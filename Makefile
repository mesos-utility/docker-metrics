default: help

HOST_GOLANG_VERSION     = $(go version | cut -d ' ' -f3 | cut -c 3-)
# this variable is used like a function. First arg is the minimum version, Second arg is the version to be checked.
ALLOWED_GO_VERSION      = $(test '$(/bin/echo -e "$(1)\n$(2)" | sort -V | head -n1)' == '$(1)' && echo 'true')

NAME := docker-metrics
COMMIT := $(shell git rev-parse HEAD 2> /dev/null || true)
GITHUB_SRC := github.com/mesos-utility
CURDIR_LINK := $(CURDIR)/vendor/$(GITHUB_SRC)
#export GOPATH := $(CURDIR)/vendor


MKDIR	= mkdir
INSTALL	= install
BIN		= $(BUILD_ROOT)
MAN		= $(BIN)
VERSION	= $(shell git describe --tags)
RELEASE	= 0
RPMSOURCEDIR	= $(shell rpm --eval '%_sourcedir')
RPMSPECDIR	= $(shell rpm --eval '%_specdir')
RPMBUILD = $(shell				\
	if [ -x /usr/bin/rpmbuild ]; then	\
		echo "/usr/bin/rpmbuild";	\
	else					\
		echo "/bin/rpm";		\
	fi )

## Make bin for docker-metrics.
bin:
	#./control build
	go build -i -ldflags "-X github.com/mesos-utility/${NAME}/g.VERSION=${VERSION}" -o ${NAME} main.go

## Get godep and restore dep.
godep:
	@go get -u github.com/tools/godep
	GO15VENDOREXPERIMENT=0 GOPATH=`godep path` godep restore

$(CURDIR_LINK):
	rm -rf $(CURDIR_LINK)
	mkdir -p $(CURDIR_LINK)
	ln -sfn $(CURDIR) $(CURDIR_LINK)/$(NAME)

## Get vet go tools.
vet:
	go get -u golang.org/x/tools/cmd/vet

# `go get github.com/golang/lint/golint`
.golint:
	go get github.com/golang/lint/golint
ifeq ($(call ALLOWED_GO_VERSION,1.5,$(HOST_GOLANG_VERSION)),true)
	golint ./...
endif

## Validate this go project.
validate:
	script/validate-gofmt
#	go vet ./...

## Run test case for this go project.
test:
	go list ./... | grep -v 'vendor' | xargs -L1 go test -v

## Clean everything (including stray volumes).
clean:
#	find . -name '*.created' -exec rm -f {} +
	-rm -rf var
	-rm -f ${NAME}
	-rm -f ${NAME}-*.tar.gz
	-rm -f ${NAME}.spec


# Rpm install action for docker-metrics rpm package build.
rpm-install:
	if [ ! -d $(BIN) ]; then $(MKDIR) -p $(BIN); fi
	$(INSTALL) -m 0755 $(NAME) $(BIN)
	$(INSTALL) -m 0644 cfg.json $(BIN)
	$(INSTALL) -m 0755 control $(BIN)
	[ -d $(MAN) ] || $(MKDIR) -p $(MAN)
	$(INSTALL) -m 0644 README.md $(MAN)

dist: clean
	sed -e "s/@@VERSION@@/$(VERSION)/g" \
		-e "s/@@RELEASE@@/$(RELEASE)/g" \
		< docker-metrics.spec.in > docker-metrics.spec
	rm -f $(NAME)-$(VERSION)
	rm -f cfg.json
	cp cfg.example.json cfg.json
	ln -s . $(NAME)-$(VERSION)
	tar czvf $(NAME)-$(VERSION).tar.gz			\
		--exclude CVS --exclude .git --exclude TAGS		\
		--exclude $(NAME)-$(VERSION)/$(NAME)-$(VERSION)	\
		--exclude $(NAME)-$(VERSION).tar.gz			\
		$(NAME)-$(VERSION)/*
	rm -f $(NAME)-$(VERSION)

rpms: dist
	cp $(NAME)-$(VERSION).tar.gz $(RPMSOURCEDIR)/
	cp $(NAME).spec $(RPMSPECDIR)/
	$(RPMBUILD) -ba $(RPMSPECDIR)/$(NAME).spec

help: # Some kind of magic from https://gist.github.com/rcmachado/af3db315e31383502660
	$(info Available targets)
	@awk '/^[a-zA-Z\-\_0-9]+:/ {                                   \
		nb = sub( /^## /, "", helpMsg );                             \
		if(nb == 0) {                                                \
			helpMsg = $$0;                                             \
			nb = sub( /^[^:]*:.* ## /, "", helpMsg );                  \
		}                                                            \
		if (nb)                                                      \
			printf "\033[1;31m%-" width "s\033[0m %s\n", $$1, helpMsg; \
	}                                                              \
	{ helpMsg = $$0 }'                                             \
	width=$$(grep -o '^[a-zA-Z_0-9]\+:' $(MAKEFILE_LIST) | wc -L)  \
	$(MAKEFILE_LIST)

