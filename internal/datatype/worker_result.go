package datatype

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/zerjioang/rooftop/v2/common/io/json"
)

var (
	boldGreen = color.New(color.FgGreen, color.Bold)
	boldWhite = color.New(color.FgHiWhite, color.Bold)
	gray      = color.New(color.FgWhite, color.Faint)
	boldRed   = color.New(color.FgRed, color.Bold)
)

type Measurements struct {
	FailedToParse   bool    `json:"failed_to_parse"`
	MaxRAMKb        int     `json:"max_ram_kb"`
	ExecTimeUs      int     `json:"exec_time_us"`
	ExecTimeMs      int     `json:"exec_time_ms"`
	ExecTimeS       float64 `json:"exec_time_s"`
	ExitStatus      int     `json:"exit_status"`
	AvgCPUPercent   float64 `json:"avg_cpu_percent"`
	Instructions    int     `json:"instructions"`
	CPUCycles       int     `json:"cpu_cycles"`
	ContextSwitches int     `json:"context_switches"`
	PageFaults      int     `json:"page_faults"`
	BranchMisses    int     `json:"branch_misses"`
}

type fileInfo struct {
	App      string `json:"app"`
	SampleId string `json:"sample_id"`
	Filename string `json:"filename"`
}

// Result represents the result of the Docker command
type Result struct {
	// Input data
	Task Task `json:"task"`
	// Raw output of the executed command
	Output []byte `json:"output"`
	// Error output of the executed command
	OutputErr        []byte        `json:"output_error"`
	Error            error         `json:"error"`
	TotalElapsedTime time.Duration `json:"total_elapsed_time"`
	Measurements     Measurements  `json:"measurements"`
	// Parsed output of the executed command
	ParsedOutput *ScanResult `json:"parsed_output"`
	// Additional reference files
	Files []fileInfo `json:"files"`
}

func (r *Result) Headers() []string {
	return []string{
		"app",
		"sample_id",
		"status",
		"error",
		"max_ram_kb",
		"exec_time_ms",
		"exec_time_s",
		"exit_code",
		"avg_cpu_percent",
		"scan_result",
	}
}

func (r *Result) Rows() []string {
	errstr := "-"
	status := boldGreen.Sprint("✔ success")
	if r.Error != nil {
		errstr = r.Error.Error()
		status = boldRed.Sprint("❌ errored")
	}
	if r.Measurements.ExitStatus != 0 {
		errstr = fmt.Sprintf("exit status %d", r.Measurements.ExitStatus)
		status = boldRed.Sprint("❌ errored")
	}
	parsedResult := "pendiente de parsear"
	if r.Task.Failed() {
		errMsg := "tool failed"
		if r.OutputErr != nil {
			errMsg = r.Error.Error()
			if len(r.OutputErr) > 40 {
				errMsg = string(r.OutputErr[0:40]) + "..."
			}
		}
		parsedResult = errMsg
	} else if r.ParsedOutput != nil {
		parsedResult = r.ParsedOutput.String()
	} else {
		parsedResult = "parser not implemented"
	}
	return []string{
		boldWhite.Sprint(r.Task.ID().app),
		boldWhite.Sprint(r.Task.TrackerId()),
		status,
		errstr,
		fmt.Sprintf("%d", r.Measurements.MaxRAMKb),
		fmt.Sprintf("%d", r.Measurements.ExecTimeMs),
		fmt.Sprintf("%.6f", r.Measurements.ExecTimeS),
		fmt.Sprintf("%d", r.Measurements.ExitStatus),
		fmt.Sprintf("%.2f", r.Measurements.AvgCPUPercent),
		parsedResult,
	}
}

// String returns the string representation of the result
func (r *Result) String() string {
	raw, _ := json.Marshal(r)
	return string(raw)
}

// Filename returns the filename of the result
func (r *Result) Filename() string {
	return fmt.Sprintf("result_%s_%s.json", r.Task.ID().App(), r.Task.TrackerId())
}

// AddFileReference adds a reference file to the result so it can be processed later
func (r *Result) AddFileReference(app string, id string, filename string) {
	r.Files = append(r.Files, fileInfo{
		App:      app,
		SampleId: id,
		Filename: filename,
	})
}
