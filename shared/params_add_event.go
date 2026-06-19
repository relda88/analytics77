package shared

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/tudorhulban/analytics77/domain/analytics"
)

type DayMonth uint8 // 1 - 31

func (d DayMonth) IsValid() bool {
	return d >= 1 && d <= 31
}

type HourDay uint8 // 0 - 23

func (h HourDay) IsValid() bool {
	return h <= 23
}

type ParamsAddEvent struct {
	SiteKey string
	Country string
	City    string

	DayOfMonth      DayMonth
	HourOfDay       HourDay
	IP              netip.Addr
	Browser         analytics.Browser
	ASNOrganization string

	TimestampUNIX int64
	OffsetUTC     int64
	IsPrivateIP   bool
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

	if !e.DayOfMonth.IsValid() {
		errs = append(
			errs,
			fmt.Errorf(
				"day index %d out of bounds (1-31)",
				e.DayOfMonth,
			),
		)
	}

	if !e.HourOfDay.IsValid() {
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
