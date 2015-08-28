# Package configuration
PROJECT = srcd-rovers
COMMANDS = rovers
DEPENDENCIES = github.com/aktau/github-release
INTERNAL = srcd-domain

# Github configuration
GITHUB_USER := tyba
GITHUB_REPO := $(PROJECT)
GITHUB_TOKEN := 08763897c930b3ff7f7cebf8da45935350a96b7d

# Environment
BASE_PATH := $(shell pwd)
BUILD_PATH := $(BASE_PATH)/build
SHA1 := $(shell git log --format='%H' -n 1 | cut -c1-10)
BUILD := $(shell date)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

# Packages content
PKG_OS = darwin linux
PKG_ARCH = amd64
PKG_CONTENT =
PKG_TAG = build

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get
GOTEST = $(GOCMD) test
GHRELEASE = github-release

# CircleCI
ifneq ($(origin CI), undefined)
	GITHUB_USER := $(CIRCLE_PROJECT_USERNAME)
	GITHUB_REPO := $(CIRCLE_PROJECT_REPONAME)
	SHA1 := $(shell echo $(CIRCLE_SHA1) | cut -c1-10)
	BRANCH := $(CIRCLE_BRANCH)
	GOPATH = /home/ubuntu/.go_workspace/
endif

# Exports
export GITHUB_TOKEN
export GITHUB_USER
export GITHUB_REPO

# Rules
all: clean upload

dependencies:
	for i in $(DEPENDENCIES); do $(GOGET) $$i; done

	for pkg in $(INTERNAL); do \
		go get github.com/tyba/$${pkg}/...; \
		cd $(GOPATH)/src/github.com/tyba/$${pkg} && git checkout $(BRANCH); \
	done; \
	$(GOGET) -d -v ./...


test:
	go test -v ./...

packages: dependencies
	for os in $(PKG_OS); do \
		for arch in $(PKG_ARCH); do \
			cd $(BASE_PATH); \
			mkdir -p $(BUILD_PATH)/$(PROJECT)_$${os}_$${arch}; \
			for cmd in $(COMMANDS); do \
				GOOS=$${os} GOARCH=$${arch} $(GOCMD) build -ldflags "-X main.version $(SHA1) -X main.build \"$(BUILD)\"" -o $(BUILD_PATH)/$(PROJECT)_$${os}_$${arch}/$${cmd} $${cmd}.go ; \
			done; \
			for content in $(PKG_CONTENT); do \
				cp -rf $${content} $(BUILD_PATH)/$(PROJECT)_$${os}_$${arch}/; \
			done; \
			cd  $(BUILD_PATH) && tar -cvzf $(BUILD_PATH)/$(PROJECT)_$${os}_$${arch}.tar.gz $(PROJECT)_$${os}_$${arch}/; \
		done; \
	done;

upload: packages
	cd $(BASE_PATH); \
	$(GHRELEASE) delete --tag $(PKG_TAG); \
	$(GHRELEASE) release --tag $(PKG_TAG) --name "$(PKG_TAG) ($(SHA1))"; \
	for os in $(PKG_OS); do \
		for arch in $(PKG_ARCH); do \
			$(GHRELEASE) upload \
		    --tag $(PKG_TAG) \
				--name "$(PROJECT)_$${os}_$${arch}.tar.gz" \
				--file $(BUILD_PATH)/$(PROJECT)_$${os}_$${arch}.tar.gz; \
		done; \
	done;

clean:
	rm -rf $(BUILD_PATH)
	$(GOCLEAN) .
