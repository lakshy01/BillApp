package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type DisGroup struct {
	GroupName string
}

func findGroup(gname string, groups []Group) []User {
	for _, group := range groups {
		if group.GroupName == gname {
			return group.Members
		}
	}
	return []User{}
}

func DisplayGroup(w http.ResponseWriter, r *http.Request) {
	groups, err := CollectGroupRecords()
	if err != nil {
		log.Println("Error while collecting the groups")
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("Error while reading the body")
		return
	}
	var msg DisGroup
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Println("Error while copying the request body")
		return
	}
	users := findGroup(msg.GroupName, groups)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
