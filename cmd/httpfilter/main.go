package main

import (
	"crypto/tls"
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
	flag.BoolVar(&ssl, "ssl", false, "Use SSL? Using SSL will also listen on port 80 to redirect HTTP traffic to HTTPS.")
	flag.StringVar(&cert, "cert", "./cert.pem", "The path of the SSL certificate.")
	flag.StringVar(&key, "key", "./key.pem", "The path of the private key.")
}

func main() {
	flag.Parse()
	fs := httpfilter.NewServer(root, filter)
	server := http.Server{
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)), //disable HTTP/2
		Addr:         port,
		Handler:      fs,
	}
	if ssl {
		go func() {
			log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)))
		}()
		log.Fatal(server.ListenAndServeTLS(cert, key))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}
