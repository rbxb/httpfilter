package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/rbxb/httpfilter"
	"github.com/rbxb/httpfilter/request"
	"github.com/rbxb/httpfilter/reverseproxy"
)

var port string
var root string

func init() {
	flag.StringVar(&port, "port", ":8080", "The address and port the fileserver listens at.")
	flag.StringVar(&root, "root", "./root", "The directory serving files.")
}

func main() {
	flag.Parse()
	sv := httpfilter.NewServer(root, request.Ops, reverseproxy.Ops)
	log.Fatal(http.ListenAndServe(port, sv))
}