package datatype

import (
	"evalevm/internal/export"
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

func (s *ScanResult) WithGraph(dot string) {

	if dot == "" {
		return
	}
	s.DotGraph = dot
	s.CFGCreated = len(dot) > 0

	var nodesRead int
	nodesRead, s.OrphanNodesInGraph, _ = CountOrphanNodes(dot)
	if s.NodesDetected == 0 {
		s.NodesDetected = nodesRead
	}
	if s.Coverage == nil || *s.Coverage == 0 {
		// calculate coverage based on .dot file data instead of tool result
		readCoverage := float64((nodesRead-s.OrphanNodesInGraph)*100) / float64(nodesRead)
		s.Coverage = &readCoverage
	}
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
