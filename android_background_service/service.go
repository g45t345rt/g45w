//go:build !android
// +build !android

package android_background_service

func Start() error {
	return nil
}

func Stop() error {
	return nil
}

func IsRunning() (bool, error) {
	return false, nil
}

func IsAvailable() bool {
	return false
}
