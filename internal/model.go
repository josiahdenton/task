package internal

import "slices"

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
	Start       string    `json:"start"` // I don't know the data type... some epoch time??
	End         string    `json:"end"`   // I don't know the data type...
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

// TODO: func to set times...
