package core

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VividCortex/ewma"
)

// Alert ...
type Alert struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// TrafficMonitor ...
type TrafficMonitor struct {
	Average      ewma.MovingAverage
	ReqThreshold float64
	Interval     time.Duration
	sigint       chan os.Signal
	Sink         chan *Alert
}

// Start ...
func (t *TrafficMonitor) Start() {
	go func() {
		for {
			select {
			case <-time.After(t.Interval * time.Second):
				t.Average.Add(float64(GlobalStats.TotalRequests.Int64()))
				if t.Average.Value() > t.ReqThreshold {
					t.Sink <- &Alert{
						Message: fmt.Sprintf("High traffic generated an alert - hits = %.2f, triggered at %s",
							t.Average.Value(), time.Now().Format(layout)),
					}
				}
			}
		}
	}()

	<-t.sigint
	close(t.Sink)
	return
}

// NewTrafficMonitor ...
func NewTrafficMonitor(reqThreshold, interval int) *TrafficMonitor {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return &TrafficMonitor{
		Average:      ewma.NewMovingAverage(),
		ReqThreshold: float64(interval),
		Interval:     time.Duration(interval * 60),
		sigint:       sigint,
		Sink:         make(chan *Alert),
	}
}
