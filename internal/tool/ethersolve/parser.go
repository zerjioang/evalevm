package ethersolve

import (
	"bufio"
	"evalevm/internal/datatype"
	"fmt"
	"log"
	"strings"
)

// parseEthersolveOutput parses the output of the ethersolve tool and returns a ScanResult
func parseEthersolveOutput(output *datatype.Result) error {
	// count how many nodes are defined in the CFG
	// example: 119 [label="119: EXIT BLOCK\l" fillcolor=crimson ];
	outStr := string(output.Output)
	nodesDetected := strings.Count(outStr, " [label=")

	// count how many edges are defined in the CFG
	// example: 119 -> 118 [label="119 -> 118\l" ];
	edgesDetected := strings.Count(outStr, " -> ")

	// check if at least one analyzer is executed
	findingsDetected := false
	txOriginVulnerable := false
	reEntrancyVulnerable := false
	var dotGraph string
	// reentrancy or tx-origin findings. parse them
	offsetData, err := parseOffsetData(outStr)
	if err != nil {
		log.Printf("failed to parse offset data: %v", err)
	}
	for _, datum := range offsetData {
		// a valid vunerability finding has at least 2 lines: header + at least one finding
		// if only the header is present, it means no findings were detected
		// example with findings:
		// >>> Analysis_re-entrancy
		// offset,opcode,detection
		// ...
		txOriginVulnerable = txOriginVulnerable || (datum.IsTxOrigin && len(datum.Content) > 1)
		reEntrancyVulnerable = reEntrancyVulnerable || (datum.IsReEntrancy && len(datum.Content) > 1)
		if datum.IsDotFile {
			dotGraph = strings.Join(datum.Content, "\n")
		}
	}
	findingsDetected = txOriginVulnerable || reEntrancyVulnerable
	log.Println("findings detected: ", findingsDetected)

	//var asPtrBool = func(b bool) *bool { return &b }
	output.ParsedOutput = &datatype.ScanResult{
		Vulnerable:           nil, //asPtrBool(findingsDetected),
		Error:                nil,
		EdgesDetected:        edgesDetected,
		NodesDetected:        nodesDetected,
		TxOriginVulnerable:   nil, //asPtrBool(txOriginVulnerable),
		ReEntrancyVulnerable: nil, //asPtrBool(reEntrancyVulnerable),
	}
	if err := output.ParsedOutput.WithGraph(dotGraph, "", output); err != nil {
		return fmt.Errorf("failed to store .dot graph: %w", err)
	}

	return nil
}

type parsedOffsetFile struct {
	File         string
	Content      []string // raw lines
	IsTxOrigin   bool
	IsReEntrancy bool
	IsDotFile    bool
}

func parseOffsetData(data string) ([]parsedOffsetFile, error) {
	var results []parsedOffsetFile
	var current *parsedOffsetFile

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, ">>> "):
			// Finish previous file if any
			if current != nil {
				results = append(results, *current)
			}
			filename := strings.TrimSpace(strings.TrimPrefix(line, ">>> "))
			current = &parsedOffsetFile{
				File:         filename,
				IsTxOrigin:   strings.Contains(filename, "tx-origin"),
				IsReEntrancy: strings.Contains(filename, "re-entrancy"),
				IsDotFile:    strings.Contains(filename, "dot"),
			}

		case strings.HasPrefix(line, "<<<"):
			if current != nil {
				results = append(results, *current)
				current = nil
			}

		default:
			if current != nil && len(strings.TrimSpace(line)) > 0 {
				current.Content = append(current.Content, line)
			}
		}
	}

	// Add last file if not closed with <<<
	if current != nil {
		results = append(results, *current)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
