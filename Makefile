REPONAME = monitors
VERSION = dev

.PHONY: all build_release 

all: build_release

build_release:
	$(MAKE) -C ./environment/cpu
	$(MAKE) -C ./dbms/sqlquery/mysql

build_container:
	$(MAKE) -C ./environment/cpu build_container
	$(MAKE) -C ./dbms/sqlquery/mysql build_container

test: build_release
	$(MAKE) -C ./environment/cpu test
	$(MAKE) -C ./dbms/sqlquery/mysql test