package internal

// BUG: add subtasks, edit parent task children disappear

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	contentStyle            = lipgloss.NewStyle().MarginLeft(2).MarginTop(1)
	formKeyStyle            = lipgloss.NewStyle().Foreground(SecondaryGrayColor).Bold(true)
	listTitleStyle          = lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)
	helpKeyStyle            = lipgloss.NewStyle().Foreground(SecondaryGrayColor).Bold(true)
	helpKeyDescriptionStyle = lipgloss.NewStyle().Foreground(SecondaryGrayColor)
)

type Repository interface {
	AllTasks() ([]Task, error)
	AllTasksWithIds([]int) ([]Task, error)
	EditTask(m *Task) error
	AddTask(m *Task) (*Task, error)
	DeleteTask(id int) (*Task, error)
}

type RefreshTasksMsg struct{}

func refreshTasks() tea.Cmd {
	return func() tea.Msg {
		return RefreshTasksMsg{}
	}
}

type Model struct {
	keys         KeyMapList
	help         help.Model
	repository   Repository
	tasks        list.Model
	toast        tea.Model
	form         tea.Model
	modeForm     bool
	modeFocus    bool
	modeArchived bool
	taskFocused  *Task
	// focusHistory tracks each focus then pops these back to taskFocused when returning
	focusHistory []*Task
	deleted      []*Task
	height       int
	width        int
}

