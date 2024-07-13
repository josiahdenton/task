# task

A simple CLI tool to track tasks. Goes with with [mark](https://github.com/josiahdenton/mark)

## Install

Assuming go is already installed, simply run
```
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
<img width="691" alt="image" src="https://github.com/user-attachments/assets/e8126a78-f568-4549-8360-3896ed806c1a">





#### Progress

Task can
- [x] show task description
- [ ] display due date (relative)
- [x] enable recursive subtasks
- [x] archive tasks to move to a separate view without losing them
- [ ] auto export completed / archived tasks to help with performance reviews
- [ ] show help / controls at bottom
