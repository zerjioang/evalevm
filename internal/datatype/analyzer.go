package datatype

import (
	"fmt"
	"os"
)

type ResultParser interface {
	ParseOutput(output string) *ScanResult
}
type Analyzer interface {
	ResultParser
	Renderable
	Name() string
	URL() string
	CreateTask(uid string, bytecode string) []Task
	DockerfilePath() (string, error)
	DockerPlatform() string
}

type Task interface {
	ID() TaskId
	Command() []string
	FinishChan() chan struct{}
	Result() *Result
	WithResult(result *Result)
	WithResultParser(parser ResultParser)
	Parse()
	Failed() bool
}

type TaskSet []Task

type commonAnalyzerFields struct {
	AppName    string
	WebsiteUrl string
	Desc       string
	Deprecated bool
	LastCommit string
	Language   string
	PaperURL   string
	Creator    string
	// docker platform: linux/amd64 or "linux/arm64
	Platform string
	WorkDir  string
}

func (r commonAnalyzerFields) Name() string {
	return r.AppName
}

func (r commonAnalyzerFields) URL() string {
	return r.WebsiteUrl
}

func (r commonAnalyzerFields) Description() string {
	return r.Desc
}

type BytecodeScanOpts struct {
	ForceRemoveHexPrefix bool
	ForceSplitRuntime    bool
}

type BytecodeAnalyzer struct {
	commonAnalyzerFields
	Options    BytecodeScanOpts
	Dockerfile string
}

func (a BytecodeAnalyzer) SetupDockerPlatform() BytecodeAnalyzer {
	a.commonAnalyzerFields.Platform = "linux/arm64"
	return a
}

func (a BytecodeAnalyzer) Headers() []string {
	return []string{"name", "url", "deprecated", "last commit", "language"}
}

func (a BytecodeAnalyzer) Rows() []string {
	return []string{
		a.Name(), a.URL(), fmt.Sprintf("%v", a.Deprecated), a.LastCommit, a.Language,
	}
}

func (a BytecodeAnalyzer) CreateTaskId(uid string) TaskId {
	appuid := TaskId{
		app:        a.AppName,
		identifier: uid,
	}
	return appuid
}

func (a BytecodeAnalyzer) DockerfilePath() (string, error) {
	// Create a temporary file in the default temp directory
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.Dockerfile", a.AppName))
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Write "hello" into the file
	if _, err := tmpFile.WriteString(a.Dockerfile); err != nil {
		return "", err
	}

	// Return the full path to the temporary file
	return tmpFile.Name(), nil
}

func (a BytecodeAnalyzer) DockerPlatform() string {
	return a.commonAnalyzerFields.Platform
}

func (a BytecodeAnalyzer) ParseOutput(output string) *ScanResult {
	return nil
}
