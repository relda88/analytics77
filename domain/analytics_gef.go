package domain

import (
	"net/netip"
	"sync"
	"sync/atomic"
)

type Metric struct {
	RecordsPerPeriod atomic.Uint32
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
