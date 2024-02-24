# TODO

Sample list of things to do.

## Functions to implement

* Download module release
* Fetch module release
* List modules
* Fetch module

* Try to use pdk to generate testing modules
  From the puppet documentation on ubuntu 22.04

    wget <https://apt.puppet.com/puppet-tools-release-jammy.deb>
    sudo dpkg -i puppet-tools-release-jammy.deb
    sudo apt-get update
    sudo apt-get install pdk

* Try to use json.Marshal() instead of toJson()

* Remove all asJSON() functions when done

* Investigate if get_v3_releases_module_result() is the same as
  "Fetch Module" from the API specification.

## Documentation

* Update README.md with Q and A
* Update with build instructions
* Update with usage instructions

## Known bugs

* In pagination the parameter first, previous, current, next does not honor all
  url.Query parameters

* Inital extraction of metadata.json from the module.tar.gz file is implemented.
  Needs more debugging and better tests.
