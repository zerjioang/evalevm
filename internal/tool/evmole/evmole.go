package evmole

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
)

type EVMole struct {
	datatype.BytecodeAnalyzer
}

var (
	//go:embed Dockerfile
	evmoleDockerfile string
)
var _ datatype.Analyzer = (*EVMole)(nil)

func NewEVMole() EVMole {
	app := EVMole{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "evmole"
	app.WebsiteUrl = "https://github.com/cdump/evmole"
	app.Desc = `Extracts function selectors, arguments, state mutability and storage layout from EVM bytecode, even for unverified contracts`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = true
	app.LastCommit = "last week"
	app.Language = "rust"
	app.Dockerfile = evmoleDockerfile
	return app
}

func (scan EVMole) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			[]string{
				// docker run command already defined. customize the flags here
				"local/evmole", "-c",
				fmt.Sprintf(`./measure.sh bash -c 'python3 /opt/evmole/run.py %s'`, bytecode),
			},
		),
	}
}

func (scan EVMole) ParseOutput(output *datatype.Result) error {
	return nil
}
