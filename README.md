# Тестовое задание на стажировку Импульс 2024
## Вакансия: Инженер по разработке ПО для базовых станций (Go)
### Запуск в контейнере
Необходимо собрать докер контейнер. Для этого, находясь в директории с Dockerfile, ввести в терминале: 
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
    ```bash
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
    ```

</details>

### Запуск без контейнера
Программа на вход принимает путь к файлу-сценарию с расширением .txt. 
Пример: go run src/main/main scenaries/scenario.txt

<details>
    <summary>Вывод</summary>
    ```bash
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
    ```

</details>