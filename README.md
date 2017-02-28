# rovers

[![Build Status](http://drone.srcd.host/api/badges/src-d/rovers/status.svg)](http://drone.srcd.host/src-d/rovers)
[![codecov.io](https://codecov.io/github/src-d/rovers/coverage.svg?branch=master&token=ObiptJsBpW)](https://codecov.io/github/src-d/rovers?branch=master)

**rovers** is a service to retrieve repository URLs from multiple repository
hosting providers.

Type `help` fore commands info.

## Supported Providers

### GitHub

Uses the GitHub API to get new repositories. Requires a GitHub API token.

### Bitbucket

Uses the Bitbucket API to get new repositories as an anonymous user.

### cgit

cgit is a popular service to create your own self-hosted Git repository provider.
Rovers uses Bing search to discover cgit instances online and tracks them to get
new repositories. Requires a Bing API key.

## Installation

```
go get github.com/src-d/rovers/...
```

## Usage

Run `rovers --help` to get help about the supported commands and their options.

## Test

This service uses PostgreSQL and RabbitMQ.

To execute test locally you need to run RabbitMQ and PostgreSQL too:

```shell
  docker run -d -p 5672:5672 rabbitmq:3
  docker run --name some-postgres -d -p 5432:5432 library/postgres
  go test ./...
```
