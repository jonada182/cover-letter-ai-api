package util

//go:generate mockgen -destination=../mocks/mock_util.go -package=mocks github.com/jonada182/cover-letter-ai-api/util Util
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Util interface {
	LoadEnvFile(string) error
}

// LoadEnvFile sets all the valid environment variables from the given filename
func LoadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			return fmt.Errorf("invalid line in .env file: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
