package datacenter

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/TudorHulban/analytics77/helpers"
)

type ResponseRecordsPerSite map[string]uint32

func (r ResponseRecordsPerSite) String() string {
	if len(r) == 0 {
		return "{}"
	}

	keys := make([]string, 0, len(r))

	for k := range r {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var builder strings.Builder
	builder.WriteString("{")

	for ix, key := range keys {
		builder.WriteString(
			fmt.Sprintf(
				"%s: %d",
				key,
				r[key],
			),
		)

		if ix < len(keys)-1 {
			builder.WriteString(", ")
		}
	}

	builder.WriteString("}")

	return builder.String()
}

func (dc *DataCenter) GetLastHourRecordsPerSite(offsets *helpers.TimestampOffsets) ResponseRecordsPerSite {
	ixDay, ixHour := helpers.ExtractDayAndHour(
		time.Now().Unix(),
		offsets,
	)

	dc.mu.Lock()

	result := make(map[string]uint32, len(dc.data))

	for siteKey, registry := range dc.data {
		result[siteKey] = registry.
			MonthCurrent[ixDay][ixHour].RecordsPerPeriod.Load()
	}

	dc.mu.Unlock()

	return result
}
