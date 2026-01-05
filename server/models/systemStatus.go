package models

type SystemStatus struct {
	Cpu    string   `json:"cpu"`
	Memory string   `json:"memory"`
	Swap   string   `json:"swap"`
	Disk   []string `json:"disk"`
}

func (s *SystemStatus) NewSystemStatus(cpu, memory, swap string, disk []string) *SystemStatus {
	return &SystemStatus{
		Cpu:    cpu,
		Memory: memory,
		Swap:   swap,
		Disk:   disk,
	}
}
