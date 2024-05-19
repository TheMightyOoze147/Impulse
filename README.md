# Тестовое задание на стажировку Импульс 2024

Версия GO: 1.19. Запуск производился на ОС Ubuntu.
### Запуск в контейнере

[![dockgo](https://skillicons.dev/icons?i=docker,go)](https://skillicons.dev)

Необходимо собрать docker контейнер. Для этого, находясь в директории с Dockerfile, ввести в терминале: 

```bash
docker build -t app .
```

По завершении сборки необходимо запустить контейнер, передав ему флаг -v, предназначенный для обмена данными между контейнером и хостом.
Запуск производить из корневой папки проекта.

```bash
docker run --rm -v $(pwd)/scenaries/scenario.txt:/app/scenaries/scenario.txt app ./scenaries/scenario.txt
```

В данном случае аргументом командной строки является 
```bash
./scenaries/scenario.txt
```

Обозначающий путь к файлу.

<details>
    <summary>Вывод</summary>
    
    09:00
    08:48 1 client1
    08:48 13 NotOpenYet
    09:41 1 client1
    09:48 1 client2
    09:52 3 client1
    09:52 13 ICanWaitNoLonger!
    09:54 2 client1 1
    10:25 2 client2 2
    10:58 1 client3
    10:59 2 client3 3
    11:30 1 client4
    11:35 2 client4 2
    11:35 13 PlaceIsBusy
    11:45 3 client4
    12:33 4 client1
    12:33 12 client4 1
    12:43 4 client2
    15:52 4 client4
    19:00 11 client3
    19:00
    1 70 05:58
    2 30 02:18
    3 90 08:01

</details>

### Запуск без контейнера

[![go](https://skillicons.dev/icons?i=go)](https://skillicons.dev)

Программа на вход принимает путь к файлу-сценарию с расширением .txt. 
Пример: 
```bash
go run src/main/main scenaries/scenario.txt
```
Это запустит программу.

Либо же использовать команду для сборки проекта в исходный файл: 
```bash
go build -o ./build/main src/main/main.go
```
Это создаст исполняемый файл main.exe, который можно запустить с помощью: 
```bash
./build/main scenaries/scenario.txt
```

<details>
    <summary>Вывод</summary>
    
        09:00
        08:48 1 client1
        08:48 13 NotOpenYet
        09:41 1 client1
        09:48 1 client2
        09:52 3 client1
        09:52 13 ICanWaitNoLonger!
        09:54 2 client1 1
        10:25 2 client2 2
        10:58 1 client3
        10:59 2 client3 3
        11:30 1 client4
        11:35 2 client4 2
        11:35 13 PlaceIsBusy
        11:45 3 client4
        12:33 4 client1
        12:33 12 client4 1
        12:43 4 client2
        15:52 4 client4
        19:00 11 client3
        19:00
        1 70 05:58
        2 30 02:18
        3 90 08:01

</details>

### Результаты тестов

<details>
    <summary>Пакет datafromfile</summary>
    
        go test -v -cover
        === RUN   TestReadFile
        --- PASS: TestReadFile (0.00s)
        === RUN   TestParsePCNumber
        --- PASS: TestParsePCNumber (0.00s)
        === RUN   TestParseTimeRange
        --- PASS: TestParseTimeRange (0.00s)
        === RUN   TestParsePrice
        --- PASS: TestParsePrice (0.00s)
        PASS
        coverage: 75.9% of statements
        ok      github.com/TheMightyOoze147/main/src/datafromfile       0.002s

</details>

<details>
    <summary>Пакет events</summary>

        go test -v -cover
        === RUN   TestNewEvent
        --- PASS: TestNewEvent (0.00s)
        === RUN   TestEventDatabase
        --- PASS: TestEventDatabase (0.00s)
        === RUN   TestFindFreeTables
        --- PASS: TestFindFreeTables (0.00s)
        === RUN   TestTableDatabase
        --- PASS: TestTableDatabase (0.00s)
        === RUN   TestClientDatabase
        --- PASS: TestClientDatabase (0.00s)
        === RUN   TestFindInClientDBByName
        --- PASS: TestFindInClientDBByName (0.00s)
        === RUN   TestFindInClientDBByTableNumber
        --- PASS: TestFindInClientDBByTableNumber (0.00s)
        === RUN   TestClientArrived
        12:00 1 John
        --- PASS: TestClientArrived (0.00s)
        === RUN   TestClientTakeASeat
        12:00 2 John 1
        --- PASS: TestClientTakeASeat (0.00s)
        === RUN   TestClientIsWaiting
        12:00 3 John
        12:00 13 ICanWaitNoLonger!
        --- PASS: TestClientIsWaiting (0.00s)
        === RUN   TestClientLeaved
        12:00 4 John
        --- PASS: TestClientLeaved (0.00s)
        === RUN   TestEventProcessing
        12:00 1 John
        --- PASS: TestEventProcessing (0.00s)
        === RUN   TestServiceClosed
        14:00 11 Jane
        14:00 11 John
        --- PASS: TestServiceClosed (0.00s)
        PASS
        coverage: 67.5% of statements
        ok      github.com/TheMightyOoze147/main/src/events     0.005s

</details>
