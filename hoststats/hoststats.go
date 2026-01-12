// Package hoststats provides functions to collect host statistics.
package hoststats

import (
	"context"
	"fmt"
	"math"

	"github.com/alertbingo/alertbingo/api"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

// Config holds the common configuration for host stats checks
type Config struct {
	Dashboard        string
	Site             string
	Service          string
	Message          string
	InactiveExpire   string
	InactiveEscalate string
	Highlighted      string
}

// Collect gathers memory, uptime, CPU, and disk statistics and returns check payloads
func Collect(ctx context.Context, cfg Config) ([]api.CheckPayload, error) {
	var checks []api.CheckPayload

	// Memory check
	memCheck, err := collectMemory(ctx, cfg)
	if err != nil {
		return nil, err
	}
	checks = append(checks, memCheck)

	// Uptime check
	uptimeCheck, err := collectUptime(ctx, cfg)
	if err != nil {
		return nil, err
	}
	checks = append(checks, uptimeCheck)

	// CPU check
	cpuCheck, err := collectCPU(ctx, cfg)
	if err != nil {
		return nil, err
	}
	checks = append(checks, cpuCheck)

	// Disk Used check
	diskUsedCheck, diskInodesCheck, err := collectDisk(ctx, cfg)
	if err != nil {
		return nil, err
	}
	checks = append(checks, diskUsedCheck, diskInodesCheck)

	return checks, nil
}

func collectMemory(ctx context.Context, cfg Config) (api.CheckPayload, error) {
	vmem, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return api.CheckPayload{}, fmt.Errorf("failed to get virtual memory: %w", err)
	}
	memPercent := int(math.Ceil(vmem.UsedPercent))
	memAlertLevel := 0 // ok
	memMessage := cfg.Message
	if memPercent >= 95 {
		memAlertLevel = 1 // warn
		memMessage = appendAlertReason(cfg.Message, "Memory % over 95")
	}
	return api.CheckPayload{
		Dashboard:        cfg.Dashboard,
		Site:             cfg.Site,
		Service:          cfg.Service,
		Name:             "Memory",
		AlertLevel:       memAlertLevel,
		Value:            fmt.Sprintf("%d%%", memPercent),
		Message:          memMessage,
		InactiveExpire:   cfg.InactiveExpire,
		InactiveEscalate: cfg.InactiveEscalate,
		Highlighted:      cfg.Highlighted,
	}, nil
}

func collectUptime(ctx context.Context, cfg Config) (api.CheckPayload, error) {
	hostInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		return api.CheckPayload{}, fmt.Errorf("failed to get host info: %w", err)
	}
	uptimeDays := int(hostInfo.Uptime / 86400) // seconds to days
	uptimeAlertLevel := 0                      // ok
	uptimeMessage := cfg.Message
	if uptimeDays < 1 {
		uptimeAlertLevel = 1 // warn
		uptimeMessage = appendAlertReason(cfg.Message, "Uptime less than 1 day")
	}
	return api.CheckPayload{
		Dashboard:        cfg.Dashboard,
		Site:             cfg.Site,
		Service:          cfg.Service,
		Name:             "Uptime",
		AlertLevel:       uptimeAlertLevel,
		Value:            fmt.Sprintf("%dd", uptimeDays),
		Message:          uptimeMessage,
		InactiveExpire:   cfg.InactiveExpire,
		InactiveEscalate: cfg.InactiveEscalate,
		Highlighted:      cfg.Highlighted,
	}, nil
}

func collectCPU(ctx context.Context, cfg Config) (api.CheckPayload, error) {
	loadAvg, err := load.AvgWithContext(ctx)
	if err != nil {
		return api.CheckPayload{}, fmt.Errorf("failed to get load average: %w", err)
	}
	cpuCount, err := cpu.CountsWithContext(ctx, true) // logical cores
	if err != nil {
		return api.CheckPayload{}, fmt.Errorf("failed to get CPU count: %w", err)
	}
	cpuPercent := int(math.Ceil((loadAvg.Load1 / float64(cpuCount)) * 100))
	if cpuPercent > 100 {
		cpuPercent = 100 // cap at 100%
	}
	cpuAlertLevel := 0 // ok
	cpuMessage := cfg.Message
	if cpuPercent == 100 {
		cpuAlertLevel = 1 // warn
		cpuMessage = appendAlertReason(cfg.Message, "CPU % at 100")
	}
	return api.CheckPayload{
		Dashboard:        cfg.Dashboard,
		Site:             cfg.Site,
		Service:          cfg.Service,
		Name:             "CPU",
		AlertLevel:       cpuAlertLevel,
		Value:            fmt.Sprintf("%d%%", cpuPercent),
		Message:          cpuMessage,
		InactiveExpire:   cfg.InactiveExpire,
		InactiveEscalate: cfg.InactiveEscalate,
		Highlighted:      cfg.Highlighted,
	}, nil
}

func collectDisk(ctx context.Context, cfg Config) (api.CheckPayload, api.CheckPayload, error) {
	diskUsage, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		return api.CheckPayload{}, api.CheckPayload{}, fmt.Errorf("failed to get disk usage: %w", err)
	}

	// Disk Used check
	diskUsedPercent := int(math.Ceil(diskUsage.UsedPercent))
	diskUsedAlertLevel := 0 // ok
	diskUsedMessage := cfg.Message
	if diskUsedPercent > 99 {
		diskUsedAlertLevel = 2 // alert
		diskUsedMessage = appendAlertReason(cfg.Message, "Disk Used % over 99")
	} else if diskUsedPercent > 95 {
		diskUsedAlertLevel = 1 // warn
		diskUsedMessage = appendAlertReason(cfg.Message, "Disk Used % over 95")
	}
	diskUsedCheck := api.CheckPayload{
		Dashboard:        cfg.Dashboard,
		Site:             cfg.Site,
		Service:          cfg.Service,
		Name:             "Disk Used",
		AlertLevel:       diskUsedAlertLevel,
		Value:            fmt.Sprintf("%d%%", diskUsedPercent),
		Message:          diskUsedMessage,
		InactiveExpire:   cfg.InactiveExpire,
		InactiveEscalate: cfg.InactiveEscalate,
		Highlighted:      cfg.Highlighted,
	}

	// Disk Inodes check
	diskInodesPercent := int(math.Ceil(diskUsage.InodesUsedPercent))
	diskInodesAlertLevel := 0 // ok
	diskInodesMessage := cfg.Message
	if diskInodesPercent > 99 {
		diskInodesAlertLevel = 2 // alert
		diskInodesMessage = appendAlertReason(cfg.Message, "Disk Inodes % over 99")
	} else if diskInodesPercent > 95 {
		diskInodesAlertLevel = 1 // warn
		diskInodesMessage = appendAlertReason(cfg.Message, "Disk Inodes % over 95")
	}
	diskInodesCheck := api.CheckPayload{
		Dashboard:        cfg.Dashboard,
		Site:             cfg.Site,
		Service:          cfg.Service,
		Name:             "Disk Inodes",
		AlertLevel:       diskInodesAlertLevel,
		Value:            fmt.Sprintf("%d%%", diskInodesPercent),
		Message:          diskInodesMessage,
		InactiveExpire:   cfg.InactiveExpire,
		InactiveEscalate: cfg.InactiveEscalate,
		Highlighted:      cfg.Highlighted,
	}

	return diskUsedCheck, diskInodesCheck, nil
}

// appendAlertReason appends an alert reason to an existing message
func appendAlertReason(message, reason string) string {
	if message == "" {
		return reason
	}
	return message + " - " + reason
}
