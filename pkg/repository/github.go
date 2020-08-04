package repository

import (
	"github.com/dathan/go-github-license-check/pkg/graphgit"
	"github.com/dathan/go-github-license-check/pkg/license"
	"github.com/davecgh/go-spew/spew"
)

type GitHubRepository struct {
	gihub *graphgit.Service
}

func NewRepository() license.Repository {
	lic := &GitHubRepository{}
	lic.gihub = graphgit.NewService()

	return lic
}

func (ghr *GitHubRepository) GetLicenses(owner, repo string) (license.LicenseCheckResults, error) {

	response, err := ghr.gihub.GetLicenses(owner, repo)
	if err != nil {
		return nil, err
	}

	var res license.LicenseCheckResults
	for _, edges := range response.Repository.DependencyGraphManifests.Edges {
		for _, node := range edges.Node.Dependencies.Nodes {
			res = append(res, license.LicenseCheckResult{
				GitHubRepo:        "github.com/" + owner + "/" + repo,
				Lang:              "todo",
				Dependency:        node.PackageName,
				DependencyLicense: node.Repository.LicenseInfo.Name,
			})
		}
	}

	return res, nil
}

func (ghr *GitHubRepository) SaveLicenses(res license.LicenseCheckResults) error {

	spew.Config.Indent = "\t"

	spew.Dump(res)

	return nil

}
