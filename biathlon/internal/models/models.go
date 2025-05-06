package models

import (
	"fmt"
	"strings"
	"time"
)

type LapData struct {
	Duration time.Duration
	Speed    float64
}

type PenaltyData struct {
	Duration time.Duration
	Speed    float64
}

type CompetitorStatus string

const (
	StatusNotStarted  CompetitorStatus = "NotStarted"
	StatusNotFinished CompetitorStatus = "NotFinished"
	StatusFinished    CompetitorStatus = "Finished"
)

type Competitor struct {
	ID             int
	Registered     bool
	StartPlanned   time.Time
	StartActual    time.Time
	StartDelta     time.Duration
	LapStarts      []time.Time
	LapTimes       []LapData
	PenaltyStarts  []time.Time
	PenaltyTimes   []PenaltyData
	HitCount       int
	ShotCount      int
	Status         CompetitorStatus
	TotalTime      time.Duration
	Comment        string
	LastFiringHits map[int]bool // to avoid double count
}

func NewCompetitor(id int) *Competitor {
	return &Competitor{
		ID:             id,
		LastFiringHits: make(map[int]bool),
	}
}

func (c *Competitor) RecordLap(start, end time.Time, lapLen float64) {
	duration := end.Sub(start)
	speed := lapLen / duration.Seconds()
	c.LapTimes = append(c.LapTimes, LapData{
		Duration: duration,
		Speed:    speed,
	})
}

func (c *Competitor) RecordPenalty(start, end time.Time, penaltyLen float64) {
	duration := end.Sub(start)
	speed := penaltyLen / duration.Seconds()
	c.PenaltyTimes = append(c.PenaltyTimes, PenaltyData{
		Duration: duration,
		Speed:    speed,
	})
}

func (c *Competitor) MarkNotStarted() {
	c.Status = StatusNotStarted
}

func (c *Competitor) MarkNotFinished(comment string) {
	c.Status = StatusNotFinished
	c.Comment = comment
}

func (c *Competitor) MarkFinished(end time.Time) {
	c.TotalTime = end.Sub(c.StartPlanned) + c.StartDelta
	c.Status = StatusFinished
}

func (c *Competitor) ResultSummary() string {
	status := ""
	switch c.Status {
	case StatusNotStarted:
		status = "[NotStarted]"
	case StatusNotFinished:
		status = "[NotFinished]"
	default:
		status = fmt.Sprintf("[%s]", c.TotalTime)
	}

	laps := ""
	for _, lap := range c.LapTimes {
		laps += fmt.Sprintf("{%v, %.3f} ", lap.Duration, lap.Speed)
	}

	var penaltyDur time.Duration
	var totalPenaltyLen float64
	for _, p := range c.PenaltyTimes {
		penaltyDur += p.Duration
		totalPenaltyLen += p.Speed * p.Duration.Seconds()
	}

	avgPenaltySpeed := 0.0
	if penaltyDur > 0 {
		avgPenaltySpeed = totalPenaltyLen / penaltyDur.Seconds()
	}

	return fmt.Sprintf("%s %d [%s] {%v, %.3f} %d/%d", status, c.ID, strings.TrimSpace(laps), penaltyDur, avgPenaltySpeed, c.HitCount, c.ShotCount)
}
