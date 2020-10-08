package main

import (
	"errors"
	"github.com/google/uuid"
	"io"
)

type CommandJob struct {
	id     string
	state  int
	runner *CommandRunner
	onFinishCallback func()
}

func (c *CommandJob) StdWrite(line string) {
	_, _ = io.WriteString(c.runner.inputWriter, line)
}

func (c *CommandJob) OnStdOut(f func(line string)) {
	c.runner.OnSTDOut = f
}

func (c *CommandJob) OnStdErr(f func(line string)) {
	c.runner.OnSTDErr = f
}

func (c CommandJob) GetId() string {
	return c.id
}

func (c CommandJob) GetState() int {
	return c.state
}

func (c CommandJob) GetType() int {
	return JOBTYPE_COMMAND
}

func (c *CommandJob) start() error {
	if c.state != STATE_WAITING {
		return errors.New("Only waiting jobs can be started")
	}

	c.runner.OnFinish = func() {
		c.state = STATE_STOPPED
		c.onFinishCallback()
	}
	c.state = STATE_RUNNING
	c.runner.Start()
	return nil
}

func (c *CommandJob) onFinish(callback func()) {
	c.onFinishCallback = callback
}

func (c *CommandJob) kill() error {
	if c.state != STATE_RUNNING {
		return errors.New("Only running jobs may be killed")
	}

	c.runner.Kill()
	return nil
}

func NewCommandJob(command []string) Job {
	cj := CommandJob{
		id:    uuid.New().String(),
		state: STATE_WAITING,
	}

	cj.runner = NewCommandRunner(command)

	return &cj
}
