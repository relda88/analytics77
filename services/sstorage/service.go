package sstorage

import (
	"strconv"

	"github.com/TudorHulban/analytics77/domain"
	"github.com/prologic/bitcask"
	"github.com/shamaton/msgpack/v3"
)

type ServiceStorage struct {
	db *bitcask.Bitcask
}

func NewServiceStorage(path string) (*ServiceStorage, error) {
	db, errCrBitcaskDB := bitcask.Open(path)
	if errCrBitcaskDB != nil {
		return nil,
			errCrBitcaskDB
	}

	return &ServiceStorage{
			db: db,
		},
		nil
}

// func (s *ServiceStorage) GetIPGeo(ip netip.Addr) (*ResponseGetIPGeo, error) {
// 	values := map[string]*ResponseGetIPGeo{
// 		"127.0.0.1": {
// 			Country: "LOC",
// 			City:    "localhost",
// 			ASN:     "localhost",
// 		},

// 		"82.77.237.37": {
// 			Country: "ROU",
// 			City:    "Iasi",
// 			ASN:     "Digi",
// 		},

// 		"82.77.237.38": {
// 			Country: "ROU",
// 			City:    "Iasi",
// 			ASN:     "Digi",
// 		},

// 		"82.77.237.39": {
// 			Country: "ROU",
// 			City:    "Iasi",
// 			ASN:     "Vodafone",
// 		},
// 	}

// 	if value, exists := values[ip.String()]; exists {
// 		return value, nil
// 	}

// 	return nil,
// 		ErrIPNotFound
// }

func (s *ServiceStorage) GetIPGeo(ip string) (*domain.GeoIP, error) {
	key := []byte("ip:" + ip)

	dbValue, errGet := s.db.Get(key)
	if errGet != nil {
		if errGet == bitcask.ErrKeyNotFound {
			return nil,
				ErrIPNotFound
		}

		return nil, errGet
	}

	var result domain.GeoIP

	if errUnmarshal := msgpack.Unmarshal(dbValue, &result); errUnmarshal != nil {
		return nil,
			errUnmarshal
	}

	// Enrich with city + ASN entity
	city, _ := s.getCity(result.CityID)
	asn, _ := s.getASN(result.ASN)
	entity, _ := s.getEntity(asn.EntityID)

	return &domain.GeoIP{
		Country: result.CountryID,
		City:    city.Name,
		ASN:     entity.Name,
	}, nil
}

func (s *ServiceStorage) getCity(id string) (*CityRecord, error) {
	raw, err := s.db.Get([]byte("city:" + id))
	if err != nil {
		return nil, err
	}
	var rec CityRecord
	return &rec, msgpack.Unmarshal(raw, &rec)
}

func (s *ServiceStorage) getASN(asn string) (*ASNRecord, error) {
	raw, err := s.db.Get([]byte("asn:" + asn))
	if err != nil {
		return nil, err
	}
	var rec ASNRecord
	return &rec, msgpack.Unmarshal(raw, &rec)
}

func (s *ServiceStorage) getEntity(id int) (*EntityRecord, error) {
	raw, err := s.db.Get([]byte("entity:" + strconv.Itoa(id)))
	if err != nil {
		return nil, err
	}
	var rec EntityRecord
	return &rec, msgpack.Unmarshal(raw, &rec)
}
