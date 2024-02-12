PROGRAM=goforge
FORGE_CACHE=cache
WGET_ARGS=--no-verbose


all: test build

test:
	@echo 'NOTE: Run "make cache" at least once before running "make test"'
	@echo '      to populate the cache with some puppet modules.'
	go test .

build:
	go build -o $(PROGRAM)

cache_init:
	echo "Creating and populating goforge:s modules cache"
	mkdir -p $(FORGE_CACHE)
	mkdir -p $(FORGE_CACHE)/p/puppetlabs
	mkdir -p $(FORGE_CACHE)/p/pdxcat

cache: cache_init
	wget $(WGET_ARGS) -O $(FORGE_CACHE)/p/puppetlabs/puppetlabs-stdlib-9.4.1.tar.gz \
	    https://forge.puppet.com/v3/files/puppetlabs-stdlib-9.4.1.tar.gz
	wget $(WGET_ARGS) -O $(FORGE_CACHE)/p/puppetlabs/puppetlabs-stdlib-9.4.0.tar.gz \
	    https://forge.puppet.com/v3/files/puppetlabs-stdlib-9.4.0.tar.gz
	wget $(WGET_ARGS) -O $(FORGE_CACHE)/p/puppetlabs/puppetlabs-concat-9.0.2.tar.gz \
	    https://forge.puppet.com/v3/files/puppetlabs-concat-9.0.2.tar.gz
	wget $(WGET_ARGS) -O $(FORGE_CACHE)/p/pdxcat/pdxcat-nrpe-2.1.1.tar.gz \
	    https://forge.puppet.com/v3/files/pdxcat-nrpe-2.1.1.tar.gz

clean:
	rm -f $(PROGRAM)

distclean: clean
	rm -rf $(FORGE_CACHE)

