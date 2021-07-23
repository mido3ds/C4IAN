package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type API struct {
	dbManager *DatabaseManager
}

func NewAPI(dbManager *DatabaseManager) *API {
	return &API{dbManager: dbManager}
}

func (api *API) Start(port int) {
	// Initialize router
	router := mux.NewRouter()

	// API endpoints
	router.HandleFunc("/api/packets", api.getPackets).Methods(http.MethodGet)
	router.HandleFunc("/api/packet-data/{hash:.*}", api.getPacketData).Methods(http.MethodGet)

	router.Use(api.jsonContentType)

	// Listen for HTTP requests
	c := cors.New(cors.Options{
		OptionsPassthrough: false,
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
	})
	handler := c.Handler(router)
	address := ":" + strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(address, handler))
}

func (api *API) jsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (api *API) getPackets(w http.ResponseWriter, r *http.Request) {
	packets := api.dbManager.GetPackets()
	json.NewEncoder(w).Encode(packets)
}

func (api *API) getPacketData(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	data := api.dbManager.GetPacketData(hash)
	if data == nil {
		w.WriteHeader(404)
	} else {
		json.NewEncoder(w).Encode(data)
	}
}
