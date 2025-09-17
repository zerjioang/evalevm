package datatype

import (
	"bufio"
	"regexp"
	"strings"
)

// CountOrphanNodes reads a Graphviz digraph file and returns the number of nodes
// that are not connected to any other node.
func CountOrphanNodes(dot string) (int, int, error) {
	declared := make(map[string]struct{})
	connected := make(map[string]struct{})

	// Regex patterns
	// Matches: 123[label="..."]
	nodeDecl := regexp.MustCompile(`^\s*"?([0-9A-Za-z_]+)"?\s*\[`)
	// Matches: 123 -> 456
	edgeDecl := regexp.MustCompile(`^\s*"?([0-9A-Za-z_]+)"?\s*->\s*"?([0-9A-Za-z_]+)"?`)
	inLegend := false
	scanner := bufio.NewScanner(strings.NewReader(dot))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}

		// Detect and skip legend subgraph
		if strings.HasPrefix(line, "subgraph cluster_legend") {
			inLegend = true
			continue
		}
		if inLegend {
			if strings.HasPrefix(line, "}") {
				inLegend = false
			}
			continue
		}

		if m := edgeDecl.FindStringSubmatch(line); m != nil {
			src, dst := m[1], m[2]
			connected[src] = struct{}{}
			connected[dst] = struct{}{}
			continue
		}

		if m := nodeDecl.FindStringSubmatch(line); m != nil {
			node := m[1]
			declared[node] = struct{}{}
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	// Collect orphans
	var orphans []string
	for n := range declared {
		if _, ok := connected[n]; !ok {
			orphans = append(orphans, n)
		}
	}
	totalNodes := len(declared)
	return totalNodes, len(orphans), nil
}
