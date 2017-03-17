package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/dmmcquay/fotia"
	"github.com/prometheus/client_golang/prometheus"
)

var addr = flag.String("addr", ":8080", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/sleep", fotia.Sleep)
	http.HandleFunc("/down/", fotia.Down)
	http.HandleFunc("/up/", fotia.Up)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
