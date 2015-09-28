package main

import "os/exec"

var DefaultPaths = []string{
	"/var/log/*.log",
	"/var/log/*/*.log",
}

func openBrowser(url string) error {
	open, err := exec.LookPath("xdg-open")
	if err != nil {
		open, err = exec.LookPath("x-www-browser")
		if err != nil {
			return nil
		}
	}

	return exec.Command(open, url).Run()
}
