package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/dathan/go-github-license-check/pkg/csv"
	"github.com/dathan/go-github-license-check/pkg/gitrepos"
	"github.com/dathan/go-github-license-check/pkg/graphgit"
	"github.com/dathan/go-github-license-check/pkg/license"
	"github.com/dathan/go-github-license-check/pkg/sheets"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// GitHubRepository is a wrapper around various services to abstract the list and save to repository from the caller
type GitHubRepository struct {
	gihub  *graphgit.Service
	sheets *sheets.Service
	csv    *csv.Service
}

// NewRepository returns the service wrapper
func NewRepository() *GitHubRepository {

	lic := &GitHubRepository{}
	lic.gihub = graphgit.NewService()
	lic.sheets = sheets.NewService()
	lic.csv = csv.NewService()

	return lic
}

// GetLicenses returns the graphql license result or an error
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

// SaveLicenses abstracts licenses
func (ghr *GitHubRepository) SaveLicenses(res license.LicenseCheckResults) error {
	/*
		spew.Config.Indent = "\t"
		spew.Dump(res)
	*/

	/* uncomment if you want to save to a spreadsheet
	if err := ghr.sheets.Save(res); err != nil {
		return err
	}
	*/

	if err := ghr.csv.Save(res); err != nil {
		return err
	}

	return nil

}

// GetRepos return the Repos struct the dependant GetRepos in github has an issue with getting the foll data
func (ghr *GitHubRepository) GetRepos(org string) (gitrepos.Repos, error) {

	log.Infof("Getting recent non-archived repos for the ORG: %s", org)
	repoArg := graphgit.ReposArg{
		Org:   org,
		After: "",
	}
	res, err := ghr.gihub.GetRepos(repoArg)
	if err != nil {
		return nil, err
	}

	//TODO: Remove hack to fill in missing repos
	for i := 2013; i <= 2021; i++ {
		if repoArg.Since, err = time.Parse("2006-01-01", fmt.Sprintf("%d-01-01", i)); err != nil {
			return nil, err
		}

		res2, err := ghr.gihub.GetRepos(repoArg)

		if err != nil {
			return nil, err
		}

		res.Repos.Edges = append(res.Repos.Edges, res2.Repos.Edges...)
	}

	var resp gitrepos.Repos
	var dedupe map[string]string = make(map[string]string)
	for _, repo := range res.Repos.Edges {

		if _, ok := dedupe[repo.Node.Name]; !ok {
			// I could use the input but let's be explicit since we are converting one type to another
			orgSplit := strings.Split(repo.Node.NameWithOwner, "/")

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
