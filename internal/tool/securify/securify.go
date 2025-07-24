package securify

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
)

type Securify struct {
	datatype.BytecodeAnalyzer
}

var (
	//go:embed Dockerfile
	securifyDockerfile string
)

var _ datatype.Analyzer = (*Securify)(nil)

func NewSecurify() Securify {
	app := Securify{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "securify"
	app.WebsiteUrl = "https://github.com/eth-sri/securify"
	app.Desc = `Security Scanner for Ethereum Smart Contracts`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: true,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = true
	app.LastCommit = "6 years ago"
	app.Language = "java"
	app.PaperURL = "https://files.sri.inf.ethz.ch/website/papers/ccs18-securify.pdf"
	app.Creator = "https://www.chainsecurity.com/"
	app.Dockerfile = securifyDockerfile
	app.Platform = "linux/amd64"
	return app
}

func (scan Securify) CreateTask(uid string, bytecode string) []datatype.Task {
	// docker run --rm -it --platform linux/amd64 --entrypoint=bash docker.io/local/securify:latest -c 'echo 0x60008080803473d7c02a75d24e5a0f8140488877874cd880dafe155af1602457600080fd5b00 > contract.evm && java -jar /securify_jar/securify.jar -fh contract.evm
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				"run", "--rm", "--cap-add=SYS_ADMIN",
				"--platform", scan.Platform, "--entrypoint=bash", "local/securify", "-c",
				fmt.Sprintf(`echo %s > code.evm && ./measure.sh bash -c 'java -Xms512m -Xmx2048m -jar /securify_jar/securify.jar -fh code.evm'`, bytecode),
			},
		),
	}
}
