package main

import "log"

func testJobs()  {
	manager := NewJobManager(2)

functionalJob := manager.ScheduleFunction(func() {
	// do your heavy task
})

}
