package sgeo

import (
	"errors"

	"github.com/TudorHulban/analytics77/domain"
	"github.com/TudorHulban/analytics77/services/sstorage"
	lru "github.com/TudorHulban/hx-lru"
)

// 1. Check LRU
// 2. If hit → return
// 3. If miss → call ServiceStorage.GetIPGeo(ip)
// 4. If storage hit → store in LRU → return
// 5. If storage miss → call geo provider
// 6. Persist provider result into Bitcask
// 7. Store in LRU
// 8. Return

type ServiceGeo struct {
	cache          *lru.CacheOneLRU[string, *domain.GeoIP]
	serviceStorage *sstorage.ServiceStorage
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

func (s *ServiceGeo) GetIPGeo(ip string) (*sstorage.ResponseGetIPGeo, error) {
	// 1. Hot cache
	if v, exists := s.cache.Get(ip); exists {
		return v, nil
	}

	// 2. Persistent store (Bitcask)
	if v, err := s.serviceStorage.GetIPGeo(ip); err == nil {
		s.cache.Put(ip, v)
		return v, nil
	}

	// 3. Provider (cold path)
	v, err := s.provider.Lookup(ip)
	if err != nil {
		return nil, err
	}

	// 4. Persist + cache
	_ = s.serviceStorage.PutIPGeo(ip, v)
	s.cache.Put(ip, v)

	return v, nil
}
