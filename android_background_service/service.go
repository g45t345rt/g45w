//go:build !android
// +build !android

package android_background_service

func StartForeground() error {
	return nil
}

func StopForeground() error {
	return nil
}

func Start() error {
	return nil
}

func Stop() error {
	return nil
}

func IsForegroundRunning() (bool, error) {
	return false, nil
}

func IsRunning() (bool, error) {
	return false, nil
}

func IsAvailable() bool {
	return false
}
