# goforge

Simple clone of the puppetlabs puppetforge system.

Project selected to get a feeling for the golang programming language.

Target usage is as a local puppetforge where all modules are locally cached.

## Usage

### Populate the cache

All puppet modules are stored in a local cache-directory structure as follows

    ${FORGE_CACHE}/{module_hash}/module_owner/{module_owner}-{modulename}-{version}.tar.gz

The `module_hash` is the first character of the module_owner string.

#### Example

    ${CACHE_ROOT}/p/puppetlabs/puppetlabs-stdlib-1.0.0.tar.gz
    ${CACHE_ROOT}/p/puppetlabs/puppetlabs-stdlib-1.1.0.tar.gz
    ${CACHE_ROOT}/p/puppetlabs/puppetlabs-concat-1.0.0.tar.gz

The script `scripts/forgecli.sh` may be used to simplify downloading modules from
<https://forgeapi.puppetlabs.com/> to populate the local cache.

## Building

Build the code using make.

    make cache  # to create a local cache and populate it with some forge modules
    make
    
## Change default puppet settings

Change the default puppet config settings to point to goforge:s
address.

    puppet config print module_repository    # = https://forgeapi.puppet.com
    puppet config set module_repository http://localhost:8080

Start goforge and try to download and install some module

    puppet module install puppetlabs-stdlib

## Running the application

Populate the local cache as described above.
Then select what ip and port the application should bind to using environment variables
FORGE_IP (default localhost), FORGE_PORT (default 8080), FORGE_CACHE (default cache).

Then start the application

    ./goforge

The application logs to stdout.

## Implementation notes

Implementation based on the specification of Puppet Forge v3 API (29), as found at
<https://forgeapi.puppet.com/>.

Code only implements a fraction of Puppet Forge API:s.

No state except the data that can be extracted from the locally cached
<owner-module-version.tar.gz> files exist.

## Q and A

Q: Why create a forgeapi clone?

A: In environments without direct access to the internet, some puppet-tools cant be used.
   This application tries to resolve some of these problems. This code should make it
   possible to use "puppet module install" and "r10k" work.

Q: Why golang?

A: This project have been used to begin learning golang. Golang also seem to have some
   native support for transportable static binaries.

Q: Why is TLS missing?

A: The idea is to have this application behind some web frontend that is responsible
   for the TLS stuff.
