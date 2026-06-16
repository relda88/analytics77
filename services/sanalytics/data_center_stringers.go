package sanalytics

import (
	"fmt"
	"strings"

	"github.com/TudorHulban/analytics77/domain"
)

func monthActiveString(label string, month []domain.DayActive, b *strings.Builder) {
	for ixDay := range month {
		day := &month[ixDay] // pointer, no copy

		for ixHour := range day {
			m := &day[ixHour] // pointer, no copy

			noRecords := m.RecordsPerPeriod.Load()
			if noRecords == 0 {
				continue
			}

			fmt.Fprintf(
				b,
				"  %-10s day%02d hour%02d  records:%d\n",

				label,
				ixDay,
				ixHour,
				noRecords,
			)
		}
	}
}

func monthArchivedString(label string, month []domain.DayArchived, b *strings.Builder) {
	for ixDay := range month {
		day := &month[ixDay] // pointer, no copy

		for ixHour := range day {
			m := &day[ixHour] // pointer, no copy

			noRecords := m.RecordsPerPeriod
			if noRecords == 0 {
				continue
			}

			fmt.Fprintf(
				b,
				"  %-10s day%02d hour%02d  records:%d\n",

				label,
				ixDay,
				ixHour,
				noRecords,
			)
		}
	}
}

func registryString(r *domain.Registry, b *strings.Builder) {
	monthActiveString("current", r.MonthCurrent[:], b)
	monthActiveString("previous", r.MonthPrevious[:], b)

	for ix, month := range r.History {
		monthArchivedString(
			fmt.Sprintf("history[%d]", ix),
			month[:],
			b,
		)
	}
}
