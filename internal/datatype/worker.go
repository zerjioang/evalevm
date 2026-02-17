package datatype

import (
	"bytes"
	"context"
	"errors"
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

func (t TaskId) UID() string {
	return t.identifier
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

func (wp *WorkerPool) Run(ctx context.Context) {
	for i := 0; i < wp.WorkerCount; i++ {
		go wp.worker(ctx, i)
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

func (wp *WorkerPool) worker(ctx context.Context, workerID int) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-wp.Tasks:
			if !ok {
				return
			}
			wp.executeTask(ctx, workerID, task)
		}
	}
}

// executeTask runs a single task with panic recovery.
// Isolating each task in its own function ensures that:
// - A panic in one task doesn't kill the worker goroutine
// - wg.Done() is always called even after a panic
// - Context cancel is always called (no leak)
func (wp *WorkerPool) executeTask(ctx context.Context, workerID int, task Task) {
	defer wp.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Worker %d] recovered from panic in task %s: %v", workerID, task.ID(), r)
		}
	}()

	log.Printf("[Worker %d] Task %s started", workerID, task.ID())

	// Use the parent context for cancellation, but add a timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	var stdout, stderr bytes.Buffer
	var err error

	start := time.Now() // measure time as late as possible

	var stats *ContainerStats

	if len(task.Command()) > 0 {
		containerName := task.ContainerName()
		cmdArgs := []string{
			"run",
			"--name", containerName,
			"--cap-add=SYS_ADMIN",
			"--entrypoint=bash",
			"--network", "none", // disable network access for the container
		}
		cmdArgs = append(cmdArgs, task.Command()...)
		log.Printf("[Worker %d] Executing command: %s in container=%s", workerID, cmdArgs, containerName)

		cmd := exec.CommandContext(ctx, "docker", cmdArgs...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		done := make(chan error, 1)
		go func() {
			done <- cmd.Run()
		}()

		select {
		case <-ctx.Done():
			// Timeout or cancellation occurred, collect metrics before killing
			log.Printf("[Worker %d] Context done: %v, collecting metrics for container %s", workerID, ctx.Err(), containerName)

			// Get container stats (CPU, RAM, MemPerc, I/O)
			statsCmd := exec.Command("docker", "stats", "--no-stream", "--format",
				"{{.CPUPerc}}|{{.MemUsage}}|{{.MemPerc}}|{{.NetIO}}|{{.BlockIO}}|{{.PIDs}}", containerName)
			statsOut, statsErr := statsCmd.Output()
			if statsErr != nil {
				log.Printf("[Worker %d] Failed to get stats for container %s: %v", workerID, containerName, statsErr)
			} else {
				statsLine := string(statsOut)
				stats = parseContainerStats(statsLine)
				if stats != nil {
					log.Printf("[Worker %d] ContainerStats for %s: CPU=%s, MemUsage=%s, MemPerc=%s, NetIO=%s, BlockIO=%s, PIDs=%s", workerID, containerName, stats.CPUPerc, stats.MemUsage, stats.MemPerc, stats.NetIO, stats.BlockIO, stats.PIDs)
				} else {
					log.Printf("[Worker %d] Failed to parse stats for container %s: %s", workerID, containerName, statsLine)
				}
			}

			// Now kill the container
			log.Printf("[Worker %d] Killing container %s", workerID, containerName)
			killCmd := exec.Command("docker", "kill", containerName)
			_ = killCmd.Run()
			err = ctx.Err()
		case runErr := <-done:
			if runErr != nil {
				err = runErr
			}
		}

		// Always attempt to remove the container (no --rm flag, we manage lifecycle)
		log.Printf("[Worker %d] Removing container %s", workerID, containerName)
		rmCmd := exec.Command("docker", "rm", "-f", containerName)
		_ = rmCmd.Run()
	}

	elapsed := time.Since(start)

	// Capture result
	task.WithResult(&Result{
		Task:             task,
		Output:           stdout.Bytes(),
		OutputErr:        stderr.Bytes(),
		Error:            err,
		TotalElapsedTime: elapsed,
		Timeout:          errors.Is(ctx.Err(), context.DeadlineExceeded),
		Stats:            stats,
	})

	task.Parse()

	// Send completion signal non-blocking
	select {
	case task.FinishChan() <- struct{}{}:
	default:
		log.Printf("[Worker %d] Warning: finish channel blocked for task %s", workerID, task.ID())
	}

	log.Printf("[Worker %d] Task %s completed in %s", workerID, task.ID(), elapsed)
}
