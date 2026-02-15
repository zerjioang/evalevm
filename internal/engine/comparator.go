package engine

import (
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
	"log"
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

func (c Comparator) Start() {
	go c.pool.Run()
}

func (c Comparator) Submit(hexBytecode string, sampleId string) datatype.TaskSet {
	// we scan the contract bytecode with all the tools
	taskId := uuid.UUIDv4()
	var ts datatype.TaskSet
	for _, analyzer := range c.analyzerList {
		hexBytecode = c.cleanBytecode(hexBytecode)
		taskList := analyzer.CreateTask(taskId, hexBytecode, sampleId)
		for _, task := range taskList {
			ts = append(ts, task)
			task.WithResultParser(analyzer)
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
			wg.Add(1)
			task.WithResultParser(analyzer)
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
