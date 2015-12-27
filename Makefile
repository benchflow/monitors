REPONAME = monitors
VERSION = dev

.PHONY: all build_release 

all: build_release

build_release:
	$(MAKE) -C ./environment/cpu

test:
	$(MAKE) -C ./environment/cpu