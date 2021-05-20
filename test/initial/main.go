package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type user struct {
	ID   string
	Name string `json:"Name"`
}

var userMap = map[string]user{}
var id = 0

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		//
	} else {
		var newUser user
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newUser)
		if err != nil {
			panic(err)
		}
		log.Println("Ok till here")
		id++
		newUser.ID = strconv.Itoa(id)
		userMap[strconv.Itoa(id)] = newUser
		log.Println(strconv.Itoa(id) + ":" + newUser.Name)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		//
	} else {
		r.ParseForm()
		for code, value := range r.Form {
			if code == "id" {
				entry := strings.Join(value, "")
				if _, present := userMap[entry]; present {
					json.NewEncoder(w).Encode(userMap[entry])
				} else {
					log.Println("User not present")
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/create/", createUser)
	http.HandleFunc("/get/", getUser)
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
