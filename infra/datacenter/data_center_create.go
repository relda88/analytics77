package datacenter

import (
	"github.com/TudorHulban/analytics77/domain"
	"github.com/TudorHulban/analytics77/shared"
)

func (dc *DataCenter) AddEvents(events ...*shared.ParamsAddEvent) []error {
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
			registrySite = &domain.Registry{}

			dc.data[event.SiteKey] = registrySite
		}

		metricSlot := &registrySite.
			MonthCurrent[event.DayOfMonth-1][event.HourOfDay]

		metricSlot.RecordsPerPeriod.Add(1)

		metricSlot.TopIPs.Increment(event.IP.String())
		metricSlot.TopBrowsers.Increment(event.Browser)
		metricSlot.TopASN.Increment(event.ASNOrganization)

		metricSlot.TopCountries.Increment(event.Country)
		metricSlot.TopCities.Increment(event.City)
	}

	dc.mu.Unlock()

	if hasErrors {
		return errorsBatch
	}

	return nil
}
