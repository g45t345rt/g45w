//go:build !ios
// +build !ios

package cpu

import (
	"time"

	cpu_utils "github.com/shirou/gopsutil/v3/cpu"
)

func Counts(logical bool) (int, error) {
	return cpu_utils.Counts(logical)
}

func Percent(interval time.Duration, percpu bool) ([]float64, error) {
	return cpu_utils.Percent(interval, percpu)
}
