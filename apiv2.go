package main

import "net/http"

func handleApiV2GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ContextUserIDKey).(string)

	name, ok := apiV1Users[userID]
	if !ok {
		w.Write([]byte("V2: Not found user " + userID + "  name"))
		return
	}

	w.Write([]byte("V2: User ID: " + userID + " Name: " + name))
}
