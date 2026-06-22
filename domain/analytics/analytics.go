package analytics

import (
	"sync/atomic"
)

type MetricActive struct {
	RecordsPerPeriod atomic.Uint32

	TopIPs       MetaActive[string]
	TopASN       MetaActive[string]
	TopCountries MetaActive[string]
	TopCities    MetaActive[string]
	TopURL       MetaActive[string]

	TopBrowsers MetaActive[Browser]
}

type MetricArchived struct {
	TopIPs       MetaArchived[string]
	TopASN       MetaArchived[string]
	TopCountries MetaArchived[string]
	TopCities    MetaArchived[string]
	TopURL       MetaArchived[string]
	TopBrowsers  MetaArchived[Browser]

	RecordsPerPeriod uint32
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
