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
	CreateTask(uid string, bytecode string, filename string) []Task
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
	ContainerName() string // Added for robust container management
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
	Platform                string
	WorkDir                 string
	SupportsVulnerabilities bool
	SupportsCFG             bool
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
	Audit                bool
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
	return []string{"name", "url", "deprecated", "last commit", "language", "Vulns", "CFG"}
}

func (scan BytecodeAnalyzer) Rows() []string {
	return []string{
		scan.Name(),
		scan.URL(),
		fmt.Sprintf("%v", scan.Deprecated),
		scan.LastCommit,
		scan.Language,
		fmt.Sprintf("%v", scan.SupportsVulnerabilities),
		fmt.Sprintf("%v", scan.SupportsCFG),
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
	// Create a temporary Dockerfile from the embedded content.
	// Note: the caller (Builder.Build) is responsible for cleanup after docker build.
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.Dockerfile", scan.AppName))
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(scan.Dockerfile); err != nil {
		return "", err
	}

	// Ensure content is flushed to disk before the path is used by docker build
	if err := tmpFile.Sync(); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func (scan BytecodeAnalyzer) DockerPlatform() string {
	return scan.commonAnalyzerFields.Platform
}

func (scan BytecodeAnalyzer) ParseOutput(output *Result) *ScanResult {
	return nil
}
