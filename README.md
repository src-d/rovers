# source{d} rovers [![Circle CI](https://circleci.com/gh/tyba/srcd-rovers.svg?style=svg&circle-token=662647aec7bf50cd00c97487a868437f8dd0fb6e)](https://circleci.com/gh/tyba/srcd-rovers) [![Circle CI](https://circleci.com/gh/tyba/srcd-rovers/tree/dev.svg?style=svg&circle-token=662647aec7bf50cd00c97487a868437f8dd0fb6e)](https://circleci.com/gh/tyba/srcd-rovers/tree/dev)

Source{d} rovers defines (sub)commands for retrieving different external services.

## Services

- `augur`: Augur API crawler, takes all emails from sourced.people collection and fetches whatever info Augur may have about them. It works incrementally. For every email we save when was the last time we fetched its info. **WARNING:** This is a very very long process, it should take about 44 days.
- `bitbucket`: Bitbucket API crawler for Bitbucket repositories.
- `github`: web crawler for GitHub profiles. Requires `augur` data, GitHub profile URL.
- `github-api`: GitHub API crawler for GitHub repositories.
- `github-api-users`: GitHub API crawler for GitHub users.
- `linkedin`: LinkedIn company employees importer.
- `twitter`: Twitter web crawler for Twitter profiles (followers, following, tweets, location, bio, ...). Requires `augur` data, Twitter profile URL.

## Usage

- `github` and `twitter` commands require `augur` command to be run first.
- Both `github-api` and `github-api-users` require a local Mongo.

## Utilities

The `utils` directory contains scripts that **require having a Go toolchain installed** but simplify usage of a certain pipeline.
