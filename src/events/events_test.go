package events

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestNewEvent(t *testing.T) {
	eventStr := "12:00 arrival John"
	expectedEvent := Event{
		ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
		Event:       "arrival",
		Name:        "John",
		TableNum:    "0",
	}

	event := NewEvent(eventStr)

	if !reflect.DeepEqual(event, expectedEvent) {
		t.Errorf("NewEvent(%q) = %v, want %v", eventStr, event, expectedEvent)
	}
}

func TestEventDatabase(t *testing.T) {
	data := []string{"1", "2", "3", "12:00 1 John", "12:05 1 Jane"}
	expectedEventsPool := []Event{
		{
			ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
			Event:       "1",
			Name:        "John",
			TableNum:    "0",
		},
		{
			ArrivedTime: time.Date(0, 1, 1, 12, 5, 0, 0, time.UTC),
			Event:       "1",
			Name:        "Jane",
			TableNum:    "0",
		},
	}

	eventsPool := EventDatabase(data)

	if !reflect.DeepEqual(eventsPool, expectedEventsPool) {
		t.Errorf("EventDatabase(%v) = %v, want %v", data, eventsPool, expectedEventsPool)
	}
}

func TestFindFreeTables(t *testing.T) {
	tablesDB := []Table{
		{Number: "1", IsBusy: false},
		{Number: "2", IsBusy: true},
	}
	expectedTable := Table{Number: "1", IsBusy: false}

	table, err := FindFreeTables(tablesDB)

	if err != nil || !reflect.DeepEqual(table, expectedTable) {
		t.Errorf("FindFreeTables(%v) = %v, %v, want %v, nil", tablesDB, table, err, expectedTable)
	}
}

func TestTableDatabase(t *testing.T) {
	tablesNum := 2
	expectedTablesDB := []Table{
		{Number: "1", IsBusy: false},
		{Number: "2", IsBusy: false},
	}

	tablesDB := TableDatabase(tablesNum)

	if !reflect.DeepEqual(tablesDB, expectedTablesDB) {
		t.Errorf("TableDatabase(%d) = %v, want %v", tablesNum, tablesDB, expectedTablesDB)
	}
}

func TestClientDatabase(t *testing.T) {
	eventsPool := []Event{
		{
			ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
			Event:       "1",
			Name:        "John",
			TableNum:    "0",
		},
		{
			ArrivedTime: time.Date(0, 1, 1, 12, 5, 0, 0, time.UTC),
			Event:       "1",
			Name:        "Jane",
			TableNum:    "0",
		},
	}
	expectedClientsPool := []Client{
		{
			Name:      "John",
			IsVisited: false,
			Table:     Table{"0", false, time.Hour, 0},
		},
		{
			Name:      "Jane",
			IsVisited: false,
			Table:     Table{"0", false, time.Hour, 0},
		},
	}

	clientsPool := ClientDatabase(eventsPool)

	if !reflect.DeepEqual(clientsPool, expectedClientsPool) {
		t.Errorf("ClientDatabase(%v) = %v, want %v", eventsPool, clientsPool, expectedClientsPool)
	}
}

func TestFindInClientDBByName(t *testing.T) {
	clientsDB := []Client{
		{Name: "John", IsVisited: false},
		{Name: "Jane", IsVisited: false},
	}
	expectedClient := clientsDB[0]

	client, index, err := FindInClientDBByName("John", clientsDB)

	if err != nil || !reflect.DeepEqual(client, expectedClient) || index != 0 {
		t.Errorf("FindInClientDBByName(%q, %v) = %v, %d, %v, want %v, 0, nil", "John", clientsDB, client, index, err, expectedClient)
	}
}

