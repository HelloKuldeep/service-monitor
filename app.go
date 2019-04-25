package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	. "service-monitor/config"
	. "service-monitor/dao"
	. "service-monitor/models"
)

var config = Config{}
var dao = StreamsDAO{}

// GET list of streams
func AllStreamsEndPoint(w http.ResponseWriter, r *http.Request) {
	streams, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, streams)
}

// POST a new stream
func CreateStreamEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var stream Stream
	// dec := json.NewDecoder(r.Body)
	// var v map[Stream]interface{}
	if err := json.NewDecoder(r.Body).Decode(&stream); err != nil {
	// if err := dec.Decode(&v); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// for val := range v {
	// 	val.ID = bson.NewObjectId()
	// 	if err := dao.Insert(val); err != nil {
	// 		respondWithError(w, http.StatusInternalServerError, err.Error())
	// 		return
	// 	}
	// }
	stream.ID = bson.NewObjectId()
	if err := dao.Insert(stream); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, stream)
}

// GET Required Output
func GetOutputEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// responseoutput, err := dao.FindOutput(input)
	responseoutput, err := dao.FindOutput1Min(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, responseoutput)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/stream", AllStreamsEndPoint).Methods("GET")
	r.HandleFunc("/stream", CreateStreamEndPoint).Methods("POST")
	r.HandleFunc("/output", GetOutputEndPoint).Methods("GET")
	if err := http.ListenAndServe(":3001", r); err != nil {
		log.Fatal(err)
	}
}
