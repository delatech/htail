// +build !linux,!darwin

package main

var DefaultPaths = []string{}

func openBrowser(url string) error {
	return nil
}
