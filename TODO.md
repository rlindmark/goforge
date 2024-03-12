# TODO

Sample list of things to do.

* Stub out all methods

* Try to use pdk to generate testing modules
  From the puppet documentation on ubuntu 22.04

    wget <https://apt.puppet.com/puppet-tools-release-jammy.deb>
    sudo dpkg -i puppet-tools-release-jammy.deb
    sudo apt-get update
    sudo apt-get install pdk

* Investigate if get_v3_releases_module_result() is the same as
  "Fetch Module" from the API specification.

* ListModuleReleases() see if code for default values can be written in a
  better way.

* Maybe rearrange project into main, api, cache for better code separation.

* Create /api and move the fetch code into that directory

## Documentation

* Update README.md with Q and A
* Update with build instructions
* Update with usage instructions

## Known bugs

* In pagination the parameter first, previous, current, next does not honor all
  url.Query parameters

* Inital extraction of metadata.json from the module.tar.gz file is implemented.
  Needs more debugging and better tests.
