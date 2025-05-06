package processor

import (
	"biathlon/internal/config"
	"biathlon/internal/event"
	"biathlon/internal/models"
	"fmt"
	"sort"
	"strconv"
	"time"
)

type Processor struct {
	config      *config.Config
	competitors map[int]*models.Competitor
	logs        []string
}

func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{
		config:      cfg,
		competitors: make(map[int]*models.Competitor),
	}
}

func (p *Processor) Process(events []event.Event) ([]string, []string) {
	for _, e := range events {
		c := p.getOrCreateCompetitor(e.CompetitorID)
		switch e.EventID {
		case 1: // Registered
			c.Registered = true
			p.logf(e.Time, "The competitor(%d) registered", c.ID)

		case 2: // Start time set by draw
			startTime, _ := time.Parse("15:04:05.000", e.Extra)
			c.StartPlanned = startTime
			p.logf(e.Time, "The start time for the competitor(%d) was set by a draw to %s", c.ID, e.Extra)

		case 3: // On start line
			p.logf(e.Time, "The competitor(%d) is on the start line", c.ID)

		case 4: // Started
			c.StartActual = e.Time
			c.StartDelta = e.Time.Sub(c.StartPlanned)
			p.logf(e.Time, "The competitor(%d) has started", c.ID)

		case 5: // On firing range
			line, _ := strconv.Atoi(e.Extra)
			c.LastFiringHits = make(map[int]bool)
			p.logf(e.Time, "The competitor(%d) is on the firing range(%d)", c.ID, line)

		case 6: // Target hit
			targetID, _ := strconv.Atoi(e.Extra)
			if !c.LastFiringHits[targetID] {
				c.HitCount++
				c.LastFiringHits[targetID] = true
			}
			c.ShotCount++
			p.logf(e.Time, "The target(%d) has been hit by competitor(%d)", targetID, c.ID)

		case 7: // Left firing range
			p.logf(e.Time, "The competitor(%d) left the firing range", c.ID)

		case 8: // Entered penalty laps
			c.PenaltyStarts = append(c.PenaltyStarts, e.Time)
			p.logf(e.Time, "The competitor(%d) entered the penalty laps", c.ID)

		case 9: // Left penalty laps
			if len(c.PenaltyStarts) > 0 {
				start := c.PenaltyStarts[len(c.PenaltyStarts)-1]
				c.RecordPenalty(start, e.Time, float64(len(c.LastFiringHits))*p.config.PenaltyLen)
			}
			p.logf(e.Time, "The competitor(%d) left the penalty laps", c.ID)

		case 10: // Ended main lap
			if len(c.LapStarts) == 0 {
				c.LapStarts = append(c.LapStarts, c.StartActual)
			}
			start := c.LapStarts[len(c.LapStarts)-1]
			c.RecordLap(start, e.Time, p.config.LapLen)
			c.LapStarts = append(c.LapStarts, e.Time)
			p.logf(e.Time, "The competitor(%d) ended the main lap", c.ID)

			if len(c.LapTimes) == p.config.Laps {
				c.MarkFinished(e.Time)
				p.logf(e.Time.Add(1*time.Second), "The competitor(%d) has finished", c.ID)
			}

		case 11: // Can't continue
			c.MarkNotFinished(e.Extra)
			p.logf(e.Time, "The competitor(%d) can`t continue: %s", c.ID, e.Extra)
		}
	}

	// Check who never started
	for _, c := range p.competitors {
		if c.Registered && c.StartActual.IsZero() {
			c.MarkNotStarted()
			p.logf(c.StartPlanned.Add(p.config.Delta), "The competitor(%d) is disqualified", c.ID)
		}
	}

	// Generate final sorted report
	var report []string
	var sorted []*models.Competitor
	for _, c := range p.competitors {
		sorted = append(sorted, c)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ID < sorted[j].ID
	})
	for _, c := range sorted {
		report = append(report, c.ResultSummary())
	}

	return p.logs, report
}

func (p *Processor) getOrCreateCompetitor(id int) *models.Competitor {
	if _, ok := p.competitors[id]; !ok {
		p.competitors[id] = models.NewCompetitor(id)
	}
	return p.competitors[id]
}

func (p *Processor) logf(t time.Time, format string, args ...any) {
	entry := fmt.Sprintf("[%s] %s", t.Format("15:04:05.000"), fmt.Sprintf(format, args...))
	p.logs = append(p.logs, entry)
}
