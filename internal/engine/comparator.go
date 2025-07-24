package engine

import (
	"evalevm/internal/datatype"
	"evalevm/internal/tool/byteinspector"
	"evalevm/internal/tool/ethersolve"
	"evalevm/internal/tool/evm_cfg"
	"evalevm/internal/tool/evm_cfg_builder"
	"evalevm/internal/tool/evmole"
	"evalevm/internal/tool/honeybadger"
	"evalevm/internal/tool/paper"
	"evalevm/internal/tool/rattle"
	"evalevm/internal/tool/securify"
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

func NewComparator() Comparator {
	analyzerList := []datatype.Analyzer{
		// tool.NewAcuaricaEVM(),
		byteinspector.NewByteInspector(),
		//conkas.NewConkas(),
		// tool.NewDefectChecker(),
		ethersolve.NewEthersolveCreator(),
		ethersolve.NewEthersolveRuntime(),
		evm_cfg.NewEvmCFG(),
		evm_cfg_builder.NewEvmCFGBuilder(),
		//evmlisa.NewEvmLisa(),
		evmole.NewEVMole(),
		// tool.NewGigaHorse(),
		// tool.NewHeimdall(),
		honeybadger.NewHoneyBadger(),
		// tool.NewMadMax(),
		// tool.NewMaian(),
		// tool.NewManticore(),
		// tool.NewMythril(),
		// tool.NewOctopus(),
		// tool.NewOsiris(),
		// tool.NewOyente(),
		// tool.NewPakala(),
		// tool.NewPanoramix(),
		// tool.NewPanoramixPalkeo(),
		paper.NewPaper(),
		// tool.NewPorosity(),
		// tool.NewPyevmasm(),
		rattle.NewRattle(),
		securify.NewSecurify(),
		// tool.NewSlither(),
		// tool.NewTeether(),
		// tool.NewVandal(),
	}
	if false {
		analyzerList = []datatype.Analyzer{
			evmole.NewEVMole(),
			evm_cfg.NewEvmCFG(),
		}
	}
	return Comparator{
		analyzerList: analyzerList,
		threads:      runtime.NumCPU(),
		pool:         datatype.NewWorkerPool(),
	}
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

func (c Comparator) Submit(hexBytecode string) datatype.TaskSet {
	// we scan the contract bytecode with all the tools
	taskId := uuid.UUIDv4()
	var ts datatype.TaskSet
	for _, analyzer := range c.analyzerList {
		hexBytecode = c.cleanBytecode(hexBytecode)
		taskList := analyzer.CreateTask(taskId, hexBytecode)
		for _, task := range taskList {
			ts = append(ts, task)
			task.WithResultParser(analyzer)
			c.pool.Submit(task)
		}
	}
	return ts
}

func (c Comparator) SubmitAndWait(hexBytecode string) datatype.TaskSet {
	log.Println("submit and wait")
	var wg sync.WaitGroup

	var ts datatype.TaskSet

	// we scan the contract bytecode with all the tools
	taskId := uuid.UUIDv4()
	for _, analyzer := range c.analyzerList {
		hexBytecode = c.cleanBytecode(hexBytecode)
		taskList := analyzer.CreateTask(taskId, hexBytecode)
		for _, task := range taskList {
			ts = append(ts, task)
			wg.Add(1)
			task.WithResultParser(analyzer)
			c.pool.Submit(task)
			go func() {
				<-task.FinishChan()
				wg.Done()
			}()
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
