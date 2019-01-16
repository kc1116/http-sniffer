package core

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gosuri/uilive"

	"github.com/mbndr/figlet4go"
	"github.com/olekukonko/tablewriter"
)

var (
	banner    string
	top       *tablewriter.Table
	reqStats  *tablewriter.Table
	respStats *tablewriter.Table

	timestampFMT  = "\nStats Last Captured @Time: %s  (unix: %v) - %s  (unix: %v)"
	topSectionFMT = "\n\nTop Section: %s - Hits: %s\n"
	totalReqFMT   = "Total Requests: %s - Current Average: %.2f\n"
	totalBytesFMT = "Total Bytes Processed: %s - Current Average: %.2f \n"
	respCodesFMT  = "\nResponse Code Stats\n"
	reqMethodsFMT = "\nRequest Method Stats\n"
)

// Console ...
type Console struct {
	Ready       bool
	sigint      chan os.Signal
	writer      *uilive.Writer
	alertWriter *uilive.Writer
	cache       string
	lastAlert   string
}

// Render ...
func (c *Console) Render() {
	fmt.Printf("%s\n", banner)
	c.writer = uilive.New()
	c.writer.RefreshInterval = time.Duration(1)
	c.writer.Start()

	c.alertWriter = uilive.New()
	c.alertWriter.RefreshInterval = time.Duration(1)
	c.alertWriter.Start()

	// fmt.Fprint(c.alertWriter, "No Alerts")
	// fmt.Println()
}

// Alert ...
func (c *Console) Alert(alert *Alert) {
	fmt.Fprint(c.alertWriter, alert.Message)
}

// Publish ...
func (c *Console) Publish(payload *ReadPayload, alert *Alert) {
	c.writer.Flush()
	if payload != nil && payload.Stats != nil {
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf(timestampFMT, payload.StartTime.Format(layout), payload.StartTime.Unix(), payload.EndTime.Format(layout), payload.EndTime.Unix()))
		totalReq := payload.Stats.TotalRequests.String()
		totalBytes := payload.Stats.BodySizeStats.TotalSize.String()
		topSection := payload.Stats.SectionStats.TopSection.Section
		topSectionHits := payload.Stats.SectionStats.TopSection.Hits.String()
		codes := payload.Stats.ResponCodeStats.Codes
		methods := payload.Stats.ReqMethodStats.Methods

		builder.WriteString(fmt.Sprintf(topSectionFMT, topSection, topSectionHits))
		builder.WriteString(fmt.Sprintf(totalReqFMT, totalReq, GlobalStats.RequestsAverage.Value()))
		builder.WriteString(fmt.Sprintf(totalBytesFMT, totalBytes, GlobalStats.BodySizeStats.Average.Value()))

		builder.WriteString(respCodesFMT)
		for _, val := range codes {
			builder.WriteString(fmt.Sprintf("  Code: %s Hits: %s\n", val.Code, val.Hits.String()))
		}
		builder.WriteString(reqMethodsFMT)
		for _, val := range methods {
			builder.WriteString(fmt.Sprintf("  Method: %s Hits: %s\n", val.Method, val.Hits.String()))
		}

		c.cache = builder.String()
	}

	s := fmt.Sprintf("Total Requests Overall %v - Current Request Rate (running average): %.2f \t Total Bytes Overall %v - Average Request Bytes (running average): %.2f\n",
		GlobalStats.TotalRequests.Int64(),
		GlobalStats.RequestsAverage.Value(),
		GlobalStats.BodySizeStats.TotalSize.Int64(),
		GlobalStats.BodySizeStats.Average.Value())

	if alert != nil {
		c.lastAlert = alert.Message
	}

	fmt.Fprintf(c.writer, "%s%s%s", c.cache, s, c.lastAlert)
}

// NewConsoleOutput ...
func NewConsoleOutput() *Console {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return &Console{
		sigint:    sigint,
		lastAlert: "No alerts detected :)",
	}
}

func init() {
	ascii := figlet4go.NewAsciiRender()

	// Adding the colors to RenderOptions
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		// Colors can be given by default ansi color codes...
		figlet4go.ColorGreen,
		figlet4go.ColorYellow,
		figlet4go.ColorCyan,
	}

	banner, _ = ascii.RenderOpts("HTTP-SNIFFER", options)
}
