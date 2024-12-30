package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserData struct {
	NickName string `json:"nickname"`
}

func getUser(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session-name")
	userName, _ := session.Values["name"].(string)
	userId, _ := session.Values["id"].(uint)

	log.Println(userName, userId)

	userData := UserData{
		NickName: userName,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(userData)
	if err != nil {
		log.Println(err)
	}
}
