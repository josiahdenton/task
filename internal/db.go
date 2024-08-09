package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Document struct {
	Id     int    `db:"doc_id"`
	Blob   string `db:"blob"`
	Hidden int    `db:"hidden"`
}

var schema = `
CREATE TABLE IF NOT EXISTS tasks (
    doc_id INTEGER PRIMARY KEY,
    blob TEXT,
    hidden INTEGER
);
`

/**
we can use the following in our refactor to setup the many to many connection

CREATE TABLE IF NOT EXISTS tasks_to_tasks (
  id SERIAL PRIMARY KEY,
  task_id_parent INTEGER NOT NULL,
  task_id_child INTEGER NOT NULL
  )

we would combine this with the above schema
*/

// hidden is a flag that determines if it should be returned when getting all tasks

// FIXME: could probably just have the list of subtasks encoded instead of the full object

func ConnectToDB(dbName string) (*TaskDatabase, error) {
	// use :memory: for testing
	db, err := sqlx.Connect("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s, %w", dbName, err)
	}

	// exec the schema or fail;
	db.MustExec(schema)

	return &TaskDatabase{
		db: db,
	}, nil
}

type TaskDatabase struct {
	db *sqlx.DB
}

func (m *TaskDatabase) AllTasksWithIds(ids []int) ([]Task, error) {
	var documents []Document
	err := m.db.Select(&documents, "SELECT * FROM tasks ORDER BY doc_id ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to select all tasks %w", err)
	}

	var tasks []Task
	for _, document := range documents {
		if slices.Contains(ids, document.Id) {
			var t Task
			err := json.Unmarshal([]byte(document.Blob), &t)
			if err != nil {
				return []Task{}, err
			}
			tasks = append(tasks, t)
		}
	}

	return tasks, nil
}

func (m *TaskDatabase) AllTasks() ([]Task, error) {
	var documents []Document
	err := m.db.Select(&documents, "SELECT * FROM tasks WHERE hidden=FALSE ORDER BY doc_id ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to select all tasks %w", err)
	}

	var tasks []Task
	for _, document := range documents {
		var t Task
		err := json.Unmarshal([]byte(document.Blob), &t)
		if err != nil {
			return []Task{}, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (m *TaskDatabase) maxId() (int, error) {
	var documents []Document
	err := m.db.Select(&documents, "SELECT * FROM tasks WHERE doc_id = (SELECT MAX(doc_id) FROM tasks)")
	if err != nil {
		return 0, fmt.Errorf("could not get max id %w", err)
	}

	if len(documents) < 1 {
		return 0, nil
	}

	return documents[0].Id, nil
}

func (m *TaskDatabase) EditTask(task *Task) error {
	blob, err := json.Marshal(task)
	if err != nil {
		return err
	}
	d := Document{Id: task.Id, Blob: string(blob)}

	_, err = m.db.Exec("UPDATE tasks SET blob=$1 WHERE doc_id=$2", d.Blob, d.Id)
	if err != nil {
		return fmt.Errorf("EditTask failed: %w", err)
	}
	return nil
}

func (m *TaskDatabase) AddTask(task *Task) (*Task, error) {
	maxId, err := m.maxId()
	if err != nil {
		log.Fatalf("could not get max id %v", err)
	}
	// our ID should be max + 1
	if task.Id == 0 { // != --> then it's an undo
		task.Id = maxId + 1
		task.Priority = task.Id
	}

	blob, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	d := Document{Id: task.Id, Blob: string(blob)}
	if task.IsSubTask {
		d.Hidden = 1
	}

	_, err = m.db.Exec("INSERT INTO tasks (doc_id, blob, hidden) VALUES ($1, $2, $3)", d.Id, d.Blob, d.Hidden)
	if err != nil {
		return nil, fmt.Errorf("failed to add task: %w", err)
	}
	return task, nil
}

func (m *TaskDatabase) DeleteTask(id int) (*Task, error) {
	task, err := m.task(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task for delete: %w", err)
	}
	_, err = m.db.Exec("DELETE FROM tasks WHERE doc_id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete task: %w", err)
	}
	return task, nil
}

func (m *TaskDatabase) task(id int) (*Task, error) {
	var documents []Document
	err := m.db.Select(&documents, "SELECT * FROM tasks WHERE doc_id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch task %w", err)
	}

	if len(documents) == 0 {
		return nil, fmt.Errorf("no task found with id %d", id)
	}

	d := documents[0]
	var t Task
	err = json.Unmarshal([]byte(d.Blob), &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
