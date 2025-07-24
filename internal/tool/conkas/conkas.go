package conkas

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
)

var (
	//go:embed Dockerfile
	conkasDockerfile string
)

type Conkas struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*Conkas)(nil)

func NewConkas() Conkas {
	app := Conkas{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "conkas"
	app.WebsiteUrl = "https://github.com/nveloso/conkas"
	app.Desc = `Ethereum Virtual Machine (EVM) Bytecode or Solidity Smart Contract static analysis tool based on symbolic execution	`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "3 years ago"
	app.Language = "python"
	app.Dockerfile = conkasDockerfile
	return app
}

func (scan Conkas) CreateTask(uid string, bytecode string) []datatype.Task {
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
				"local/conkas",
				"-c",
				fmt.Sprintf(`echo %s > code.evm && ./measure.sh bash -c 'python /opt/conkas/conkas.py -fav code.evm'`, bytecode),
			},
		),
	}
}

func (scan Conkas) ParseOutput(output string) *datatype.ScanResult {
	// TODO pending
	return nil
}
