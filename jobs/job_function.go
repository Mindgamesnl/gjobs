package jobs

import (
	"errors"
	"github.com/google/uuid"
)

type FunctionJob struct {
	id string
	state int
	onFinishCallback func()
	task func()
}

func (f FunctionJob) StdWrite(line string) {
	panic("Function Job's don't have a Std handler")
}

func (f FunctionJob) OnStdOut(f2 func(line string)) {
	panic("Function Job's don't have a Std handler")
}

func (f FunctionJob) OnStdErr(f2 func(line string)) {
	panic("Function Job's don't have a Std handler")
}

func (f *FunctionJob) onFinish(callback func()) {
	f.onFinishCallback = callback
}

func (f *FunctionJob) start() error {
	if f.state != STATE_WAITING {
		return errors.New("Only waiting jobs can be started")
	}

	f.state = STATE_RUNNING
	go f.task()
	return nil
}

func (f FunctionJob) Kill() error {
	return errors.New("Functional job's cant be killed")
}

func (f FunctionJob) GetId() string {
	return f.id
}

func (f FunctionJob) GetState() int {
	return f.state
}

func (f FunctionJob) GetType() int {
	return JOBTYPE_FUNCTION
}

func NewFunctionJob(task func()) Job {
	fj := FunctionJob{
		id: uuid.New().String(),
		state: STATE_WAITING,
		task: task,
	}

	return &fj
}
