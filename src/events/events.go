package events

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	ArrivedTime time.Time
	Event       string
	Name        string
	TableNum    string
}

type Table struct {
	Number   string
	IsBusy   bool
	InWork   time.Duration
	Earnings int
}

type Client struct {
	Name        string
	ArrivedTime time.Time
	LeavedTime  time.Time
	IsVisited   bool
	Table       Table
}

// Функция для разбития строки ивента в структуру Event с необходимыми полями
func NewEvent(event string) Event {
	parts := strings.Split(event, " ")
	if len(parts) < 3 || len(parts) > 4 {
		log.Fatal(fmt.Errorf("bad format: %s", event))
	}

	arrivedTime, err := time.Parse("15:04", parts[0])
	if err != nil {
		log.Fatal(err)
	}

	if len(parts) == 4 {
		return Event{
			ArrivedTime: arrivedTime,
			Event:       parts[1],
			Name:        parts[2],
			TableNum:    parts[3],
		}
	} else {
		return Event{
			ArrivedTime: arrivedTime,
			Event:       parts[1],
			Name:        parts[2],
			TableNum:    "0",
		}
	}
}

// Функция для создания слайса Ивентов
func EventDatabase(data []string) (eventsPool []Event) {
	for i := 3; i <= len(data)-1; i++ {
		eventsPool = append(eventsPool, NewEvent(data[i]))
	}

	return
}

// "Конструктор" Стола
func NewTable(num string) Table {
	return Table{
		Number: num,
		IsBusy: false,
	}
}

// Функция для обхода слайса Столов и поиска первого свободного
func FindFreeTables(tablesDB []Table) (Table, error) {
	for _, table := range tablesDB {
		if !table.IsBusy {
			return table, nil
		}
	}

	return Table{}, errors.New("no free tables")
}

// Функция для создания слайса столов
func TableDatabase(tablesNum int) (tablesDB []Table) {
	for i := 0; i < tablesNum; i++ {
		tablesDB = append(tablesDB, NewTable(strconv.Itoa(i+1)))
	}

	return
}

// "Конструктор" Клиента
func NewClient(event Event) Client {
	return Client{
		Name:      event.Name,
		IsVisited: false,
		Table:     Table{"0", false, time.Hour, 0},
	}
}

// Поиск Клиента в слайсе по имени (так же возвращает его индекс)
func FindInClientDBByName(name string, DB []Client) (Client, int, error) {
	for id, client := range DB {
		if client.Name == name {
			return client, id, nil
		}
	}

	return Client{}, -1, errors.New("client not found")
}

// Поиск Клиента в слайсе по номеру Стола (так же возвращает его индекс)
func FindInClientDBByTableNumber(tableNum string, DB []Client) (Client, int, error) {
	for id, client := range DB {
		if client.Table.Number == tableNum {
			return client, id, nil
		}
	}

	return Client{}, -1, errors.New("table not found")
}

// Создание базы данных Клиентов
func ClientDatabase(eventsPool []Event) (clientsPool []Client) {
	for _, event := range eventsPool {
		_, _, err := FindInClientDBByName(event.Name, clientsPool)
		if err != nil {
			clientsPool = append(clientsPool, NewClient(event))
		}
	}

	return
}

// Событие "Клиент пришёл"
func ClientArrived(event Event, client Client, startTime time.Time, endTime time.Time) Client {
	fmt.Println(event.ArrivedTime.Format("15:04"), event.Event, event.Name)
	// Если клиент пришёл во временной промежуток от открытия до закрятия, а так же если он ещё не приходил ранее
	if event.ArrivedTime.After(startTime) && event.ArrivedTime.Before(endTime) && !client.IsVisited {
		client.ArrivedTime = event.ArrivedTime
		client.IsVisited = true
	} else if client.IsVisited {
		fmt.Println(event.ArrivedTime.Format("15:04"), "13", "YouShallNotPass")
	} else {
		fmt.Println(event.ArrivedTime.Format("15:04"), "13", "NotOpenYet")
	}

	return client
}

// Событие "Клиент занял стол"
func ClientTakeASeat(event Event, client Client, tablesDB []Table, isFree error) Client {
	fmt.Println(event.ArrivedTime.Format("15:04"), event.Event, event.Name, event.TableNum)
	// Проверяем, пришёл ли клиент в клуб
	if client.IsVisited {
		// Проверяем, не занято ли место, которое он хочет занять
		if isFree != nil {
			tableID, _ := strconv.Atoi(event.TableNum)
			tablesDB[tableID-1].IsBusy = true
			client.ArrivedTime = event.ArrivedTime
			client.Table = tablesDB[tableID-1]
		} else {
			fmt.Println(event.ArrivedTime.Format("15:04"), "13", "PlaceIsBusy")
		}
	} else {
		fmt.Println(event.ArrivedTime.Format("15:04"), "13", "ClientUnknown")
	}

	return client
}

