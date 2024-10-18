package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"
)

const (
	saveFile = "tasks.json"
)

func main() {
	fmt.Println("Welcome to ToGo - the shittiest todo list of all time!")
	fmt.Println("Type 'help' to see a list of commands.")

	s := NewFileStorage(saveFile)
	err := s.Load()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(fmt.Errorf("failed to load tasks from file: %w", err))
		}
		fmt.Println("No save file found, creating new ToGo list ...")
	} else {
		fmt.Printf("Loaded ToGo list from %s\n", saveFile)
	}

	_ = s.Add("this is some random task")
	_ = s.Add("this is another task")

	input := make(chan string)
	go readCmd(input)
	for {
		fmt.Print("> ")
		select {
		case cmd := <-input:
			err := handleCmd(s, cmd)
			if err != nil {
				panic(err)
			}
		case <-time.After(60 * time.Second):
			fmt.Println("Inactive for 1 minute. Terminating ...")
			os.Exit(0)
		}
	}
}

func handleCmd(s Storage, cmd string) error {
	re := regexp.MustCompile(`"([^"]*)"|(\S+)`)
	matches := re.FindAllStringSubmatch(cmd, -1)

	var parts []string
	for _, match := range matches {
		for _, group := range match[1:] {
			if group == "" {
				continue
			}
			parts = append(parts, group)
		}
	}

	switch parts[0] {
	case "list":
		return handleCmdList(s)
	case "add":
		return handleCmdAdd(s, parts)
	case "exit":
		return handleCmdExit(s)
	default:
		return fmt.Errorf("command not supported: %s\n", cmd)
	}
}

func handleCmdList(s Storage) error {
	listTasks(s)
	return nil
}

func handleCmdAdd(s Storage, parts []string) error {
	if len(parts) != 2 {
		return fmt.Errorf("invalid cmd format: %+v", parts)
	}
	return s.Add(parts[1])
}

func handleCmdExit(s Storage) error {
	err := s.Save()
	if err != nil {
		return fmt.Errorf("failed to save task file: %w", err)
	}
	os.Exit(0)
	return nil
}

func readCmd(input chan<- string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		u, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Errorf("failed to read input: %w", err))
		}
		u = u[:len(u)-1]
		input <- u
	}
}

func listTasks(s Storage) {
	for _, task := range s.GetTasks() {
		fmt.Printf("%s: %s\n", task.ID, task.Text)
	}
}
