package sqlmigrator

import (
	"bufio"
	"os"
)

func CreateFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if content == "" {
		return nil
	}

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	return lines, nil
}
