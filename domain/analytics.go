package domain

import (
	"errors"
	"net/netip"
	"sync"
	"time"
)

type Metric struct {
	RecordsPerPeriod uint32
	TopIPs           Meta[netip.Addr]
	TopBrowsers      Meta[Browser]
	TopASN           Meta[AsnEntity]
	TopCountries     Meta[string]
	TopCities        Meta[string]
	TopURL           Meta[string]
}

type Day [24]Metric

type Registry struct {
	MonthPrevious [31]Day
	MonthCurrent  [31]Day

	History [7][31]Day
}

type DataCenter struct {
	data map[string]*Registry
	mu   sync.Mutex
}

func NewDataCenter() *DataCenter {
	return &DataCenter{
		data: map[string]*Registry{},
	}
}

type ParamsAddEvent struct {
	SiteKey string
	DayIdx  int // 1-31
	HourIdx int // 0-23
	IP      netip.Addr
	Browser Browser
	ASN     AsnEntity
	Country string
	City    string
}

func (dc *DataCenter) AddEvents(events ...*ParamsAddEvent) []error {
	errorsBatch := make([]error, len(events))
	indexesNoError := make([]int, 0, len(events))

	var hasErrors bool

	for ix, event := range events {
		if event.DayIdx < 1 || event.DayIdx > 31 {
			hasErrors = true

			errorsBatch[ix] = errors.New("day index out of bounds (1-31)")

			continue
		}

		if event.HourIdx < 0 || event.HourIdx > 23 {
			hasErrors = true

			errorsBatch[ix] = errors.New("hour index out of bounds (0-23)")

			continue
		}

		indexesNoError = append(indexesNoError, ix)
	}

	dc.mu.Lock()

	for _, eventIx := range indexesNoError {
		// Create a local reference to keep the pointer code readable
		event := events[eventIx]

		// O(1) Fetch or Allocate the Registry for this site/datacenter
		registry, exists := dc.data[event.SiteKey]
		if !exists {
			registry = &Registry{}

			dc.data[event.SiteKey] = registry
		}

		// Drill down directly to the exact memory address in the matrix
		// We do dayIdx-1 to map human days (1-31) to zero-based array indices (0-30)
		metricSlot := &registry.MonthCurrent[event.DayIdx-1][event.HourIdx]

		// Update the flat metrics using pure pointer execution
		metricSlot.RecordsPerPeriod++

		metricSlot.TopIPs.Increment(event.IP)
		metricSlot.TopBrowsers.Increment(event.Browser)
		metricSlot.TopASN.Increment(event.ASN)

		if event.Country != "" {
			metricSlot.TopCountries.Increment(event.Country)
		}

		if event.City != "" {
			metricSlot.TopCities.Increment(event.City)
		}
	}

	dc.mu.Unlock()

	if hasErrors {
		return errorsBatch
	}

	return nil
}

func (dc *DataCenter) GetLastHourRecordsPerSite() map[string]uint32 {
	now := time.Now()

	// Map human time to zero-based array indices
	currentDayIdx := now.Day() - 1
	currentHourIdx := now.Hour()

	dc.mu.Lock()

	// Pre-allocate the map based on current registry size to avoid reallocations
	result := make(map[string]uint32, len(dc.data))

	// Extract the granular metric directly per site
	for siteKey, registry := range dc.data {
		result[siteKey] = registry.MonthCurrent[currentDayIdx][currentHourIdx].RecordsPerPeriod
	}

	dc.mu.Unlock()

	return result
}
