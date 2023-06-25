package utils

import (
	"bufio"
	"os"
)

// GetFileAsLines returns the file present in the specified path, panics if the file is not available ( lines are limited to 64k )
func GetFileAsLines(path string) ([]string){
	file, err := os.Open(path)
    if err != nil {
        panic(err)
    }
	defer func() {
        if err := file.Close(); err != nil {
            panic(err)
        }
    }()
	lines := []string{}
	scanner := bufio.NewScanner(file)
    // optionally, resize scanner's capacity for lines over 64K
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        panic(err)
    }
	return lines
}

// SafeGetFileAsLines works as GetFileAsLines but is safe in case no file was found ( still panic in case of failed close file since it could result in memory leaks )
func SafeGetFileAsLines(path string) ([]string, error){
	file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
	defer func() {
        if err := file.Close(); err != nil {
            panic(err)
        }
    }()
	lines := []string{}
	scanner := bufio.NewScanner(file)
    // optionally, resize scanner's capacity for lines over 64K
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
	return lines, nil
}