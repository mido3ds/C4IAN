package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mido3ds/C4IAN/src/models"
	"github.com/mido3ds/C4IAN/src/unit/halapi"
	"github.com/rs/cors"
	"gopkg.in/antage/eventsource.v1"
)

type API struct {
	context     *Context
	eventSource eventsource.EventSource
}

func newAPI(context *Context) *API {
	es := eventsource.New(nil, func(req *http.Request) [][]byte {
		return [][]byte{
			[]byte(fmt.Sprintf("Access-Control-Allow-Origin: *")),
		}
	})
	return &API{eventSource: es, context: context}
}

func (api *API) start(socket string) {
	// Initialize router
	router := mux.NewRouter()

	// API endpoints
	router.HandleFunc("/api/name", api.getName).Methods(http.MethodGet)
	router.HandleFunc("/api/audio-msg", api.postAudioMsg).Methods(http.MethodPost)
	router.HandleFunc("/api/code-msg", api.postMsg).Methods(http.MethodPost)

	// SSE endpoint
	router.Handle("/events", api.eventSource)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// Use CORS handler with mux router
	c := cors.New(cors.Options{
		OptionsPassthrough: false,
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
	})
	handler := c.Handler(router)

	// Open unix socket
	if err := os.RemoveAll(socket); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("unix", socket)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer listener.Close()

	log.Println("API listening on: ", socket)
	// Serve HTTP requests over unix socket
	log.Fatal(http.Serve(listener, handler))
}

func (api *API) sendCodeMsgEvent(code int) {
	api.eventSource.SendEventMessage(strconv.Itoa(code), "CODE-EVENT", "")
}

func (api *API) sendAudioMsgEvent(body *models.Audio) {
	payload, err := json.Marshal(body)
	if err != nil {
		log.Panic(err)
	}
	api.eventSource.SendEventMessage(string(payload), "AUDIO-EVENT", "")
}

func (api *API) getName(w http.ResponseWriter, r *http.Request) {
	nameMap := make(map[string]string)
	nameMap["name"] = api.context.name
	json.NewEncoder(w).Encode(nameMap)
}

func (api *API) postAudioMsg(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("audio")
	if err != nil {
		log.Panic(err)
	}

	var buffer bytes.Buffer
	io.Copy(&buffer, file)

	api.context.onAudioMsgReceivedFromHAL(&halapi.AudioMsg{Audio: buffer.Bytes()})

	w.WriteHeader(http.StatusOK)
}

func (api *API) postMsg(w http.ResponseWriter, r *http.Request) {
	msg := halapi.CodeMsg{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Panic(err)
	}
	api.context.onCodeMsgReceivedFromHAL(&msg)
	w.WriteHeader(http.StatusOK)
}
