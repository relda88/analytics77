package sgeo

import (
	"errors"
	"net/netip"

	"github.com/TudorHulban/analytics77/services/sstorage"
)

// ServiceGeo provides geo information for the request IP.
type ServiceGeo struct {
	serviceStorage *sstorage.ServiceStorage // provides access to the geo persisted info
}

func NewServiceGeo(service *sstorage.ServiceStorage) (*ServiceGeo, error) {
	if service == nil {
		return nil,
			errors.New("passed service storage is nil")
	}

	return &ServiceGeo{
			serviceStorage: service,
		},
		nil
}

func (s *ServiceGeo) GetIPGeo(ip netip.Addr) (*sstorage.ResponseGetIPGeo, error) {
	return s.serviceStorage.GetIPGeo(ip)
}
