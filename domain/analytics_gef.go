package domain

import (
	"fmt"
	"net/netip"
	"slices"
	"strings"
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

func (dc *DataCenter) String() string {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	var b strings.Builder
	fmt.Fprintf(&b, "DataCenter: %d registr%s\n", len(dc.data),
		func() string {
			if len(dc.data) == 1 {
				return "y"
			}
			return "ies"
		}())

	keys := make([]string, 0, len(dc.data))
	for k := range dc.data {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, k := range keys {
		fmt.Fprintf(&b, "\n[%s]\n", k)
		registryString(dc.data[k], &b)
	}

	return b.String()
}

func registryString(r *Registry, b *strings.Builder) {
	monthString("current", r.MonthCurrent[:], b)
	monthString("previous", r.MonthPrevious[:], b)

	for i, month := range r.History {
		monthString(fmt.Sprintf("history[%d]", i), month[:], b)
	}
}

func monthString(label string, month []Day, b *strings.Builder) {
	for dayIdx, day := range month {
		for hourIdx, m := range day {
			n := m.RecordsPerPeriod.Load()
			if n == 0 {
				continue
			}
			fmt.Fprintf(b, "  %-10s day%02d hour%02d  records:%d\n",
				label, dayIdx, hourIdx, n)
		}
	}
}
