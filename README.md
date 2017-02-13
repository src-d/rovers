# source{d} rovers [![Build Status](http://drone.srcd.host/api/badges/src-d/rovers/status.svg)](http://drone.srcd.host/src-d/rovers) [![codecov.io](https://codecov.io/github/src-d/rovers/coverage.svg?branch=master&token=ObiptJsBpW)](https://codecov.io/github/src-d/rovers?branch=master)

Source{d} rovers defines (sub)commands for retrieving repositories urls from different external services.

Type `help` fore commands info.

## Services

- `github`: GitHub API crawler for GitHub repositories.
- `cgit`: Cgit urls scraper.

## Usage

This service uses PostgreSQL and RabbitMQ.

To execute test locally you need to run RabbitMQ and PostgreSQL too:

```shell
  docker run -d -p 5672:5672 rabbitmq:3
  docker run --name some-postgres -d -p 5432:5432 library/postgres
  go test ./...
```
