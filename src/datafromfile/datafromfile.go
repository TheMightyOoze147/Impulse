package datafromfile

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Чтение данных из файла, возвращает массив строк
func ReadFile(filePath string) (lines []string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}

// Парсинг номера - в будущем для обозначения количества столов
func ParsePCNumber(number string) int {
	num, err := strconv.Atoi(number)
	if err != nil {
		log.Fatal(err)
	}

	return num
}

// Парсинг времени открытия и времени закрытия
func ParseTimeRange(timeRange string) (time.Time, time.Time) {
	parts := strings.Split(timeRange, " ")
	if len(parts) != 2 {
		log.Fatal(fmt.Errorf("bad format: %s", timeRange))
	}

	firstTime, err := time.Parse("15:04", parts[0])
	if err != nil {
		log.Fatal(err)
	}

	secondTime, err := time.Parse("15:04", parts[1])
	if err != nil {
		log.Fatal(err)
	}

	return firstTime, secondTime
}

// Парсинг цены за один час работы за столом
func ParsePrice(value string) int {
	price, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}

	return price
}
