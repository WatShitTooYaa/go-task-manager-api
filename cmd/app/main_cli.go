package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/storage"
)

type Command struct {
	Name   string
	Regex  *regexp.Regexp
	Handle func([]string) error
}

func main_cli() {
	fileName := "storage.json"
	storage := storage.NewStorage(fileName)

	scanner := bufio.NewScanner(os.Stdin)

	var commands []Command = []Command{
		{
			Name:   "add",
			Regex:  regexp.MustCompile(`^add\s+"([^"]+)"$`),
			Handle: storage.HandleAdd(),
		},
		{
			Name:   "list",
			Regex:  regexp.MustCompile(`^list$`),
			Handle: storage.HandleList(),
		},
		{
			Name:   "update",
			Regex:  regexp.MustCompile(`^update\s+(\d+)\s+-(task|priority|status|p|t|s)\s+"([^"]+)"$`),
			Handle: storage.HandleUpdate(),
		},
		{
			Name:   "delete",
			Regex:  regexp.MustCompile(`^delete\s+(\d+)$`),
			Handle: storage.HandleDelete(),
		},
		{
			Name:   "import",
			Regex:  regexp.MustCompile(`^import\s+"([^"]+)"$`),
			Handle: storage.ImportCSV(),
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
			// fmt.Println("rawscan : ", rawScan)
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
