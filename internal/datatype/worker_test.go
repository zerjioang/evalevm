package datatype

import (
	"log"
	"testing"
)

func TestNewWorkerPool(t *testing.T) {
	pool := NewWorkerPool(4) // 4 concurrent workers
	go pool.Run()

	task := DockerTask{
		id:  "",
		cmd: []string{"run", "hello-world"},
	}
	pool.Submit(task)

	pool.Close()

	// Collect results
	for result := range pool.Results {
		if result.Error != nil {
			log.Printf("Task %d failed: %v\nOutput:\n%sTime:\n%v", result.TaskID, result.Error, result.Output, result.ElapsedTime)
		} else {
			log.Printf("Task %d succeeded:\nOutput:\n%sTime:\n%v", result.TaskID, result.Output, result.ElapsedTime)
		}
	}
}

func TestWorkerPoolMeasure(t *testing.T) {
	pool := NewWorkerPool(4) // 4 concurrent workers
	go pool.Run()

	// docker run --rm -v /path/to/measure.sh:/measure.sh:ro myimage /measure.sh somecommand arg1 arg2
	task := DockerTask{
		id: TaskId{},
		cmd: []string{
			"run", "--rm",
			"-v", "/path/to/measure.sh:/measure.sh:ro", // bind mount measure.sh
			"myimage",
			"/measure.sh", "somecommand", "arg1", "arg2",
		},
	}
	pool.Submit(task)

	pool.Close()

	// Collect results
	for result := range pool.Results {
		if result.Error != nil {
			log.Printf("Task %d failed: %v\nOutput:\n%sTime:\n%v", result.TaskID, result.Error, result.Output, result.ElapsedTime)
		} else {
			log.Printf("Task %d succeeded:\nOutput:\n%sTime:\n%v", result.TaskID, result.Output, result.ElapsedTime)
		}
	}
}
