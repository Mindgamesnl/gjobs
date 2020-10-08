package main

import (
	"io"
	"os"
	"os/exec"
)

var (
	RealHandlerSet = handlerSet{
		Out:   os.Stdout,
		Error: os.Stderr,
	}
)

type handlerSet struct {
	Out   io.Writer
	Error io.Writer
}

type WrappedWriter struct {
	OnWrite func(written []byte)

	Replaces io.Writer
}

func (w WrappedWriter) Write(p []byte) (n int, err error) {
	w.OnWrite(p)
	return w.Replaces.Write(p)
}

type CommandRunner struct {
	command     *exec.Cmd
	inputWriter io.WriteCloser
	OnFinish    func()
	OnSTDOut    func(line string)
	OnSTDErr    func(line string)
}

func (runner *CommandRunner) Start() {
	go func() {
		_ = runner.command.Run()
		runner.OnFinish()
	}()
}

func (runner *CommandRunner) Kill() {
	_ = runner.command.Process.Kill()
}

func NewCommandRunner(command []string) *CommandRunner {
	runner := CommandRunner{}

	runner.command = exec.Command(command[0], command[1:]...)

	runner.command.Stdout = WrappedWriter{
		OnWrite: func(written []byte) {
			runner.OnSTDOut(string(written))
		},
		Replaces: RealHandlerSet.Out,
	}

	runner.command.Stderr = WrappedWriter{
		OnWrite: func(written []byte) {
			runner.OnSTDErr(string(written))
		},
		Replaces: RealHandlerSet.Error,
	}

	epi, er := runner.command.StdinPipe()
	if er != nil {
		println(er)
	}
	runner.inputWriter = epi

	return &runner
}
