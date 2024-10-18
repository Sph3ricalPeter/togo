package main

import (
	"os"
	"reflect"
	"testing"
)

const (
	testFilename = "tmp/test_tasks.json"
)

func TestFileStorage(t *testing.T) {
	setupTest(t)
	defer teardownTest(t)

	expTasks := TaskList{
		{ID: "1", Text: "some task"},
		{ID: "2", Text: "some other task"},
	}

	s := NewFileStorage(testFilename)
	s.list = &expTasks

	err := s.Save()
	if err != nil {
		t.Fatalf("failed to save tasks to file: %+v", err)
	}

	err = s.Load()
	if err != nil {
		t.Fatalf("failed to load tasks from file: %+v", err)
	}

	actualTasks := s.GetTasks()

	if !reflect.DeepEqual(expTasks, actualTasks) {
		t.Fatalf("expected tasks %+v, got %+v", expTasks, actualTasks)
	}
}

func setupTest(tb testing.TB) {
	tb.Helper()

	err := os.MkdirAll("tmp", 0755)
	if err != nil {
		tb.Fatalf("failed to create test directory: %+v", err)
	}
}

func teardownTest(tb testing.TB) {
	tb.Helper()

	err := os.RemoveAll("tmp")
	if err != nil {
		tb.Fatalf("failed to remove test directory: %+v", err)
	}
}
