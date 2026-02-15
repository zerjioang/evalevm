package evmlisa

import (
	_ "embed"
	"evalevm/internal/datatype"
	"fmt"
	"log"
	"os"

	"github.com/zerjioang/rooftop/v2/common/io/json"
)

var (
	//go:embed Dockerfile
	evmLisaDockerfile string
)

type EvmLisa struct {
	datatype.BytecodeAnalyzer
	WorkDir string
}

var _ datatype.Analyzer = (*EvmLisa)(nil)

func NewEvmLisa(audit bool) EvmLisa {
	app := EvmLisa{}
	app.BytecodeAnalyzer = app.SetupDockerPlatform()
	app.AppName = "evm-lisa"
	app.WebsiteUrl = "https://github.com/lisa-analyzer/evm-lisa"
	app.Desc = `EVMLiSA: an abstract interpretation-based static analyzer for EVM bytecode`
	app.Options = datatype.BytecodeScanOpts{
		ForceRemoveHexPrefix: false,
		ForceSplitRuntime:    false,
		Audit:                audit,
	}
	app.Deprecated = false
	app.LastCommit = "2 months ago"
	app.Language = "java"
	app.Creator = ""
	app.PaperURL = "https://vincenzoarceri.github.io/papers/ftfjp2024.pdf"
	app.Dockerfile = evmLisaDockerfile
	app.SupportsVulnerabilities = true
	app.SupportsCFG = true
	return app
}

func (scan EvmLisa) CreateTask(uid string, bytecode string, filename string) []datatype.Task {
	// docker run --rm -v $(pwd)/.env:/app/.env -v $(pwd)/execution/docker:/app/execution/results docker.io/library/evm-lisa:latest --bytecode 0x366028576000600060006000303173f43febf30d4a00fa9b23e49e36e7acb5ca8591616103e8f1005b6388c2a0bf60e060020a026000526000358043116077574390036001016003023562ffffff16600452600060006024600060007306012c8cf97bead5deae237070f9587f8e7a266d6103e85a03f15b00 --stack-size 1024 --stack-set-size 1024 --checker-all
	// Define volume mounts based on the scanner working directory
	if scan.WorkDir == "" {
		path, err := scan.createTempDir("evm-lisa")
		if err != nil {
			log.Printf("evmlisa: failed to create temp dir: %v", err)
			return nil
		}
		scan.WorkDir = path
	}
	envMount := fmt.Sprintf("%s/.env:/app/.env", scan.WorkDir)
	resultsMount := fmt.Sprintf("%s/execution/docker:/app/execution/results", scan.WorkDir)

	// Assemble Docker command arguments
	dockerArgs := []string{
		// docker run command already defined. customize the flags here
		"-v", envMount,
		"-v", resultsMount,
		"local/evm-lisa",
		"-c",
	}

	cmd := fmt.Sprintf(`./helper.sh evmlisa %s`, bytecode)
	if scan.Options.Audit {
		cmd += " --checker-all"
	}
	dockerArgs = append(dockerArgs, cmd)

	return []datatype.Task{
		datatype.NewDockerTask(
			scan.CreateTaskId(uid),
			scan.Options,
			bytecode,
			filename,
			dockerArgs,
		),
	}
}

// CreateTempDir creates a temporary directory under the scanner's working directory.
// The prefix string is used to generate the directory name.
func (scan EvmLisa) createTempDir(prefix string) (string, error) {
	// Use os.MkdirTemp to create a temp directory under WorkDir
	tmpDir, err := os.MkdirTemp(scan.WorkDir, prefix)
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	return tmpDir, nil
}

func (scan EvmLisa) ParseOutput(result *datatype.Result) error {

	files, err := parseOffsetData(string(result.Output))
	if err != nil {
		return err
	}

	jsonData := files["results.json"]
	dotGraph := files["CFG.dot"]

	var dst output
	if err := json.Unmarshal(jsonData.Content, &dst); err != nil {
		return err
	}

	var asPtrBool = func(b bool) *bool { return &b }
	vulnerabilities := dst.Vulnerabilities.TxOrigin + dst.Vulnerabilities.RandomnessDependency + dst.Vulnerabilities.Reentrancy
	edges := 0
	for _, bb := range dst.BasicBlocks {
		edges += len(bb.OutgoingEdges)
	}
	result.ParsedOutput = &datatype.ScanResult{
		Vulnerable:           asPtrBool(vulnerabilities > 0),
		TxOriginVulnerable:   asPtrBool(dst.Vulnerabilities.TxOrigin > 0),
		ReEntrancyVulnerable: asPtrBool(dst.Vulnerabilities.Reentrancy > 0),
		EdgesDetected:        edges,
		NodesDetected:        len(dst.BasicBlocks),
	}

	if err := result.ParsedOutput.WithGraph(string(dotGraph.Content), "", result); err != nil {
		return fmt.Errorf("failed to store .dot graph: %w", err)
	}

	return nil
}
