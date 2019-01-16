package main

import (
	"fmt"

	"github.com/kc1116/http-sniffer/core"
	"github.com/spf13/cobra"
)

var (
	logFile          string
	statsInterval    int
	monitorInterval  int
	monitorThreshold int
	web              bool
	layout           = "02/Jan/2006:15:04:04 -0700"
)

func main() {
	root.Execute()
}

var root = &cobra.Command{
	Use:   "run",
	Short: "HTTP-SNIFFER monitors acess logs",
	Long: `HTTP-SNIFFER monitors acess logs allowing 
	flexibility with cinfigurable options such as 
	request threshhold, statistics capturing interval, and alerting`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := core.NewFileReader(statsInterval)
		reqMonitor := core.NewTrafficMonitor(monitorThreshold, monitorInterval)

		if web {
			startWebClient()
		} else {
			startConsoleOutput(reader, reqMonitor)
		}
	},
}

func printOptions() {
	fmt.Printf("Reading From LogFile: %v\n", logFile)
	fmt.Printf("Monitor Interval: %v min\n", monitorInterval)
	fmt.Printf("Request Per Second Threshold: > %vreq/%vmin\n", monitorThreshold, monitorInterval)
	fmt.Printf("Stat Publishing Interval: %v/s\n\n", statsInterval)
}

func startWebClient() {}
func startConsoleOutput(reader *core.FileReader, monitor *core.TrafficMonitor) {
	output := core.NewConsoleOutput()
	output.Render()
	printOptions()
	go monitor.Start()
	go func() {
		for {
			select {
			case payload := <-reader.Sink:
				output.Publish(payload, nil)
			case alert := <-monitor.Sink:
				output.Publish(nil, alert)
			}
		}
	}()

	err := reader.OpenFile(logFile)
	fmt.Print(err)
}

func init() {
	root.PersistentFlags().StringVarP(&logFile, "log-file", "f", "/tmp/access.log", "path to access log file")
	root.PersistentFlags().IntVarP(&statsInterval, "stats-interval", "s", 10, "how often stats should be captures ()")
	root.PersistentFlags().IntVarP(&monitorInterval, "monitor-interval", "m", 2, "how often req/s threshold should be checked")
	root.PersistentFlags().IntVarP(&monitorThreshold, "monitor-threshold", "t", 10, "req/s")
	root.PersistentFlags().BoolVarP(&web, "web", "w", false, "view output in web UI")
}
