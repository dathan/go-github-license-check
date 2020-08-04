package main

import (
	"fmt"
	"os"

	"github.com/dathan/go-github-license-check/pkg/gitrepos"
	"github.com/dathan/go-github-license-check/pkg/license"
	"github.com/dathan/go-github-license-check/pkg/repository"
)

func main() {

	var ro *repository.GitHubRepository = repository.NewRepository()
	var saving *license.Service = license.NewService(ro)
	var orgs gitrepos.Repos
	var enable bool = true
	var err error

	var listing *gitrepos.Service = gitrepos.NewService(ro)
	if enable {
		orgs, err = listing.ListRepos("WeConnect")
		if err != nil {
			fmt.Printf("TEMP ERROR: %s\n", err)
			os.Exit(-1)
		}
	}

	//orgs = gitrepos.Repos{{Org: "WeConnect", Name: "referral-web-app"}}

	for _, org := range orgs {
		//TODO call a listing service to list all the repos to perform the update
		if err := saving.Execute(org.Org, org.Name); err != nil {
			fmt.Printf("TEMP ERROR: %s\n", err)
			os.Exit(-1)
		}
	}

	os.Exit(0)

}
