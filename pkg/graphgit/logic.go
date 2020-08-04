package graphgit

import (
	"context"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

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

// use graphql to get data from github for the repository
type Service struct {
	gClient *graphql.Client
}

func NewService() *Service {
	service := &Service{}

	service.gClient = graphql.NewClient("https://api.github.com/graphql")
	service.gClient.Log = func(s string) {
		log.Printf("check.NewService()  - %s\n", s)
	}
	return service
}

// return the licenses
func (service *Service) GetLicenses(owner, repo string) (*DependencyGraphManifestsResponse, error) {

	req := graphql.NewRequest(`
	query {
		repository(name: "` + repo + `", owner: "` + owner + `") {
		  dependencyGraphManifests {
			edges {
			  cursor
			  node {
				dependencies {
				  nodes {
					packageName
					repository {
					  licenseInfo {
						name
					  }
					}
				  }
				  totalCount
				}
				dependenciesCount
			  }
			}
			pageInfo {
			  endCursor
			  hasNextPage
			  hasPreviousPage
			  startCursor
			}
			totalCount
		  } 
		}
	  }
	`)
	token := os.Getenv("GITHUB_GRAPHQL_CHECK")
	req.Header.Add("Authorization", "bearer "+token)
	req.Header.Add("Accept", "application/vnd.github.hawkgirl-preview+json")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData DependencyGraphManifestsResponse
	if err := service.gClient.Run(ctx, req, &respData); err != nil {
		//log.Fatal(errors.Wrap(err, "NewService graphql client failed - "))
		return nil, err
	}

	return &respData, nil
}
