package shared

import (
	"encoding/gob"
	"net"
	"net/netip"
	"net/url"
	"strings"

	"github.com/TudorHulban/analytics77/domain"
	"github.com/TudorHulban/analytics77/helpers"
)

func init() {
	gob.Register(&url.URL{})
}

type Request struct {
	RemoteAddr    string
	Host          string
	Method        string
	URL           *url.URL
	Header        map[string][]string
	TimestampUNIX int64
}

func (req Request) AsParamsAddEvent(offsets *helpers.TimestampOffsets) (*domain.ParamsAddEvent, error) {
	host, _, errHost := net.SplitHostPort(req.RemoteAddr)
	if errHost != nil {
		host = req.RemoteAddr
	}

	ip, errParseIP := netip.ParseAddr(host)
	if errParseIP != nil {
		return nil,
			errParseIP
	}

	userAgent := req.Header["User-Agent"]

	var uaString string

	if len(userAgent) > 0 {
		uaString = userAgent[0]
	}

	var browser domain.Browser

	switch {
	case strings.Contains(uaString, "Chrome"):
		browser = domain.Chrome

	case strings.Contains(uaString, "Safari") && !strings.Contains(uaString, "Chrome"):
		browser = domain.Safari

	case strings.Contains(uaString, "Edg"):
		browser = domain.Edge

	case strings.Contains(uaString, "Firefox"):
		browser = domain.Firefox

	case strings.Contains(uaString, "Brave"):
		browser = domain.Brave

	default:
		browser = 0
	}

	var country, city string

	if reqCountry := req.Header["Cf-Ipcountry"]; len(reqCountry) > 0 {
		country = reqCountry[0]
	}
	if reqCity := req.Header["X-Client-City"]; len(reqCity) > 0 {
		city = reqCity[0]
	}

	ixDay, ixHour := helpers.ExtractDayAndHour(
		req.TimestampUNIX,
		offsets,
	)

	return &domain.ParamsAddEvent{
			DayIdx:  ixDay,
			HourIdx: ixHour,
			IP:      ip,
			Browser: browser,
			Country: country,
			City:    city,
		},
		nil
}

type Requests []Request