func TestFindInClientDBByTableNumber(t *testing.T) {
	clientsDB := []Client{
		{Name: "John", IsVisited: false, Table: Table{Number: "1"}},
		{Name: "Jane", IsVisited: false, Table: Table{Number: "2"}},
	}
	expectedClient := clientsDB[0]

	client, index, err := FindInClientDBByTableNumber("1", clientsDB)

	if err != nil || !reflect.DeepEqual(client, expectedClient) || index != 0 {
		t.Errorf("FindInClientDBByTableNumber(%q, %v) = %v, %d, %v, want %v, 0, nil", "1", clientsDB, client, index, err, expectedClient)
	}
}

func TestClientArrived(t *testing.T) {
	event := Event{
		ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
		Event:       "1",
		Name:        "John",
		TableNum:    "0",
	}
	client := Client{
		Name:      "John",
		IsVisited: false,
		Table:     Table{"0", false, time.Hour, 0},
	}
	startTime := time.Date(0, 1, 1, 11, 59, 0, 0, time.UTC)
	endTime := time.Date(0, 1, 1, 12, 1, 0, 0, time.UTC)

	expectedClient := Client{
		Name:        "John",
		IsVisited:   true,
		ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
		Table:       Table{"0", false, time.Hour, 0},
	}

	client = ClientArrived(event, client, startTime, endTime)

	if !reflect.DeepEqual(client, expectedClient) {
		t.Errorf("ClientArrived(%v, %v, %v, %v) = %v, want %v", event, client, startTime, endTime, client, expectedClient)
	}
}

func TestClientTakeASeat(t *testing.T) {
	event := Event{
		ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
		Event:       "2",
		Name:        "John",
		TableNum:    "1",
	}
	client := Client{
		Name:      "John",
		IsVisited: true,
		Table:     Table{"0", false, time.Hour, 0},
	}
	tablesDB := []Table{
		{Number: "1", IsBusy: false},
		{Number: "2", IsBusy: true},
	}

	expectedClient := Client{
		Name:        "John",
		IsVisited:   true,
		ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
		Table:       Table{"1", true, time.Duration(0), 0},
	}
	clientsDB := make([]Client, 0)
	clientsDB = append(clientsDB, client)
	_, _, IsBusy := FindInClientDBByTableNumber(event.TableNum, clientsDB)

	client = ClientTakeASeat(event, client, tablesDB, IsBusy)

	if !reflect.DeepEqual(client, expectedClient) {
		t.Errorf("ClientTakeASeat, %v, %v, %v, %v) = %v, want %v", event, client, tablesDB, nil, client, expectedClient)
	}
}

func TestClientIsWaiting(t *testing.T) {
	event := Event{
		ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
		Event:       "3",
		Name:        "John",
		TableNum:    "0",
	}
	client := Client{
		Name:      "John",
		IsVisited: true,
		Table:     Table{"0", false, time.Duration(0), 0},
	}
	queue := []Client{}
	tablesDB := []Table{
		{Number: "1", IsBusy: false},
		{Number: "2", IsBusy: false},
		{Number: "3", IsBusy: false},
		{Number: "4", IsBusy: false},
	}

	expectedClient := Client{
		Name:      "John",
		IsVisited: true,
		Table:     Table{"0", false, time.Duration(0), 0},
	}

	client, err := ClientIsWaiting(event, client, queue, tablesDB)

	if err == nil || reflect.DeepEqual(client, expectedClient) {
		t.Errorf("ClientIsWaiting(%v, %v, %v, %v) = %v, %v, want %v, %v", event, client, queue, tablesDB, client, err, expectedClient, errors.New("cant add client to queue"))
	}
}

