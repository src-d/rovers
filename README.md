# rovers [![Build Status](https://travis-ci.org/src-d/rovers.svg?branch=master)](https://travis-ci.org/src-d/rovers) [![codecov.io](https://codecov.io/github/src-d/rovers/coverage.svg?branch=master&token=ObiptJsBpW)](https://codecov.io/github/src-d/rovers?branch=master)

**rovers** is a service to retrieve repository URLs from multiple repository
hosting providers.

Type `help` fore commands info.

# Quick start using docker images

## Download docker images

Get the last version of rovers spark image:

```bash
docker pull srcd/rovers
```

Also, you will need Postgres and RabbitMQ

```bash
docker pull postgres:9.6-alpine
docker pull rabbitmq:3-management
```

## Start everything

Start Rabbit and Postgres

```bash
docker run -d --hostname postgres --name postgres -e POSTGRES_PASSWORD=testing -p 5432:5432 -e POSTGRES_USER=testing postgres:9.6-alpine
```
```bash
docker run -d --hostname rabbitmq --name rabbitmq -p 8081:15672 -p 5672:5672 rabbitmq:3-management
```

Then, you can execute rovers using docker:
```bash
docker run --name rovers --link rabbitmq --link postgres \
  -e CONFIG_GITHUB_TOKEN=[REPLACEWITHGHKEY] \
  -e CONFIG_BING_KEY=[REPLACEWITHBINGKEY] \
  srcd/rovers /bin/sh -c "rovers initdb; rovers repos --queue=rovers"
```
After that, rovers will generate a lot of 'mentions' (git repositories found on the internet), and sending them to the 'rovers' queue in Rabbit.

Finally, you can use [Borges](https://github.com/src-d/borges) to fetch the content of these repositories.

# Supported Providers

All the supported providers are used by default. In case you need run only some of them you must use the `--provider` flag:
```bash
rovers repos --provider=github --provider=bitbucket
```

## Generate needed API keys for the providers

To be able to fetch github and cgit repositories, you should create several API keys:

- Get Github token: https://github.com/settings/tokens
- Get Bing token (Bing is the search engine used to fetch cgit repositories from internet): https://azure.microsoft.com/en-us/pricing/details/cognitive-services/search-api/web/

## GitHub

Uses the GitHub API to get new repositories. Requires a GitHub API token. You can set the token through the environment variable:

```bash
$ export CONFIG_GITHUB_TOKEN=github-token
```

## Bitbucket

Uses the Bitbucket API to get new repositories as an anonymous user.

**Note:** the first time this provider runs it may take a while to find git repositories, as the vast majority of repositories at the beginning of the API results are mercurial repositories.

### cgit

cgit is a popular service to create your own self-hosted Git repository provider.
Rovers uses Bing search to discover cgit instances online and tracks them to get
new repositories. Requires a Bing API key. You can set the key through the environment variable:

```bash
$ export CONFIG_BING_KEY=bing-api-key
```

# Installation

```
go get -u github.com/src-d/rovers/...
```

# Usage

Run `rovers --help` to get help about the supported commands and their options.


To initialize the database schemas. You need to run this command only once.
```
rovers initdb
```
To start collecting repository URLs
```
rovers repos --provider=github
```

You can configure rovers by environment variables:

Providers:
- `CONFIG_GITHUB_TOKEN` to set the github api key.
- `CONFIG_BING_KEY` to set the cgit api key.

Broker:
- `CONFIG_BROKER_URL` to set the broker url, by default `amqp://guest:guest@localhost:5672`

Database:
- `CONFIG_DBUSER`: database username,by default if not set `testing`
- `CONFIG_DBPASS`: database user password,by default if not set `testing`
- `CONFIG_DBHOST`: database host,by default if not set `0.0.0.0`
- `CONFIG_DBPORT`: database port,by default if not set `5432`
- `CONFIG_DBNAME`: database name,by default if not set `testing`
- `CONFIG_DBSSLMODE`: ssl mode to use,by default if not set `disable`
- `CONFIG_DBTIMEOUT`: connection timeout,by default if not set `30s`
- `CONFIG_DBAPPNAME`: application name

# Development

## Build

- `rm Makefile.main; rm -rf .ci` to make sure you will have the last Makefile changes.
- `make dependencies` to download vendor dependencies.
- `make packages` to generate binaries for several platforms.

You will find the built binaries under `./build`.

## Test

This service uses PostgreSQL and RabbitMQ.

To execute test locally you need to run RabbitMQ and PostgreSQL.

```bash
  docker run --hostname postgres --name postgres -e POSTGRES_PASSWORD=testing -p 5432:5432 -e POSTGRES_USER=testing -d postgres
  docker run -d --hostname rabbit --name rabbit -p 8081:15672 -p 5672:5672 rabbitmq:3-management
  go test ./...
```

# Running Rovers in Kubernetes

You can use the official [Helm](https://github.com/kubernetes/helm) [chart](https://github.com/src-d/charts/tree/master/rovers) to deploy Rovers in your kubernetes cluster.
