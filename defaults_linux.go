package main

import "os/exec"

var DefaultPaths = []string{
	"/var/log/*.log",
	"/var/log/*/*.log",
}

func openBrowser(url string) error {
	return exec.Command("xdg-open", url).Run()
}
