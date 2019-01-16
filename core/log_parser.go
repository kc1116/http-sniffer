package core

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	layout         = "02/Jan/2006:15:04:04 -0700"
	accessLogRegex = `^(\S+) (\S+) (\S+) \[([\w:/]+\s[+\-]\d{4})\] "(\S+)\s?(\S+)?\s?(\S+)?" (\d{3}|-) (\d+|-)\s?"?([^"]*)"?\s?"?([^"]*)?"?$`
)

// remotehost rfc931 authuser [date] "request" status bytes

// Log ...
type Log struct {
	Orig           string `json:"orig"`
	RemoteHost     string `json:"remote_host"`
	RequestingUser string `json:"requesting_user"`
	Timestamp      string `json:"timestamp"`
	Method         string `json:"method"`
	Request        string `json:"request"`
	Section        string `json:"section"`
	HTTPVersion    string `json:"http_version"`
	ResponseCode   int    `json:"reponse_code"`
	Size           int    `json:"size"`
}

func toInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0
}

func getSection(route string) string {
	return fmt.Sprintf("/%s", strings.Split(route, "/")[1])
}

func getTime(timestamp string) string {
	if _, err := time.Parse(layout, timestamp); err != nil {
		return timestamp
	}
	return "0"
}

// ParseLog ...
func ParseLog(line string) (Log, error) {
	if match := regexp.MustCompile(accessLogRegex).FindStringSubmatch(line); len(match) == 12 {
		return Log{
			Orig:           line,
			RemoteHost:     match[1],
			RequestingUser: match[3],
			Timestamp:      match[4],
			Method:         match[5],
			Request:        match[6],
			Section:        getSection(match[6]),
			HTTPVersion:    match[7],
			ResponseCode:   toInt(match[8]),
			Size:           toInt(match[9]),
		}, nil
	}
	return Log{}, errors.New("failed to parse log line no regex match")
}
