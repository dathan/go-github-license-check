package graphgit

import (
	"context"

	"log"
	"os"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/pkg/errors"
	logrus "github.com/sirupsen/logrus"
)

// use graphql to get data from github for the repository
type Service struct {
	gClient *graphql.Client
}

const DEBUG_ON = false

func NewService() *Service {
	service := &Service{}

	service.gClient = graphql.NewClient("https://api.github.com/graphql")
	service.gClient.Log = func(s string) {
		if DEBUG_ON {
			log.Printf("check.NewService()  - %s\n", s)
		}
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	logrus.SetReportCaller(true)

	return service
}

// return the licenses
func (service *Service) GetLicenses(owner, repo, after string) (*DependencyGraphManifestsResponse, error) {

	afterStr := ""
	if len(after) > 0 {
		afterStr = ", after:\"" + after + "\""
	}
	req := graphql.NewRequest(`
	query {
		repository(name: "` + repo + `", owner: "` + owner + `"` + afterStr + `) {
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
								  primaryLanguage {
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

	req.Header.Add("Accept", "application/vnd.github.hawkgirl-preview+json")

	// run it and capture the response
	var respData DependencyGraphManifestsResponse
	if err := service.execute(req, &respData); err != nil {
		if strings.Contains(err.Error(), "decoding") {
			skipErrorNextTime(repo)
			return nil, nil
		}
		return nil, err
	}

	if respData.Repository.DependencyGraphManifests.PageInfo.HasNextPage {
		data, err := service.GetLicenses(owner, repo, respData.Repository.DependencyGraphManifests.PageInfo.EndCursor)
		if err != nil {
			return nil, err
		}
		respData.Repository.DependencyGraphManifests.Edges = append(respData.Repository.DependencyGraphManifests.Edges, data.Repository.DependencyGraphManifests.Edges...)
	}

	depdendancySize := len(respData.Repository.DependencyGraphManifests.Edges)
	if depdendancySize == 0 { //hack
		skipErrorNextTime(repo)
	}
	return &respData, nil
}

func skipErrorNextTime(repo string) {
	logrus.Infof("Skipping this repo due to some error: %s\n", repo)
	fd, _ := os.Create("./data/" + repo + ".csv")
	defer fd.Close()
}

func (service *Service) execute(req *graphql.Request, respData interface{}) error {
	token := os.Getenv("GITHUB_GRAPHQL_CHECK")
	req.Header.Add("Authorization", "bearer "+token)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	if err := service.gClient.Run(ctx, req, respData); err != nil {

		// catch timeout errors
		if strings.Contains(err.Error(), "timedout") != false || strings.Contains(err.Error(), "loading") != false {
			logrus.Infof("WARNING - recovering from a timeout: %s\n", err)
			return service.execute(req, respData)
		}

		log.Println(errors.Wrap(err, "NewService graphql.Client.Run() failed"))
		return err
	}

	return nil

}

func (service *Service) GetRepos(org string, after string) (*GithubRepositoriesResponse, error) {
	afterStr := ""
	if len(after) > 0 {
		afterStr = ", after:\"" + after + "\""
	}
	req := graphql.NewRequest(`
	{
		repos: search(query: "org:` + org + ` archived:false pushed:>=2020-02-03", type: REPOSITORY, first: 100` + afterStr + `) {
		  repositoryCount
		  pageInfo { endCursor startCursor hasNextPage }
		  edges {
			node {
			  ... on Repository {
				name
				nameWithOwner
				primaryLanguage {
				  name
				} 
			  }
		  }
		}
	  }
	  }
	`)

	var respData GithubRepositoriesResponse
	if err := service.execute(req, &respData); err != nil {
		return nil, err
	}

	if respData.Repos.PageInfo.HasNextPage == true {
		data, err := service.GetRepos(org, respData.Repos.PageInfo.EndCursor)
		if err != nil {
			return nil, err
		}

		respData.Repos.Edges = append(respData.Repos.Edges, data.Repos.Edges...)

	}
	return &respData, nil

}
