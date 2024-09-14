package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// simulate database by using map. key is userID, value is user's name
var apiV1Users = make(map[string]string)

func handleApiV1GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ContextUserIDKey).(string)

	name, ok := apiV1Users[userID]
	if !ok {
		w.Write([]byte("Not found user " + userID + "  name"))
		return
	}

	w.Write([]byte("User ID: " + userID + " Name: " + name))
}

func handleApiV1PostUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ContextUserIDKey).(string)

	var reqBody struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	apiV1Users[userID] = reqBody.Name

	log.Printf("Saving user %s name %s", userID, apiV1Users[userID])
	w.Write([]byte("User Name: " + apiV1Users[userID]))
}

func handleApiV1GetUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	name, ok := apiV1Users[userID]
	if !ok {
		w.Write([]byte("Not found user " + userID + "  name"))
		return
	}

	w.Write([]byte("User ID: " + userID + " Name: " + name))
}
