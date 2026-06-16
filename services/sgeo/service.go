package sgeo

import (
	"errors"
	"net/http"
	"time"

	"github.com/TudorHulban/analytics77/domain"
	requestgeo "github.com/TudorHulban/analytics77/infra/request-geo"
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
	apiKeyGeolocation string

	cache          *lru.CacheOneLRU[string, domain.GeoIP]
	serviceStorage *sstorage.ServiceStorage
	httpClient     *http.Client
}

type ParamsNewServiceGeo struct {
	APIKeyGeolocation string
}

func NewServiceGeo(params *ParamsNewServiceGeo, service *sstorage.ServiceStorage) (*ServiceGeo, error) {
	if service == nil {
		return nil,
			errors.New("passed service storage is nil")
	}

	return &ServiceGeo{
			apiKeyGeolocation: params.APIKeyGeolocation,

			cache: lru.NewCacheOneLRU[string, domain.GeoIP](
				&lru.ParamsNewCacheLRU{
					TTL:      14 * 24 * time.Hour,
					Capacity: 5000,
				},
			),
			serviceStorage: service,
			httpClient: &http.Client{
				Timeout: 5 * time.Second,
			},
		},
		nil
}

func (s *ServiceGeo) GetIPGeo(ip string) (*domain.GeoIP, error) {
	// 1. Hot cache
	if cacheValue, errGetLRU := s.cache.Get(ip); errGetLRU == nil {
		return cacheValue,
			nil
	}

	// 2. Persistent store (Bitcask)
	if kvValue, errGetPersistence := s.serviceStorage.GetIPGeo(ip); errGetPersistence == nil {
		s.cache.Put(ip, *kvValue)

		return kvValue,
			nil
	}

	// 3. Provider (cold path)
	providerValue, errGetProvider := requestgeo.GetLocationByIP(
		&requestgeo.ParamsGetLocationByIP{
			Client:    s.httpClient,
			APIKey:    s.apiKeyGeolocation,
			IPAddress: ip,
		},
	)
	if errGetProvider != nil {
		return nil,
			errGetProvider
	}

	// 4. Persist + cache
	_ = s.serviceStorage.PutGeoIP(providerValue)
	s.cache.Put(ip, *providerValue)

	return providerValue,
		nil
}
