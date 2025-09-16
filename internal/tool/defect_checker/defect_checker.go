package defect_checker

import (
	"evalevm/internal/datatype"
)

type DefectChecker struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*DefectChecker)(nil)

func NewDefectChecker() DefectChecker {
	app := DefectChecker{}
	app.AppName = "defect-checker"
	app.WebsiteUrl = "https://github.com/Jiachi-Chen/DefectChecker"
	app.Desc = `Automated Smart Contract Defect Detection by Analyzing EVM Bytecode`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "5 years ago"
	app.Language = "python"
	app.PaperURL = "https://www4.comp.polyu.edu.hk/~csxluo/DEFECTCHECKER.pdf"
	return app
}

func (scan DefectChecker) CreateTask(uid string, bytecode string) []datatype.Task {
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

func (scan DefectChecker) ParseOutput(output *datatype.Result) error {
	return nil
}
