package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	flag "github.com/spf13/pflag"
)

func main() {
	var (
		listenAddress string
		metricsPath string
	)
	flag.StringVar(&listenAddress, "listenAddress", ":8778", "Address on which to expose metrics")
	flag.StringVar(&metricsPath, "metricsPath", "/metrics", "Path under which to expose metrics")
	flag.Parse()
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc(metricsPath, func(w http.ResponseWriter, r *http.Request) {
		HandleMetrics(w, r)
	})
	corsHandler := cors.AllowAll().Handler(router)
	ctxHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsHandler.ServeHTTP(w, r)
	})
	http.ListenAndServe(listenAddress, ctxHandler)
	return
}

