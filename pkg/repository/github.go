package repository

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/dathan/go-github-license-check/pkg/gitrepos"
	"github.com/dathan/go-github-license-check/pkg/graphgit"
	"github.com/dathan/go-github-license-check/pkg/license"
	"github.com/pkg/errors"
)

type GitHubRepository struct {
	gihub *graphgit.Service
}

func NewRepository() *GitHubRepository {
	lic := &GitHubRepository{}
	lic.gihub = graphgit.NewService()

	return lic
}

func (ghr *GitHubRepository) GetLicenses(owner, repo string) (license.LicenseCheckResults, error) {

	filename := "./data/" + repo + ".csv"
	// better semiphore is needed
	if fileExists(filename) {
		fmt.Printf("Warning Skipping get - %s exists..skipping\n", filename)
		return nil, nil
	}

	response, err := ghr.gihub.GetLicenses(owner, repo, "")
	if err != nil {
		return nil, err
	}

	if response == nil {
		fmt.Printf("Warning REPO: %s Data does not exist\n")
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
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (ghr *GitHubRepository) SaveLicenses(res license.LicenseCheckResults) error {
	/*
		spew.Config.Indent = "\t"
		spew.Dump(res)
	*/
	//TODO put this in its own domain so you can save to sql,csv,google-sheets
	splitResult := strings.Split(res[0].GitHubRepo, "/")
	//TODO make this configurable
	filename := "./data/" + splitResult[len(splitResult)-1] + ".csv"
	// better semiphore is needed
	if fileExists(filename) {
		fmt.Printf("Warning - %s exists..skipping\n", filename)
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

	res, err := ghr.gihub.GetRepos(org, "")
	if err != nil {
		return nil, err
	}

	var resp gitrepos.Repos
	for _, repo := range res.Repos.Edges {
		orgSplit := strings.Split(repo.Node.NameWithOwner, "/")

		resp = append(resp, gitrepos.Repo{
			Org:  orgSplit[0],
			Name: repo.Node.Name,
			Lang: repo.Node.PrimaryLanguage.Name,
		})
	}

	return resp, nil
}
