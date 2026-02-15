package ethersolve

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
)

var (
	//go:embed Dockerfile
	ethersolveDockerfile string
)

type Ethersolve struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Ethersolve)(nil)

func NewEthersolve(runMode string, audit bool) Ethersolve {
	app := Ethersolve{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "ethersolve"
	app.WebsiteUrl = "https://github.com/SeUniVr/EtherSolve"
	app.Desc = `EtherSolve is a tool for Control Flow Graph (CFG) reconstruction and static analysis of Solidity smart-contracts from Ethereum bytecode.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: true,
		ForceSplitRuntime:    true,
		RunMode:              runMode,
		Audit:                audit,
	}
	app.Deprecated = false
	app.LastCommit = "2 years ago"
	app.Language = "java"
	app.Dockerfile = ethersolveDockerfile
	app.SupportsVulnerabilities = true
	app.SupportsCFG = true
	return app
}

func (scan Ethersolve) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	var modes []string
	if scan.Options.RunMode == "any" || scan.Options.RunMode == "" {
		modes = []string{"runtime", "creator"}
	} else {
		modes = []string{scan.Options.RunMode}
	}

	var tasks []datatype.Task
	for _, mode := range modes {
		// Use a distinct ID for each mode if multiple are running, or just append mode to app name?
		// TaskId struct has app field. I can fake it or just rely on sampleId?
		// WorkerPool uses TaskId to deduplicate?
		// Let's rely on the list of tasks.
		// NOTE: If I return multipl
		taskId := scan.CreateTaskId(uid)
		// We can't easily change TaskId app name without breaking things maybe?
		// Analyzer interface returns []Task.
		// TaskId is {app, identifier}.
		// If I run multiple, I should probably distinguish them.

		// Append mode to filename/sampleId to avoid collision in result filenames
		// If filename is empty, NewDockerTask would use hash. We can pass mode as prefix or suffix.
		// If filename is provided (directory scan), we append mode.
		// If not provided (single scan), we append mode to make it distinct.
		distinctFilename := filename
		if distinctFilename == "" {
			distinctFilename = mode
		} else {
			distinctFilename = fmt.Sprintf("%s_%s", filename, mode)
		}

		t := datatype.NewDockerTask(
			taskId,
			scan.Options,
			bytecode,
			distinctFilename,
			[]string{
				"local/ethersolve",
				"-c",
				fmt.Sprintf(`./helper.sh ethersolve %s %v %s`, mode, scan.Options.Audit, bytecode),
			},
		)
		tasks = append(tasks, t)
	}
	return tasks
}

func (scan Ethersolve) ParseOutput(output *datatype.Result) error {
	return parseEthersolveOutput(output)
}
