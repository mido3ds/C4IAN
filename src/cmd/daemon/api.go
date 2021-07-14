package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mido3ds/C4IAN/src/models"
	"gopkg.in/antage/eventsource.v1"
)

type API struct {
	dbManager   *DatabaseManager
	eventSource eventsource.EventSource
}

func NewAPI(dbManager *DatabaseManager) *API {
	es := eventsource.New(nil, func(req *http.Request) [][]byte {
		return [][]byte{
			[]byte("Access-Control-Allow-Origin: *"),
		}
	})
	return &API{dbManager: dbManager, eventSource: es}
}

func (api *API) Start(port int) {
	router := mux.NewRouter()
	router.HandleFunc("/api/audio-msg/{ip}", api.postAudioMsg).Methods(http.MethodPost)
	router.HandleFunc("/api/msg/{ip}", api.postMsg).Methods(http.MethodPost)
	router.HandleFunc("/api/audio-msgs/{ip}", api.getAudioMsgs).Methods(http.MethodGet)
	router.HandleFunc("/api/msgs/{ip}", api.getMsgs).Methods(http.MethodGet)
	router.HandleFunc("/api/videos/{ip}", api.getVideos).Methods(http.MethodGet)
	router.HandleFunc("/api/sensors-data/{ip}", api.getSensorsData).Methods(http.MethodGet)
	router.Handle("/events", api.eventSource)
	address := ":" + strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(address, router))
}

func (api *API) SendEvent(body models.Event) {
	payload, err := json.Marshal(body)
	if err != nil {
		log.Fatalln(err)
	}
	api.eventSource.SendEventMessage(string(payload), body.EventType(), "")
}

func (api *API) postAudioMsg(w http.ResponseWriter, r *http.Request) {
	audioMsg := models.Audio{}
	ip := mux.Vars(r)["ip"]
	json.NewDecoder(r.Body).Decode(&audioMsg)
	// TODO: Add to database and send to unit(s)
	log.Println("Sending audio message: ", audioMsg.Body, ", to: ", ip)
}

func (api *API) postMsg(w http.ResponseWriter, r *http.Request) {
	msg := models.Message{}
	ip := mux.Vars(r)["ip"]
	json.NewDecoder(r.Body).Decode(&msg)
	// TODO: Add to database and send to unit(s)
	log.Println("Sending message: ", msg.Code, ", to: ", ip)
}

func (api *API) getAudioMsgs(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	log.Println("Fetching audio message history for: ", ip)
	// TODO: Fetch audio messages
}

func (api *API) getMsgs(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	log.Println("Fetching message history for: ", ip)
	// TODO: Fetch messages
}

func (api *API) getVideos(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	log.Println("Fetching video history for: ", ip)
	// TODO: Fetch videos
}

func (api *API) getSensorsData(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	log.Println("Fetching sensors data history for: ", ip)
	// TODO: Fetch sensors
}
