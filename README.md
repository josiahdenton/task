# task

A simple CLI tool to track tasks. Goes with with [mark]()

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


#### Progress

Task can
- [x] show task description
- [ ] display due date (relative)
- [x] enable recursive subtasks
- [x] archive tasks to move to a separate view without losing them
- [ ] auto export completed / archived tasks to help with performance reviews
