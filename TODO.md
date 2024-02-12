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

* Found some problem with Marshal of metadata. As the type is defined as string
  the conversion will add \"somedata\" to the output. Need to test a lot of
  json.Unmarshal for metadata.json file. If that works out change the type of metadata
  in the puppetmodule struct.
  NOTE: In the specification it states that metadata is "Verbatim contents of
        release's metadata.json file".

* Try to use json.Marshal() istead of toJson()

## Documentation

* Update README.md with Q and A
* Update with build instructions
* Update with usage instructions

## Known bugs

* In pagination the parameter first, previous, current, next does not honor all
  url.Query parameters

* Inital extraction of metadata.json from the module.tar.gz file is implemented.
  Needs more debugging and better tests.
