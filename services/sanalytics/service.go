package sanalytics

import (
	"fmt"

	"github.com/TudorHulban/analytics77/domain"
	"github.com/TudorHulban/analytics77/helpers"
	"github.com/TudorHulban/analytics77/shared"
)

// TODO: add methods to update the DST moments
type ServiceAnalytics struct {
	DC      *domain.DataCenter
	offsets *helpers.TimestampOffsets
}

func NewServiceAnalytics(dataCenter *domain.DataCenter, offsets *helpers.TimestampOffsets) *ServiceAnalytics {
	return &ServiceAnalytics{
		DC:      dataCenter,
		offsets: offsets,
	}
}

// RecordEvents returns transformation and validation / processing errors.
func (s *ServiceAnalytics) RecordEvents(events shared.Requests) ([]error, []error) {
	errorsTransformation := make([]error, 0, len(events))
	validEvents := make([]*domain.ParamsAddEvent, 0, len(events))

	for ix, event := range events {
		param, errTransformation := event.AsParamsAddEvent(s.offsets)
		if errTransformation != nil {
			errorsTransformation = append(
				errorsTransformation,
				fmt.Errorf(
					"transformation error for event %d:%w",
					ix,
					errTransformation,
				),
			)

			continue
		}

		validEvents = append(validEvents, param)
	}

	if len(validEvents) == 0 {
		return errorsTransformation,
			nil
	}

	errorsProcess := s.DC.AddEvents(validEvents...)

	return errorsTransformation, errorsProcess
}
