# task

A simple CLI tool to track tasks. Pairs well with [mark](https://github.com/josiahdenton/mark)

## Requirements

- `go` is installed
- nerd font for icons

## Install

```
git clone https://github.com/josiahdenton/task.git
cd task
go install .
```
This installs `task` in `~/go/bin/`. To call `task` from anywher
make sure to add `~/go/bin/` to your path.

## Usage

Run in any terminal
```
task
```

##### Basic task list
<img width="800" alt="image" src="https://github.com/user-attachments/assets/fcffae72-e88a-41ed-be1b-012fd502659d">

##### Focus mode on one task with 3 subtasks
<img width="800" alt="image" src="https://github.com/user-attachments/assets/110c8ef1-2b9a-45cf-9e1b-5faca34271af">

##### Help
<img width="800" alt="image" src="https://github.com/user-attachments/assets/40767ed0-8fd1-4b46-af0b-dac59c5b833f">

##### Archived View
<img width="800" alt="image" src="https://github.com/user-attachments/assets/03158540-6653-4a29-8e4f-e50dc9275db4">


#### Progress

- [ ] improve error logs in `log` file
- [x] show total time spent working on task in focused view
- [x] show task description
- [x] auto-complete when all sub-tasks are marked as done
- [x] enable recursive subtasks
- [x] archive tasks to move to a separate view without losing them
- [x] auto export completed / archived tasks to help with performance reviews
- [x] show help / controls at bottom

