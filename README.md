# gjobs
Simple Go library to schedule command or functional jobs with concurrent execution limits, install with
```
go get github.com/Mindgamesnl/gjobs
```

# Job manager
You need to get started by making a job manager, the only argument is `maxConcurrentJobs` which specifies the limit of allowed concurrent jobs
```go
manager := gjobs.NewJobManager(2)
```

# Jobs
The job manager is what contains and manages the queue, each job has a few public api functions, these are
- `Kill()` Kills the process (only if its a command and actually running)
- `GetId() string` returns the job id
- `GetState() int` returns the current job state, which is `STATE_WAITING`, `STATE_RUNNING`, or `STATE_STOPPED`
- `GetType() int` returns the job type, which is `JOBTYPE_COMMAND` or `JOBTYPE_FUNCTION`
- `StdWrite(line string)` writes your input string to the StdIn of the child process (only if the job is a command type)
- `OnStdOut(func(line string))` and `OnStdErr(func(line string))` handle both std channels respectively on a line-by-line bases (only if the job is a command type)

### Example command job
```go
commandJob := manager.ScheduleCommand([]string{"ping", "google.com"})

commandJob.OnStdErr(func(line string) {
	log.Println("Error: " + line)
})

commandJob.OnStdOut(func(line string) {
	log.Println("Out: " + line)
})
```

### Example functional job
```go
functionalJob := manager.ScheduleFunction(func() {
	// do your heavy task
})
```