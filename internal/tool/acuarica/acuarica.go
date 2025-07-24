package acuarica

import (
	"evalevm/internal/datatype"
)

type AcuaricaEVM struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*AcuaricaEVM)(nil)

func NewAcuaricaEVM() AcuaricaEVM {
	app := AcuaricaEVM{}
	app.AppName = "acuarica-evm"
	app.WebsiteUrl = "https://github.com/acuarica/evm"
	app.Desc = `A Symbolic Ethereum Virtual Machine (EVM) interpreter and decompiler, along with several other utils for programmatically extracting information from bytecode.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = true
	app.LastCommit = "6 months ago"
	app.Language = "typescript"
	return app
}

func (scan AcuaricaEVM) CreateTask(uid string, bytecode string) []datatype.Task {
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
