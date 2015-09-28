package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"

	"github.com/ActiveState/tail"
)

type Line struct {
	File string `json:"file"`
	Text string `json:"text"`
	Time int64  `json:"time"`
}

type Tailer struct {
	files   map[string]*tail.Tail
	outputs []Output
	lines   chan Line
}

type Output interface {
	WriteLine(Line) error
}

func NewTailer() *Tailer {
	return &Tailer{
		files:   make(map[string]*tail.Tail),
		outputs: make([]Output, 0),
		lines:   make(chan Line),
	}
}

func (t *Tailer) AddFile(filename string) error {
	if _, ok := t.files[filename]; ok {
		return nil
	}

	config := tail.Config{Location: &tail.SeekInfo{Whence: os.SEEK_END}, Follow: true}
	tf, err := tail.TailFile(filename, config)
	if err != nil {
		return err
	}
	t.files[filename] = tf
	return nil
}

func (t *Tailer) AddReader(name string, r io.Reader) {
	go func() {
		br := bufio.NewReader(r)
		for {
			l, err := br.ReadString('\n')
			if l == "" && err == io.EOF {
				return
			}

			line := Line{
				File: name,
				Text: l,
				Time: time.Now().Unix(),
			}

			t.lines <- line

			if err != nil {
				if err != io.EOF {
					log.Printf("Error while reading from %s: %s\n", name, err)
				}
				return
			}
		}
	}()
}

func (t *Tailer) AddOutput(out Output) {
	t.outputs = append(t.outputs, out)
}

func (t *Tailer) readTail(file string, ft *tail.Tail) {
	for line := range ft.Lines {
		l := Line{
			File: file,
			Text: line.Text,
			Time: line.Time.Unix(),
		}

		t.lines <- l
	}
}

func (t *Tailer) Run() {
	for i := range t.files {
		go t.readTail(i, t.files[i])
	}
	go t.output()
}

func (t *Tailer) output() {
	for l := range t.lines {
		for i := range t.outputs {
			t.outputs[i].WriteLine(l)
		}
	}
}

func (t *Tailer) Close() {
	for i := range t.files {
		t.files[i].Tomb.Kill(nil)
	}
	close(t.lines)
}
