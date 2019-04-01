package main

import (
	"flag"
	"github.com/rbxb/fileserve"
	"log"
	"net/http"
)

var port string          // the port the fileserver listens at
var fileDirectory string // the directory serving files

func main() {
	flag.Parse()
	server := fileserve.NewServer(fileDirectory, nil)
	log.Fatal(http.ListenAndServe(port, server))
}

func init() {
	flag.StringVar(&port, "port", ":8080", "The address and port the fileserver listens at. (:8080)")
	flag.StringVar(&fileDirectory, "directory", "./root", "The directory serving files. (./root)")
}
