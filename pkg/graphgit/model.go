package graphgit

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
								PrimaryLanguage struct {
									Name string `json:"name"`
								} `json:"primaryLanguage"`
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

type GithubRepositoriesResponse struct {
	Repos struct {
		RepositoryCount int `json:"repositoryCount"`
		PageInfo        struct {
			EndCursor   string `json:"endCursor"`
			StartCursor string `json:"startCursor"`
			HasNextPage bool   `json:"hasNextPage"`
		} `json:"pageInfo"`
		Edges []struct {
			Node struct {
				Name            string `json:"name"`
				NameWithOwner   string `json:"nameWithOwner"`
				PrimaryLanguage struct {
					Name string `json:"name"`
				} `json:"primaryLanguage"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"repos"`
}
