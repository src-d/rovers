# source{d} rovers [![Circle CI](https://circleci.com/gh/tyba/srcd-rovers.svg?style=svg&circle-token=662647aec7bf50cd00c97487a868437f8dd0fb6e)](https://circleci.com/gh/tyba/srcd-rovers) [![Circle CI](https://circleci.com/gh/tyba/srcd-rovers/tree/dev.svg?style=svg&circle-token=662647aec7bf50cd00c97487a868437f8dd0fb6e)](https://circleci.com/gh/tyba/srcd-rovers/tree/dev)

Source{d} rovers defines (sub)commands for retrieving different external services.

## Services

- `augur`: Augur API crawler, given a list of emails it fetches information about them.
- `augur-emails`: Inserts email address from either a Mongo database or a file for later processing via the `augur` command.
- `bitbucket`: Bitbucket API crawler for Bitbucket repositories.
- `github`: web crawler for GitHub profiles. Requires `augur` data, GitHub profile URL.
- `github-api`: GitHub API crawler for GitHub repositories.
- `github-api-users`: GitHub API crawler for GitHub users.
- `twitter`: Twitter web crawler for Twitter profiles (followers, following, tweets, location, bio, ...). Requires `augur` data, Twitter profile URL.

## Environment

- `augur` command requires `augur-emails` command to be run first.
- `github` and `twitter` commands require `augur` command to be run first.
- Both `github-api` and `github-api-users` require a local Mongo.
