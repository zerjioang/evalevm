package evm_cfg_builder

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
)

var (
	//go:embed Dockerfile
	evmCfgBuilderDockerfile string
)

type EvmCFGBuilder struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*EvmCFGBuilder)(nil)

func NewEvmCFGBuilder() EvmCFGBuilder {
	app := EvmCFGBuilder{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "evm_cfg_builder"
	app.WebsiteUrl = "https://github.com/crytic/evm_cfg_builder"
	app.Desc = `evm_cfg_builder is used to extract a control flow graph (CFG) from EVM bytecode. It is used by Ethersplay, Manticore, and other tools from Trail of Bits. It is a reliable foundation to build program analysis tools for EVM.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = true
	app.LastCommit = "4 years ago"
	app.Language = "python"
	app.Dockerfile = evmCfgBuilderDockerfile
	return app
}

func (scan EvmCFGBuilder) CreateTask(uid string, bytecode string) []datatype.Task {
	// docker run --rm -it --entrypoint=bash local/evm_cfg_builder -c 'echo 606060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806361461954146044575b600080fd5b3415604e57600080fd5b6054606a565b6040518082815260200191505060405180910390f35b6000806073606a565b8101905080806001019150915050905600a165627a7a723058201cff09a7222fbd72f9f18386a0a03a1a1f02313950b8306cbdb5ce84ed7749c40029 > contract.evm && evm-cfg-builder contract.evm --export-dot out && cat out/contract.evm_-FULL_GRAPH.dot'
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				"run", "--rm", "--cap-add=SYS_ADMIN", "--entrypoint=bash", "local/evm_cfg_builder", "-c",
				fmt.Sprintf(`echo %s > code.evm && ./measure.sh bash -c 'evm-cfg-builder code.evm --export-dot out && cat out/code.evm_-FULL_GRAPH.dot'`, bytecode),
			},
		),
	}
}

func (scan EvmCFGBuilder) ParseOutput(output string) *datatype.ScanResult {
	// TODO pending
	return nil
}
