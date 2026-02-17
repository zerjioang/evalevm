package engine

import (
	"context"
	"encoding/json"
	"evalevm/internal/datatype"
	"evalevm/internal/tool/bytespector"
	"evalevm/internal/tool/ethersolve"
	"evalevm/internal/tool/evm_cfg"
	"evalevm/internal/tool/evm_cfg_builder"
	"evalevm/internal/tool/evmlisa"
	"evalevm/internal/tool/evmole"
	"evalevm/internal/tool/octopus"
	"evalevm/internal/tool/paper"
	"evalevm/internal/tool/rattle"
	"evalevm/internal/tool/vandal"
	"evalevm/internal/uuid"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

type Comparator struct {
	analyzerList []datatype.Analyzer
	threads      int
	pool         *datatype.WorkerPool
}

func NewComparator(audit bool, runMode string) Comparator {
	analyzerList := []datatype.Analyzer{
		//ethir.NewEthIR(),
		bytespector.NewByteInspector(),
		// conkas.NewConkas(),
		// tool.NewDefectChecker(),
		ethersolve.NewEthersolve(runMode, audit),
		evm_cfg.NewEvmCFG(),
		evm_cfg_builder.NewEvmCFGBuilder(),
		evmlisa.NewEvmLisa(audit),
		evmole.NewEVMole(),
		// tool.NewGigaHorse(),
		// tool.NewHeimdall(),
		// honeybadger.NewHoneyBadger(),
		// tool.NewMadMax(),
		// tool.NewMaian(),
		// tool.NewManticore(),
		// tool.NewMythril(),
		// tool.NewMythril(),
		octopus.NewOctopus(),
		// tool.NewOsiris(),
		// tool.NewOsiris(),
		// tool.NewOyente(),
		// tool.NewPakala(),
		// tool.NewPanoramix(),
		// tool.NewPanoramixPalkeo(),
		paper.NewPaper(),
		// tool.NewPorosity(),
		// tool.NewPyevmasm(),
		rattle.NewRattle(),
		// securify.NewSecurify(),
		// tool.NewSlither(),
		// tool.NewTeether(),
		vandal.NewVandal(),
	}
	return Comparator{
		analyzerList: analyzerList,
		threads:      runtime.NumCPU(),
		pool:         datatype.NewWorkerPool(),
	}
}

// FilterByTools filters the analyzer list to only include analyzers whose
// Name() matches one of the provided names. If names is empty, no filtering
// is applied. This allows the --tools flag to select specific tools.
func (c *Comparator) FilterByTools(names []string) {
	if len(names) == 0 {
		return
	}
	for _, n := range names {
		if n == "all" {
			return
		}
	}
	allowed := make(map[string]bool, len(names))
	for _, n := range names {
		allowed[n] = true
	}
	filtered := make([]datatype.Analyzer, 0, len(names))
	for _, a := range c.analyzerList {
		if allowed[a.Name()] {
			filtered = append(filtered, a)
		}
	}
	c.analyzerList = filtered
}

func (c Comparator) Analyzers() []datatype.Analyzer {
	return c.analyzerList
}

func (c Comparator) Threads() int {
	return c.threads
}

func (c Comparator) Start(ctx context.Context) {
	go c.pool.Run(ctx)
}

func (c Comparator) Submit(hexBytecode string, sampleId string) datatype.TaskSet {
	// we scan the contract bytecode with all the tools
	taskId := uuid.UUIDv4()
	var ts datatype.TaskSet
	for _, analyzer := range c.analyzerList {
		hexBytecode = c.cleanBytecode(hexBytecode)
		taskList := analyzer.CreateTask(taskId, hexBytecode, sampleId)
		for _, task := range taskList {
			task.WithResultParser(analyzer)

			// Try to load cached result first
			if cached, err := c.tryLoadCachedResult(task); err == nil {
				log.Printf("Using cached result for %s", task.ID())
				// Link the loaded result to the task and vice versa
				task.WithResult(cached)
				// Send a signal to finish channel so consumers don't block
				select {
				case task.FinishChan() <- struct{}{}:
				default:
				}
				ts = append(ts, task)
				continue
			}

			ts = append(ts, task)
			c.pool.Submit(task)
		}
	}
	return ts
}

func (c Comparator) SubmitAndWait(hexBytecode string, sampleId string) datatype.TaskSet {
	log.Println("submit and wait")
	var wg sync.WaitGroup

	var ts datatype.TaskSet

	// we scan the contract bytecode with all the tools
	taskId := uuid.UUIDv4()
	for _, analyzer := range c.analyzerList {
		hexBytecode = c.cleanBytecode(hexBytecode)
		taskList := analyzer.CreateTask(taskId, hexBytecode, sampleId)
		for _, task := range taskList {
			ts = append(ts, task)
			task.WithResultParser(analyzer)

			// Try to load cached result first
			if cached, err := c.tryLoadCachedResult(task); err == nil {
				log.Printf("Using cached result for %s", task.ID())
				// Link the loaded result to the task and vice versa
				task.WithResult(cached)
				// Signal completion immediately
				select {
				case task.FinishChan() <- struct{}{}:
				default:
				}
				// Do NOT add to waitgroup or submit to pool
				continue
			}

			wg.Add(1)
			c.pool.Submit(task)
			go func(t datatype.Task) {
				<-t.FinishChan()
				wg.Done()
			}(task)
		}
	}
	wg.Wait()
	return ts
}

// cleanBytecode removes the 0x prefix from bytecode hex string
func (c Comparator) cleanBytecode(bytecode string) string {
	if strings.HasPrefix(bytecode, "0x") {
		return bytecode[2:]
	}
	return bytecode
}

// tryLoadCachedResult attempts to find and load a previously saved result file
// for the given task. Returns nil if file not found or load fails.
func (c *Comparator) tryLoadCachedResult(task datatype.Task) (*datatype.Result, error) {
	// Reconstruct the expected filename: result_<App>_<TrackerId>.json
	filename := fmt.Sprintf("result_%s_%s.json", task.ID().App(), task.TrackerId())

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal into a shadow struct to handle interface fields (Task, Error)
	// which json.Unmarshal cannot handle directly.
	type cachedResult struct {
		datatype.Result
		Task  json.RawMessage `json:"task"`
		Error json.RawMessage `json:"error"`
	}

	var shadow cachedResult
	if err := json.Unmarshal(content, &shadow); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached result: %w", err)
	}

	// Reconstruct the valid Result
	res := &shadow.Result
	res.Task = task // Link back to the live task
	res.Error = nil // Error interface cannot be fully restored from JSON generic object, assume nil or rely on other fields

	return res, nil
}
