package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mido3ds/C4IAN/src/models"
	"github.com/rs/cors"
	"gopkg.in/antage/eventsource.v1"
)

const (
	M3U8Name = "index.m3u8"
)

type API struct {
	unitsPort   int
	dbManager   *DatabaseManager
	netManager  *NetworkManager
	eventSource eventsource.EventSource
	VideosPath  string
}

func NewAPI() *API {
	es := eventsource.New(nil, func(req *http.Request) [][]byte {
		return [][]byte{
			[]byte("Access-Control-Allow-Origin: *"),
		}
	})
	return &API{eventSource: es}
}

func (api *API) Start(port int, unitsPort int, VideosPath string, dbManager *DatabaseManager, netManager *NetworkManager) {
	// Initialize members
	api.unitsPort = unitsPort
	api.netManager = netManager
	api.dbManager = dbManager
	api.VideosPath = VideosPath

	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/api/audio-msg/{ip}", api.postAudioMsg).Methods(http.MethodPost)
	router.HandleFunc("/api/msg/{ip}", api.postMsg).Methods(http.MethodPost)
	router.HandleFunc("/api/units", api.getUnits).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/units-names", api.getUnitsNames).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/groups", api.getGroups).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/memberships", api.getMemberships).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/audio-msgs/{ip}", api.getAudioMsgs).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/msgs/{ip}", api.getMsgs).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/videos/{ip}", api.getVideos).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/sensors-data/{ip}", api.getSensorsData).Methods(http.MethodGet, http.MethodOptions)
	router.Handle("/events", api.eventSource)

	router.HandleFunc("/api/stream/{ip}/{mId:[0-9]+}", api.StreamHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/stream/{ip}/{mId:[0-9]+}/{segName:index[0-9]+.ts}", api.StreamHandler).Methods(http.MethodGet)

	router.Use(api.jsonContentType)

	// Listen for HTTP requests
	address := fmt.Sprintf(":%d", port)
	allowedOrigin := fmt.Sprintf("http://localhost%s", address)
	fmt.Println(allowedOrigin)
	c := cors.New(cors.Options{
		OptionsPassthrough: false,
		AllowedOrigins:     []string{allowedOrigin},
		AllowCredentials:   false,
	})
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(address, handler))
	// log.Fatal(http.ListenAndServe(address, cors.Default().Handler(router)))
}

func (api *API) StreamHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Printf("Request: %s\n", request.RequestURI)
	vars := mux.Vars(request)
	mId, err := strconv.Atoi(vars["mId"])
	ip, ok := vars["ip"]
	fmt.Printf("mId: %d\n", mId)
	fmt.Printf("ip: %s\n", ip)
	if err != nil || !ok {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("passed")
	segName, ok := vars["segName"]
	fmt.Printf("segName: %s\n", segName)
	if !ok {
		api.serveHlsM3u8(response, request, ip, mId)
	} else {
		fmt.Println(segName)
		api.serveHlsTs(response, request, ip, mId, segName)
	}
}

func (api *API) serveHlsM3u8(w http.ResponseWriter, r *http.Request, ip string, mId int) {
	mediaFile := fmt.Sprintf("%s/%s/%d/%s", api.VideosPath, ip, mId, M3U8Name)
	fmt.Printf("Get file: %s\n", mediaFile)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "application/x-mpegURL")
}

func (api *API) serveHlsTs(w http.ResponseWriter, r *http.Request, ip string, mId int, segName string) {
	mediaFile := fmt.Sprintf("%s/%s/%d/%s", api.VideosPath, ip, mId, segName)
	fmt.Printf("Get file: %s\n", mediaFile)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "video/MP2T")
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
	ip := mux.Vars(r)["ip"]
	file, _, err := r.FormFile("audio")
	if err != nil {
		log.Panic(err)
	}

	var buffer bytes.Buffer
	io.Copy(&buffer, file)

	audioMsg := models.Audio{}
	audioMsg.Body = buffer.Bytes()
	audioMsg.Time = time.Now().Unix()

	go api.dbManager.AddSentAudio(ip, &audioMsg)

	if isMulticastOrBroadcast(ip) {
		go api.netManager.SendUDP(ip, api.unitsPort, audioMsg)
	} else {
		go api.netManager.SendTCP(ip, api.unitsPort, audioMsg)
	}
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

func (api *API) getUnitsNames(w http.ResponseWriter, r *http.Request) {
	unitsNames := api.dbManager.GetUnitsNames()
	json.NewEncoder(w).Encode(unitsNames)
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

func isMulticastOrBroadcast(ip string) bool {
	parsedIP := net.ParseIP(ip).To4()
	isBroadcast := parsedIP[0] == 255 && parsedIP[1] == 255
	return isBroadcast || parsedIP.IsMulticast()
}
