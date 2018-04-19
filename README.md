# rovers [![Build Status](https://travis-ci.org/src-d/rovers.svg?branch=master)](https://travis-ci.org/src-d/rovers) [![codecov.io](https://codecov.io/github/src-d/rovers/coverage.svg?branch=master&token=ObiptJsBpW)](https://codecov.io/github/src-d/rovers?branch=master)

**rovers** is a service to retrieve repository URLs from multiple repository
hosting providers.

Type `help` fore commands info.

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
