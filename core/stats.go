package core

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/VividCortex/ewma"
)

var (
	// GlobalStats ...
	GlobalStats *LogStats
	once        sync.Once
)

// LogStats ...
type LogStats struct {
	TotalRequests   *big.Int
	RequestsAverage ewma.MovingAverage
	SectionStats
	ResponCodeStats
	ReqMethodStats
	BodySizeStats
}

// UpdateTotalRequests ...
func (l *LogStats) UpdateTotalRequests(n int) *LogStats {
	l.TotalRequests.Add(l.TotalRequests, big.NewInt(int64(n)))
	l.RequestsAverage.Add(float64(n))
	return l
}

// UpdateSectionStats ...
func (l *LogStats) UpdateSectionStats(section string, n int) *LogStats {
	if sectionStat, ok := l.SectionStats.Sections[section]; ok {
		sectionStat.Hits.Add(sectionStat.Hits, big.NewInt(int64(n)))
	} else {
		avg := ewma.NewMovingAverage()
		avg.Add(float64(n))
		l.SectionStats.Sections[section] = &SectionStat{
			Section: section,
			Hits:    big.NewInt(int64(n)),
		}
	}

	if l.SectionStats.Sections[section].Hits.Cmp(l.SectionStats.TopSection.Hits) > 0 {
		l.SectionStats.TopSection = l.SectionStats.Sections[section]
	}

	return l
}

// UpdateResponseCodeStats ...
func (l *LogStats) UpdateResponseCodeStats(code string, n int) *LogStats {
	if codeStat, ok := l.ResponCodeStats.Codes[code]; ok {
		codeStat.Hits.Add(codeStat.Hits, big.NewInt(int64(n)))
	} else {
		avg := ewma.NewMovingAverage()
		avg.Add(float64(n))
		l.ResponCodeStats.Codes[code] = &ResponCodeStat{
			Code: code,
			Hits: big.NewInt(int64(n)),
		}
	}
	return l
}

// UpdateReqMethodtats ...
func (l *LogStats) UpdateReqMethodtats(method string, n int) *LogStats {
	if methodStat, ok := l.ReqMethodStats.Methods[method]; ok {
		methodStat.Hits.Add(methodStat.Hits, big.NewInt(int64(n)))
	} else {
		avg := ewma.NewMovingAverage()
		avg.Add(float64(n))
		l.ReqMethodStats.Methods[method] = &ReqMethodStat{
			Method: method,
			Hits:   big.NewInt(int64(n)),
		}
	}
	return l
}

// UpdateBodySizeStat ...
func (l *LogStats) UpdateBodySizeStat(size int) *LogStats {
	l.BodySizeStats.TotalSize.Add(l.BodySizeStats.TotalSize, big.NewInt(int64(size)))
	l.BodySizeStats.Average.Add(float64(size))
	return l
}

// SectionStats ...
type SectionStats struct {
	TopSection *SectionStat `json:"top_section"`
	Sections   map[string]*SectionStat
}

// SectionStat ...
type SectionStat struct {
	Section string   `json:"section"`
	Hits    *big.Int `json:"hits"`
}

// ResponCodeStats ...
type ResponCodeStats struct {
	Codes map[string]*ResponCodeStat `json:"codes"`
}

// ResponCodeStat ...
type ResponCodeStat struct {
	Code string   `json:"code"`
	Hits *big.Int `json:"hits"`
}

// ReqMethodStats ...
type ReqMethodStats struct {
	Methods map[string]*ReqMethodStat `json:"request_methods"`
}

// ReqMethodStat ...
type ReqMethodStat struct {
	Method string   `json:"method"`
	Hits   *big.Int `json:"hits"`
}

// BodySizeStats ...
type BodySizeStats struct {
	TotalSize *big.Int `json:"total_size"`
	Average   ewma.MovingAverage
}

// NewStats ...
func NewStats() *LogStats {
	return &LogStats{
		TotalRequests:   big.NewInt(0),
		RequestsAverage: ewma.NewMovingAverage(),
		SectionStats: SectionStats{
			TopSection: &SectionStat{Section: "", Hits: big.NewInt(0)},
			Sections:   make(map[string]*SectionStat),
		},
		ResponCodeStats: ResponCodeStats{Codes: make(map[string]*ResponCodeStat)},
		ReqMethodStats:  ReqMethodStats{Methods: make(map[string]*ReqMethodStat)},
		BodySizeStats:   BodySizeStats{TotalSize: big.NewInt(0), Average: ewma.NewMovingAverage()},
	}
}

// Append ...
func (l *LogStats) Append(stats *LogStats) *LogStats {
	l.UpdateTotalRequests(int(stats.TotalRequests.Int64()))
	for _, val := range stats.SectionStats.Sections {
		l.UpdateSectionStats(val.Section, int(val.Hits.Int64()))
	}

	for _, val := range stats.ResponCodeStats.Codes {
		l.UpdateResponseCodeStats(val.Code, int(val.Hits.Int64()))
	}

	for _, val := range stats.ReqMethodStats.Methods {
		l.UpdateReqMethodtats(val.Method, int(val.Hits.Int64()))
	}

	l.UpdateBodySizeStat(int(stats.BodySizeStats.TotalSize.Int64()))

	return l
}

// CaptureStat ...
func CaptureStat(l Log, stats *LogStats) *LogStats {
	stats.
		UpdateTotalRequests(1).
		UpdateSectionStats(l.Section, 1).
		UpdateResponseCodeStats(fmt.Sprintf("%v", l.ResponseCode), 1).
		UpdateReqMethodtats(l.Method, 1).
		UpdateBodySizeStat(l.Size)

	return stats
}

func init() {
	once.Do(func() {
		GlobalStats = NewStats()
	})
}
