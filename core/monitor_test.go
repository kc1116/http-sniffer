package core_test

import (
	"testing"
	"time"

	"github.com/kc1116/http-sniffer/core"
)

func TestAlert(t *testing.T) {
	reqMonitor := core.NewTrafficMonitor(1, 1)
	reqMonitor.Interval = time.Duration(1)

	go reqMonitor.Start()
	core.GlobalStats.UpdateTotalRequests(10)
	time.Sleep(1 * time.Second)

	var alert *core.Alert
	select {
	case alert = <-reqMonitor.Sink:
		break
	case <-time.After(reqMonitor.Interval * time.Second):
		break
	}

	if alert == nil {
		t.Errorf("expected alert to be published instead got nil %v", alert)
	}

	t.Log(alert.Message)
}
