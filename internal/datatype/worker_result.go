package datatype

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/zerjioang/rooftop/v2/common/io/json"
)

var (
	boldGreen  = color.New(color.FgGreen, color.Bold)
	boldPurple = color.New(color.FgMagenta, color.Bold)
	boldWhite  = color.New(color.FgHiWhite, color.Bold)
	gray       = color.New(color.FgWhite, color.Faint)
	boldRed    = color.New(color.FgRed, color.Bold)
)

type Measurements struct {
	FailedToParse   bool    `json:"failed_to_parse"`
	MaxRAMKb        int64   `json:"max_ram_kb"`
	ExecTimeUs      int64   `json:"exec_time_us"`
	ExecTimeMs      int64   `json:"exec_time_ms"`
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
	Timeout          bool          `json:"timeout"`
	TotalElapsedTime time.Duration `json:"total_elapsed_time"`
	Measurements     Measurements  `json:"measurements"`
	// Parsed output of the executed command
	ParsedOutput *ScanResult `json:"parsed_output"`
	// Additional reference files
	Files []fileInfo      `json:"files"`
	Stats *ContainerStats `json:"stats"`
}

func (r *Result) Headers() []string {
	return []string{
		"sota tool",
		"sample id",
		"status",
		"max ram",
		"time",
		"avg cpu",
		"nodes",
		"edges",
		"connected",
	}
}

func (r *Result) Rows() []string {
	status := boldGreen.Sprint("✔")
	if r.Error != nil {
		status = boldRed.Sprint("❌")
	}
	if r.Measurements.ExitStatus != 0 {
		status = boldRed.Sprint("❌")
	}
	results := &ScanResult{}
	if r.ParsedOutput != nil {
		results = r.ParsedOutput
	}
	appName := boldWhite.Sprint(r.Task.ID().app)
	if r.Task.ID().app == "paper" {
		appName = boldPurple.Sprint(r.Task.ID().app)
	}
	return []string{
		boldWhite.Sprint(appName),
		boldWhite.Sprint(r.Task.TrackerId()),
		status,
		beautifyRAM(r.Measurements.MaxRAMKb),
		fmt.Sprintf("%s (%d ms)", beautifyTimeWithUnits(r.Measurements.ExecTimeMs), r.Measurements.ExecTimeMs),
		fmt.Sprintf("%.2f %%", r.Measurements.AvgCPUPercent),
		fmt.Sprintf("%d", results.NodesDetected),
		fmt.Sprintf("%d", results.EdgesDetected),
		fmt.Sprintf("%s %%", toFloat(results.Coverage)),
	}
}

func toBool(v *bool) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%v", *v)
}

func toFloat(v *float64) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%.2f", *v)
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

// beautifyRAM converts an integer representing bytes into a human-readable string.
func beautifyRAM(kb int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	// note: profiler returns the data as kb
	bytes := kb * 1024
	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// BeautifyTimeWithUnits formats milliseconds into a human-readable string with units.
func beautifyTimeWithUnits(ms int64) string {
	totalSeconds := ms / 1000
	seconds := totalSeconds % 60
	minutes := (totalSeconds / 60) % 60
	hours := totalSeconds / 3600

	result := ""
	if hours > 0 {
		result += fmt.Sprintf("%dh ", hours)
	}
	if minutes > 0 || hours > 0 { // show minutes if there are hours
		result += fmt.Sprintf("%dm ", minutes)
	}
	result += fmt.Sprintf("%ds", seconds)

	return result
}
