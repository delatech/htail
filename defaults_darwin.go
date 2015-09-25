package main

import "flag"

var DefaultPaths = []string{
	"/usr/local/var/log",  // homebrew log locations
	"/var/log/system.log", // default system log
}

func init() {
	flag.BoolVar(&cfg.NoOpen, "no-open", false, "Disable automatic browser open")
}
