package domain

import (
	"errors"
	"fmt"
	"net/netip"
	"sort"
	"strings"
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
	mu   sync.RWMutex
}

func NewDataCenter() *DataCenter {
	return &DataCenter{
		data: map[string]*Registry{},
	}
}

type ParamsAddEvent struct {
	SiteKey string
	Country string
	City    string

	DayIdx  int // 1-31
	HourIdx int // 0-23
	IP      netip.Addr
	Browser Browser
	ASN     AsnEntity
}

func (e *ParamsAddEvent) Validate() []error {
	var errs []error

	if len(e.SiteKey) == 0 {
		errs = append(errs, errors.New("site key cannot be empty"))
	}

	if len(e.Country) == 0 {
		errs = append(errs, errors.New("country cannot be empty"))
	}

	if len(e.City) == 0 {
		errs = append(errs, errors.New("city cannot be empty"))
	}

	if e.DayIdx < 1 || e.DayIdx > 31 {
		errs = append(
			errs,
			fmt.Errorf(
				"day index %d out of bounds (1-31)",
				e.DayIdx,
			),
		)
	}

	if e.HourIdx < 0 || e.HourIdx > 23 {
		errs = append(
			errs,
			fmt.Errorf(
				"hour index %d out of bounds (0-23)",
				e.HourIdx,
			),
		)
	}

	if !e.IP.IsValid() {
		errs = append(
			errs,
			errors.New("invalid or missing IP address"),
		)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func (dc *DataCenter) AddEvents(events ...*ParamsAddEvent) []error {
	errorsBatch := make([]error, 0)
	indexesNoError := make([]int, 0, len(events))

	var hasErrors bool

	for ix, event := range events {
		if errorsValidation := event.Validate(); errorsValidation != nil {
			hasErrors = true

			errorsBatch = append(errorsBatch, errorsValidation...)

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

type ResponseRecordsPerSite map[string]uint32

func (r ResponseRecordsPerSite) String() string {
	if len(r) == 0 {
		return "{}"
	}

	keys := make([]string, 0, len(r))

	for k := range r {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var builder strings.Builder
	builder.WriteString("{")

	for i, key := range keys {
		builder.WriteString(fmt.Sprintf("%s: %d", key, r[key]))

		if i < len(keys)-1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString("}")

	return builder.String()
}

func (dc *DataCenter) GetLastHourRecordsPerSite() ResponseRecordsPerSite {
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
