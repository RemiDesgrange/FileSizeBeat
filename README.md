FileSizeBeat
============

*You know, for checking file, folder size and content*
Current status: **beta release**.

pull-request are welcome

## Requirements

filesizebeat require beat from [here](https://github.com/elastic/beats/)

## Install

You need [Go](https://golang.org/doc/install).

Fetch this package via `git clone https://github.com/remidesgrange/filesizebeat`. The
resulting binary will be in `$GOPATH/bin`

## Usage

There is a (configuration file)[etc/filesizebeat-example.yml] in `etc` folder.
On Windows you can put the path like that :
`C:/Program Files`

## Note on time mgmt

You can potentially scan *any* folder, if you have the right to of course. But
be careful, If you put a period of 2 sec and try to scan `C:`or `/` this will
*not* work. Unsuccessful scan are logged.



## TODO
* need test
* use fnotify ? maybe it can be an option ?
* use different checking time for each file describe in `path`
* elasticsearch template
