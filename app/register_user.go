package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type User struct {
	UserName string
	Balance  int
}

func addUser(name string, bal int, users []User) ([]User, error) {
	users = append(users, User{UserName: name, Balance: bal})
	content, err := json.MarshalIndent(users, "", "	")
	if err != nil {
		log.Println("Not able to connect the users to json object")
		return users, err
	}
	err = ioutil.WriteFile("Users.json", content, 0644)
	if err != nil {
		log.Println("Error while adding the user")
		return users, err
	}
	return users, nil
}

func collectUserRecords() ([]User, error) {
	content, err := ioutil.ReadFile("Users.json")
	if err != nil {
		log.Println("Error while getting the content out of the Users.json")
		return []User{}, err
	}
	users := []User{}
	err = json.Unmarshal(content, &users)
	if err != nil {
		log.Println("Error while copying the users to the local slice")
		return []User{}, err
	}
	return users, nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("Users.json")
	users := []User{}
	if err == nil {
		users, err = collectUserRecords()
		if err != nil {
			log.Println("Error while collecting the previous users")
			return
		}

	}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("Error while referring to the body of the request")
		return
	}
	var msg User
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Println("Error while unmarshalling the request body")
		return
	}
	users, _ = addUser(msg.UserName, msg.Balance, users)
	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(users)
}
