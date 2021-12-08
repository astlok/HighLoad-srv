package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"highload-srv/httputils"
	"highload-srv/metric"
	"log"
	"math/rand"
	"net/http"
)

func main() {
	port := ":8080"

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

	num := uint64(1)

	for i := uint64(1); i < 10000; i++ {
		for j := 1; j < 1000; j++ {
			num += i
		}
	}

	answer := fmt.Sprintf("get api handler %d", num)

	httputils.Respond(w, r, reqID, http.StatusOK, answer)
}
