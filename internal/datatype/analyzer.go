package datatype

import (
	"encoding/json"
	"fmt"
	"os"
)

type ResultParser interface {
	ParseOutput(output *Result) error
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
	json.Marshaler
	ID() TaskId
	TrackerId() string
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

func (scan BytecodeAnalyzer) SetupDockerPlatform() BytecodeAnalyzer {
	scan.commonAnalyzerFields.Platform = "linux/arm64"
	return scan
}

func (scan BytecodeAnalyzer) Headers() []string {
	return []string{"name", "url", "deprecated", "last commit", "language"}
}

func (scan BytecodeAnalyzer) Rows() []string {
	return []string{
		scan.Name(), scan.URL(), fmt.Sprintf("%v", scan.Deprecated), scan.LastCommit, scan.Language,
	}
}

func (scan BytecodeAnalyzer) CreateTaskId(uid string) TaskId {
	appuid := TaskId{
		app:        scan.AppName,
		identifier: uid,
	}
	return appuid
}

func (scan BytecodeAnalyzer) DockerfilePath() (string, error) {
	// Create scan temporary file in the default temp directory
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.Dockerfile", scan.AppName))
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Write "hello" into the file
	if _, err := tmpFile.WriteString(scan.Dockerfile); err != nil {
		return "", err
	}

	// Return the full path to the temporary file
	return tmpFile.Name(), nil
}

func (scan BytecodeAnalyzer) DockerPlatform() string {
	return scan.commonAnalyzerFields.Platform
}

func (scan BytecodeAnalyzer) ParseOutput(output *Result) *ScanResult {
	return nil
}
