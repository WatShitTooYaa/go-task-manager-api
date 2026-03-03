package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type ArgsError struct {
	Message string
	// Err     error
}

func (ie *ArgsError) Error() string {
	return ie.Message
}

type Command struct {
	Name   string
	Regex  *regexp.Regexp
	Handle func([]string) error
}

func main() {
	fileName := "storage.json"
	storage := NewStorage(fileName)

	scanner := bufio.NewScanner(os.Stdin)

	var commands []Command = []Command{
		{
			Name:   "add",
			Regex:  regexp.MustCompile(`^add\s+"([^"]+)"$`),
			Handle: storage.handleAdd(),
		},
		{
			Name:   "list",
			Regex:  regexp.MustCompile(`^list$`),
			Handle: storage.handleList(),
		},
		{
			Name:   "update",
			Regex:  regexp.MustCompile(`^update\s+(\d+)\s+-(task|priority|status|p|t|s)\s+"([^"]+)"$`),
			Handle: storage.handleUpdate(),
		},
		{
			Name:   "delete",
			Regex:  regexp.MustCompile(`^delete\s+(\d+)$`),
			Handle: storage.handleDelete(),
		},
	}

mainLoop:
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "break" {
			break mainLoop
		}

		var rawScan []string
		var found = false
		for _, cmd := range commands {
			rawScan = cmd.Regex.FindStringSubmatch(input)
			if len(rawScan) > 0 {
				err := cmd.Handle(rawScan)
				if err != nil {
					fmt.Println("error : ", err.Error())
				}
				name := cmd.Name
				if name == "add" || name == "update" || name == "delete" {
					fmt.Printf("perintah %s berhasil dijalankan\n", name)
				}
				found = true
				break
			}
		}
		if !found {
			fmt.Println("command tidak ditemukan")
		}
	}

}
