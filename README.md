# rovers
[![Build Status](https://travis-ci.com/src-d/rovers.svg?branch=master)](https://travis-ci.com/src-d/rovers)
[![codecov](https://codecov.io/gh/src-d/rovers/branch/master/graph/badge.svg)](https://codecov.io/gh/src-d/rovers)

**rovers** is a service to retrieve repository URLs from multiple repository
hosting providers and store the URLs and metadata to a PostgreSQL database while
issuing messages to a RabbitMQ queue so you can connect other processes to it.

# Quick start using docker images

## Download docker images

Get the latest rovers image:

```bash
docker pull srcd/rovers
```

### Start everything using docker-compose

Install [docker-compose](https://docs.docker.com/compose/install/).

Start Rabbit and Postgres

```bash
docker-compose -d up rovers-postgres rovers-rabbitmq
```

Export as environment variables the API keys ([see section](#supported-providers)). Then, you can execute rovers:
```bash
docker-compose up --no-deps rovers
```

If you need to run just some of the available supported-providers, you can use this command replacing the flag `--provider` with those providers you want to use:
```bash
docker-compose run --rm --no-deps --service-ports rovers /bin/sh -c "rovers initdb; rovers repos --provider=github --provider=bitbucket"
```

After that, rovers will generate a lot of 'mentions' (git repositories found on the internet), and sending them to the 'rovers' queue in Rabbit.

To stop and remove all the containers running `docker-compose down`.

Finally, you can use [Borges](https://github.com/src-d/borges) to fetch the content of these repositories.

# Supported Providers

All the supported providers are used by default. In case you need to run only some of them you must use the `--provider` flag:
```bash
rovers repos --provider=github --provider=bitbucket
```

## Generate needed API keys for the providers

To be able to fetch github and cgit repositories, you should create several API keys:

- Get Github token: https://github.com/settings/tokens
- Get Bing token (Bing is the search engine used to fetch cgit repositories from internet): https://azure.microsoft.com/en-us/pricing/details/cognitive-services/search-api/web/

## GitHub

Uses the GitHub API to get new repositories. Requires a GitHub API token, which will only need repository read access. You can set the token through the environment variable:

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
The commands `initdb`, `replay` and `repos` also have their own `--help` switch
for specific options related to them.


To initialize the database schemas. You need to run this command only once.
```
rovers initdb
```
To start collecting repository URLs
```
rovers repos --provider=github
```

You can configure rovers using environment variables:

Providers:
- `CONFIG_GITHUB_TOKEN` to set the github api key.
- `CONFIG_BING_KEY` to set the cgit api key.

Broker:
- `CONFIG_BROKER_URL` to set the broker url, by default `amqp://guest:guest@localhost:5672`

Database:
- `CONFIG_DBUSER`: database username, by default if not set `testing`
- `CONFIG_DBPASS`: database user password, by default if not set `testing`
- `CONFIG_DBHOST`: database host, by default if not set `0.0.0.0`
- `CONFIG_DBPORT`: database port, by default if not set `5432`
- `CONFIG_DBNAME`: database name, by default if not set `testing`
- `CONFIG_DBSSLMODE`: ssl mode to use, by default if not set `disable`
- `CONFIG_DBTIMEOUT`: connection timeout, by default if not set `30s`
- `CONFIG_DBAPPNAME`: application name

# Development

## Build

- `rm Makefile.main; rm -rf .ci` to make sure you will have the last Makefile changes.
- `make dependencies` to download vendor dependencies.
- `make packages` to generate binaries for several platforms.

You will find the built binaries under `./build`.

## Test

This service uses PostgreSQL and RabbitMQ.

To execute tests locally you need to run RabbitMQ and PostgreSQL.

```bash
  docker-compose up -d rovers-postgres rovers-rabbitmq
  make test
```

# Running Rovers in Kubernetes

You can use the official [Helm](https://github.com/kubernetes/helm) [chart](https://github.com/src-d/charts/tree/master/stable/rovers) to deploy Rovers in your Kubernetes cluster.
