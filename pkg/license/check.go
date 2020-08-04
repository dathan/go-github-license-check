package license

//DependencyGraphResponse is for the specific query performmed, anon struct types are used for the types.
type DependencyGraphManifestsResponse struct {
	Repository struct {
		DependencyGraphManifests struct {
			Edges []struct {
				Cursor string `json:"cursor"`
				Node   struct {
					Dependencies struct {
						Nodes []struct {
							PackageName string `json:"packageName"`
							Repository  struct {
								LicenseInfo struct {
									Name string `json:"name"`
								} `json:"licenseInfo"`
							} `json:"repository"`
						} `json:"nodes"`
						TotalCount int `json:"totalCount"`
					} `json:"dependencies"`
					DependenciesCount int `json:"dependenciesCount"`
				} `json:"node"`
			} `json:"edges"`
			PageInfo struct {
				EndCursor       string `json:"endCursor"`
				HasNextPage     bool   `json:"hasNextPage"`
				HasPreviousPage bool   `json:"hasPreviousPage"`
				StartCursor     string `json:"startCursor"`
			} `json:"pageInfo"`
			TotalCount int `json:"totalCount"`
		} `json:"dependencyGraphManifests"`
	} `json:"repository"`
}

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

func NewService(ro Repository) *Service {

	service := Service{}
	service.repository = ro
	return &service

}

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
