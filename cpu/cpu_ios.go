//go:build ios
// +build ios

package cpu

import (
	"errors"
	"time"
)

func Counts(logical bool) (int, error) {
	return 0, errors.New("ios: func not available")
}

func Percent(interval time.Duration, percpu bool) ([]float64, error) {
	return make([]float64, 0), errors.New("ios: func not available")
}
