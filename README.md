# source{d} rovers [![Build Status](http://drone.srcd.host/api/badges/src-d/rovers/status.svg)](http://drone.srcd.host/src-d/rovers) [![codecov.io](https://codecov.io/github/src-d/rovers/coverage.svg?branch=master&token=ObiptJsBpW)](https://codecov.io/github/src-d/rovers?branch=master)

Source{d} rovers defines (sub)commands for retrieving repositories urls from different external services.

Type `help` fore commands info.

## Services

- `github`: GitHub API crawler for GitHub repositories.
- `cgit`: Cgit urls scraper.

## Usage

This service uses MongoDB and Beanstalk.

To execute test locally you need to run Beanstalk and MongoDB too:

```shell
  docker run -d -p 11300:11300 schickling/beanstalkd
  docker run --name some-mongo -d -p 27017:27017 mongo
  go test ./...
```
