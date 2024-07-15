package internal

import (
	"reflect"
	"slices"
	"time"
)

// the DB will literally just be ID and JSON blob

type TaskState int

const (
	Ready TaskState = iota
	Focused
	Hold
	Completed
	Urgent
	TotalStates
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	IsSubTask   bool      `json:"isSubTask"`
	IsArchived  bool      `json:"isArchived"`
	State       TaskState `json:"state"`
	Priority    int       `json:"priority"`
	SubTasks    []int     `json:"sub_tasks"`
	Start       time.Time `json:"start"` // I don't know the data type... some epoch time??
	End         time.Time `json:"end"`   // I don't know the data type...
}

func (t *Task) FilterValue() string {
	return t.Description
}

func (t *Task) RemoveSubTask(id int) bool {
	if len(t.SubTasks) == 0 {
		return false
	}

	for i, stId := range t.SubTasks {
		if stId == id {
			t.SubTasks = slices.Delete(t.SubTasks, i, i+1)
			return true
		}
	}

	return false
}

func (t *Task) Open() {
	// can only be started once
	if reflect.ValueOf(t.Start).IsZero() {
		t.Start = time.Now()
	}
}

func (t *Task) Close() {
	// allow for re-opening a task
	t.End = time.Now()
}
