package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/websocket"
)

// Config holds the program configuration
type Config struct {
	Paths    []string
	Host     string
	NoStdout bool
	NoOpen   bool
	NoHTTP   bool
	Verbose  bool
}

var cfg = &Config{Paths: DefaultPaths}

func init() {
	log.SetOutput(os.Stderr)
	flag.StringVar(&cfg.Host, "host", "localhost:48245", "Address on which the htail will listen for connections")
	flag.BoolVar(&cfg.NoStdout, "no-stdout", false, "Disable output to stdout")
	flag.BoolVar(&cfg.NoHTTP, "no-http", false, "Disable output to http")
	flag.BoolVar(&cfg.Verbose, "v", false, "Verbose output")
}

func scanDir(dir string) []string {
	files := make([]string, 0)
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("Error while reading directory %#v: %s\n", dir, err)
		return files
	}

	for _, fi := range fis {
		if fi.IsDir() {
			f := scanDir(path.Join(dir, fi.Name()))
			files = append(files, f...)
			continue
		}
		files = append(files, path.Join(dir, fi.Name()))
	}
	return files
}

func ScanPaths(paths []string) []string {
	files := make([]string, 0)

	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			log.Printf("Pattern %#v is not valid\n", path)
			continue
		}
		for i := range matches {
			p := matches[i]
			fi, err := os.Stat(p)
			if err != nil && !os.IsNotExist(err) {
				logrus.Errorf("Error while stating %#v: %s", p, err)
				continue
			}

			if !fi.IsDir() {
				files = append(files, p)
				continue
			}

			f := scanDir(p)
			files = append(files, f...)
		}
	}
	return files
}

func main() {
	flag.Parse()
	paths := flag.Args()
	// if there's nothing on the argument line, we check for the environment variable
	if len(paths) == 0 {
		paths = strings.Split(os.Getenv("HTAIL_PATH"), ":")

		// if there's nothing in the env var, we switch to the default
		if len(paths) == 1 && paths[0] == "" {
			paths = DefaultPaths
		}
	}

	cfg.Paths = paths
	files := ScanPaths(cfg.Paths)

	tailer := NewTailer()
	tailer.AddReader("stdin", os.Stdin)
	for _, f := range files {
		if cfg.Verbose {
			log.Printf("Tailing %s\n", f)
		}
		if err := tailer.AddFile(f); err != nil {
			log.Println(err)
			continue
		}
	}

	if !cfg.NoStdout {
		tailer.AddOutput(NewOutputStdout())
	}

	if !cfg.NoHTTP {
		hub := newWebsocketHub()
		tailer.AddOutput(hub)

		http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
			tmpl, err := template.New("index").Parse(templateIndex)
			if err != nil {
				log.Panicln(err)
			}

			rw.Header().Set("Content-Type", "text/html; charset=UTF-8")
			tmpl.Execute(rw, cfg)
		})

		http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
			hub.addConn(ws)
			select {}
		}))

		go func() {
			if cfg.Verbose {
				log.Println("Listening to", cfg.Host)
			}
			go func() {
				err := http.ListenAndServe(cfg.Host, nil)
				if err != nil {
					log.Fatalf("Error starting HTTP server: %s\n", err)
				}
			}()
			if !cfg.NoOpen {
				if err := openBrowser(fmt.Sprintf("http://%s", cfg.Host)); err != nil {
					log.Printf("Error while trying to open browser: %s\n", err)
				}
			}
		}()
	}

	tailer.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	_ = <-c
	tailer.Close()
}
