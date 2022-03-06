package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Group struct {
	GroupName    string
	LeaderName   string
	TotalBalance int
	Members      []User
}

func CollectGroupRecords() ([]Group, error) {
	content, err := ioutil.ReadFile("Groups.json")
	if err != nil {
		log.Println("Error while getting the content of the group.json")
		return []Group{}, nil
	}
	groups := []Group{}
	err = json.Unmarshal(content, &groups)
	if err != nil {
		log.Println("Error while copying the groups in the local groups slice")
		return []Group{}, nil
	}
	return groups, nil
}

func createGroups(gname string, lname string, groups []Group) ([]Group, error) {
	mems := []User{User{UserName: lname}}
	groups = append(groups, Group{GroupName: gname, LeaderName: lname, Members: mems})
	content, err := json.MarshalIndent(groups, "", "	")
	if err != nil {
		log.Println("Not able to connect to the groups in the Groups.json")
		return []Group{}, err
	}
	err = ioutil.WriteFile("Groups.json", content, 0644)
	if err != nil {
		log.Println("Errro while adding the group")
		return []Group{}, err
	}
	return groups, nil
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("Groups.json")
	groups := []Group{}
	if err == nil {
		groups, err = CollectGroupRecords()
		if err != nil {
			log.Println("Error while collecting the previous groups")
			return
		}
	}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("Error while referring the body")
		return
	}
	var msg Group
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Println("Error while unmarshalling the request body")
		return
	}
	groups, _ = createGroups(msg.GroupName, msg.LeaderName, groups)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}
