package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mido3ds/C4IAN/src/models"
)

func serveRequests(port int) {
	router := mux.NewRouter()
	router.HandleFunc("/api/audio-msg/{ip}", postAudioMsg).Methods(http.MethodPost)
	router.HandleFunc("/api/msg/{ip}", postMsg).Methods(http.MethodPost)
	router.HandleFunc("/api/audio-msgs/{ip}", getAudioMsgs).Methods(http.MethodGet)
	router.HandleFunc("/api/msgs/{ip}", getMsgs).Methods(http.MethodGet)
	router.HandleFunc("/api/videos/{ip}", getVideos).Methods(http.MethodGet)
	router.HandleFunc("/api/sensors-data/{ip}", getSensorsData).Methods(http.MethodGet)

	address := ":" + strconv.Itoa(port)
	http.ListenAndServe(address, router)
}

func postAudioMsg(w http.ResponseWriter, r *http.Request) {
	audioMsg := models.Audio{}
	ip := mux.Vars(r)["ip"]
	json.NewDecoder(r.Body).Decode(&audioMsg)
	// TODO: Add to database and send to unit(s)
	log.Println("Sending audio message: ", audioMsg.Body, ", to: ", ip)
}

func postMsg(w http.ResponseWriter, r *http.Request) {
	msg := models.Message{}
	ip := mux.Vars(r)["ip"]
	json.NewDecoder(r.Body).Decode(&msg)
	// TODO: Add to database and send to unit(s)
	log.Println("Sending message: ", msg.Code, ", to: ", ip)
}

func getAudioMsgs(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	log.Println("Fetching audio message history for: ", ip)
	// TODO: Fetch audio messages
}

func getMsgs(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	log.Println("Fetching message history for: ", ip)
	// TODO: Fetch messages
}

func getVideos(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	log.Println("Fetching video history for: ", ip)
	// TODO: Fetch videos
}

func getSensorsData(w http.ResponseWriter, r *http.Request) {
	ip := mux.Vars(r)["ip"]
	log.Println("Fetching sensors data history for: ", ip)
	// TODO: Fetch sensors
}
