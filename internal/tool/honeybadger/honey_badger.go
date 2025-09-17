package honeybadger

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	//go:embed Dockerfile
	honeyBadgerDockerfile string
)

type HoneyBadger struct {
	datatype.BytecodeAnalyzer
}

var _ datatype.Analyzer = (*HoneyBadger)(nil)

func NewHoneyBadger() HoneyBadger {
	app := HoneyBadger{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "honeybadger"
	app.WebsiteUrl = "https://github.com/christoftorres/HoneyBadger"
	app.Desc = `A tool that detects honeypots in Ethereum smart contracts 🍯 (USENIX 2019).`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
	}
	app.Deprecated = false
	app.LastCommit = "2 years ago"
	app.Language = "python"
	app.Dockerfile = honeyBadgerDockerfile
	app.PaperURL = "https://arxiv.org/pdf/1902.06976"
	app.Platform = "linux/amd64"
	return app
}

func (scan HoneyBadger) CreateTask(uid string, bytecode string) []datatype.Task {
	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			[]string{
				// docker run command already defined. customize the flags here
				"local/honeybadger",
				"-c",
				fmt.Sprintf(`echo %s > code.evm && ./measure.sh bash -c 'python honeybadger/honeybadger.py -b -s code.evm'`, bytecode),
			},
		),
	}
}

func (scan HoneyBadger) ParseOutput(output *datatype.Result) error {
	_, err := scan.parse(output.Output)
	if err != nil {
		return &datatype.ScanResult{
			Error: err,
		}
	}
	return &datatype.ScanResult{}
}

type SymExecResult struct {
	Coverage    float64
	Results     map[string]bool
	ScanTimeSec float64
}

func (scan HoneyBadger) parse(response string) (*SymExecResult, error) {
	lines := strings.Split(response, "\n")
	result := &SymExecResult{
		Results: make(map[string]bool),
	}

	coverageRegex := regexp.MustCompile(`EVM code coverage:\s+([\d.]+)%`)
	keyValueRegex := regexp.MustCompile(`^\s*(.+?):\s+(True|False)$`)
	timeRegex := regexp.MustCompile(`---\s+([\d.]+)\s+seconds\s+---`)

	for _, line := range lines {
		line = strings.TrimSpace(strings.TrimPrefix(line, "INFO:symExec:"))
		if line == "" {
			continue
		}

		// Match coverage
		if match := coverageRegex.FindStringSubmatch(line); match != nil {
			coverage, err := strconv.ParseFloat(match[1], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid coverage value: %v", err)
			}
			result.Coverage = coverage
			continue
		}

		// Match boolean key-value lines
		if match := keyValueRegex.FindStringSubmatch(line); match != nil {
			key := match[1]
			val := match[2] == "True"
			result.Results[key] = val
			continue
		}

		// Match scan time
		if match := timeRegex.FindStringSubmatch(line); match != nil {
			scanTime, err := strconv.ParseFloat(match[1], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid scan time: %v", err)
			}
			result.ScanTimeSec = scanTime
			continue
		}
	}

	return result, nil
}
