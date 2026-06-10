package shared

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strings"

	"github.com/TudorHulban/analytics77/domain"
	"github.com/TudorHulban/analytics77/helpers"
	"github.com/TudorHulban/analytics77/services/sgeo"

	"github.com/tudorhulban/hxhelpers"
)

func init() {
	gob.Register(&url.URL{})
}

type Request struct {
	RemoteAddr string
	Host       string
	Method     string
	URL        *url.URL
	Header     map[string][]string

	TimestampUNIX int64
	OffsetUTC     int64
}

type PiersAsParamsAddEvent struct {
	Offsets    *helpers.TimestampOffsets
	ServiceGeo *sgeo.ServiceGeo
}

func (req Request) AsParamsAddEvent(piers *PiersAsParamsAddEvent) (*domain.ParamsAddEvent, error) {
	if piers.Offsets == nil {
		return nil, errors.New(
			"AsParamsAddEvent - passed offsets is nil",
		)
	}
	if piers.ServiceGeo == nil {
		return nil, errors.New(
			"AsParamsAddEvent - passed ServiceGeo is nil",
		)
	}

	host, _, errHost := net.SplitHostPort(req.RemoteAddr)
	if errHost != nil {
		host = req.RemoteAddr
	}

	ip, errParseIP := netip.ParseAddr(host)
	if errParseIP != nil {
		return nil,
			fmt.Errorf(
				"parsing IP: %s: %w",
				ip,
				errParseIP,
			)
	}

	userAgent := req.Header["User-Agent"]

	var uaString string

	if len(userAgent) > 0 {
		uaString = userAgent[0]
	}

	var browser domain.Browser

	switch {
	case strings.Contains(uaString, "Safari") && !strings.Contains(uaString, "Chrome"):
		browser = domain.Safari

	case strings.Contains(uaString, "Edg"):
		browser = domain.Edge

	case strings.Contains(uaString, "Firefox"):
		browser = domain.Firefox

	case strings.Contains(uaString, "Brave"):
		browser = domain.Brave

	case strings.Contains(uaString, "Chrome"):
		browser = domain.Chrome

	default:
		browser = 0
	}

	responseGeo, errGeo := piers.ServiceGeo.GetIPGeo(ip)
	if errGeo != nil {
		return nil,
			fmt.Errorf(
				"geo call for IP: %s: %w",
				ip,
				errGeo,
			)
	}

	offsetUTC := hxhelpers.Ternary(
		req.OffsetUTC > 0,

		req.OffsetUTC,
		piers.Offsets.OffsetUTC,
	)

	ixDay, ixHour := helpers.ExtractDayAndHour(
		req.TimestampUNIX,
		&helpers.TimestampOffsets{
			TimestampDSTWinter: piers.Offsets.TimestampDSTWinter,
			TimestampDSTSpring: piers.Offsets.TimestampDSTSpring,

			OffsetUTC: offsetUTC,
		},
	)

	return &domain.ParamsAddEvent{
			SiteKey: host,
			Country: responseGeo.Country,
			City:    responseGeo.City,

			DayOfMonth: domain.DayMonth(ixDay),
			HourOfDay:  domain.HourDay(ixHour),
			IP:         ip,
			Browser:    browser,

			ASN: domain.AsnEntity{
				Name: responseGeo.ASN,
			},

			OffsetUTC:     offsetUTC,
			TimestampUNIX: req.TimestampUNIX,
		},
		nil
}

type Requests []Request
