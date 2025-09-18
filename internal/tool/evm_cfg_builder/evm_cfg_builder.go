package evm_cfg_builder

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
	"log"
	"strings"
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

func (scan EvmCFGBuilder) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			[]string{
				// docker run command already defined. customize the flags here
				"local/evm_cfg_builder", "-c",
				fmt.Sprintf(`echo 0x%s > code.evm && ./measure.sh bash -c 'evm-cfg-builder code.evm --export-dot out && cat out/code.evm_-FULL_GRAPH.dot'`, bytecode),
			},
		),
	}
}

func (scan EvmCFGBuilder) ParseOutput(output *datatype.Result) error {
	outStr := string(output.Output)
	output.ParsedOutput = &datatype.ScanResult{
		EdgesDetected: strings.Count(outStr, "->"),
		NodesDetected: strings.Count(outStr, "[label="),
		DotGraph:      outStr,
	}
	output.ParsedOutput.WithGraph(outStr)
	filename := fmt.Sprintf("cfg_%s_%s.svg", output.Task.ID().App(), output.Task.TrackerId())
	if err := output.ParsedOutput.SaveGraph(outStr, filename); err != nil {
		log.Println("failed to save graph: ", err)
		return err
	}
	output.AddFileReference(output.Task.ID().App(), output.Task.TrackerId(), filename)
	return nil
}
