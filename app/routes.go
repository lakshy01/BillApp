package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	router := mux.NewRouter()
	router.HandleFunc("/register_user", CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/register_group", CreateGroup).Methods(http.MethodPost)
	router.HandleFunc("/add_members", AddMembers).Methods(http.MethodPost)
	router.HandleFunc("/disp_group", DisplayGroup).Methods(http.MethodPost)
	router.HandleFunc("/disp_user", DisplayUser).Methods(http.MethodPost)
	router.HandleFunc("/split_bill", SplitBill).Methods(http.MethodPost)
	fmt.Println("Running")
	err := http.ListenAndServe("localhost:8000", router)
	if err != nil {
		log.Println(err)
	}
}
