package domain

import (
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
