# goforge

Simple clone of the puppetlabs puppetforge system.

Project selected to get a feeling for golang.

Target usage of this code is for use as a local puppetforge where all modules
are locally cached.

## Usage

### Populate the cache

All puppet modules are stored in a local cache-directory structure as follows

    ${CACHE_ROOT}/{module_hash}/module_owner/{module_owner}-{modulename}-version.tar.gz

The `module_hash` is the first character of the module_owner string.

#### Example

    ${CACHE_ROOT}/p/puppetlabs/puppetlabs-stdlib-1.0.0.tar.gz
    ${CACHE_ROOT}/p/puppetlabs/puppetlabs-stdlib-1.1.0.tar.gz
    ${CACHE_ROOT}/p/puppetlabs/puppetlabs-concat-1.0.0.tar.gz

## Building

Build the code using make.

## Change default puppet settings

Change the default puppet config settings to point to goforge:s
address.

    puppet config print module_repository    # = https://forgeapi.puppet.com
    puppet config set module_repository http://localhost:8080

Start goforge and try to download and install some module

    puppet module install puppetlabs-stdlib

## Implementation notes

Specifications loosely based on Puppet Forge v3 API (29) as found at
<https://forgeapi.puppet.com/>

Code does not implement any POST Puppet Forge API:s as the cache is managed outside of this program.

Currently implements the following api endpoints.

* </v3/files/{filename}>
