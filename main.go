package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/yalue/merged_fs"
)

func main() {
	var infc string
	var port int
	flag.StringVar(&infc, "i", "0.0.0.0", "interface to listen on (default: 0.0.0.0)")
	flag.IntVar(&port, "p", 8080, "port to listen on (default: 8080)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("No valid directories to serve")
	}

	var servedFS fs.FS
	servedFS = os.DirFS(args[0])
	rest := args[1:]
	for _, f := range rest {
		servedFS = merged_fs.NewMergedFS(servedFS, os.DirFS(f))
	}

	tmpl, err := template.New("index.html").Parse(defaultIndex)
	if err != nil {
		log.Fatal(err)
	}

	httpFS := autoIndexedFS{fs: http.FS(servedFS), template: tmpl}

	lstn := fmt.Sprintf("%s:%d", infc, port)
	fmt.Println("Listening on:", lstn)
	log.Fatal(http.ListenAndServe(lstn, AccessLogger(http.FileServer(httpFS))))
}
