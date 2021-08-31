package repository

import (
	"github.com/dathan/go-github-license-check/pkg/csv"
	"github.com/dathan/go-github-license-check/pkg/gitrepos"
	"github.com/dathan/go-github-license-check/pkg/graphgit"
	"github.com/dathan/go-github-license-check/pkg/license"
	"github.com/dathan/go-github-license-check/pkg/sheets"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type GitHubRepository struct {
	gihub  *graphgit.Service
	sheets *sheets.Service
	csv    *csv.Service
}

func NewRepository() *GitHubRepository {

	lic := &GitHubRepository{}
	lic.gihub = graphgit.NewService()
	lic.sheets = sheets.NewService()
	lic.csv = csv.NewService()

	return lic
}

func (ghr *GitHubRepository) GetLicenses(owner, repo string) (license.LicenseCheckResults, error) {

	filename := "./data/" + repo + ".csv"
	// better semiphore is needed
	if ghr.csv.FileExists(filename) {
		logrus.Warningf("repository.GetLicenses() Filename: %s exists..skipping", filename)
		return nil, nil
	}

	response, err := ghr.gihub.GetLicenses(owner, repo, "")
	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, nil
	}

	var res license.LicenseCheckResults
	for _, edges := range response.Repository.DependencyGraphManifests.Edges {
		for _, node := range edges.Node.Dependencies.Nodes {
			res = append(res, license.LicenseCheckResult{
				GitHubRepo:        "github.com/" + owner + "/" + repo,
				Lang:              node.Repository.PrimaryLanguage.Name,
				Dependency:        node.PackageName,
				DependencyLicense: node.Repository.LicenseInfo.Name,
			})
		}
	}

	return res, nil
}

func (ghr *GitHubRepository) SaveLicenses(res license.LicenseCheckResults) error {
	/*
		spew.Config.Indent = "\t"
		spew.Dump(res)
	*/
	/*
		if err := ghr.sheets.Save(res); err != nil {
			return err
		}
	*/

	if err := ghr.csv.Save(res); err != nil {
		return err
	}

	return nil

}

func (ghr *GitHubRepository) GetRepos(org string) (gitrepos.Repos, error) {

	log.Infof("Getting recent non-archived repos for the ORG: %s", org)

	res, err := ghr.gihub.GetRepos(org, "")
	if err != nil {
		return nil, err
	}

	res2, err := ghr.gihub.GetReposSince(org, "", "2020-01-01")
	if err != nil {
		return nil, err
	}

	res.Repos.Edges = append(res.Repos.Edges, res2.Repos.Edges...)
	res2, err = ghr.gihub.GetReposSince(org, "", "2019-01-01")
	if err != nil {
		return nil, err
	}

	res.Repos.Edges = append(res.Repos.Edges, res2.Repos.Edges...)
	res2, err = ghr.gihub.GetReposSince(org, "", "2018-01-01")
	if err != nil {
		return nil, err
	}

	res.Repos.Edges = append(res.Repos.Edges, res2.Repos.Edges...)
	res2, err = ghr.gihub.GetReposSince(org, "", "2017-01-01")
	if err != nil {
		return nil, err
	}

	var dedupe map[string]string = make(map[string]string)
	for _, repo := range res.Repos.Edges {
		if _, ok := dedupe[repo.Node.Name]; !ok {
			resp = append(resp, gitrepos.Repo{
				Org:  orgSplit[0],
				Name: repo.Node.Name,
				Lang: repo.Node.PrimaryLanguage.Name,
			})
			dedupe[repo.Node.Name] = repo.Node.Name
		}
	}

	log.Infof("Total Number of Repos : %d -->", len(resp))
	return resp, nil
}
