package datatype

import "bytes"

type ContainerStats struct {
	CPUPerc  string
	MemUsage string
	MemPerc  string
	NetIO    string
	BlockIO  string
	PIDs     string
}

func parseContainerStats(statsLine string) *ContainerStats {
	// statsLine format: "CPU%|MemUsage|MemPerc|NetIO|BlockIO|PIDs"
	parts := bytes.Split([]byte(statsLine), []byte("|"))
	if len(parts) != 6 {
		return nil
	}
	return &ContainerStats{
		CPUPerc:  string(parts[0]),
		MemUsage: string(parts[1]),
		MemPerc:  string(parts[2]),
		NetIO:    string(parts[3]),
		BlockIO:  string(parts[4]),
		PIDs:     string(parts[5]),
	}
}
