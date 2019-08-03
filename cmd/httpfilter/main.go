package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/rbxb/httpfilter"
)

var port string
var root string

func init() {
	flag.StringVar(&port, "port", ":8080", "The address and port the fileserver listens at. (:8080)")
	flag.StringVar(&root, "root", "./root", "The directory serving files. (./root)")
}

func main() {
	flag.Parse()
	sv := httpfilter.NewServer(root)
	log.Fatal(http.ListenAndServe(port, sv))
}