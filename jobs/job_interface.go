package jobs

const (
	STATE_WAITING = iota
	STATE_RUNNING = iota
	STATE_STOPPED = iota
)

const (
	JOBTYPE_COMMAND  = iota
	JOBTYPE_FUNCTION = iota
)

type Job interface {
	GetId() string
	GetState() int
	GetType() int
	StdWrite(line string)
	OnStdOut(func(line string))
	OnStdErr(func(line string))
	onFinish(callback func())
	start() error
	Kill() error
}
