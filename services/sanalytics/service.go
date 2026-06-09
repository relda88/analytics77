package sanalytics

import (
	"github.com/TudorHulban/analytics77/domain"
	"github.com/TudorHulban/analytics77/helpers"
	"github.com/TudorHulban/analytics77/shared"
)

// TODO: add methods to update the DST moments
type ServiceAnalytics struct {
	DC      *domain.DataCenter
	offsets *helpers.TimestampOffsets
}

func NewServiceAnalytics(dataCenter *domain.DataCenter) *ServiceAnalytics {
	return &ServiceAnalytics{
		DC: dataCenter,
	}
}

func (s *ServiceAnalytics) RecordEvents(events shared.Requests) ([]error, []error) {
	errorsValidation := make([]error, 0, len(events))
	validEvents := make([]*domain.ParamsAddEvent, 0, len(events))

	for _, request := range events {
		param, errTransformation := request.AsParamsAddEvent(s.offsets)
		if errTransformation != nil {
			errorsValidation = append(errorsValidation, errTransformation)

			continue
		}

		validEvents = append(validEvents, param)
	}

	if len(validEvents) == 0 {
		return errorsValidation,
			nil
	}

	errorsProcess := s.DC.AddEvents(validEvents...)

	return errorsValidation, errorsProcess
}
