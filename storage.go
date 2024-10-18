package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
)

const (
	idBytesLen = 4
)

type Task struct {
	ID          string
	Text        string
	IsCompleted bool
}

func NewTask(text string) *Task {
	return &Task{
		ID:   generateUniqueId(),
		Text: text,
	}
}

type TaskList []Task

type Storage interface {
	Load() error
	Save() error
	Find(taskId string) (*Task, error)
	Add(text string) error
	GetTasks() TaskList
}

type FileStorage struct {
	filename string
	list     *TaskList
}

func NewFileStorage(filename string) *FileStorage {
	return &FileStorage{
		filename: filename,
		list:     &TaskList{},
	}
}

func (s *FileStorage) Load() error {
	file, err := os.Open(s.filename)
	if err != nil {
		return fmt.Errorf("failed to open task file: %w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(s.list)
	if err != nil {
		return fmt.Errorf("failed to decode task list: %w", err)
	}

	return nil
}

func (s *FileStorage) Save() error {
	file, err := os.Create(s.filename)
	if err != nil {
		return fmt.Errorf("failed to create task file: %w", err)
	}

	e := json.NewEncoder(file)
	e.SetIndent("", "  ")
	err = e.Encode(s.list)
	if err != nil {
		return fmt.Errorf("failed to encode task list: %w", err)
	}

	return nil
}

func (s *FileStorage) Find(taskId string) (*Task, error) {
	for _, task := range *s.list {
		if task.ID == taskId {
			return &task, nil
		}
	}
	return nil, fmt.Errorf("no task with id: %s", taskId)
}

func (s *FileStorage) Add(text string) error {
	newTask := NewTask(text)

	if task, _ := s.Find(newTask.ID); task != nil {
		return fmt.Errorf("task with id %s already exists", newTask.ID)
	}

	*s.list = append(*s.list, *newTask)
	return nil
}

func (s *FileStorage) GetTasks() TaskList {
	return *s.list
}

func generateUniqueId() string {
	id := make([]byte, idBytesLen)
	_, err := rand.Read(id)
	if err != nil {
		panic(fmt.Errorf("failed to generate id: %w", err))
	}
	return fmt.Sprintf("%x", id)
}
