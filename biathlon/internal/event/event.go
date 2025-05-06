package event

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Event struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	Extra        string
	RawLine      string
}

func ParseEvents(filename string) ([]Event, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		ts, err := time.Parse("15:04:05.000", strings.Trim(parts[0], "[]"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %v", err)
		}
		var eid, cid int
		fmt.Sscanf(parts[1], "%d", &eid)
		fmt.Sscanf(parts[2], "%d", &cid)

		extra := ""
		if len(parts) > 3 {
			extra = strings.Join(parts[3:], " ")
		}

		events = append(events, Event{
			Time:         ts,
			EventID:      eid,
			CompetitorID: cid,
			Extra:        extra,
			RawLine:      line,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
