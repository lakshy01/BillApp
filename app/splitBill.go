package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Member struct {
	UserName    string
	AmountSpent int
	AmountOwned int
}

type BillRequest struct {
	GroupName    string
	TotalAmount  int
	NumberOfUser int
	SplitEqual   bool
	Members      []Member
}

func helpGroupBill(groups []Group, uname string, gname string, bill int) []Group {
	for i, group := range groups {
		if group.GroupName == gname {
			for j, mem := range groups[i].Members {
				for mem.UserName == uname {
					groups[i].Members[j].Balance = bill
					break
				}
			}
		}
	}
	return groups
}

func updateUserBills(groups []Group, users []User) []User {
	for i, user := range users {
		users[i].Balance = 0
		for _, group := range groups {
			for _, mem := range group.Members {
				if mem.UserName == user.UserName {
					users[i].Balance += mem.Balance
				}

			}
		}
	}
	return users
}

func updateGroupBills(gname string, mems []Member, ta int, groups []Group) []Group {
	for _, mem := range mems {
		uname := mem.UserName
		bill := mem.AmountSpent - mem.AmountOwned
		if bill < 0 {
			bill *= (-1)
		}
		groups = helpGroupBill(groups, uname, gname, bill)
	}
	return groups
}

func updateUser(users []User) error {
	content, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		log.Println("Error while converting to the json object")
		return err
	}
	err = ioutil.WriteFile("Users.json", content, 0644)
	if err != nil {
		log.Println("Error while updating")
		return err
	}
	return nil
}

func SplitBill(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("Error while reading the request body")
		return
	}

	var msg BillRequest
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Println("Error while copying the body to the local")
		return
	}
	if msg.SplitEqual == true {
		for i, _ := range msg.Members {
			msg.Members[i].AmountOwned = msg.TotalAmount / msg.NumberOfUser
		}
	}

	groups, err := CollectGroupRecords()
	if err != nil {
		log.Println("Error while getting the groups")
		return
	}

	users, err := collectUserRecords()
	if err != nil {
		log.Println("Error while getting the users")
		return
	}

	groups = updateGroupBills(msg.GroupName, msg.Members, msg.TotalAmount, groups)
	users = updateUserBills(groups, users)
	err = updateGroup(groups)
	if err != nil {
		log.Println("Cannot able to update the group")
		return
	}
	err = updateUser(users)
	if err != nil {
		log.Println("Cannot able to update the user")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}
