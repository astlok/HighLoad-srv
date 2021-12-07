package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"highload-srv/httputils"
	"highload-srv/metric"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("port")

	router := mux.NewRouter()

	router.Methods(http.MethodGet).Path("/metrics").Handler(promhttp.Handler())

	router.HandleFunc("/api", apiHandler).Methods(http.MethodGet)

	metric.New()
	defer metric.Destroy()

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	log.Printf("Starting Server at port %v", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	reqID := rand.Uint64()
	time.Sleep(time.Millisecond * time.Duration(rand.Int31n(1000)))

	answer := "get api handler"

	httputils.Respond(w, r, reqID, http.StatusOK, answer)
}