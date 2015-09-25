package main

import (
	"fmt"
	"os"
)

type outputStdout struct {
}

func NewOutputStdout() Output {
	return &outputStdout{}
}

func (s *outputStdout) WriteLine(l Line) error {
	_, err := fmt.Fprintf(os.Stdout, "%s\t%s\n", l.File, l.Text)
	return err
}
