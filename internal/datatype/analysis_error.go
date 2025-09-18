package datatype

import (
	"evalevm/internal/export"
	"fmt"
	"log"
	"time"

	"github.com/zerjioang/rooftop/v2/common/io/json"
)

type ScanResult struct {
	// Vulnerable flag set to true if the analysis found at least one vulnerability
	Vulnerable           *bool    `json:"vulnerable,omitempty"`
	Error                error    `json:"error,omitempty"`
	EdgesDetected        int      `json:"edges,omitempty"`
	NodesDetected        int      `json:"nodes,omitempty"`
	OrphanNodesInGraph   int      `json:"orphan_nodes,omitempty"`
	TxOriginVulnerable   *bool    `json:"tx_origin_vulnerable,omitempty"`
	ReEntrancyVulnerable *bool    `json:"re_entrancy_vulnerable,omitempty"`
	CFGCreated           bool     `json:"cfg_created,omitempty"`
	DotGraph             string   `json:"dot_graph,omitempty"`
	Coverage             *float64 `json:"coverage,omitempty"`
}

func (s *ScanResult) WithGraph(dot string, fileId string, result *Result) error {

	if dot == "" {
		return nil
	}

	task := result.Task
	filename := fmt.Sprintf("cfg_%s_%s_%s.svg", task.ID().App(), task.TrackerId(), fileId)
	if err := s.SaveGraph(dot, filename); err != nil {
		log.Println("failed to save graph: ", err)
		return err
	}
	result.AddFileReference(task.ID().App(), task.TrackerId(), filename)

	s.DotGraph = dot
	s.CFGCreated = len(dot) > 0

	if s.Coverage == nil || *s.Coverage == 0 {
		var nodesRead int
		nodesRead, s.OrphanNodesInGraph, _ = CountOrphanNodes(dot)
		if s.NodesDetected == 0 {
			s.NodesDetected = nodesRead
		}
		// calculate coverage based on .dot file data instead of tool result
		readCoverage := float64((nodesRead-s.OrphanNodesInGraph)*100) / float64(nodesRead)
		s.Coverage = &readCoverage
	}

	return nil
}

func (s *ScanResult) SaveGraph(dot string, filename string) error {
	log.Println("Saving graph to file")
	return export.GenerateSVGFromDot(dot, filename, 10*time.Second)
}

func (s *ScanResult) String() string {
	if s == nil {
		return ""
	}
	dot := s.DotGraph
	s.DotGraph = ""
	raw, _ := json.Marshal(s)
	s.DotGraph = dot
	return string(raw)
}

type ScanErrorDetails struct {
	Name    string
	Message string
}

var _ Renderable = (*ScanErrorDetails)(nil)

func (s ScanErrorDetails) Headers() []string {
	return []string{"name", "error message"}
}

func (s ScanErrorDetails) Rows() []string {
	return []string{
		boldRed.Sprint(s.Name),
		s.Message,
	}
}
