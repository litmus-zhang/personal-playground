package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}
type List []item

func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CompletedAt: time.Time{},
		CreatedAt:   time.Now(),
	}
	*l = append(*l, t)
}

func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d not in the list", i)
	}
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()
	return nil
}

func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d not in the list", i)
	}
	*l = append(ls[:i-1], ls[i:]...)
	return nil
}

func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, js, 0644)
}

func (l *List) Get(filename string) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
	}
	if len(f) == 0 {
		return nil
	}
	return json.Unmarshal(f, l)
}

func (l *List) String() string {
	formatted := ""

	for k, t := range *l {
		prefix := " "
		if t.Done {
			prefix = "✅ "
		}

		formatted += fmt.Sprintf("%s%d. %s\n", prefix, k+1, t.Task)
	}
	return formatted
}

