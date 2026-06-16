package initialization

import (
	"log"
	"os"

	"github.com/TudorHulban/analytics77/helpers"
	"github.com/TudorHulban/analytics77/services/sanalytics"
	"github.com/TudorHulban/analytics77/services/sgeo"
	"github.com/TudorHulban/analytics77/services/sstorage"
)

type ParamsServices struct {
	Offsets           helpers.TimestampOffsets
	APIKeyGeolocation string
}

func Services(params *ParamsServices) *sanalytics.ServiceAnalytics {
	log.Println("Initializing Analytics Application...")

	serviceStorage, errCrServiceStorage := sstorage.NewServiceStorage(".")
	if errCrServiceStorage != nil {
		log.Printf(
			"service geo creation: %s\n",
			errCrServiceStorage.Error(),
		)

		os.Exit(10)
	}

	serviceGeo, errCrServiceGeo := sgeo.NewServiceGeo(
		&sgeo.ParamsNewServiceGeo{
			APIKeyGeolocation: params.APIKeyGeolocation,
		},
		serviceStorage,
	)
	if errCrServiceGeo != nil {
		log.Printf(
			"service geo creation: %s",
			errCrServiceGeo.Error(),
		)

		os.Exit(11)
	}

	return sanalytics.NewServiceAnalytics(
		&sanalytics.PiersNewServiceAnalytics{
			ServiceGeo: serviceGeo,
		},
		&params.Offsets,
	)
}
