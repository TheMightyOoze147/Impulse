package main

import (
	"fmt"
	"os"

	data "github.com/TheMightyOoze147/main/src/datafromfile"
	events "github.com/TheMightyOoze147/main/src/events"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Использование: go run main.go <путь к файлу>")

		return
	}

	filepath := args[1]
	var dataFromFile []string = data.ReadFile(filepath)

	pcNumber := data.ParsePCNumber(dataFromFile[0])
	open, closed := data.ParseTimeRange(dataFromFile[1])
	price := data.ParsePrice(dataFromFile[2])

	tablesDB := events.TableDatabase(pcNumber)
	eventsDB := events.EventDatabase(dataFromFile)
	clientsDB := events.ClientDatabase(eventsDB)
	queue := make([]events.Client, 0)

	fmt.Println(open.Format("15:04"))
	for event := range eventsDB {
		_, queue = events.EventProcessing(eventsDB[event], clientsDB, queue, tablesDB, open, closed)
	}

	clientsDB, tablesDB = events.ServiceClosed(clientsDB, tablesDB, closed, price)

	fmt.Println(closed.Format("15:04"))
	for id, table := range tablesDB {
		fmt.Printf("%d %d %02d:%02d\n", id+1, table.Earnings, int(table.InWork.Hours()), int(table.InWork.Minutes())%60)
	}

}
