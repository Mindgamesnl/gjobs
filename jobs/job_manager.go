package jobs

import (
	"github.com/orcaman/concurrent-map"
	"time"
)

type JobManager struct {
	maxConcurrentJobs   int
	taskMap             cmap.ConcurrentMap
}

func (manager *JobManager) ScheduleFunction(f func()) Job {
	job := NewFunctionJob(f)
	manager.taskMap.Set(job.GetId(), job)
	return job
}

func (manager *JobManager) ScheduleCommand(command []string) Job {
	job := NewCommandJob(command)
	manager.taskMap.Set(job.GetId(), job)
	return job
}

func (manager *JobManager) findWaitingJob() Job {
	keys := manager.taskMap.Keys()
	for i := range keys {
		q, found := manager.taskMap.Get(keys[i])
		if found {
			j, _ := q.(Job)
			if j.GetState() == STATE_WAITING {
				return j
			}
		}
	}
	return nil
}

func (manager *JobManager) countRunningJobs() int {
	r := 0

	keys := manager.taskMap.Keys()
	for i := range keys {
		q, found := manager.taskMap.Get(keys[i])
		if found {
			j, _ := q.(Job)
			if j.GetState() == STATE_RUNNING {
				r++
			}
		}
	}

	return r
}

func (manager *JobManager) dispatchJobs(jobsToDispatch int) {
	if jobsToDispatch < 1 {
		return
	}

	for i := 0; i < jobsToDispatch; i++ {
		waitingJob := manager.findWaitingJob()
		if waitingJob != nil {
			waitingJob.start()
			manager.taskMap.Set(waitingJob.GetId(), waitingJob)
		}
	}
}

func (manager *JobManager) clean() {
	keys := manager.taskMap.Keys()
	for i := range keys {
		q, found := manager.taskMap.Get(keys[i])
		if found {
			j, _ := q.(Job)
			if j.GetState() == STATE_STOPPED {
				manager.taskMap.Remove(j.GetId())
			}
		}
	}
}

func NewJobManager(maxConcurrentJobs int) *JobManager {
	manager := JobManager{
		maxConcurrentJobs:   maxConcurrentJobs,
		taskMap:             cmap.New(),
	}

	go func() {
		for range time.NewTicker(time.Second * 1).C {
			manager.dispatchJobs(maxConcurrentJobs - manager.countRunningJobs())
		}
	}()

	return &manager
}
