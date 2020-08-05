## Purpose

Dump the same data you can get from
https://github.com/orgs/[ORG]/insights/dependencies but export the results
to a csv or google sheet


# How to run
```
export GITHUB_GRAPHQL_CHECK=<GITHUB SSO enabled token which has permissions to view your org>
```

```
make run
```


## Algorithm

* create authorized https client to the graphql endpoint
* get all repos added to in the last 6 months which is not archived
* if the repo contains the lang supported by the github dependencies get the
  license dump
* otherwise mark for manual intervention

## OUTPUT
 Service , github repo, lang, lib, license

