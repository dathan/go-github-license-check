package graphgit

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

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

	return service
}

// return the licenses
func (service *Service) GetLicenses(owner, repo, after string) (*DependencyGraphManifestsResponse, error) {

	afterStr := ""
	if len(after) > 0 {
		afterStr = ", after:\"" + after + "\""
	}

	query := fmt.Sprintf(`query {
		repository(name:"%s", owner:"%s" %s) {
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
  	}`, repo, owner, afterStr)

	logrus.Warn(query)
	req := graphql.NewRequest(query)

	// this is a beta version graphql call
	req.Header.Add("Accept", "application/vnd.github.hawkgirl-preview+json")

	// run it and capture the response
	var respData DependencyGraphManifestsResponse
	if err := service.execute(req, &respData); err != nil {
		/*
			if strings.Contains(err.Error(), "decoding") {
				skipErrorNextTime(repo)
				return nil, nil
			}
		*/
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
	logrus.Infof("Skipping this repo due to some error: %s", repo)
	fd, _ := os.Create("./data/" + repo + ".csv")
	defer fd.Close()
}

func (service *Service) execute(req *graphql.Request, respData interface{}) error {
	token := os.Getenv("GITHUB_GRAPHQL_CHECK")

	// verify requirements
	if token == "" {
		logrus.Error("GITHUB_GRAPHQL_CHECK is not set")
		return errors.New("GITHUB_GRAPHQL_CHECK is not set")
	}

	req.Header.Add("Authorization", "bearer "+token)

	// define a Context for the request
	ctx := context.Background()

	//service.gClient.Log = func(s string) { logrus.Warn(s) }
	// run it and capture the response
	if err := service.gClient.Run(ctx, req, respData); err != nil {

		// catch timeout errors
		if strings.Contains(err.Error(), "timedout") || strings.Contains(err.Error(), "loading") {
			logrus.Warningf("recovering from a timeout: %s", err)
			time.Sleep(3 * time.Minute)
			return service.execute(req, respData)
		}

		logrus.Warn(err)
		if strings.Contains(err.Error(), "invalid character") != false {
			return nil
		}
		return err
	}

	return nil

}

// GHRequestJSON a helper function to build a graphql search query in string format
// The same query can be run at https://docs.github.com/en/graphql/overview/explorer
func (service *Service) GHRequestJSON(org, after, since string) string {
	// need to change the query from after to created since we need to walk the index
	// https://github.community/t/graphql-github-api-how-to-get-more-than-1000-pull-requests/13838/11
	afterStr := ""
	if len(after) > 0 {
		// going to use a cursort for the base query, then call another method for filling in the gaps
		afterStr = fmt.Sprintf("after:\"%s\",", after)
	}
	created := ""
	if len(since) > 0 {
		created = fmt.Sprintf(" created:>%s", since)
	}

	requestJSON := fmt.Sprintf(`
	{
		repos: search(query: "org:%s archived:false %s sort:updated-asc", %s type: REPOSITORY, first: 100) {
		  repositoryCount
		  pageInfo { endCursor startCursor hasNextPage }
		  edges {
			node {
			  ... on Repository {
				name
				nameWithOwner
				createdAt
				updatedAt
				primaryLanguage {
				  name
				} 
			  }
		  }
		}
	  }
	  }
	`, org, created, afterStr)
	return requestJSON
}

func (service *Service) GetRepos(org string, after string) (*GithubRepositoriesResponse, error) {

	requestJSON := service.GHRequestJSON(org, after, "")
	logrus.Infof("About to make a request: %s", requestJSON)

	respData, err := service.getGHResponse(requestJSON)
	if err != nil {
		logrus.Warnf("Error for GHResponse: %s", err)
		return nil, err
	}

	if respData.Repos.PageInfo.HasNextPage {
		//pos := len(respData.Repos.Edges) - 1
		//pos = 0
		//data, err := service.GetRepos(org, respData.Repos.Edges[pos].Node.CreatedAt.Format("2006-01-02"))
		data, err := service.GetRepos(org, respData.Repos.PageInfo.EndCursor)
		if err != nil {
			logrus.Warnf("ERROR: %s", err)
			return nil, err
		}

		respData.Repos.Edges = append(respData.Repos.Edges, data.Repos.Edges...)

	}

	return respData, nil

}

func (service *Service) GetReposSince(org string, after string, since string) (*GithubRepositoriesResponse, error) {

	requestJSON := service.GHRequestJSON(org, after, since)
	logrus.Infof("About to make a request: %s", requestJSON)

	respData, err := service.getGHResponse(requestJSON)
	if err != nil {
		logrus.Warnf("Error for GHResponse: %s", err)
		return nil, err
	}

	if respData.Repos.PageInfo.HasNextPage {
		data, err := service.GetReposSince(org, respData.Repos.PageInfo.EndCursor, since)
		if err != nil {
			logrus.Warnf("ERROR: %s", err)
			return nil, err
		}

		respData.Repos.Edges = append(respData.Repos.Edges, data.Repos.Edges...)

	}

	return respData, nil

}

func (service *Service) getGHResponse(requestJSON string) (*GithubRepositoriesResponse, error) {
	req := graphql.NewRequest(requestJSON)

	var respData GithubRepositoriesResponse
	if err := service.execute(req, &respData); err != nil {

		logrus.Warnf("Received an error: %s", err)
		return nil, err
	}

	startPos, _ := base64.StdEncoding.DecodeString(respData.Repos.PageInfo.StartCursor)
	endPos, _ := base64.StdEncoding.DecodeString(respData.Repos.PageInfo.EndCursor)

	logrus.Infof("Repos: %d Paginfo Struct (%+v) Edges %d out of %d for start: %s and end: %s", respData.Repos.RepositoryCount, respData.Repos.PageInfo, respData.Repos.RepositoryCount, len(respData.Repos.Edges), string(startPos), string(endPos))
	return &respData, nil
}
