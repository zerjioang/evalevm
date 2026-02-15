package render

import (
	"evalevm/internal/datatype"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func ScanResults(items datatype.TaskSet, transpose bool) error {
	switch len(items) {
	case 0:
		fmt.Println("no result found")
	default:
		if transpose {
			renderTransposedTable(items)
		} else {
			withRenderTable(len(items), func(i int) datatype.Renderable {
				return items[i].Result()
			})
		}
	}
	return nil
}

func ScanSuccess(item datatype.ScanSuccess) error {
	withRenderTable(1, func(i int) datatype.Renderable {
		return item
	})
	return nil
}

func ScanError(item datatype.ScanErrorDetails) error {
	withRenderTable(1, func(i int) datatype.Renderable {
		return item
	})
	return nil
}

func RenderAnalyzers(items []datatype.Analyzer) {
	switch len(items) {
	case 0:
		fmt.Println("no analyzers found")
	default:
		withRenderTable(len(items), func(i int) datatype.Renderable {
			return items[i]
		})
	}
}

func renderTransposedTable(items datatype.TaskSet) {
	if len(items) == 0 {
		return
	}

	// 1. Find the best tool according to performance and quality metrics
	winnerIdx := -1
	maxScore := -1.0
	for i, task := range items {
		score := calculateToolScore(task.Result())
		if score > maxScore {
			maxScore = score
			winnerIdx = i
		}
	}

	// Get headers (features) from the first item
	firstResult := items[0].Result()
	headers := firstResult.Headers()

	// Prepare columns: First column is Feature Name, subsequent are Tool Names
	data := make([][]string, len(headers))
	for i := range data {
		data[i] = make([]string, len(items)+1)
		data[i][0] = headers[i] // Feature Name
	}

	tableHeaders := make([]string, len(items)+1)
	tableHeaders[0] = "METRIC / TOOL"

	for colIdx, task := range items {
		res := task.Result()
		toolName := res.Task.ID().App()
		rows := res.Rows()

		// If this is the winner, highlight header
		if colIdx == winnerIdx {
			tableHeaders[colIdx+1] = color.GreenString(toolName)
		} else {
			tableHeaders[colIdx+1] = toolName
		}

		for rowIdx := 0; rowIdx < len(headers) && rowIdx < len(rows); rowIdx++ {
			val := rows[rowIdx]
			if colIdx == winnerIdx {
				val = color.GreenString(val)
			}
			data[rowIdx][colIdx+1] = val
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header(tableHeaders)
	table.Bulk(data)
	table.Render()
}

func calculateToolScore(res *datatype.Result) float64 {
	// 0. Base checks: failure results in zero score
	if res == nil || res.Error != nil || res.Measurements.ExitStatus != 0 || res.Timeout {
		return 0
	}

	score := 0.0

	// 1. Performance - Time (Lower is better) - Max 100 points
	// Assuming 5s (5000ms) as a slow baseline
	timePoints := 100.0 - (float64(res.Measurements.ExecTimeMs) / 50.0)
	if timePoints < 0 {
		timePoints = 0
	}
	score += timePoints

	// 2. Performance - RAM (Lower is better) - Max 100 points
	// Assuming 1GB as a high consumption baseline for small tasks
	ramPoints := 100.0 - (float64(res.Measurements.MaxRAMKb) / 10240.0)
	if ramPoints < 0 {
		ramPoints = 0
	}
	score += ramPoints

	// 3. Performance - CPU (Balanced) - Max 50 points
	// We favor lower overall cpu usage for efficiency if tasks are simple
	cpuPoints := 50.0 - (res.Measurements.AvgCPUPercent / 2.0)
	if cpuPoints < 0 {
		cpuPoints = 0
	}
	score += cpuPoints

	// 4. Quality - CFG Metrics (Higher is better)
	if res.ParsedOutput != nil && res.ParsedOutput.Metrics != nil {
		m := res.ParsedOutput.Metrics
		// Coverage is king (0-100% -> 0-200 points)
		score += m.CodeCoverage * 2.0

		// Nodes and Edges density (Complexity handled)
		score += float64(m.NodeCount) * 0.5
		score += float64(m.EdgeCount) * 0.5

		// Island penalty
		score -= float64(m.NumDisconnectedSets) * 5.0
	}

	return score
}

func withRenderTable(size int, pickerFn func(i int) datatype.Renderable) {
	table := tablewriter.NewWriter(os.Stdout)
	var rows [][]string
	if size > 0 {
		headers := pickerFn(0).Headers()
		table.Header(headers)
	}
	for i := 0; i < size; i++ {
		item := pickerFn(i)
		if item != nil {
			rows = append(rows, item.Rows())
		}
	}
	_ = table.Bulk(rows)
	_ = table.Render()
}
