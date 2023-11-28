# goforge

Simple clone of the puppetlabs puppetforge system.

Project selected to get a feeling for golang.

Target usage of this code is to use as a local puppetforge where all modules
are locally cached.



# Usage

## Populate the cache

All puppet modules are stored in a local cache-directory structure as follows

    ${CACHE_ROOT}/{module_hash}/module_owner/{module_owner}-{modulename}-version.tar.gz

The `module_hash` is the first character in the module_owner.

### Example

    ${CACHE_ROOT}/p/puppetlabs/puppetlabs-stdlib-1.0.0.tar.gz
    ${CACHE_ROOT}/p/puppetlabs/puppetlabs-stdlib-1.1.0.tar.gz
    ${CACHE_ROOT}/p/puppetlabs/puppetlabs-concat-1.0.0.tar.gz


# Building

Build the code using make.


# Implementation notes

Specifications loosely based on Puppet Forge v3 API (29) as found in
https://forgeapi.puppet.com/
