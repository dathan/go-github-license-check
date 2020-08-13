package csv

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/dathan/go-github-license-check/pkg/license"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Service struct {
}

func NewService() *Service {

	return &Service{}

}

func (s *Service) FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func (s *Service) SkipIfExecuted(input license.LicenseCheckResults) bool {

	filename := getFileNameFromLicense(input)
	// better semiphore is needed
	if s.FileExists(filename) {
		logrus.Warningf("respository.SaveLicenses(): Filename: %s exists...skipping", filename)
		return true
	}
	return false

}

func getFileNameFromLicense(input license.LicenseCheckResults) string {
	splitResult := strings.Split(input[0].GitHubRepo, "/")
	//TODO make this configurable
	filename := "./data/" + splitResult[len(splitResult)-1] + ".csv"
	return filename

}

func (s *Service) Save(res license.LicenseCheckResults) error {

	if s.SkipIfExecuted(res) {
		return nil
	}

	filename := getFileNameFromLicense(res)
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "creating a file - improve this")
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range res {
		err := writer.Write([]string{value.GitHubRepo, value.Dependency, value.DependencyLicense, value.Lang})
		if err != nil {
			return err
		}
	}

	return nil
}
