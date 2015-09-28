package main

import (
	"flag"
	"os/exec"
)

var DefaultPaths = []string{
	"/usr/local/var/log",  // homebrew log locations
	"/var/log/system.log", // default system log
}

func openBrowser(url string) error {
	return exec.Command("open", url).Run()
}

func init() {
	flag.BoolVar(&cfg.NoOpen, "no-open", false, "Disable automatic browser open")
}
