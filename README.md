## Purpose

Dump the same data you can get from
https://github.com/orgs/[ORG]/insights/dependencies but export the results
to a csv or google sheet

# How to run

```
export GITHUB_GRAPHQL_CHECK=<GITHUB SSO enabled token which has permissions to view your org>
```

```
* Setup a googleapp following these directions https://developers.google.com/sheets/api/quickstart/go
* Make sure the credentials.json file runs in the same directory as the binary
```

```
make run
```

## Algorithm

- create authorized https client to the graphql endpoint
- get all repos added to in the last 6 months which is not archived
- if the repo contains the lang supported by the github dependencies get the
  license dump
- otherwise mark for manual intervention

## OUTPUT

Service , github repo, lang, lib, license

## DOCKER

Github builds and pushes the docker images to docker.pkg.github.com.
Splitting testing bettween builds can be sped up with a cache in the same github
action run, not between builds.`

## Known Problems

- [Github limits any search to 1000](https://github.community/t/graphql-github-api-how-to-get-more-than-1000-pull-requests/13838/11)
- [Search limitation](https://docs.github.com/en/graphql/reference/queries#searchresultitemconnection)
