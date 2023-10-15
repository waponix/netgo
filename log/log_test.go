package log

import (
	"bufio"
	"os"
	"reflect"
	"regexp"
	"testing"
)

func TestLoggerType(t *testing.T) {
	l := Logger()
	loggerType := reflect.TypeOf(l)
	expectedType := reflect.TypeOf(&logger{})

	if loggerType != expectedType {
		t.Fatalf(`Logger() func should be returning type %v returned %v`, expectedType, loggerType)
	}
}

func TestInfoShouldWrite(t *testing.T) {
	l := Logger()
	l.Filename = "test.log"

	defer os.Remove(l.Filename)

	want := regexp.MustCompile(`\b` + "." + INFO + ": This is an info log" + `\b`)

	l.Info("This is an info log")

	lastLine, err := ReadLastLine(l.Filename)

	if !want.MatchString(lastLine) || err != nil {
		t.Fatalf(`l.Info("This is an info log") = %q, %v, want match for %#q, nil`, lastLine, err, want)
	}
}

func TestInfoShouldNotWrite(t *testing.T) {
	l := Logger()
	l.Filename = "test.log"

	defer os.Remove(l.Filename)

	l.LogLevels = []string{DEBUG, NOTICE, ERROR, FATAL}

	l.Info("This is an info log")

	lastLine, _ := ReadLastLine(l.Filename)

	if fileExists(l.Filename) && lastLine != "" {
		t.Fatalf("Log file should not be created or at least be empty")
	}
}

func TestDebugShouldWrite(t *testing.T) {
	l := Logger()
	l.Filename = "test.log"

	defer os.Remove(l.Filename)

	want := regexp.MustCompile(`\b` + "." + DEBUG + ": This is a debug log" + `\b`)

	l.Debug("This is a debug log")

	lastLine, err := ReadLastLine(l.Filename)

	if !want.MatchString(lastLine) || err != nil {
		t.Fatalf(`l.Debug("This is a debug log") = %q, %v, want match for %#q, nil`, lastLine, err, want)
	}
}

func TestDebugShouldNotWrite(t *testing.T) {
	l := Logger()
	l.Filename = "test.log"

	defer os.Remove(l.Filename)

	l.LogLevels = []string{INFO, NOTICE, ERROR, FATAL}

	l.Debug("This is a debug log")

	lastLine, _ := ReadLastLine(l.Filename)

	if fileExists(l.Filename) && lastLine != "" {
		t.Fatalf("Log file should not be created or at least be empty")
	}
}

func TestNoticeShouldWrite(t *testing.T) {
	l := Logger()
	l.Filename = "test.log"

	defer os.Remove(l.Filename)

	want := regexp.MustCompile(`\b` + "." + NOTICE + ": This is a notice log" + `\b`)

	l.Notice("This is a notice log")

	lastLine, err := ReadLastLine(l.Filename)

	if !want.MatchString(lastLine) || err != nil {
		t.Fatalf(`l.Notice("This is a notice log") = %q, %v, want match for %#q, nil`, lastLine, err, want)
	}
}

func TestNoticeShouldNotWrite(t *testing.T) {
	l := Logger()
	l.Filename = "test.log"

	defer os.Remove(l.Filename)

	l.LogLevels = []string{INFO, DEBUG, ERROR, FATAL}

	l.Notice("This is a notice log")

	lastLine, _ := ReadLastLine(l.Filename)

	if fileExists(l.Filename) && lastLine != "" {
		t.Fatalf("Log file should not be created or at least be empty")
	}
}

func TestErrorShouldWrite(t *testing.T) {
	l := Logger()
	l.Filename = "test.log"

	defer os.Remove(l.Filename)

	want := regexp.MustCompile(`\b` + "." + ERROR + ": This is an error log" + `\b`)

	l.Error("This is an error log")

	lastLine, err := ReadLastLine(l.Filename)

	if !want.MatchString(lastLine) || err != nil {
		t.Fatalf(`l.Error("This is an error log") = %q, %v, want match for %#q, nil`, lastLine, err, want)
	}
}

func TestErrorShouldNotWrite(t *testing.T) {
	l := Logger()
	l.Filename = "test.log"

	defer os.Remove(l.Filename)

	l.LogLevels = []string{INFO, NOTICE, DEBUG, FATAL}

	l.Error("This is an Error log")

	lastLine, _ := ReadLastLine(l.Filename)

	if fileExists(l.Filename) && lastLine != "" {
		t.Fatalf("Log file should not be created or at least be empty")
	}
}

func TestFatalShouldWrite(t *testing.T) {
	l := Logger()
	l.Filename = "test.log"

	defer os.Remove(l.Filename)

	want := regexp.MustCompile(`\b` + "." + FATAL + ": This is a fatal log" + `\b`)

	l.Fatal("This is a fatal log")

	lastLine, err := ReadLastLine(l.Filename)

	if !want.MatchString(lastLine) || err != nil {
		t.Fatalf(`l.Fatal("This is a fatal log") = %q, %v, want match for %#q, nil`, lastLine, err, want)
	}
}

func TestFatalShouldNotWrite(t *testing.T) {
	l := Logger()
	l.Filename = "test.log"

	defer os.Remove(l.Filename)

	l.LogLevels = []string{INFO, NOTICE, DEBUG, ERROR}

	l.Fatal("This is a fatal log")

	lastLine, _ := ReadLastLine(l.Filename)

	if fileExists(l.Filename) && lastLine != "" {
		t.Fatalf("Log file should not be created or at least be empty")
	}
}

func ReadLastLine(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return lastLine, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
