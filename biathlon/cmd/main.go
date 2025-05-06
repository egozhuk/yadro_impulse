package main

import (
	"fmt"
	"log"

	"biathlon/internal/config"
	"biathlon/internal/event"
	"biathlon/internal/processor"
)

func main() {
	conf, err := config.LoadConfig("../sunny_5_skiers/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	events, err := event.ParseEvents("../sunny_5_skiers/events")
	if err != nil {
		log.Fatalf("Failed to parse events: %v", err)
	}

	proc := processor.NewProcessor(conf)
	logs, finalReport := proc.Process(events)

	fmt.Println("==== EVENT LOG ====")
	for _, line := range logs {
		fmt.Println(line)
	}

	fmt.Println("\n==== FINAL REPORT ====")
	for _, line := range finalReport {
		fmt.Println(line)
	}
}
