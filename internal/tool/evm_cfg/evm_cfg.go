package evm_cfg

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
)

var (
	//go:embed Dockerfile
	evmCFGDockerfile string
)

type EvmCFG struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*EvmCFG)(nil)

func NewEvmCFG() EvmCFG {
	app := EvmCFG{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "evm-cfg"
	app.WebsiteUrl = "https://github.com/plotchy/evm-cfg"
	app.Desc = `Symbolic stack CFG generator for EVM`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = true
	app.LastCommit = "5 months ago"
	app.Language = "rust"
	app.Dockerfile = evmCFGDockerfile
	return app
}

func (scan EvmCFG) CreateTask(uid string, bytecode string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				"run", "--rm", "--cap-add=SYS_ADMIN", "--entrypoint=bash", "local/evm-cfg", "-c",
				fmt.Sprintf(`./measure.sh bash -c '/opt/evm-cfg/evm-cfg %s -o cfg.dot && cat cfg.dot'`, bytecode),
			},
		),
	}
}

func (scan EvmCFG) ParseOutput(output string) *datatype.ScanResult {
	// TODO pending
	return nil
}
