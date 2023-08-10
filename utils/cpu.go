//go:build !ios
// +build !ios

package utils

import (
	"time"

	cpu_utils "github.com/shirou/gopsutil/v3/cpu"
)

func CPU_Counts(logical bool) (int, error) {
	return cpu_utils.Counts(logical)
}

func CPU_Percent(interval time.Duration, percpu bool) ([]float64, error) {
	return cpu_utils.Percent(interval, percpu)
}
