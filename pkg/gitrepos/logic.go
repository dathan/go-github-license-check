package gitrepos

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

// interface for the core to save
type Repository interface {
	GetRepos(string) (Repos, error)
}

// A github repo
type Repo struct {
	Org  string
	Name string
	Lang string
}

// type declarations
type Repos []Repo

// return the repository service
func NewService(ro Repository) *Service {

	service := Service{}
	service.repository = ro
	return &service

}

// generic algorithm for check to get the results and save them
func (service *Service) ListRepos(org string) (Repos, error) {

	results, err := service.repository.GetRepos(org)

	if err != nil {

		return nil, err

	}

	return results, nil
}
