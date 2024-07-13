package internal

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	formTitleStyle = lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)
	formLabelStyle = lipgloss.NewStyle().Foreground(SecondaryColor)
)

type FormIndex int

// -- results --
type CloseFormMsg struct{}

func closeForm() tea.Cmd {
	return func() tea.Msg {
		return CloseFormMsg{}
	}
}

type TaskModifiedMsg struct {
	task Task
}

func taskModified(t Task) tea.Cmd {
	return func() tea.Msg {
		return TaskModifiedMsg{t}
	}
}

type TaskCreatedMsg struct {
	task Task
}

func taskCreated(t Task) tea.Cmd {
	return func() tea.Msg {
		return TaskCreatedMsg{t}
	}
}

// -- actions --
type EditTaskMsg struct {
	task *Task
}

func editTask(t *Task) tea.Cmd {
	return func() tea.Msg {
		return EditTaskMsg{
			task: t,
		}
	}
}

type AddTaskMsg struct {
	task Task
}

func addTask(t Task) tea.Cmd {
	return func() tea.Msg {
		return AddTaskMsg{
			task: t,
		}
	}
}

func NewForm() *FormModel {
	description := textinput.New()
	description.Focus()
	description.Width = 60
	description.CharLimit = 60
	description.Prompt = "Task: "
	description.PromptStyle = formLabelStyle
	description.Placeholder = "..."

	description.Validate = func(s string) error {
		if len(strings.Trim(s, " \n")) < 1 {
			return fmt.Errorf("name missing")
		}
		return nil
	}
	return &FormModel{
		description: description,
		keys:        DefaultKeyMapForm(),
	}
}

type FormModel struct {
	description textinput.Model
	t           Task
	activeInput FormIndex
	keys        KeyMapForm
}

func (f *FormModel) Init() tea.Cmd {
	return nil
}

func (f *FormModel) View() string {
	var b strings.Builder
	b.WriteString("\n\n")
	b.WriteString(f.description.View())
	return b.String()
}

func (f *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// -- init --
	switch msg := msg.(type) {
	case EditTaskMsg:
		f.t.Id = msg.task.Id
		f.t.Description = msg.task.Description
		f.t.SubTasks = msg.task.SubTasks
		f.t.State = msg.task.State
		f.t.Priority = msg.task.Priority
		f.t.IsSubTask = msg.task.IsSubTask
		f.t.Start = msg.task.Start
		f.t.End = msg.task.End

		f.description.SetValue(msg.task.Description)
	}

	// -- key messags --
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, f.keys.Close):
			f.reset()
			cmds = append(cmds, closeForm())
		case key.Matches(msg, f.keys.Submit):
			// assign form values and signal parent component
			f.t.Description = f.description.Value()

			// send off
			if f.t.Id != 0 { // existing mark
				log.Printf("modifying %d, now is task = %+v", f.t.Id, f.t)
				cmds = append(cmds, taskModified(f.t), closeForm())
			} else { // new mark
				cmds = append(cmds, taskCreated(f.t), closeForm())
			}
			f.reset() // clear form + task
		}

	}

	// form field input
	// -- name --
	f.description, cmd = f.description.Update(msg)
	cmds = append(cmds, cmd)

	return f, tea.Batch(cmds...)
}

func (f *FormModel) reset() {
	f.t = Task{}
	f.description.Reset()
}

func validateForm(errs ...error) tea.Cmd {
	for _, err := range errs {
		if err != nil {
			return ShowToast(fmt.Sprintf("%v", err), ToastWarn)
		}
	}
	return nil
}
