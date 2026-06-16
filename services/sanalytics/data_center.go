package sanalytics

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/TudorHulban/analytics77/domain"
)

type DataCenter struct {
	data map[string]*domain.Registry
	mu   sync.RWMutex
}

func NewDataCenter() *DataCenter {
	return &DataCenter{
		data: map[string]*domain.Registry{},
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
