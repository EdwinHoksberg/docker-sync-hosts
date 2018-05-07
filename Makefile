BUILDPATH=$(CURDIR)
MANPATH="/usr/local/share/man/man1"
GO=$(shell which go)

# Compile time values
PROGRAM=docker-sync-hosts
VERSION=$(shell printf "%s [%s]" `git describe --abbrev=0 --tags 2> /dev/null || echo 'from_source'` `git rev-parse --verify HEAD 2> /dev/null || echo 'no_commit'`)

# Interpolate the variable values using go link flags
LDFLAGS=-ldflags "-s -w -X 'main.Name=${PROGRAM}' -X 'main.Version=${VERSION}'"

build:
	@if [ ! -d $(BUILDPATH)/bin ] ; then mkdir -p $(BUILDPATH)/bin ; fi
	$(GO) build ${LDFLAGS} -o $(BUILDPATH)/bin/$(PROGRAM)

install:
	@if [ ! -d $(MANPATH) ] ; then mkdir -p $(MANPATH) ; fi

	cp $(BUILDPATH)/bin/$(PROGRAM) /usr/sbin/$(PROGRAM)
	cp $(PROGRAM).service /lib/systemd/system/$(PROGRAM).service
	gzip -c $(PROGRAM).man | tee $(MANPATH)/$(PROGRAM).1.gz > /dev/null

clean:
	rm -f $(BUILDPATH)/bin/$(PROGRAM)
