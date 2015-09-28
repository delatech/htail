# htail

An utility that displays your common logs, to stdout and in the browser

By firing up one command, htail gives you in your terminal and in your browser:

- Your applications logs
- Your webserver access/error logs
- Your database logs
- Your system logs

![how-it-works](https://i.gyazo.com/439e1ce9b4156661a52b2fe869418209.gif)

## Why

I was tired of searching for all my log files when debugging issues that could
span multiple applications. I wanted to fire up one command and be able to read
all log files that were affected by an action. A few hours of coding later, I
came up to that. Other solutions exists, but I wanted it to be platform
independent and be easily deployable.

## Usage

`htail` is configurable via the command line.

- `htail`: Run htails with the default options
- `htail /my/log1 /my/log2`: tail the two provided files
- `htail -h`: Display help

`htail` also parses the `HTAIL_PATH` environment variable for directories or
log files to parses. The format is the same as the `PATH` environment variable:

    export HTAIL_PATH="/var/log:/usr/local/var/log:/my/path/to/projects/*/logs/dev.log"
    htail
    # Similar to
    htail /var/log /usr/local/var/log '/my/path/to/projects/*/logs/dev.log'

## Installation

Provided that you have Go installed, run:

- `go get github.com/delatech/htail`

## Credits

Created by Marc Weistroff for [DeLaTech](http://www.delatech.net)
