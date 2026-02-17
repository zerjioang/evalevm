package cmd

import (
	"encoding/csv"
	"evalevm/internal/datatype"
	"fmt"
	"os"
)

// csvExportHeaders defines the column headers for the exported CSV.
var csvExportHeaders = []string{
	"sample_id",
	"status",
	"max_ram_kb",
	"exec_time_ms",
	"avg_cpu_percent",
	"failed_to_parse",
	"exit_status",
	"nodes",
	"edges",
	"mccabe",
	"density",
	"max_out_degree",
	"max_in_degree",
	"avg_out_degree",
	"avg_in_degree",
	"branching_nodes",
	"merge_nodes",
	"entry_points",
	"exit_points",
	"islands",
	"cycles",
	"detected_entry",
	"reachable_nodes",
	"total_nodes",
	"coverage_percent",
	"max_depth",
	"orphan_nodes",
}

// ResultStreamWriter handles streaming export of results to CSV files.
type ResultStreamWriter struct {
	writers map[string]*csv.Writer
	files   map[string]*os.File
}

// NewResultStreamWriter initializes a new writer for streaming results.
func NewResultStreamWriter() *ResultStreamWriter {
	return &ResultStreamWriter{
		writers: make(map[string]*csv.Writer),
		files:   make(map[string]*os.File),
	}
}

// Write streams a batch of results (typically from one contract) to the appropriate CSV files.
func (rsw *ResultStreamWriter) Write(tasks datatype.TaskSet) error {
	// Group tasks by tool name
	grouped := make(map[string][]datatype.Task, len(tasks))
	for _, t := range tasks {
		name := t.ID().App()
		grouped[name] = append(grouped[name], t)
	}

	for toolName, toolTasks := range grouped {
		filename := fmt.Sprintf("evalevm_%s.csv", toolName)

		// Get or create writer for this tool
		w, exists := rsw.writers[toolName]
		if !exists {
			// Check if file exists to decide whether to write headers
			writeHeader := true
			if info, err := os.Stat(filename); err == nil && info.Size() > 0 {
				writeHeader = false
			}

			file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("failed to open CSV file %s: %w", filename, err)
			}
			rsw.files[toolName] = file

			w = csv.NewWriter(file)
			rsw.writers[toolName] = w

			if writeHeader {
				if err := w.Write(csvExportHeaders); err != nil {
					return fmt.Errorf("failed to write CSV headers for %s: %w", toolName, err)
				}
			}
		}

		// Write rows
		for _, task := range toolTasks {
			row := buildCSVRow(task)
			if err := w.Write(row); err != nil {
				return fmt.Errorf("failed to write CSV row for %s: %w", toolName, err)
			}
		}
		w.Flush()
		if err := w.Error(); err != nil {
			return fmt.Errorf("error flushing CSV for %s: %w", toolName, err)
		}
	}
	return nil
}

// Close closes all open file handles.
func (rsw *ResultStreamWriter) Close() error {
	var errs []error
	for name, w := range rsw.writers {
		w.Flush()
		if err := w.Error(); err != nil {
			errs = append(errs, fmt.Errorf("flush error %s: %w", name, err))
		}
	}
	for name, f := range rsw.files {
		if err := f.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close error %s: %w", name, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors closing writers: %v", errs)
	}
	return nil
}

func buildCSVRow(task datatype.Task) []string {
	r := task.Result()

	status := "ok"
	if task.Failed() {
		status = "fail"
	}
	if r.Timeout {
		status = "timeout"
	}

	row := []string{
		task.TrackerId(),
		status,
		fmt.Sprintf("%d", r.Measurements.MaxRAMKb),
		fmt.Sprintf("%d", r.Measurements.ExecTimeMs),
		fmt.Sprintf("%.2f", r.Measurements.AvgCPUPercent),
		fmt.Sprintf("%v", r.Measurements.FailedToParse),
		fmt.Sprintf("%d", r.Measurements.ExitStatus),
	}

	if r.ParsedOutput != nil && r.ParsedOutput.Metrics != nil {
		m := r.ParsedOutput.Metrics
		hasCycles := "false"
		if m.HasCycles {
			hasCycles = "true"
		}
		row = append(row,
			fmt.Sprintf("%d", m.NodeCount),
			fmt.Sprintf("%d", m.EdgeCount),
			fmt.Sprintf("%d", m.CyclomaticComplexity),
			fmt.Sprintf("%.4f", m.GraphDensity),
			fmt.Sprintf("%d", m.MaxOutDegree),
			fmt.Sprintf("%d", m.MaxInDegree),
			fmt.Sprintf("%.2f", m.AvgOutDegree),
			fmt.Sprintf("%.2f", m.AvgInDegree),
			fmt.Sprintf("%d", m.BranchingNodes),
			fmt.Sprintf("%d", m.MergeNodes),
			fmt.Sprintf("%d", m.EntryPoints),
			fmt.Sprintf("%d", m.ExitPoints),
			fmt.Sprintf("%d", m.NumDisconnectedSets),
			hasCycles,
			m.DetectedEntryPoint,
			fmt.Sprintf("%d", m.ReachableNodes),
			fmt.Sprintf("%d", m.NodeCount),
			fmt.Sprintf("%.2f", m.CodeCoverage),
			fmt.Sprintf("%d", m.MaxDepth),
			fmt.Sprintf("%d", len(m.OrphanNodes)),
		)
	} else {
		// Fill with empty values for tasks without metrics
		empty := make([]string, len(csvExportHeaders)-7) // 7 base columns already written
		for i := range empty {
			empty[i] = ""
		}
		row = append(row, empty...)
	}

	return row
}
