package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type NodeOverview struct {
	Hostname      string      `json:"hostname"`
	OS            string      `json:"os"`
	Kernel        string      `json:"kernel"`
	Architecture  string      `json:"architecture"`
	Uptime        string      `json:"uptime"`
	LoadAverage   string      `json:"loadAverage"`
	CpuPercent    float64     `json:"cpuPercent"`
	MemoryTotalMB int         `json:"memoryTotalMB"`
	MemoryUsedMB  int         `json:"memoryUsedMB"`
	MemoryPercent float64     `json:"memoryPercent"`
	DiskTotal     string      `json:"diskTotal"`
	DiskUsed      string      `json:"diskUsed"`
	DiskPercent   float64     `json:"diskPercent"`
	LastReportAt  string      `json:"lastReportAt"`
	Mounts        []NodeMount `json:"mounts"`
}

type NodeMount struct {
	Name       string  `json:"name"`
	FileSystem string  `json:"fileSystem"`
	Size       string  `json:"size"`
	Used       string  `json:"used"`
	Usage      string  `json:"usage"`
	Percent    float64 `json:"percent"`
}

const nodeOverviewCommand = `
HOSTNAME=$(hostname 2>/dev/null)
OS=$(awk -F= '/^PRETTY_NAME=/{gsub(/"/,"",$2); print $2}' /etc/os-release 2>/dev/null)
KERNEL=$(uname -r 2>/dev/null)
ARCH=$(uname -m 2>/dev/null)
UPTIME=$(uptime -p 2>/dev/null | sed 's/^up //')
LOAD=$(awk '{printf "%s / %s / %s", $1, $2, $3}' /proc/loadavg 2>/dev/null)
CPU_PERCENT=$(vmstat 1 2 2>/dev/null | tail -1 | awk '{printf "%.2f", 100 - $15}')
MEMORY=$(free -m 2>/dev/null | awk '/^Mem:/ {printf "%s|%s|%.2f", $2, $3, ($3 * 100 / $2)}')
ROOT_DISK=$(df -hP / 2>/dev/null | awk 'NR==2 {gsub("%", "", $5); printf "%s|%s|%s", $2, $3, $5}')

echo "HOSTNAME=${HOSTNAME}"
echo "OS=${OS}"
echo "KERNEL=${KERNEL}"
echo "ARCH=${ARCH}"
echo "UPTIME=${UPTIME}"
echo "LOAD=${LOAD}"
echo "CPU_PERCENT=${CPU_PERCENT}"
echo "MEMORY=${MEMORY}"
echo "ROOT_DISK=${ROOT_DISK}"

df -hPT 2>/dev/null | awk 'NR>1 && $2 != "tmpfs" && $2 != "devtmpfs" && $2 != "overlay" && $2 != "squashfs" {
	gsub("%", "", $6)
	printf "MOUNT=%s|%s|%s|%s|%s%%|%s\n", $7, $2, $3, $4, $6, $6
}'
`

func (s *SSHService) GetNodeOverview(nodeID uint) (*NodeOverview, error) {
	output, err := s.RunCommand(nodeID, nodeOverviewCommand)
	if err != nil {
		return nil, err
	}
	return parseNodeOverview(output)
}

func parseNodeOverview(output string) (*NodeOverview, error) {
	overview := &NodeOverview{
		LastReportAt: time.Now().Format("2006-01-02 15:04:05"),
		Mounts:       make([]NodeMount, 0),
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "HOSTNAME":
			overview.Hostname = value
		case "OS":
			overview.OS = value
		case "KERNEL":
			overview.Kernel = value
		case "ARCH":
			overview.Architecture = value
		case "UPTIME":
			overview.Uptime = value
		case "LOAD":
			overview.LoadAverage = value
		case "CPU_PERCENT":
			overview.CpuPercent = parseFloatValue(value)
		case "MEMORY":
			memoryParts := strings.Split(value, "|")
			if len(memoryParts) != 3 {
				continue
			}
			overview.MemoryTotalMB = parseIntValue(memoryParts[0])
			overview.MemoryUsedMB = parseIntValue(memoryParts[1])
			overview.MemoryPercent = parseFloatValue(memoryParts[2])
		case "ROOT_DISK":
			diskParts := strings.Split(value, "|")
			if len(diskParts) != 3 {
				continue
			}
			overview.DiskTotal = strings.TrimSpace(diskParts[0])
			overview.DiskUsed = strings.TrimSpace(diskParts[1])
			overview.DiskPercent = parseFloatValue(diskParts[2])
		case "MOUNT":
			mountParts := strings.Split(value, "|")
			if len(mountParts) != 6 {
				continue
			}

			overview.Mounts = append(overview.Mounts, NodeMount{
				Name:       strings.TrimSpace(mountParts[0]),
				FileSystem: strings.TrimSpace(mountParts[1]),
				Size:       strings.TrimSpace(mountParts[2]),
				Used:       strings.TrimSpace(mountParts[3]),
				Usage:      strings.TrimSpace(mountParts[4]),
				Percent:    parseFloatValue(mountParts[5]),
			})
		}
	}

	if overview.OS == "" && overview.Kernel == "" && overview.Hostname == "" {
		return nil, fmt.Errorf("解析节点信息失败")
	}

	return overview, nil
}

func parseIntValue(value string) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0
	}
	return parsed
}

func parseFloatValue(value string) float64 {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0
	}
	return parsed
}
