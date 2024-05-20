# TODO

Sample list of things to do.

* Change type name Module to Release


* FetchUser()
  Stub out query parameters

* FetchModule()
  * Need some new type definitions to Marshal json
  * Testing
  * Stub out query parameters

* Write test for GetAllUsers (or GetUser() )

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

* Implement exclude_field in query string
    Sample code has exclude_fields=readme,changelog,license,reference,tags

## Documentation

## Known bugs

* In pagination the parameter first, previous, current, next does not honor all
  url.Query parameters
  
* Tests are missing and failing  
