PROGRAM=goforge
FORGE_CACHE=cache
WGET_ARGS=--no-verbose


all: test build

test:
	go test .

build:
	go build -o $(PROGRAM)

cache_init:
	echo "Creating and populating goforge:s modules cache"
	mkdir -p $(FORGE_CACHE)
	mkdir -p cache/p/puppetlabs

cache: cache_init
	wget $(WGET_ARGS) -O $(FORGE_CACHE)/p/puppetlabs/puppetlabs-stdlib-9.4.1.tar.gz \
	    https://forge.puppet.com/v3/files/puppetlabs-stdlib-9.4.1.tar.gz
	wget $(WGET_ARGS) -O $(FORGE_CACHE)/p/puppetlabs/puppetlabs-stdlib-9.4.0.tar.gz \
	    https://forge.puppet.com/v3/files/puppetlabs-stdlib-9.4.0.tar.gz

clean:
	rm -f $(PROGRAM)

distclean: clean
	rm -rf $(FORGE_CACHE)

