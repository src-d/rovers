# rovers [![Build Status](https://travis-ci.org/src-d/rovers.svg?branch=master)](https://travis-ci.org/src-d/rovers) [![codecov.io](https://codecov.io/github/src-d/rovers/coverage.svg?branch=master&token=ObiptJsBpW)](https://codecov.io/github/src-d/rovers?branch=master)

**rovers** is a service to retrieve repository URLs from multiple repository
hosting providers.

Type `help` fore commands info.

## Quick start using docker images

### Generate needed API keys for the providers

To be able to fetch github and cgit repositories, you should create several API keys:

- Get Github token: https://github.com/settings/tokens
- Get Bing token (Bing is the search engine used to fetch cgit repositories from internet): https://azure.microsoft.com/en-us/pricing/details/cognitive-services/search-api/web/

### Download docker images

Get the last version of rovers spark image:

```bash
docker pull quay.io/srcd/rovers
```

Also, you will need Postgres and RabbitMQ

```bash
docker pull postgres:9.6-alpine
docker pull rabbitmq:3-management
```

### Start everything

Start Rabbit and Postgres

```bash
docker run -d --name postgres -e POSTGRES_PASSWORD=testing -p 5432:5432 -e POSTGRES_USER=testing postgres:9.6-alpine
```
```bash
docker run -d --hostname rabbitmq --name rabbitmq -p 8081:15672 -p 5672:5672 rabbitmq:3-management
```

Then, you can execute rovers using docker:
```bash
docker run --name rovers --link rabbitmq --link postgres -e CONFIG_DBUSER=testing -e CONFIG_DBPASS=testing -e CONFIG_DBHOST=postgres -e CONFIG_DBNAME=testing -e CONFIG_BROKER_URL=amqp://guest:guest@rabbitmq:5672/ -e CONFIG_GITHUB_TOKEN=[REPLACEWITHGHKEY] -e CONFIG_BING_KEY=[REPLACEWITHBINGKEY] quay.io/srcd/rovers /bin/sh -c "rovers initdb; rovers repos --queue=rovers"
```
After that, rovers will generate a lot of 'mentions' (git repositories found on the internet), and sending them to the 'rovers' queue in Rabbit.

Finally, you can use [Borges](https://github.com/src-d/borges) to fetch the content of these repositories.

## Supported Providers

### GitHub

Uses the GitHub API to get new repositories. Requires a GitHub API token. You can set the token through the environment variable:

```bash
$ export CONFIG_GITHUB_TOKEN=github-token
```

### Bitbucket

Uses the Bitbucket API to get new repositories as an anonymous user.

### cgit

cgit is a popular service to create your own self-hosted Git repository provider.
Rovers uses Bing search to discover cgit instances online and tracks them to get
new repositories. Requires a Bing API key. You can set the key through the environment variable:

```bash
$ export CONFIG_BING_KEY=bing-api-key
```

## Installation

```
go get -u github.com/src-d/rovers/...
```

## Usage

Run `rovers --help` to get help about the supported commands and their options.

## Test

This service uses PostgreSQL and RabbitMQ.

To execute test locally you need to run RabbitMQ and PostgreSQL too. To set broker's URL for RabbitMQ you can do it through the environment variable:

```bash
$ export CONFIG_BROKER_URL=url
```

By default this URL is set to `amqp://guest:guest@localhost:5672/`. To run tests:

```bash
  docker run --name some-postgres -e POSTGRES_PASSWORD=testing -p 5432:5432 -e POSTGRES_USER=testing -d postgres
  docker run -d --hostname rabbit --name rabbit -p 8081:15672 -p 5672:5672 rabbitmq:3-management
  go test ./...
```

# Running Rovers in Kubernetes

You can use the official [Helm](https://github.com/kubernetes/helm) [chart](https://github.com/src-d/charts/tree/master/rovers) to deploy Rovers in your kubernetes cluster.
