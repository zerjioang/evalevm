package datatype

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"
)

type TaskId struct {
	app        string
	identifier string
}

func (t TaskId) App() string {
	return t.app

}

type WorkerPool struct {
	WorkerCount int
	Tasks       chan Task
	wg          sync.WaitGroup
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		WorkerCount: 1, // one single container running each test
		Tasks:       make(chan Task, 1000),
	}
}

func (wp *WorkerPool) Run() {
	for i := 0; i < wp.WorkerCount; i++ {
		go wp.worker(i)
	}
	wp.wg.Wait()
}

func (wp *WorkerPool) Submit(task Task) {
	wp.wg.Add(1)
	wp.Tasks <- task
}

func (wp *WorkerPool) Close() {
	close(wp.Tasks)
}

func (wp *WorkerPool) worker(workerID int) {
	for task := range wp.Tasks {
		log.Printf("[Worker %d] Task %s started", workerID, task.ID())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel() // ensure context is cancelled to release resources

		var stdout, stderr bytes.Buffer
		var err error

		start := time.Now() // measure time as late as possible

		if len(task.Command()) > 0 {
			log.Printf("[Worker %d] Executing command: %s", workerID, task.Command())
			cmd := exec.CommandContext(ctx, "docker", task.Command()...)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err = cmd.Run()

			// Check for context timeout
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				err = fmt.Errorf("task timed out after 10 minutes")
			}
		}

		elapsed := time.Since(start)
		output := stdout.String() + stderr.String()

		// Capture result
		task.WithResult(&Result{
			Task:             task,
			Output:           output,
			Error:            err,
			TotalElapsedTime: elapsed,
		})

		task.Parse()

		// Send completion signal non-blocking
		select {
		case task.FinishChan() <- struct{}{}:
		default:
			log.Printf("[Worker %d] Warning: FinishChan blocked for task %s", workerID, task.ID())
		}

		log.Printf("[Worker %d] Task %s completed in %s", workerID, task.ID(), elapsed)
		wp.wg.Done()
	}
}
