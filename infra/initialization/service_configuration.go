package initialization

import (
	"log"
	"os"

	"github.com/TudorHulban/analytics77/helpers"
	"github.com/TudorHulban/hxerrors/goerrors"
)

func Configuration(path string) map[string]any {
	result, errParse := helpers.ParseJSONFile(path)
	if errParse != nil {
		log.Printf(
			"configuration error: %s",
			errParse.Error(),
		)

		os.Exit(goerrors.OSExitForConfigurationIssues)
	}

	return result
}
