package translation

import (
	"io/ioutil"
	"log"

	translationDomain "github.com/eujoy/openrecipe/tools/mdgen/internal/domain/translation"
	"gopkg.in/yaml.v2"
)

// Service describes the translations service.
type Service struct {
	translationsFile string
}

// New creates and returns a translation service.
func New(translationsFile string) *Service {
	return &Service{translationsFile: translationsFile}
}

// Prepare retrieves the translations from the provided file, prepares the respective struct and returns it.
func (s *Service) Prepare() translationDomain.Details {
	xLationsFile, xLationsReadFileErr := ioutil.ReadFile(s.translationsFile)
	if xLationsReadFileErr != nil {
		log.Printf("Failed to read translations file with error: #%v ", xLationsReadFileErr)
	}

	var xLations translationDomain.Details
	xLationsUnmarshalErr := yaml.Unmarshal(xLationsFile, &xLations)
	if xLationsUnmarshalErr != nil {
		log.Fatalf("Failed to unmarshal translations file with error: %v", xLationsUnmarshalErr)
	}

	return xLations
}
