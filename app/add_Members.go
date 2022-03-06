package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type mem struct {
	GroupName string
	Members   []User
}

func addMember(gname string, members []User, groups []Group) ([]Group, error) {
	for i, group := range groups {
		if group.GroupName == gname {
			members = append(members, User{UserName: group.LeaderName})
			groups[i].Members = members
			break
		}
	}
	return groups, nil
}

func updateGroup(groups []Group) error {
	content, err := json.MarshalIndent(groups, "", " ")
	if err != nil {
		log.Println("Error while converting to the json object")
		return err
	}
	err = ioutil.WriteFile("Groups.json", content, 0644)
	if err != nil {
		log.Println("Error while updating")
		return err
	}
	return nil
}

func AddMembers(w http.ResponseWriter, r *http.Request) {
	groups, err := CollectGroupRecords()
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("Error while reading the body")
		return
	}
	var msg mem
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Println("Error while referencing the message")
		return
	}
	groups, err = addMember(msg.GroupName, msg.Members, groups)
	if err != nil {
		log.Println("Error while adding the memebers in the respective group")
		return
	}
	err = updateGroup(groups)
	if err != nil {
		log.Println("Cannot able to update the group")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}
