package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/rbxb/httpfilter"
)

var port string
var root string
var filter string
var ssl bool
var cert string
var key string

func init() {
	flag.StringVar(&port, "port", ":8080", "The address and port the fileserver listens at.")
	flag.StringVar(&root, "root", "./root", "The directory serving files.")
	flag.StringVar(&filter, "filter", "", "The filter file to use (optional).")
	flag.BoolVar(&ssl, "ssl", false, "Use SSL?")
	flag.StringVar(&cert, "cert", "./cert.pem", "The path of the SSL certificate.")
	flag.StringVar(&key, "key", "./key.pem", "The path of the private key.")
}

func main() {
	flag.Parse()
	sv := httpfilter.NewServer(root, filter)
	if ssl {
		log.Fatal(http.ListenAndServeTLS(port, cert, key, sv))
	} else {
		log.Fatal(http.ListenAndServe(port, sv))
	}
}
