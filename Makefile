PROGRAM=goforge
FORGE_CACHE=cache

all: test build

test:
	go test .

build:
	go build -o $(PROGRAM)

cache:
	echo "Creating and populating goforge:s cache"
	@mkdir -p $(FORGE_CACHE)
	@mkdir -p cache/p/puppetlabs
	@curl -o $(FORGE_CACHE)/p/puppetlabs/puppetlabs-stdlib-9.4.1.tar.gz \
	    https://forge.puppet.com/v3/files/puppetlabs-stdlib-9.4.1.tar.gz
	@curl -o $(FORGE_CACHE)/p/puppetlabs/puppetlabs-stdlib-9.4.0.tar.gz \
	    https://forge.puppet.com/v3/files/puppetlabs-stdlib-9.4.0.tar.gz

clean:
	rm -f $(PROGRAM)

distclean: clean
	rm -r $(FORGE_CACHE)

