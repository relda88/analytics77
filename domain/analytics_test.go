package domain

import (
	"net/netip"
	"testing"

	"github.com/TudorHulban/analytics77/helpers"
	"github.com/stretchr/testify/require"
)

func TestE2E(t *testing.T) {
	dc := NewDataCenter()
	require.NotNil(t, dc)

	siteKey := "xxx.eu"

	p1 := ParamsAddEvent{
		SiteKey: siteKey,
		Country: "ROU",
		City:    "Iasi",

		DayOfMonth: 1,
		HourOfDay:  1,
		IP:         netip.IPv4Unspecified(),
		Browser:    Brave,
		ASN: AsnEntity{
			ID:   1,
			Name: "Digi",
		},
	}
	require.Empty(t, p1.Validate())

	errsAddEvent := dc.AddEvents(&p1)
	require.Empty(t, errsAddEvent)

	p2 := ParamsAddEvent{
		SiteKey: siteKey,
		Country: "ROU",
		City:    "Iasi",

		DayOfMonth: 1,
		HourOfDay:  2,
		IP:         netip.IPv4Unspecified(),
		Browser:    Brave,
		ASN: AsnEntity{
			ID:   1,
			Name: "Digi",
		},
	}
	require.Empty(t, p2.Validate())

	p3 := ParamsAddEvent{
		SiteKey: siteKey,
		Country: "ROU",
		City:    "Iasi",

		DayOfMonth: 1,
		HourOfDay:  3,
		IP:         netip.IPv4Unspecified(),
		Browser:    Brave,
		ASN: AsnEntity{
			ID:   1,
			Name: "Digi",
		},
	}

	errsAddEvents := dc.AddEvents(&p2, &p3)
	require.Empty(t, errsAddEvents)

	offsetsROU := helpers.TimestampOffsets{
		OffsetUTC: 3,
	}

	records := dc.GetLastHourRecordsPerSite(&offsetsROU)
	require.Len(t, records, 1)
	require.EqualValues(t,
		3,
		records[siteKey],
		records.String(),
	)
}