// Событие "Клиент ожидает"
func ClientIsWaiting(event Event, client Client, queue []Client, tablesDB []Table) (Client, error) {
	fmt.Println(event.ArrivedTime.Format("15:04"), event.Event, event.Name)
	// Проверяем, не больше ли "ждунов" в очереди, чем мест в заведении
	if len(queue) >= len(tablesDB) {
		fmt.Println(event.ArrivedTime.Format("15:04"), "11", event.Name)
		client.IsVisited = false
		client.LeavedTime = event.ArrivedTime
	} else if !client.IsVisited { // Проверка на то, пришёл ли Клиент
		fmt.Println(event.ArrivedTime.Format("15:04"), "13", "ClientUnknown")
	} else {
		_, err := FindFreeTables(tablesDB) // Если есть свободные столы, появляется ошибка
		if err == nil {
			fmt.Println(event.ArrivedTime.Format("15:04"), "13", "ICanWaitNoLonger!")
		} else {
			return client, nil
		}
	}

	return Client{}, errors.New("cant add client to queue")
}

// Событие "Клиент уходит"
func ClientLeaved(event Event, client Client, clientsDB []Client, queue []Client, tablesDB []Table) (Client, []Client) {
	fmt.Println(event.ArrivedTime.Format("15:04"), event.Event, event.Name)
	if !client.IsVisited { // Проверка, был ли клиент в заведении
		fmt.Println(event.ArrivedTime.Format("15:04"), "13", "ClientUnknown")
	} else {
		if len(queue) != 0 { // Если в очереди кто-то есть, он займёт освободившееся место
			fmt.Println(event.ArrivedTime.Format("15:04"), "12", queue[0].Name, client.Table.Number)
			queue[0].ArrivedTime = event.ArrivedTime
			queue[0].Table = client.Table

			client.LeavedTime = event.ArrivedTime
			client.IsVisited = false
		} else { // Если в очереди никого, то освободившееся место будет пустым
			tableNum, _ := strconv.Atoi(client.Table.Number)
			tablesDB[tableNum-1].IsBusy = false
			client.LeavedTime = event.ArrivedTime
			client.IsVisited = false
		}
	}

	return client, queue
}

// Глобальная обработка событий
func EventProcessing(event Event, clientsDB []Client, queue []Client, tablesDB []Table, startTime time.Time, endTime time.Time) (Client, []Client) {
	switch event.Event {
	case "1": // "Клиент пришёл"
		findedClient, id, _ := FindInClientDBByName(event.Name, clientsDB)
		eventResponse := ClientArrived(event, findedClient, startTime, endTime)
		clientsDB[id] = eventResponse

		return eventResponse, queue
	case "2": // "Клиент занял место"
		findedClient, id, _ := FindInClientDBByName(event.Name, clientsDB)
		_, _, IsBusy := FindInClientDBByTableNumber(event.TableNum, clientsDB)
		eventResponse := ClientTakeASeat(event, findedClient, tablesDB, IsBusy)
		clientsDB[id] = eventResponse

		return eventResponse, queue
	case "3": // "Клиент ожидает"
		findedClient, _, _ := FindInClientDBByName(event.Name, clientsDB)
		eventResponse, err := ClientIsWaiting(event, findedClient, queue, tablesDB)
		if err == nil {
			queue = append(queue, eventResponse)
		}

		return eventResponse, queue
	case "4": // "Клиент ушёл"
		findedClient, id, _ := FindInClientDBByName(event.Name, clientsDB)
		eventResponse, queue := ClientLeaved(event, findedClient, clientsDB, queue, tablesDB)
		clientsDB[id] = eventResponse
		if len(queue) != 0 {
			_, id, _ := FindInClientDBByName(queue[0].Name, clientsDB)
			clientsDB[id] = queue[0]
			queue = queue[1:]
		}

		return eventResponse, queue
	default:

		return Client{}, queue
	}
}

// Событие "Рабочий день окончен"
func ServiceClosed(clientsDB []Client, tablesDB []Table, endTime time.Time, price int) ([]Client, []Table) {
	// Сортировка клиентов
	sort.Slice(clientsDB, func(i, j int) bool {
		return clientsDB[i].Name < clientsDB[j].Name
	})

	// Все оставшиеся клиенты уходят
	for i, client := range clientsDB {
		if client.IsVisited {
			fmt.Println(endTime.Format("15:04"), "11", client.Name)
			client.IsVisited = false
			client.LeavedTime = endTime
		}

		// Освобождение столов и подсчёт прибыли
		if client.Table.IsBusy {
			intTableNumber, _ := strconv.Atoi(client.Table.Number)
			table := &tablesDB[intTableNumber-1]

			duration := client.LeavedTime.Sub(client.ArrivedTime)
			table.InWork += duration
			table.IsBusy = false

			hours := int(duration.Hours())
			minutes := int(duration.Minutes()) % 60
			if minutes > 0 {
				hours++
			}
			table.Earnings += hours * price

			client.Table = *table
			clientsDB[i] = client
		}
	}

	return clientsDB, tablesDB
}
