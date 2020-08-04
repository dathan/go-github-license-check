package license

//Service is a struct to execute methods
type Service struct {
	repository Repository
}

// LicenseCheckResult contains the information needed to generate the output we desire.
type LicenseCheckResult struct {
	GitHubRepo        string
	Lang              string
	Dependency        string
	DependencyLicense string
}

// List of licensechecks
type LicenseCheckResults []LicenseCheckResult

//interface for the core to save
type Repository interface {
	SaveLicenses(LicenseCheckResults) error
	GetLicenses(owner, repo string) (LicenseCheckResults, error)
}

// return the repository service
func NewService(ro Repository) *Service {

	service := Service{}
	service.repository = ro
	return &service

}

// generic algorithm for check to get the results and save them
func (service *Service) Execute(owner, repo string) error {

	lcr, err := service.repository.GetLicenses(owner, repo)
	if err != nil {
		return nil
	}

	if err := service.repository.SaveLicenses(lcr); err != nil {
		return err
	}

	return nil
}
