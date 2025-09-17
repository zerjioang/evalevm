package render

import (
	"evalevm/internal/datatype"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func ScanResults(items datatype.TaskSet) error {
	switch len(items) {
	case 0:
		fmt.Println("no result found")
	default:
		withRenderTable(len(items), func(i int) datatype.Renderable {
			return items[i].Result()
		})
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

func withRenderTable(size int, pickerFn func(i int) datatype.Renderable) {
	table := tablewriter.NewWriter(os.Stdout)
	var rows [][]string
	table.Header(pickerFn(0).Headers())
	for i := 0; i < size; i++ {
		item := pickerFn(i)
		if item != nil {
			rows = append(rows, item.Rows())
		}
	}
	_ = table.Bulk(rows)
	_ = table.Render()
}
