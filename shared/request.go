package shared

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strings"

	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/analytics77/services/sgeo"

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

func (req Request) AsParamsAddEvent(piers *PiersAsParamsAddEvent) (*ParamsAddEvent, error) {
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

	geoInfo := analytics.GeoIP{
		IsPrivate: true,
	}

	if !ip.IsPrivate() && !ip.IsLoopback() {
		responseGeo, errGeo := piers.ServiceGeo.GetIPGeo(ip.String())
		if errGeo != nil {
			return nil,
				fmt.Errorf(
					"geo call for IP: %s: %w",
					ip,
					errGeo,
				)
		}

		geoInfo = *responseGeo
	}

	userAgent := req.Header["User-Agent"]

	var uaString string

	if len(userAgent) > 0 {
		uaString = userAgent[0]
	}

	var browser analytics.Browser

	switch {
	case strings.Contains(uaString, "Safari") && !strings.Contains(uaString, "Chrome"):
		browser = analytics.Safari

	case strings.Contains(uaString, "Edg"):
		browser = analytics.Edge

	case strings.Contains(uaString, "Firefox"):
		browser = analytics.Firefox

	case strings.Contains(uaString, "Brave"):
		browser = analytics.Brave

	case strings.Contains(uaString, "Chrome"):
		browser = analytics.Chrome

	default:
		browser = 0
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

	return &ParamsAddEvent{
			SiteKey: hxhelpers.Ternary(
				len(req.Host) == 0,

				host,
				req.Host,
			),
			Country: geoInfo.Location.CountryCode,
			City:    geoInfo.Location.City,

			DayOfMonth: DayMonth(ixDay),
			HourOfDay:  HourDay(ixHour),
			IP:         ip,
			Browser:    browser,

			ASNOrganization: geoInfo.ASN.Organization,

			OffsetUTC:     offsetUTC,
			TimestampUNIX: req.TimestampUNIX,

			IsPrivateIP: geoInfo.IsPrivate,
		},
		nil
}

type Requests []Request
