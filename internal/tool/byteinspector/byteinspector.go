package byteinspector

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
)

var (
	//go:embed Dockerfile
	byteInspectorDockerfile string
)

type ByteInspector struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*ByteInspector)(nil)

func NewByteInspector() ByteInspector {
	app := ByteInspector{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "byte-inspector"
	app.WebsiteUrl = "https://github.com/franck44/evm-dis"
	app.Desc = `This project provides an EVM bytecode disassembler and Control Flow Graph (CFG) generator. ByteSpector can verify the CFGs by generating a Dafny file that encodes the semantics of the EVM bytecode. The Dafny file can be verified with Dafny. If a CFG is successfully verified, we obtain the following guarantees on the CFG and the bytecode`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "7 months ago"
	app.Language = "dafny"
	app.Dockerfile = byteInspectorDockerfile
	app.Platform = "linux/amd64"
	return app
}

func (scan ByteInspector) CreateTask(uid string, bytecode string) []datatype.Task {
	// docker run --rm -it --entrypoint=bash local/byte-inspector -c "echo 0x366028576000600060006000303173f43febf30d4a00fa9b23e49e36e7acb5ca8591616103e8f1005b6388c2a0bf60e060020a026000526000358043116077574390036001016003023562ffffff16600452600060006024600060007306012c8cf97bead5deae237070f9587f8e7a266d6103e85a03f15b00 > code.evm && /tacas25/evm-dis/measure.sh bash -c '/tacas25/evm-dis/makeCFG.sh code.evm && cat build/dot/code.evm/code.evm.dot'"
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				"run",
				"--rm",
				"--platform",
				"linux/amd64",
				"--cap-add=SYS_ADMIN",
				"--entrypoint=bash",
				"local/byte-inspector",
				"-c",
				fmt.Sprintf(`echo %s > code.evm && ./measure.sh bash -c 'cd /tacas25/evm-dis/ && ./makeCFG.sh code.evm && cat build/dot/code.evm/code.evm.dot'`, bytecode),
			},
		),
	}
}

func (scan ByteInspector) ParseOutput(output *datatype.Result) error {
	return nil
}
