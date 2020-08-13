package repository

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/dathan/go-github-license-check/pkg/gitrepos"
	"github.com/dathan/go-github-license-check/pkg/graphgit"
	"github.com/dathan/go-github-license-check/pkg/license"
	"github.com/dathan/go-github-license-check/pkg/sheets"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type GitHubRepository struct {
	gihub  *graphgit.Service
	sheets *sheets.Service
}

func NewRepository() *GitHubRepository {
	lic := &GitHubRepository{}
	lic.gihub = graphgit.NewService()
	lic.sheets = sheets.NewService()

	return lic
}

func (ghr *GitHubRepository) GetLicenses(owner, repo string) (license.LicenseCheckResults, error) {

	filename := "./data/" + repo + ".csv"
	// better semiphore is needed
	if fileExists(filename) {
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func (ghr *GitHubRepository) SaveLicenses(res license.LicenseCheckResults) error {
	/*
		spew.Config.Indent = "\t"
		spew.Dump(res)
	*/

	if err := ghr.sheets.Save(res); err != nil {
		return err
	}

	//TODO put this in its own domain so you can save to sql,csv,google-sheets
	splitResult := strings.Split(res[0].GitHubRepo, "/")
	//TODO make this configurable
	filename := "./data/" + splitResult[len(splitResult)-1] + ".csv"
	// better semiphore is needed
	if fileExists(filename) {
		logrus.Warningf("respository.SaveLicenses(): Filename: %s exists...skipping", filename)
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "creating a file - improve this")
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range res {
		err := writer.Write([]string{value.GitHubRepo, value.Dependency, value.DependencyLicense, value.Lang})
		if err != nil {
			return err
		}
	}

	return nil

}

func (ghr *GitHubRepository) GetRepos(org string) (gitrepos.Repos, error) {

	log.Infof("Getting recent non-archived repos for the ORG: %s", org)

	res, err := ghr.gihub.GetRepos(org, "")
	if err != nil {
		return nil, err
	}

	var resp gitrepos.Repos
	for _, repo := range res.Repos.Edges {
		// I could use the input but let's be explicit since we are converting one type to another
		orgSplit := strings.Split(repo.Node.NameWithOwner, "/")

		resp = append(resp, gitrepos.Repo{
			Org:  orgSplit[0],
			Name: repo.Node.Name,
			Lang: repo.Node.PrimaryLanguage.Name,
		})
	}

	log.Infof("Total Number of Repos : %d", len(resp))
	return resp, nil
}
