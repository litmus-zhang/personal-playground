package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"todo-cli-app"
)

var todoFile = ".todo.json"

func main() {

	if os.Getenv("TODO_FILENAME") == "" {
		todoFile = os.Getenv("TODO_FILENAME")
	}
	add := flag.Bool("add", false, "Add a new task to the list")
	list := flag.Bool("list", false, "List all tasks")

	complete := flag.Int("complete", 0, "Item to be completed")
	delete := flag.Int("delete", 0, "Item to be deleted")
	verbose := flag.Bool("verbose", false, "Print verbose output")
	notDone := flag.Bool("notdone", false, "List all tasks not done")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s tool:\n developed by Lukmon Abdulsalam (Litmus Zhang)\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2024\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	l := &todo.List{}
	if err := l.Get(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch {
	case *list:
		fmt.Print(l)
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(t)
		if err := l.Save(todoFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > 0:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *notDone:
		for i, t := range *l {
			if !t.Done {
				fmt.Printf("%d. %s\n", i+1, t.Task)
			}
		}
	case *verbose:
		// cmdPath := os.Getwd()
		fmt.Println("Verbose mode")
		// print verbose information about the todoFile and the list
		fmt.Println("Todo file:", todoFile)
		for i, t := range *l {
			fmt.Printf("%d. %s %v %s %s\n", i+1, t.Task, t.Done, t.CompletedAt, t.CreatedAt)
		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be empty")
	}
	return s.Text(), nil
}
