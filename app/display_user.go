package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type DispUser struct {
	UserName string
}

func findUserBill(uname string, users []User) int {
	for _, user := range users {
		if user.UserName == uname {
			return user.Balance
		}
	}
	return 0
}

func DisplayUser(w http.ResponseWriter, r *http.Request) {
	users, err := collectUserRecords()
	if err != nil {
		log.Println("Not able to collect the users")
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("Error while reading the body")
		return
	}

	var msg DispUser
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Println("Error while copying the user to the local")
		return
	}
	cost := findUserBill(msg.UserName, users)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cost)
}
