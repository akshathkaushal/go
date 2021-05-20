package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// use time.now().nano unixnano is correct

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890~!@#$%^&*")
var uniqueKeyLength = 10
var id = 0 // id for the incoming user requests
var IDlock sync.Mutex

func generateRandomString() string {
	b := make([]rune, uniqueKeyLength)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// unique id UniQ is not exported
type userData struct {
	UniQ      string
	ID        int
	Name      string         `json:"Name"`
	Address   string         `json:"Address"`
	PhNo      string         `json:"PhNo"`
	Type      string         `json:"Type"`
	Desc      string         `json:"Desc"`
	ReqID     map[int]string `json:"ReqID"`
	PendReqID map[int]string `json:"PendReqID"`
	ConID     map[int]string `json:"ConID"`
}

type userKey struct {
	UniqKey string `json:"UniqKey"`
}

var userDataMap = map[string]userData{} // contains the mapping of unique user ids to userdata struct
//var IDTypeMap = map[string]string{}     // contains mapping from user id to user type for identification

// Requested users ids:			// map of ids for the users that you have requested
// Pending request user ids:	// map of ids of the users that have requested you
// connected users ids:			// map of ids of users that you are connected to

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// method of sending secret code is ?code=<code>
	} else {
		var userUniqueCode userKey
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userUniqueCode)
		if err != nil {
			panic(err)
		}
		if _, present := userDataMap[userUniqueCode.UniqKey]; present {
			log.Println("User present")
			json.NewEncoder(w).Encode(userDataMap[userUniqueCode.UniqKey])
		} else {
			log.Println("User not present")
		}
	}
}

func create(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		panic("Error, wrong request method")
	} else {
		//
		var newUser userData
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newUser)
		if err != nil {
			panic(err)
		}

		// concurrency
		IDlock.Lock()
		id++
		newUser.ID = id
		uqKey := generateRandomString()
		IDlock.Unlock()

		if val, present := userDataMap[uqKey]; present {
			log.Println(val)
			create(w, r) // again create the user with a different unique string
		} else {
			newUser.UniQ = uqKey
			userDataMap[uqKey] = newUser
			//IDTypeMap[strconv.Itoa(newUser.ID)] = newUser.Type
		}

		json.NewEncoder(w).Encode(newUser)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		//
	} else {
		var userUniqueCode userKey
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userUniqueCode)
		if err != nil {
			panic(err)
		}

		if _, present := userDataMap[userUniqueCode.UniqKey]; present {
			delete(userDataMap, userUniqueCode.UniqKey)
			log.Println("User deleted successfully")
		} else {
			log.Println("User not present")
		}

	}
}

func update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		panic("Error, wrong request method")
	} else {
		//

		var updatedUser userData
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&updatedUser)
		if err != nil {
			panic(err)
		}

		if _, present := userDataMap[updatedUser.UniQ]; present {
			// update here

			updatedUser.ID = userDataMap[updatedUser.UniQ].ID
			userDataMap[updatedUser.UniQ] = updatedUser
			json.NewEncoder(w).Encode(updatedUser)
		} else {
			log.Println("User not present")
		}

	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		//
	} else {
		var entry = -1
		r.ParseForm()
		for code, value := range r.Form {
			if code == "id" {
				entry, _ = strconv.Atoi(strings.Join(value, ""))
			}
		}
		for _, user := range userDataMap {
			if user.ID == entry {
				toSendUser := user
				toSendUser.UniQ = ""
				toSendUser.ID = -1
				json.NewEncoder(w).Encode(toSendUser)
				break
			}
		}

	}
}

func getAllPatient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// implement the check that the requesting user is a donor
	} else {
		var userUniqueCode userKey
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userUniqueCode)
		if err != nil {
			panic(err)
		}

		if _, present := userDataMap[userUniqueCode.UniqKey]; present {
			if userDataMap[userUniqueCode.UniqKey].Type == "Donor" {
				log.Println("Donor detected")
				retUserArray := []userData{}
				for _, user := range userDataMap {
					if user.Type == "Patient" {
						toSendUser := user
						toSendUser.UniQ = ""
						toSendUser.ID = -1
						retUserArray = append(retUserArray, toSendUser)
					}
				}
				json.NewEncoder(w).Encode(retUserArray)
			} else {
				log.Println("You need to be a donor to see all the patients")
			}
		}
	}
}

func getAllDonor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// implement the check that the requesting user is a donor
	} else {
		var userUniqueCode userKey
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userUniqueCode)
		if err != nil {
			panic(err)
		}

		if _, present := userDataMap[userUniqueCode.UniqKey]; present {
			if userDataMap[userUniqueCode.UniqKey].Type == "Patient" {
				retUserArray := []userData{}
				for _, user := range userDataMap {
					if user.Type == "Donor" {
						toSendUser := user
						toSendUser.UniQ = ""
						toSendUser.ID = -1
						retUserArray = append(retUserArray, toSendUser)
					}
				}
				json.NewEncoder(w).Encode(retUserArray)
			} else {
				log.Println("You need to be a patient to see all the donors")
			}
		}
	}
}

func sendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// ?name=<name>
	} else {
		var userUniqueCode userKey
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userUniqueCode)
		if err != nil {
			panic(err)
		}

		r.ParseForm()
		if _, present := userDataMap[userUniqueCode.UniqKey]; present {
			//log.Println("Requesting user found while sending request")
			for key, value := range r.Form {
				if key == "name" {
					log.Println(key + ": " + strings.Join(value, ""))
					for _, user := range userDataMap {
						if user.Name == strings.Join(value, "") {
							//log.Println("Requested user found while sending request")
							userDataMap[userUniqueCode.UniqKey].ReqID[user.ID] = user.Name // update the map for the requesting user
							user.PendReqID[userDataMap[userUniqueCode.UniqKey].ID] = userDataMap[userUniqueCode.UniqKey].Name
							//fmt.Printf("Details added successfully to both users")

						}
					}
				}
			}
		}
	}
}

func acceptRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// ?user=<name>
	} else {
		var userUniqueCode userKey
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userUniqueCode)
		if err != nil {
			panic(err)
		}

		r.ParseForm()
		if _, present := userDataMap[userUniqueCode.UniqKey]; present {
			for key, value := range r.Form { // value is the requested user's name
				if key == "name" {
					// value is the patient's name
					for id, name := range userDataMap[userUniqueCode.UniqKey].PendReqID {
						if name == strings.Join(value, "") {
							userDataMap[userUniqueCode.UniqKey].ConID[id] = name      // update the connected map for the donor
							delete(userDataMap[userUniqueCode.UniqKey].PendReqID, id) // delete the entry from the donor's pending request table
						}
					}

					//log.Println("All okay till here 1")
					// for the patient update the connected map and delete the entry from the requested map
					for _, user := range userDataMap {
						if user.Name == strings.Join(value, "") {
							user.ConID[userDataMap[userUniqueCode.UniqKey].ID] = userDataMap[userUniqueCode.UniqKey].Name // fill the donor in connected map of patient
							for rid, rname := range user.ReqID {
								if rid == userDataMap[userUniqueCode.UniqKey].ID {
									delete(user.ReqID, rid)
									log.Println(rname)
								}
							}
						}
					}
					//log.Println("All okay till here 2")
				}
			}
		}
	}
}

func CancelConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		// ?user=<name>
	} else {
		var userUniqueCode userKey
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userUniqueCode)
		if err != nil {
			panic(err)
		}

		r.ParseForm()
		for key, value := range r.Form {

			if key == "name" {
				// delete the entry from donor's and patient's connection map
				for _, user := range userDataMap {
					if user.Name == strings.Join(value, "") {
						delete(userDataMap[userUniqueCode.UniqKey].ConID, user.ID) // deleted from the user deleting
						delete(user.ConID, userDataMap[userUniqueCode.UniqKey].ID) // deleted the deleting user entry
					}
				}
			}
		}

	}
}

func cancelRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		// ?user=<name>
	} else {
		var userUniqueCode userKey
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userUniqueCode)
		if err != nil {
			panic(err)
		}

		r.ParseForm()
		for key, value := range r.Form {

			if key == "name" {
				// delete the entry from donor's pending request and patient's requesting map
				// here value is the name of user whose request is to be deleted
				for _, user := range userDataMap {
					if user.Name == strings.Join(value, "") {

						if userDataMap[userUniqueCode.UniqKey].Type == "Donor" { // if donor is deleting request
							delete(userDataMap[userUniqueCode.UniqKey].PendReqID, user.ID) // deleted from the user deleting
							delete(user.ReqID, userDataMap[userUniqueCode.UniqKey].ID)     // deleted the deleting user entry
						} else { // if patient is deleting request
							delete(userDataMap[userUniqueCode.UniqKey].ReqID, user.ID)     // deleted from the user deleting
							delete(user.PendReqID, userDataMap[userUniqueCode.UniqKey].ID) // deleted the deleting user entry
						}

					}
				}
			}
		}

	}
}

func main() {

	http.HandleFunc("/login/", login)
	http.HandleFunc("/create/", create)
	http.HandleFunc("/delete/", deleteUser)
	http.HandleFunc("/update/", update)
	http.HandleFunc("/getu/", getUser)
	http.HandleFunc("/getad/", getAllDonor)
	http.HandleFunc("/getap/", getAllPatient)
	http.HandleFunc("/sendreq/", sendRequest)
	http.HandleFunc("/acceptreq/", acceptRequest)
	http.HandleFunc("/cancelcon/", CancelConnection)
	http.HandleFunc("/cancelreq/", cancelRequest)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
