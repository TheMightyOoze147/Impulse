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

func NewEvent(event string) Event {
	parts := strings.Split(event, " ")
	if len(parts) < 3 || len(parts) > 4 {
		log.Fatal(fmt.Errorf("Bad format: %s", event))
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

func EventDatabase(data []string) (eventsPool []Event) {
	for i := 3; i <= len(data)-1; i++ {
		eventsPool = append(eventsPool, NewEvent(data[i]))
	}

	return
}

func NewTable(num string) Table {
	return Table{
		Number: num,
		IsBusy: false,
	}
}

func FindFreeTables(tablesDB []Table) (Table, error) {
	for _, table := range tablesDB {
		if table.IsBusy == false {
			return table, nil
		}
	}

	return Table{}, errors.New("no free tables")
}

func TableDatabase(tablesNum int) (tablesDB []Table) {
	for i := range tablesNum {
		tablesDB = append(tablesDB, NewTable(strconv.Itoa(i+1)))
	}

	return
}

func NewClient(event Event) Client {
	return Client{
		Name:      event.Name,
		IsVisited: false,
		Table:     Table{"0", false, time.Hour, 0},
	}
}

func FindInClientDBByName(name string, DB []Client) (Client, int, error) {
	for id, client := range DB {
		if client.Name == name {
			return client, id, nil
		}
	}

	return Client{}, -1, errors.New("client not found")
}

func FindInClientDBByTableNumber(tableNum string, DB []Client) (Client, int, error) {
	for id, client := range DB {
		if client.Table.Number == tableNum {
			return client, id, nil
		}
	}

	return Client{}, -1, errors.New("table not found")
}

func ClientDatabase(eventsPool []Event) (clientsPool []Client) {
	for _, event := range eventsPool {
		_, _, err := FindInClientDBByName(event.Name, clientsPool)
		if err != nil {
			clientsPool = append(clientsPool, NewClient(event))
		}
	}

	return
}

func ClientArrived(event Event, client Client, startTime time.Time, endTime time.Time) Client {
	fmt.Println(event.ArrivedTime.Format("15:04"), event.Event, event.Name)
	if event.ArrivedTime.After(startTime) && event.ArrivedTime.Before(endTime) && client.IsVisited == false {
		client.ArrivedTime = event.ArrivedTime
		client.IsVisited = true
	} else if client.IsVisited == true {
		fmt.Println(event.ArrivedTime.Format("15:04"), "13", "YouShallNotPass")
	} else {
		fmt.Println(event.ArrivedTime.Format("15:04"), "13", "NotOpenYet")
	}

	return client
}

func ClientTakeASeat(event Event, client Client, tablesDB []Table, isFree error) Client {
	fmt.Println(event.ArrivedTime.Format("15:04"), event.Event, event.Name, event.TableNum)
	if client.IsVisited == true {
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

func ClientIsWaiting(event Event, client Client, queue []Client, tablesDB []Table) (Client, error) {
	fmt.Println(event.ArrivedTime.Format("15:04"), event.Event, event.Name)
	if len(queue) >= len(tablesDB) {
		fmt.Println(event.ArrivedTime.Format("15:04"), "11", event.Name)
		client.IsVisited = false
		client.LeavedTime = event.ArrivedTime
	} else if client.IsVisited == false {
		fmt.Println(event.ArrivedTime.Format("15:04"), "13", "ClientUnknown")
	} else {
		_, err := FindFreeTables(tablesDB)
		if err == nil {
			fmt.Println(event.ArrivedTime.Format("15:04"), "13", "ICanWaitNoLonger!")
		} else {
			return client, nil
		}
	}

	return Client{}, errors.New("cant add client to queue")
}

func ClientLeaved(event Event, client Client, clientsDB []Client, queue []Client, tablesDB []Table) (Client, []Client) {
	fmt.Println(event.ArrivedTime.Format("15:04"), event.Event, event.Name)
	if client.IsVisited == false {
		fmt.Println(event.ArrivedTime.Format("15:04"), "13", "ClientUnknown")
	} else {
		if len(queue) != 0 {
			fmt.Println(event.ArrivedTime.Format("15:04"), "12", queue[0].Name, client.Table.Number)
			queue[0].ArrivedTime = event.ArrivedTime
			queue[0].Table = client.Table

			client.LeavedTime = event.ArrivedTime
			client.IsVisited = false
		} else {
			tableNum, _ := strconv.Atoi(client.Table.Number)
			tablesDB[tableNum-1].IsBusy = false
			client.LeavedTime = event.ArrivedTime
			client.IsVisited = false
		}
	}

	return client, queue
}

func EventProcessing(event Event, clientsDB []Client, queue []Client, tablesDB []Table, startTime time.Time, endTime time.Time) (Client, []Client) {
	switch event.Event {
	case "1":
		findedClient, id, _ := FindInClientDBByName(event.Name, clientsDB)
		eventResponse := ClientArrived(event, findedClient, startTime, endTime)
		clientsDB[id] = eventResponse

		return eventResponse, queue
	case "2":
		findedClient, id, _ := FindInClientDBByName(event.Name, clientsDB)
		_, _, IsBusy := FindInClientDBByTableNumber(event.TableNum, clientsDB)
		eventResponse := ClientTakeASeat(event, findedClient, tablesDB, IsBusy)
		clientsDB[id] = eventResponse

		return eventResponse, queue
	case "3":
		findedClient, _, _ := FindInClientDBByName(event.Name, clientsDB)
		eventResponse, err := ClientIsWaiting(event, findedClient, queue, tablesDB)
		if err == nil {
			queue = append(queue, eventResponse)
		}

		return eventResponse, queue
	case "4":
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

func ServiceClosed(clientsDB []Client, tablesDB []Table, endTime time.Time, price int) ([]Client, []Table) {
	sort.Slice(clientsDB, func(i, j int) bool {
		return clientsDB[i].Name < clientsDB[j].Name
	})

	for i, client := range clientsDB {
		if client.IsVisited {
			fmt.Println(endTime.Format("15:04"), "11", client.Name)
			client.IsVisited = false
			client.LeavedTime = endTime
		}

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

	return clientsDB, tablesDB
}
