package models

type CPUInfo struct {
	Usage       float64 `json:"usage"`
	Cores       int     `json:"cores"`
	LoadAverage float64 `json:"load_average"`
}

type MemoryInfo struct {
	Total string  `json:"total"`
	Used  string  `json:"used"`
	Free  string  `json:"free"`
	Usage float64 `json:"usage"`
}

type DiskInfo struct {
	Path   string  `json:"path"`
	Total  string  `json:"total"`
	Used   string  `json:"used"`
	Free   string  `json:"free"`
	Usage  float64 `json:"usage"`
	FSType string  `json:"fs_type"`
}

type NetInfo struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}

type ProcInfo struct {
	Count int `json:"count"`
}

type HealthData struct {
	Timestamp string      `json:"timestamp"`
	Hostname  string      `json:"hostname"`
	Platform  string      `json:"platform"`
	Uptime    uint64      `json:"uptime"`
	CPU       *CPUInfo    `json:"cpu,omitempty"`
	Memory    *MemoryInfo `json:"memory,omitempty"`
	Disk      []DiskInfo  `json:"disk,omitempty"`
	Network   []NetInfo   `json:"network,omitempty"`
	Processes *ProcInfo   `json:"processes,omitempty"`
	Alerts    []string    `json:"alerts,omitempty"`
	Healthy   bool        `json:"healthy"`
}
