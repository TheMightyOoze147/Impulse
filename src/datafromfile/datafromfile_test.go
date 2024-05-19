package datafromfile

import (
	"os"
	"testing"
	"time"
)

func TestReadFile(t *testing.T) {
	// Создаем временный файл для тестирования
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Записываем данные в файл
	_, err = tmpFile.WriteString("test line 1\ntest line 2\ntest line 3")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Тестируем функцию ReadFile
	lines := ReadFile(tmpFile.Name())
	expectedLines := []string{"test line 1", "test line 2", "test line 3"}
	if len(lines) != len(expectedLines) {
		t.Errorf("ReadFile() returned %d lines, expected %d", len(lines), len(expectedLines))
	}
	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("ReadFile() line %d = %q, expected %q", i, line, expectedLines[i])
		}
	}
}

func TestParsePCNumber(t *testing.T) {
	number := "1234"
	expected := 1234
	result := ParsePCNumber(number)
	if result != expected {
		t.Errorf("ParsePCNumber(%q) = %d, expected %d", number, result, expected)
	}
}

func TestParseTimeRange(t *testing.T) {
	timeRange := "12:00 13:00"
	expectedStartTime, _ := time.Parse("15:04", "12:00")
	expectedEndTime, _ := time.Parse("15:04", "13:00")
	start, end := ParseTimeRange(timeRange)
	if !start.Equal(expectedStartTime) || !end.Equal(expectedEndTime) {
		t.Errorf("ParseTimeRange(%q) = (%v, %v), expected (%v, %v)", timeRange, start, end, expectedStartTime, expectedEndTime)
	}
}

func TestParsePrice(t *testing.T) {
	value := "100"
	expected := 100
	result := ParsePrice(value)
	if result != expected {
		t.Errorf("ParsePrice(%q) = %d, expected %d", value, result, expected)
	}
}
