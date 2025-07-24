package datatype

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

var (
	boldGreen = color.New(color.FgGreen, color.Bold)
	boldWhite = color.New(color.FgHiWhite, color.Bold)
	gray      = color.New(color.FgWhite, color.Faint)
	boldRed   = color.New(color.FgRed, color.Bold)
)

type Measurements struct {
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

// Result represents the result of the Docker command
type Result struct {
	Task             Task
	Output           string
	Error            error
	TotalElapsedTime time.Duration
	Measurements     Measurements
}

func (r *Result) Headers() []string {
	return []string{
		"app",
		"status",
		"error",
		"max_ram_kb",
		"exec_time_ms",
		"exec_time_s",
		"exit_status",
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
	return []string{
		boldWhite.Sprint(r.Task.ID().app),
		status,
		errstr,
		fmt.Sprintf("%d", r.Measurements.MaxRAMKb),
		fmt.Sprintf("%d", r.Measurements.ExecTimeMs),
		fmt.Sprintf("%.6f", r.Measurements.ExecTimeS),
		fmt.Sprintf("%d", r.Measurements.ExitStatus),
		fmt.Sprintf("%.2f", r.Measurements.AvgCPUPercent),
		"pendiente de parsear",
	}
}
