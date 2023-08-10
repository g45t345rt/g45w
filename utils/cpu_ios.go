//go:build ios
// +build ios

package utils

import (
	"errors"
	"time"
)

func CPU_Counts(logical bool) (int, error) {
	return 0, errors.New("ios: func not available")
}

func CPU_Percent(interval time.Duration, percpu bool) ([]float64, error) {
	return make([]float64, 0), errors.New("ios: func not available")
}
