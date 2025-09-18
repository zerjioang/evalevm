package slither

import (
	"evalevm/internal/datatype"
)

type Slither struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Slither)(nil)

func NewSlither() Slither {
	app := Slither{}
	app.AppName = "slither"
	app.WebsiteUrl = "https://github.com/crytic/slither"
	app.Desc = `Slither is a Solidity & Vyper static analysis framework written in Python3. It runs a suite of vulnerability detectors, prints visual information about contract details, and provides an API to easily write custom analyses. Slither enables developers to find vulnerabilities, enhance their code comprehension, and quickly prototype custom analyses.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "3 months ago"
	app.Language = "python"
	return app
}

func (scan Slither) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			[]string{
				"run", "--rm", "--cap-add=SYS_ADMIN", "hello-world",
			},
		),
	}
}

func (scan Slither) ParseOutput(output *datatype.Result) error {
	return nil
}
