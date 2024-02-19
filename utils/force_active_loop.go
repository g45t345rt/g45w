package utils

import (
	"sync"
	"time"
)

// This is a simple timer that execute an update function.
// You have to enable after every update.
// Use this in UI render so when a layout is not rendered anymore it won't call the update func.

type ForceActiveLoop struct {
	ticker   *time.Ticker
	isActive bool
	done     chan bool
	sync.Mutex
}

func NewForceActiveLoop(d time.Duration, updateFunc func()) *ForceActiveLoop {
	activeLoop := &ForceActiveLoop{
		isActive: false,
		ticker:   time.NewTicker(d),
		done:     make(chan bool),
	}

	go func() {
		for {
			select {
			case <-activeLoop.done:
				return
			case <-activeLoop.ticker.C:
				activeLoop.Lock()
				if activeLoop.isActive {
					go updateFunc()
					activeLoop.isActive = false
				}
				activeLoop.Unlock()
			}
		}
	}()

	return activeLoop
}

func (a *ForceActiveLoop) SetActive() {
	defer a.Unlock()
	a.Lock()
	a.isActive = true
}

func (a *ForceActiveLoop) Close() {
	a.ticker.Stop()
	a.done <- true
}
