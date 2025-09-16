package ethersolve

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
)

type EthersolveCreator struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*EthersolveCreator)(nil)

func NewEthersolveCreator() EthersolveCreator {
	app := EthersolveCreator{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "ethersolve_creator"
	app.WebsiteUrl = "https://github.com/SeUniVr/EtherSolve"
	app.Desc = `EtherSolve is a tool for Control Flow Graph (CFG) reconstruction and static analysis of Solidity smart-contracts from Ethereum bytecode.`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: true,
		ForceSplitRuntime:    true,
	}
	app.Deprecated = false
	app.LastCommit = "2 years ago"
	app.Language = "java"
	app.Dockerfile = ethersolveDockerfile
	return app
}

func (scan EthersolveCreator) CreateTask(uid string, bytecode string) []datatype.Task {
	// docker run --rm --entrypoint=bash local/ethersolve -c 'java -jar /opt/ethersolve/artifact/EtherSolve.jar --runtime --tx-origin --re-entrancy --dot 608060405234801561001057600080fd5b50610179806100206000396000f3fe608060405234801561001057600080fd5b50600436106100455760003560e01c80632a1bbc3414610061578063d96073cf1461007f578063e7a721cf146100d257610046565b5b6000618888905061999981101561005e576001810190505b50005b61006961011b565b6040518082815260200191505060405180910390f35b6100b56004803603604081101561009557600080fd5b810190808035906020019092919080359060200190929190505050610125565b604051808381526020018281526020019250505060405180910390f35b6100fe600480360360208110156100e857600080fd5b8101908080359060200190929190505050610135565b604051808381526020018281526020019250505060405180910390f35b600061cccc905090565b6000808284915091509250929050565b60008082839150915091509156fea264697066735822122064dcf93052eeac5706ef56f301b5cfbd23c32458ac9b5e4c191e85bc70003d5d64736f6c63430006060033 && cat /opt/ethersolve/Analysis_*'
	return []datatype.Task{
		// creation block
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				"run",
				"--rm",
				"--cap-add=SYS_ADMIN",
				"--entrypoint=bash",
				"local/ethersolve_creator",
				"-c",
				fmt.Sprintf(`./helper.sh ethersolve_creator %s`, bytecode),

				// output order: re-entrancy, tx-origin, dot file
			},
		),
	}
}

func (scan EthersolveCreator) ParseOutput(output *datatype.Result) error {
	return parseEthersolveOutput(output)
}
