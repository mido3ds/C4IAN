package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mido3ds/C4IAN/src/models"
	"github.com/rs/cors"
	"gopkg.in/antage/eventsource.v1"
)

type API struct {
	unitsPort   int
	dbManager   *DatabaseManager
	netManager  *NetworkManager
	eventSource eventsource.EventSource
}

func NewAPI() *API {
	es := eventsource.New(nil, func(req *http.Request) [][]byte {
		return [][]byte{
			[]byte("Access-Control-Allow-Origin: *"),
		}
	})
	return &API{eventSource: es}
}

func (api *API) Start(port int, unitsPort int, dbManager *DatabaseManager, netManager *NetworkManager) {
	// Initialize members
	api.unitsPort = unitsPort
	api.netManager = netManager
	api.dbManager = dbManager

	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/api/audio-msg/{ip}", api.postAudioMsg).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/msg/{ip}", api.postMsg).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/units", api.getUnits).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/groups", api.getGroups).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/memberships", api.getMemberships).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/audio-msgs/{ip}", api.getAudioMsgs).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/msgs/{ip}", api.getMsgs).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/videos/{ip}", api.getVideos).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/sensors-data/{ip}", api.getSensorsData).Methods(http.MethodGet, http.MethodOptions)
	router.Handle("/events", api.eventSource)

	router.Use(api.jsonContentType)

	// Listen for HTTP requests
	address := ":" + strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(address, cors.Default().Handler(router)))
}

func (api *API) SendEvent(body models.Event) {
	payload, err := json.Marshal(body)
	if err != nil {
		log.Panic(err)
	}
	api.eventSource.SendEventMessage(string(payload), body.EventType(), "")
}

func (api *API) jsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (api *API) postAudioMsg(w http.ResponseWriter, r *http.Request) {
	audioMsg := models.Audio{}
	ip := mux.Vars(r)["ip"]
	fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&audioMsg)
	if err != nil {
		log.Panic(err)
	}
	audioMsg.Time = time.Now().Unix()
	go api.dbManager.AddSentAudio(ip, &audioMsg)
	go api.netManager.SendTCP(ip, api.unitsPort, audioMsg)
	w.WriteHeader(http.StatusOK)
}

func (api *API) postMsg(w http.ResponseWriter, r *http.Request) {
	msg := models.Message{}
	ip := mux.Vars(r)["ip"]
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Panic(err)
	}
	msg.Time = time.Now().Unix()
	go api.dbManager.AddSentMessage(ip, &msg)
	go api.netManager.SendTCP(ip, api.unitsPort, msg)
	w.WriteHeader(http.StatusOK)
}

func (api *API) getUnits(w http.ResponseWriter, r *http.Request) {
	units := api.dbManager.GetUnits()
	json.NewEncoder(w).Encode(units)
}

func (api *API) getGroups(w http.ResponseWriter, r *http.Request) {
	groups := api.dbManager.GetGroups()
	json.NewEncoder(w).Encode(groups)
}

func (api *API) getMemberships(w http.ResponseWriter, r *http.Request) {
	memberships := api.dbManager.GetMemberships()
	json.NewEncoder(w).Encode(memberships)
}

func (api *API) getAudioMsgs(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	audios := api.dbManager.GetReceivedAudio(ip)
	json.NewEncoder(w).Encode(audios)
}

func (api *API) getMsgs(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	msgs := api.dbManager.GetConversation(ip)
	json.NewEncoder(w).Encode(msgs)
}

func (api *API) getVideos(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	videos := api.dbManager.GetReceivedVideos(ip)
	json.NewEncoder(w).Encode(videos)
}

func (api *API) getSensorsData(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	data := api.dbManager.GetReceivedSensorsData(ip)
	json.NewEncoder(w).Encode(data)
}
