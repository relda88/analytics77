package domain

import (
	"errors"
	"fmt"
	"net/netip"
)

type ParamsAddEvent struct {
	SiteKey string
	Country string
	City    string

	DayOfMonth int // 1-31
	HourOfDay  int // 0-23
	IP         netip.Addr
	Browser    Browser
	ASN        AsnEntity

	TimestampUNIX int64
	OffsetUTC     int64
}

func (e *ParamsAddEvent) Validate() []error {
	var errs []error

	if len(e.SiteKey) == 0 {
		errs = append(
			errs,
			errors.New("site key cannot be empty"),
		)
	}

	if len(e.Country) == 0 {
		errs = append(
			errs,
			errors.New("country cannot be empty"),
		)
	}

	if len(e.City) == 0 {
		errs = append(
			errs,
			errors.New("city cannot be empty"),
		)
	}

	if e.DayOfMonth < 1 || e.DayOfMonth > 31 {
		errs = append(
			errs,
			fmt.Errorf(
				"day index %d out of bounds (1-31)",
				e.DayOfMonth,
			),
		)
	}

	if e.HourOfDay < 0 || e.HourOfDay > 23 {
		errs = append(
			errs,
			fmt.Errorf(
				"hour index %d out of bounds (0-23)",
				e.HourOfDay,
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

			errorsBatch = append(
				errorsBatch,
				errorsValidation...)

			continue
		}

		indexesNoError = append(indexesNoError, ix)
	}

	dc.mu.Lock()

	for _, eventIndex := range indexesNoError {
		event := events[eventIndex]

		registrySite, exists := dc.data[event.SiteKey]
		if !exists {
			registrySite = &Registry{}

			dc.data[event.SiteKey] = registrySite
		}

		metricSlot := &registrySite.
			MonthCurrent[event.DayOfMonth-1][event.HourOfDay]

		metricSlot.RecordsPerPeriod.Add(1)

		metricSlot.TopIPs.Increment(event.IP)
		metricSlot.TopBrowsers.Increment(event.Browser)
		metricSlot.TopASN.Increment(event.ASN)

		metricSlot.TopCountries.Increment(event.Country)
		metricSlot.TopCities.Increment(event.City)
	}

	dc.mu.Unlock()

	if hasErrors {
		return errorsBatch
	}

	return nil
}
