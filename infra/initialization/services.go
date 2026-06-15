package initialization

import (
	"log"
	"os"

	"github.com/TudorHulban/analytics77/domain"
	"github.com/TudorHulban/analytics77/helpers"
	"github.com/TudorHulban/analytics77/services/sanalytics"
	"github.com/TudorHulban/analytics77/services/sgeo"
	"github.com/TudorHulban/analytics77/services/sstorage"
)

type ParamsServices struct {
	Offsets helpers.TimestampOffsets
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

	serviceGeo, errCrServiceGeo := sgeo.NewServiceGeo(serviceStorage)
	if errCrServiceGeo != nil {
		log.Printf(
			"service geo creation: %s",
			errCrServiceGeo.Error(),
		)

		os.Exit(11)
	}

	return sanalytics.NewServiceAnalytics(
		&sanalytics.PiersNewServiceAnalytics{
			DC:         domain.NewDataCenter(),
			ServiceGeo: serviceGeo,
		},
		&params.Offsets,
	)
}
