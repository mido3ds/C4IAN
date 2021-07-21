package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mido3ds/C4IAN/src/unit/halapi"
	"github.com/rs/cors"
	"gopkg.in/antage/eventsource.v1"
)

const UIPort = 3006

type API struct {
	context     *Context
	eventSource eventsource.EventSource
}

func newAPI(context *Context) *API {
	es := eventsource.New(nil, func(req *http.Request) [][]byte {
		return [][]byte{
			[]byte(fmt.Sprintf("Access-Control-Allow-Origin: http://localhost:%v", UIPort)),
		}
	})
	return &API{eventSource: es, context: context}
}

func (api *API) start(port int) {
	// Initialize router
	router := mux.NewRouter()

	// API endpoints
	router.HandleFunc("/api/audio-msg", api.postAudioMsg).Methods(http.MethodPost)
	router.HandleFunc("/api/code-msg", api.postMsg).Methods(http.MethodPost)
	router.HandleFunc("/api/sensors-data", api.postSensorsData).Methods(http.MethodPost)

	// SSE endpoint
	router.Handle("/events", api.eventSource)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/x-msgpack")
			next.ServeHTTP(w, r)
		})
	})

	// Listen for HTTP requests
	c := cors.New(cors.Options{
		OptionsPassthrough: false,
		AllowedOrigins:     []string{"http://localhost:*"},
		AllowCredentials:   true,
	})
	handler := c.Handler(router)
	address := ":" + strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(address, handler))
}

func (api *API) sendCodeMsgEvent(code int) {
	api.eventSource.SendEventMessage(strconv.Itoa(code), "CODE-EVENT", "")
}

func (api *API) sendAudioMsgEvent(audio []byte) {
	payload, err := json.Marshal(audio)
	if err != nil {
		log.Panic(err)
	}

	api.eventSource.SendEventMessage(string(payload), "AUDIO-EVENT", "")
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
	fmt.Println(msg)
	if err != nil {
		log.Panic(err)
	}
	api.context.onCodeMsgReceivedFromHAL(&msg)
	w.WriteHeader(http.StatusOK)
}

func readSensorsData(body io.ReadCloser) halapi.SensorData {
	var vp halapi.VideoFragment
	var s halapi.SensorData
	var a halapi.AudioMsg
	var msg halapi.CodeMsg
	recvdType, err := halapi.ReadFromHAL(body, &vp, &s, &a, &msg)
	if err != nil {
		log.Panic(err)
	}
	if recvdType != halapi.SensorDataType {
		log.Panic("invalid type, expected SensorDataType")
	}
	return s
}

func (api *API) postSensorsData(w http.ResponseWriter, r *http.Request) {
	sensorsdata := readSensorsData(r.Body)
	api.context.onSensorDataReceivedFromHAL(&sensorsdata)
	w.WriteHeader(http.StatusOK)
}
