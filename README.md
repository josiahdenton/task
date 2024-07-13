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
<img width="695" alt="usage_1" src="https://github.com/user-attachments/assets/1fe9eeae-ed93-491d-bbd2-2266e18bb78c">

##### Focus mode on one task with 3 subtasks
<img width="695" alt="image" src="https://github.com/user-attachments/assets/e342d512-9faf-4b9e-beeb-d471ec590ded">




#### Progress

Task can
- [x] show task description
- [ ] display due date (relative)
- [x] enable recursive subtasks
- [x] archive tasks to move to a separate view without losing them
- [ ] auto export completed / archived tasks to help with performance reviews
- [ ] show help / controls at bottom
- [ ] check if there's a need to enable filtering
