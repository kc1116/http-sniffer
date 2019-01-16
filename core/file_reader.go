package core

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hpcloud/tail"
)

// ReadPayload ...
type ReadPayload struct {
	StartTime time.Time
	EndTime   time.Time
	Stats     *LogStats
	Logs      []Log
}

// FileReader ...
type FileReader struct {
	sigint       chan os.Signal
	readInterval time.Duration
	file         *tail.Tail
	Sink         chan *ReadPayload
}

// OpenFile ...
func (f *FileReader) OpenFile(path string) error {
	var err error
	config := tail.Config{
		Follow:    true,
		MustExist: true,
		Poll:      true,
		Logger:    nil,
		// Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
	}
	if f.file, err = tail.TailFile(path, config); err == nil {
		f.startReading()
	}
	return err
}

// CloseFile ...
func (f *FileReader) CloseFile() error {
	f.file.Cleanup()
	return f.file.Stop()
}

func (f *FileReader) startReading() {
	logs, intervalStats, startTime := f.initReadInfo()
	go func() {
		for {
			select {
			case line := <-f.file.Lines:
				if log, err := ParseLog(line.Text); err == nil {
					intervalStats = CaptureStat(log, intervalStats)
					logs = append(logs, log)
				}
			case <-time.After(f.readInterval * time.Second):
				go GlobalStats.Append(intervalStats)
				if len(logs) > 0 {
					f.Sink <- &ReadPayload{startTime, time.Now(), intervalStats, logs}
				} else {
					GlobalStats.UpdateTotalRequests(0).UpdateBodySizeStat(0)
					f.Sink <- &ReadPayload{startTime, time.Now(), nil, []Log{}}
				}

				logs, intervalStats, startTime = f.initReadInfo()
			}
		}
	}()

	<-f.sigint
	close(f.Sink)
	return
}

func (f *FileReader) initReadInfo() ([]Log, *LogStats, time.Time) {
	return make([]Log, 0), NewStats(), time.Now()
}

func (f *FileReader) watchSignal() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-sigint
	f.CloseFile()
}

// NewFileReader ...
func NewFileReader(readInterval int) *FileReader {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return &FileReader{
		sigint:       sigint,
		readInterval: time.Duration(readInterval),
		Sink:         make(chan *ReadPayload),
	}
}
