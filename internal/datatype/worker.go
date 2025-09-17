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

		var stats *ContainerStats

		if len(task.Command()) > 0 {
			containerName := task.ContainerName()
			cmdArgs := []string{
				"run",
				"--rm",
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

				// Get container start time for uptime
				inspectCmd := exec.Command("docker", "inspect", "-f", "{{.State.StartedAt}}", containerName)
				inspectOut, inspectErr := inspectCmd.Output()
				if inspectErr != nil {
					log.Printf("[Worker %d] Failed to inspect container %s: %v", workerID, containerName, inspectErr)
				} else {
					startedAt := string(inspectOut)
					log.Printf("[Worker %d] Container %s started at: %s", workerID, containerName, startedAt)
					// Optionally, parse startedAt and calculate uptime
				}

				// Now kill the container
				log.Printf("[Worker %d] Killing container %s", workerID, containerName)
				killCmd := exec.Command("docker", "kill", containerName)
				_ = killCmd.Run() // ignore error, just try to kill
				err = ctx.Err()
			case runErr := <-done:
				if runErr != nil {
					err = runErr
				}
			}

			log.Printf("[Worker %d] Context done: %v, removing container %s", workerID, ctx.Err(), containerName)
			killCmd := exec.Command("docker", "rm", containerName)
			_ = killCmd.Run() // ignore error, just try to rm
			err = ctx.Err()
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
		wp.wg.Done()
	}
}
