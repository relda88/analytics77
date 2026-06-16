package domain

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
)

type MetricActive struct {
	RecordsPerPeriod atomic.Uint32

	TopIPs       MetaActive[string]
	TopBrowsers  MetaActive[Browser]
	TopASN       MetaActive[string]
	TopCountries MetaActive[string]
	TopCities    MetaActive[string]
	TopURL       MetaActive[string]
}

type MetricArchived struct {
	RecordsPerPeriod uint32

	TopIPs       MetaArchived[string]
	TopBrowsers  MetaArchived[Browser]
	TopASN       MetaArchived[string]
	TopCountries MetaArchived[string]
	TopCities    MetaArchived[string]
	TopURL       MetaArchived[string]
}

type (
	DayActive   [24]MetricActive
	DayArchived [24]MetricArchived
)

type Registry struct {
	MonthPrevious [31]DayActive
	MonthCurrent  [31]DayActive

	History [7][31]DayArchived
}

type DataCenter struct {
	data map[string]*Registry
	mu   sync.RWMutex
}

func NewDataCenter() *DataCenter {
	return &DataCenter{
		data: map[string]*Registry{},
	}
}

func (dc *DataCenter) String() string {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	var b strings.Builder
	fmt.Fprintf(&b, "DataCenter: %d registr%s\n", len(dc.data),
		func() string {
			if len(dc.data) == 1 {
				return "y"
			}
			return "ies"
		}())

	keys := make([]string, 0, len(dc.data))
	for k := range dc.data {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		fmt.Fprintf(&b, "\n[%s]\n", k)
		registryString(dc.data[k], &b)
	}

	return b.String()
}

func monthActiveString(label string, month []DayActive, b *strings.Builder) {
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

func monthArchivedString(label string, month []DayArchived, b *strings.Builder) {
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

func registryString(r *Registry, b *strings.Builder) {
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
