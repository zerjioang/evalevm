package datatype

import (
	"testing"
)

func TestNewWorkerPool(t *testing.T) {
	pool := NewWorkerPool() // 4 concurrent workers
	go pool.Run()

	task := &DockerTask{
		id:  TaskId{},
		cmd: []string{"run", "hello-world"},
	}
	pool.Submit(task)

	pool.Close()
}

func TestWorkerPoolMeasure(t *testing.T) {
	pool := NewWorkerPool() // 4 concurrent workers
	go pool.Run()

	// docker run --rm -v /path/to/measure.sh:/measure.sh:ro myimage /measure.sh somecommand arg1 arg2
	task := &DockerTask{
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
}
