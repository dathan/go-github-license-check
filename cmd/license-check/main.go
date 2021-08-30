package main

import (
	"flag"
	"os"

	"github.com/dathan/go-github-license-check/pkg/gitrepos"
	"github.com/dathan/go-github-license-check/pkg/license"
	"github.com/dathan/go-github-license-check/pkg/repository"
	"github.com/sirupsen/logrus"
)

func main() {

	var ro *repository.GitHubRepository = repository.NewRepository()
	var saving *license.Service = license.NewService(ro)
	var repos gitrepos.Repos
	var enable bool = true
	var err error

	// validate inputs
	org := flag.String("org", "WeConnect", "Provide your github organization to crawl")
	flag.Parse()
	if len(*org) == 0 {
		flag.Usage()
		os.Exit(-1)
	}

	// hack to set up environment
	_, err = os.Stat("./data")
	if os.IsNotExist(err) {
		err := os.Mkdir("./data", 0755)
		if err != nil {
			panic("Hack environment where the csv is unable to be made. Accepting PRs")
		}
	}

	listing := gitrepos.NewService(ro)
	if enable {
		repos, err = listing.ListRepos(*org)
		if err != nil {
			logrus.Errorf("TEMP ERROR: %s", err)
			os.Exit(-1)
		}
	}
	logrus.Infof("len repos: %d", len(repos))
	//spew.Dump(repos)

	//orgs = gitrepos.Repos{{Org: "WeConnect", Name: "referral-web-app"}}

	for _, org := range repos {
		//TODO call a listing service to list all the repos to perform the update
		if err := saving.Execute(org.Org, org.Name); err != nil {
			logrus.Printf("TEMP ERROR: %s", err)
			os.Exit(-1)
		}
	}

	os.Exit(0)

}
