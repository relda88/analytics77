package sstorage

import (
	"errors"
	"net/netip"
)

type ServiceStorage struct{}

func NewServiceStorage() *ServiceStorage {
	return &ServiceStorage{}
}

type ResponseGetIPGeo struct {
	Country string
	City    string
	ASN     string
}

var ErrIPNotFound = errors.New("ip not found")

func (s *ServiceStorage) GetIPGeo(ip netip.Addr) (*ResponseGetIPGeo, error) {
	values := map[string]*ResponseGetIPGeo{
		"82.77.237.37": {
			Country: "ROU",
			City:    "Iasi",
			ASN:     "Digi",
		},

		"82.77.237.38": {
			Country: "ROU",
			City:    "Iasi",
			ASN:     "Digi",
		},

		"82.77.237.39": {
			Country: "ROU",
			City:    "Iasi",
			ASN:     "Vodafone",
		},
	}

	if value, exists := values[ip.String()]; exists {
		return value, nil
	}

	return nil,
		ErrIPNotFound
}
