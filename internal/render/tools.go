package render

import (
	"evalevm/internal/datatype"
	"fmt"
	"os"

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

		tableHeaders[colIdx+1] = toolName

		for rowIdx := 0; rowIdx < len(headers) && rowIdx < len(rows); rowIdx++ {
			data[rowIdx][colIdx+1] = rows[rowIdx]
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header(tableHeaders)
	table.Bulk(data)
	table.Render()
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
