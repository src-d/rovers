# Package configuration
PROJECT = rovers
COMMANDS = rovers
REQUIREMENTS = mongo

# Including devops Makefile
MAKEFILE = Makefile.main
DEVOPS_REPOSITORY = git@github.com:src-d/devops.git
DEVOPS_FOLDER = .devops
TRAVIS_FOLDER = .travis

$(MAKEFILE):
	git clone $(DEVOPS_REPOSITORY) $(DEVOPS_FOLDER)
	cp $(DEVOPS_FOLDER)/common/$(MAKEFILE) ./
	cp -r $(DEVOPS_FOLDER)/common/travis $(TRAVIS_FOLDER)
	rm -rf $(DEVOPS_FOLDER)

include $(MAKEFILE)
