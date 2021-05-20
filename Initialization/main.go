package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

//var rootPath = "http://localhost:9090"

// struct to store the user data
type userData struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	DOB  string `json:"DOB"`
	PhNo string `json:"PhNo"`
}

// array of structs to store users' data
type usersDataArray []userData

var FinalArray = usersDataArray{}

func createUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		var newUser userData

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&newUser)

		if err != nil {
			panic(err)
		}

		FinalArray = append(FinalArray, newUser)

		log.Println("createUser was executed successfully")
	} else {
		log.Println("Incorrect request method")
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		re, _ := regexp.Compile("/get/(.*)")
		values := re.FindStringSubmatch(r.URL.Path)

		recordID := values[1]

		for _, singleUser := range FinalArray {
			if singleUser.ID == recordID {
				json.NewEncoder(w).Encode(singleUser)
			}
		}

		log.Println("getUser was executed successfully")
	} else {
		log.Println("Incorrect request method")
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodDelete {
		re, _ := regexp.Compile("/delete/(.*)")
		values := re.FindStringSubmatch(r.URL.Path)

		recordID := values[1]

		for i, singleRecord := range FinalArray {
			if singleRecord.ID == recordID {
				FinalArray = append(FinalArray[:i], FinalArray[i+1:]...)
				fmt.Fprintf(w, "User with ID %v has been deleted successfully", recordID)
			}
		}
		log.Println("delete was executed successfully")
	} else {
		log.Println("Incorrect request method")
	}
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(FinalArray)
		log.Println("getAll was executed successfully")
	} else {
		log.Println("Incorrect request method")
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPatch {
		re, _ := regexp.Compile("/update/(.*)")
		values := re.FindStringSubmatch(r.URL.Path)

		recordID := values[1]

		var updatedUser userData

		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&updatedUser)

		if err != nil {
			panic(err)
		}

		for i, singleUser := range FinalArray {
			if singleUser.ID == recordID {
				singleUser.Name = updatedUser.Name
				singleUser.DOB = updatedUser.DOB
				singleUser.PhNo = updatedUser.PhNo

				FinalArray = append(FinalArray[:i], singleUser)
				json.NewEncoder(w).Encode(singleUser)
			}
		}
		log.Println("updateUser was executed successfully")
	} else {
		log.Println("Incorrect request method")
	}
}

func main() {

	http.HandleFunc("/create/", createUser)
	http.HandleFunc("/update/", updateUser)
	http.HandleFunc("/get/", getUser)
	http.HandleFunc("/getall/", getAllUsers)
	http.HandleFunc("/delete/", deleteUser)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
