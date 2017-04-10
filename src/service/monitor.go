package service

const (
	cpuUsage      = "cpu"
	memPercentage = "memory"
	memUsage      = "memory_usage"
	memTotal      = "memory_total"
	networkRX     = "network_rx"
	networkTX     = "network_tx"
	fsRead        = "fs_read"
	fsWrite       = "fs_write"
	fsUsage       = "fs_usage"
	fsLimit       = "fs_limit"
)

func isInArray(array []string, value string) bool {
	for _, valueInList := range array {
		if value == valueInList {
			return true
		}
	}
	return false
}