// New takes in a path and sets up a sqlite DB in that
// location to persist tasks. If successful, will return a Model.
// Panics on any error.
func New(path string) *Model {
	dbLocation := fmt.Sprintf("%s/task.db", path)
	r, err := ConnectToDB(dbLocation)
	// r, err := ConnectToDB(":memory:")
	if err != nil {
		log.Fatalf("failed to connect to DB %v", err)
	}
	// when pop list, filterOutArchived
	tasks, err := r.AllTasks()
	if err != nil {
		log.Fatalf("failed to get all tasks %v", err)
	}
	l := newList(orderTasks(filterOutArchived(tasks)))

	return &Model{
		keys:       DefaultKeyMapList(),
		repository: r,
		tasks:      l,
		toast:      NewToast(),
		form:       NewForm(),
		help:       help.New(),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m *Model) View() string {
	var builder strings.Builder
	if m.modeForm {
		builder.WriteString(m.form.View())
	} else if m.modeFocus && m.taskFocused != nil {
		builder.WriteString(renderFocusedTask(m.taskFocused))
		builder.WriteString(m.tasks.View())
	} else {
		builder.WriteString(renderHeader(m.modeArchived))
		builder.WriteString(m.tasks.View())
	}
	builder.WriteString("\n")
	builder.WriteString(m.help.View(m.keys))
	builder.WriteString("\n")
	builder.WriteString(m.toast.View())
	return contentStyle.Render(builder.String())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	// global listeners
	m.toast, cmd = m.toast.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case RefreshTasksMsg:
		if m.modeArchived {
			tasks, err := m.repository.AllTasks()
			if err != nil {
				log.Printf("AllTasks failed with reason %v", err)
			}
			tasks = filterToArchived(tasks)
			m.tasks.SetItems(transformToItems(orderTasks(tasks)))
		} else if m.modeFocus && m.taskFocused != nil {
			tasks, err := m.repository.AllTasksWithIds(m.taskFocused.SubTasks)
			if err != nil {
				log.Printf("failed to get all tasks by id %v", err)
				return m, tea.Batch(append(cmds, ShowToast("failed to get all tasks", ToastWarn))...)
			}
			m.tasks.SetItems(transformToItems(orderTasks(tasks)))
		} else {
			tasks, err := m.repository.AllTasks()
			if err != nil {
				cmds = append(cmds, ShowToast("failed to get all tasks", ToastWarn))
				return m, tea.Batch(cmds...)
			}
			// post filter for archived
			tasks = filterOutArchived(tasks)
			items := transformToItems(orderTasks(tasks))
			m.tasks.SetItems(items)
		}
		return m, tea.Batch(cmds...)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, tea.Batch(cmds...)
	// always enable quit from anywhere
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.modeForm || m.tasks.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.Export):
			if !m.modeArchived {
				return m, tea.Batch(append(cmds, ShowToast("only supports exporting archived", ToastWarn))...)
			}

			items := m.tasks.Items()
			var builder strings.Builder
			builder.WriteString("### Archived Tasks\n")
			for _, item := range items {
				task := item.(*Task)
				builder.WriteString("- [x] ")
				builder.WriteString(task.Description)
				builder.WriteString("\n")
				builder.WriteString("	- (add impact here)\n")
				builder.WriteString("	- (add resources here)\n")
			}
			err := clipboard.WriteAll(builder.String())
			if err != nil {
				log.Printf("failed to copy to clipboard: %v", err)
				return m, tea.Batch(append(cmds, ShowToast("failed to copy to clipboard", ToastWarn))...)
			}
			return m, tea.Batch(append(cmds, ShowToast("copied export to clipboard", ToastInfo))...)
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keys.IncreasePriority):
			if selected := m.tasks.SelectedItem(); selected != nil {
				task := selected.(*Task)
				index := m.tasks.Index()
				if index == 0 {
					return m, tea.Batch(cmds...)
				}
				task.Priority = index - 1
				err := m.repository.EditTask(task)
				if err != nil {
					log.Printf("failed to decrease priority %v", err)
					return m, tea.Batch(append(cmds, ShowToast("unable to decrease priority", ToastWarn))...)
				}
				above := m.tasks.Items()[index-1].(*Task)
				above.Priority = index
				err = m.repository.EditTask(above)
				if err != nil {
					log.Printf("failed to decrease priority %v", err)
					return m, tea.Batch(append(cmds, ShowToast("unable to decrease priority", ToastWarn))...)
				}

				m.tasks.CursorUp()
				return m, tea.Batch(append(cmds, refreshTasks())...)
			}
		case key.Matches(msg, m.keys.DecreasePriority):
			if selected := m.tasks.SelectedItem(); selected != nil {
				task := selected.(*Task)
				index := m.tasks.Index()
				items := m.tasks.Items()
				if index == len(items)-1 {
					return m, tea.Batch(cmds...)
				}
				task.Priority = index + 1
				err := m.repository.EditTask(task)
				if err != nil {
					log.Printf("failed to decrease priority %v", err)
					return m, tea.Batch(append(cmds, ShowToast("unable to decrease priority", ToastWarn))...)
				}

				below := items[index+1].(*Task)
				below.Priority = index
				err = m.repository.EditTask(below)
				if err != nil {
					log.Printf("failed to decrease priority %v", err)
					return m, tea.Batch(append(cmds, ShowToast("unable to decrease priority", ToastWarn))...)
				}

				m.tasks.CursorDown()
				return m, tea.Batch(append(cmds, refreshTasks())...)
			}
		case key.Matches(msg, m.keys.ArchivedTaskToggle):
			// can only archive at root
			if selected := m.tasks.SelectedItem(); selected != nil && !m.modeFocus {
				task := selected.(*Task)
				task.IsArchived = !task.IsArchived
				err := m.repository.EditTask(task)
				if err != nil {
					log.Printf("failed to archive task, reason %v", err)
					return m, tea.Batch(append(cmds, ShowToast("unable to archive task", ToastWarn))...)
				}
				return m, tea.Batch(append(cmds, ShowToast("archived task", ToastInfo), refreshTasks())...)
			}
		case key.Matches(msg, m.keys.ToggleArchived):
			if m.modeFocus { // only allow the archived filtering at the root
				break
			}
			tasks, err := m.repository.AllTasks()
			if err != nil {
				log.Printf("AllTasks failed with reason %v", err)
			}
			m.modeArchived = !m.modeArchived
			if m.modeArchived {
				tasks = filterToArchived(tasks)
			} else {
				tasks = filterOutArchived(tasks)
			}
			m.tasks.SetItems(transformToItems(orderTasks(tasks)))
		case key.Matches(msg, m.keys.MoveStateForward):
			if selected := m.tasks.SelectedItem(); selected != nil {
				task := selected.(*Task)
				task.State = (task.State + 1) % TotalStates
				err := m.repository.EditTask(task)
				if err != nil {
					return m, tea.Batch(append(cmds, ShowToast("failed to toggle task", ToastWarn))...)
				}
			}
		case key.Matches(msg, m.keys.MoveStateBackward):
			if selected := m.tasks.SelectedItem(); selected != nil {
				task := selected.(*Task)
				if task.State == 0 {
					task.State = Urgent
				} else {
					task.State--
				}
				err := m.repository.EditTask(task)
				if err != nil {
					return m, tea.Batch(append(cmds, ShowToast("failed to toggle task", ToastWarn))...)
				}
			}
		case key.Matches(msg, m.keys.Undo): // FIXME: add ability to undo archiving...
			if len(m.deleted) == 0 {
				return m, tea.Batch(append(cmds, ShowToast("no more to undo", ToastInfo))...)
			}
			lastRemoved := m.deleted[len(m.deleted)-1]
			m.deleted = m.deleted[:len(m.deleted)-1]
			m.repository.AddTask(lastRemoved)
			if m.modeFocus && m.taskFocused != nil {
				m.taskFocused.SubTasks = append(m.taskFocused.SubTasks, *&lastRemoved.Id)
				err := m.repository.EditTask(m.taskFocused)
				if err != nil {
					log.Printf("failed to undo, reason: %v", err)
					return m, tea.Batch(append(cmds, ShowToast("failed to undo", ToastInfo))...)
				}
			}
			return m, tea.Batch(append(cmds, ShowToast("re-added deleted mark!", ToastInfo), refreshTasks())...)
		case key.Matches(msg, m.keys.Add):
			m.modeForm = true
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keys.Edit):
			if selected := m.tasks.SelectedItem(); selected != nil {
				task := selected.(*Task)
				m.modeForm = true
				return m, tea.Batch(append(cmds, editTask(task))...)
			}
		case key.Matches(msg, m.keys.Delete):
			if selected := m.tasks.SelectedItem(); selected != nil {
				task := selected.(*Task)
				if m.modeFocus && m.taskFocused != nil {
					ok := m.taskFocused.RemoveSubTask(task.Id)
					if !ok {
						log.Printf("weird, sub task not removed as ID not in parent, p:%d, st:%d", m.taskFocused.Id, task.Id)
					}
					task, err := m.repository.DeleteTask(task.Id)
					if err != nil {
						log.Printf("%v", err)
						return m, tea.Batch(append(cmds, ShowToast("failed to delete task", ToastWarn))...)
					}
					m.deleted = append(m.deleted, task)
					err = m.repository.EditTask(m.taskFocused)
					if err != nil {
						log.Printf("%v", err)
						return m, tea.Batch(append(cmds, ShowToast("failed to delete task", ToastWarn))...)
					}
				} else {
					task, err := m.repository.DeleteTask(task.Id)
					if err != nil {
						log.Printf("%v", err)
						return m, tea.Batch(append(cmds, ShowToast("failed to delete task", ToastWarn))...)
					}
					m.deleted = append(m.deleted, task)
				}
				return m, tea.Batch(append(cmds, refreshTasks(), ShowToast("deleted task!", ToastInfo))...)
			}
		case key.Matches(msg, m.keys.Return):
			if m.modeFocus && len(m.focusHistory) == 0 {
				m.modeFocus = false
				m.taskFocused = nil
				return m, tea.Batch(append(cmds, refreshTasks())...)
			} else if m.modeFocus && len(m.focusHistory) > 0 {
				last := m.focusHistory[len(m.focusHistory)-1]
				m.focusHistory = m.focusHistory[:len(m.focusHistory)-1]
				m.taskFocused = last
				return m, tea.Batch(append(cmds, refreshTasks())...)
			}
		case key.Matches(msg, m.keys.Focus):
			if m.taskFocused != nil {
				m.focusHistory = append(m.focusHistory, m.taskFocused)
			}
			if selected := m.tasks.SelectedItem(); selected != nil {
				task := selected.(*Task)
				m.taskFocused = task
				m.modeFocus = true
				tasks, err := m.repository.AllTasksWithIds(task.SubTasks)
				if err != nil {
					log.Printf("failed to get all tasks by id %v", err)
					return m, tea.Batch(append(cmds, ShowToast("failed to focus", ToastWarn))...)
				}
				m.tasks.SetItems(transformToItems(tasks))
				return m, tea.Batch(cmds...)
			}
		case key.Matches(msg, m.keys.Copy):
			if selected := m.tasks.SelectedItem(); selected != nil {
				task := selected.(*Task)
				err := clipboard.WriteAll(task.Description)
				if err != nil {
					log.Printf("failed to copy to clipboard: %v", err)
					return m, tea.Batch(append(cmds, ShowToast("failed to copy to clipboard", ToastWarn))...)
				}
				return m, tea.Batch(append(cmds, ShowToast("copied to clipboard!", ToastInfo))...)
			}
		}
	}

	switch msg := msg.(type) {
	case CloseFormMsg:
		log.Printf("m.showForm = %v", m.modeForm)
		m.modeForm = false
	case TaskCreatedMsg:
		log.Printf("task created %+v", msg.task)
		if m.modeFocus && m.taskFocused != nil {
			msg.task.IsSubTask = true
			msg.task.Priority = len(m.tasks.Items())
			task, err := m.repository.AddTask(&msg.task)
			if err != nil {
				log.Printf("%v", err)
				return m, tea.Batch(append(cmds, ShowToast("failed to add task", ToastWarn))...)
			}
			m.taskFocused.SubTasks = append(m.taskFocused.SubTasks, task.Id)
			err = m.repository.EditTask(m.taskFocused)
			if err != nil {
				log.Printf("%v", err)
				return m, tea.Batch(append(cmds, ShowToast("failed to add task", ToastWarn))...)
			}
		} else {
			_, err := m.repository.AddTask(&msg.task)
			if err != nil {
				log.Printf("%v", err)
				return m, tea.Batch(append(cmds, ShowToast("failed to add task", ToastWarn))...)
			}
		}

		return m, tea.Batch(append(cmds, refreshTasks())...)
	case TaskModifiedMsg:
		log.Printf("modifying task %+v", msg.task)
		err := m.repository.EditTask(&msg.task)
		if err != nil {
			log.Printf("%v", err)
			return m, tea.Batch(append(cmds, ShowToast("failed to edit mark", ToastWarn))...)
		}
		return m, tea.Batch(append(cmds, refreshTasks())...)
	}

	// these all are mutually exclusive...
	if m.modeForm {
		m.form, cmd = m.form.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	} else {
		m.tasks, cmd = m.tasks.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func newList(tasks []Task) list.Model {
	items := transformToItems(tasks)
	l := list.New(items, delegate{}, 30, 15)
	l.Styles.Title = listTitleStyle
	l.Title = ""
	l.SetShowStatusBar(false)
	l.DisableQuitKeybindings()
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	return l
}

func transformToItems(tasks []Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, mark := range tasks {
		item := &mark
		items[i] = item
	}
	return items
}

// filterByArchiveStatus removes archived tasks
func filterOutArchived(tasks []Task) []Task {
	var active []Task
	for _, task := range tasks {
		if !task.IsArchived {
			active = append(active, task)
		}
	}
	return active
}

// filterByArchiveStatus removes archived tasks
func filterToArchived(tasks []Task) []Task {
	var archived []Task
	for _, task := range tasks {
		if task.IsArchived {
			archived = append(archived, task)
		}
	}
	return archived
}

func orderTasks(tasks []Task) []Task {
	slices.SortFunc(tasks, func(a Task, b Task) int {
		return a.Priority - b.Priority
	})
	return tasks
}
