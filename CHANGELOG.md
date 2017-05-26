# Changelog

## v2.6 May 26, 2017

* Escape shell output (@erichs)

## v2.5 April 10, 2017

* Add support for AWS, Azure, Digital Ocean, Google, and OpenStack cloud metadata (@erichs)
* Munge non-portable characters for shell output (@erichs)

## v2.4 April 2, 2017

* Prevent HTTP server from exiting when a 500 error is returned (@erichs)

## v2.3 April 1, 2017

* Support specifying format in HTTP parameters (@erichs)

## v2.2 March 30, 2017

* Fixed issue with HTTP path recognition

## v2.1 December 19, 2016

* Added support for shell output

## v2.0 December 16, 2016

* Switched to govendor for versioning
* All elements/facts are now obtained solely by `github.com/shirou/gopsutil`
* Changed `System` and `External` back to `system` and `external`

## v0.2 June 4, 2016

* Godep
* Changed `system` and `external` to `System` and `External`.

## v0.1 April 30, 2016

Initial release
