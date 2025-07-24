package rattle

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
)

type Rattle struct {
	datatype.BytecodeAnalyzer
}

var (
	//go:embed Dockerfile
	rattleDockerfile string
)

var _ datatype.Analyzer = (*Rattle)(nil)

func NewRattle() Rattle {
	app := Rattle{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "rattle"
	app.WebsiteUrl = "https://github.com/crytic/rattle"
	app.Desc = `Rattle is an EVM binary static analysis framework designed to work on deployed smart contracts. Rattle takes EVM byte strings, uses a flow-sensitive analysis to recover the original control flow graph, lifts the control flow graph into an SSA/infinite register form, and optimizes the SSA – removing DUPs, SWAPs, PUSHs, and POPs. The conversion from a stack machine to SSA form removes 60%+ of all EVM instructions and presents a much friendlier interface to those who wish to read the smart contracts they’re interacting with.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "2 years ago"
	app.Language = "python"
	app.Dockerfile = rattleDockerfile
	return app
}

func (scan Rattle) CreateTask(uid string, bytecode string) []datatype.Task {
	// docker run --rm -it --entrypoint=bash local/rattle -c 'echo $code > contract.evm && python /opt/rattle/rattle-cli.py --no-split-functions --optimize --input contract.evm
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				"run",
				"--rm",
				"--cap-add=SYS_ADMIN",
				"--entrypoint=bash",
				"local/rattle",
				"-c",
				fmt.Sprintf(`echo %s > code.evm && ./measure.sh bash -c 'python /opt/rattle/rattle-cli.py --no-split-functions --optimize --input code.evm'`, bytecode),
			},
		),
	}
}

func (scan Rattle) ParseOutput(output string) *datatype.ScanResult {
	// TODO pending
	return nil
}
