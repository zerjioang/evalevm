package vandal

import (
	"evalevm/internal/datatype"
)

type Vandal struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Vandal)(nil)

func NewVandal() Vandal {
	app := Vandal{}
	app.AppName = "vandal"
	app.WebsiteUrl = "https://github.com/usyd-blockchain/vandal"
	app.Desc = `Vandal is a static program analysis framework for Ethereum smart contract bytecode, developed at The University of Sydney. It decompiles an EVM bytecode program to an equivalent intermediate representation that encodes the program's control flow graph. This representation removes all stack operations, thereby exposing data dependencies that are otherwise obscured. This information is then fed, with a Datalog specification, into the Souffle analysis engine for the extraction of program properties.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "5 years ago"
	app.Language = "python"
	return app
}

func (scan Vandal) CreateTask(uid string, bytecode string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				"run", "--rm", "--cap-add=SYS_ADMIN", "hello-world",
			},
		),
	}
}

func (scan Vandal) ParseOutput(output *datatype.Result) error {
	return nil
}