func TestClientLeaved(t *testing.T) {
	event := Event{
		ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
		Event:       "4",
		Name:        "John",
		TableNum:    "0",
	}
	client := Client{
		Name:        "John",
		IsVisited:   true,
		ArrivedTime: time.Date(0, 1, 1, 11, 59, 0, 0, time.UTC),
		Table:       Table{"1", true, time.Hour, 0},
	}
	clientsDB := []Client{client}
	queue := []Client{}
	tablesDB := []Table{
		{Number: "1", IsBusy: true},
		{Number: "2", IsBusy: true},
	}
	expectedClient := Client{
		Name:        "John",
		IsVisited:   false,
		ArrivedTime: time.Date(0, 1, 1, 11, 59, 0, 0, time.UTC),
		LeavedTime:  time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
		Table:       Table{"1", true, time.Hour, 0},
	}

	client, queue = ClientLeaved(event, client, clientsDB, queue, tablesDB)

	if !reflect.DeepEqual(client, expectedClient) {
		t.Errorf("ClientLeaved(%v, %v, %v, %v, %v) = %v, %v, want %v, []Client{}", event, client, clientsDB, queue, tablesDB, client, queue, expectedClient)
	}
}

func TestEventProcessing(t *testing.T) {
	event := Event{
		ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
		Event:       "1",
		Name:        "John",
		TableNum:    "0",
	}
	clientsDB := []Client{
		{Name: "John", IsVisited: false},
	}
	queue := []Client{}
	tablesDB := []Table{
		{Number: "1", IsBusy: false},
		{Number: "2", IsBusy: true},
	}
	startTime := time.Date(0, 1, 1, 11, 59, 0, 0, time.UTC)
	endTime := time.Date(0, 1, 1, 12, 1, 0, 0, time.UTC)
	expectedClient := Client{
		Name:        "John",
		IsVisited:   true,
		ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
		Table:       Table{"", false, time.Duration(0), 0},
	}

	client, queue := EventProcessing(event, clientsDB, queue, tablesDB, startTime, endTime)

	if !reflect.DeepEqual(client, expectedClient) {
		t.Errorf("EventProcessing(%v, %v, %v, %v, %v, %v) = %v, %v, want %v, []Client{}", event, clientsDB, queue, tablesDB, startTime, endTime, client, queue, expectedClient)
	}
}

func TestServiceClosed(t *testing.T) {
	clientsDB := []Client{
		{
			Name:        "John",
			IsVisited:   true,
			ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
			LeavedTime:  time.Date(0, 1, 1, 13, 0, 0, 0, time.UTC),
			Table:       Table{"1", true, time.Duration(0), 0},
		},
		{
			Name:        "Jane",
			IsVisited:   true,
			ArrivedTime: time.Date(0, 1, 1, 12, 00, 0, 0, time.UTC),
			LeavedTime:  time.Date(0, 1, 1, 14, 00, 0, 0, time.UTC),
			Table:       Table{"2", true, time.Duration(0), 0},
		},
	}
	tablesDB := []Table{
		{Number: "1", IsBusy: true},
		{Number: "2", IsBusy: true},
	}
	endTime := time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)
	price := 10
	expectedClientsDB := []Client{
		{
			Name:        "Jane",
			IsVisited:   false,
			ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
			LeavedTime:  time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC),
			Table:       Table{"2", false, 2 * time.Hour, 20},
		},
		{
			Name:        "John",
			IsVisited:   false,
			ArrivedTime: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
			LeavedTime:  time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC),
			Table:       Table{"1", false, 2 * time.Hour, 20},
		},
	}
	expectedTablesDB := []Table{
		{Number: "1", IsBusy: false, InWork: 2 * time.Hour, Earnings: 20},
		{Number: "2", IsBusy: false, InWork: 2 * time.Hour, Earnings: 20},
	}

	clientsDB, tablesDB = ServiceClosed(clientsDB, tablesDB, endTime, price)

	if !reflect.DeepEqual(clientsDB, expectedClientsDB) || !reflect.DeepEqual(tablesDB, expectedTablesDB) {
		t.Errorf("ServiceClosed(%v, %v, %v, %d) = %v, %v, want %v, %v", clientsDB, tablesDB, endTime, price, clientsDB, tablesDB, expectedClientsDB, expectedTablesDB)
	}
}
