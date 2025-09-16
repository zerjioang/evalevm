package panoramix_palkeo

import (
	"evalevm/internal/datatype"
)

type PanoramixPalkeo struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*PanoramixPalkeo)(nil)

func NewPanoramixPalkeo() PanoramixPalkeo {
	app := PanoramixPalkeo{}
	app.AppName = "panoramix-palkeo"
	app.WebsiteUrl = "https://github.com/palkeo/panoramix"
	app.Desc = `Ethereum decompiler based on panoramix original implementation`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "2 years ago"
	app.Language = "python"
	return app
}

func (scan PanoramixPalkeo) CreateTask(uid string, bytecode string) []datatype.Task {
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

func (scan PanoramixPalkeo) ParseOutput(output *datatype.Result) error {
	return nil
}
