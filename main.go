package main

import (
	"crud/servidor"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter() //contém todas as rotas
	router.HandleFunc("/add", servidor.CreateToodo).Methods(http.MethodPost)
	router.HandleFunc("/get-all", servidor.GetAllToodos).Methods(http.MethodGet)
	router.HandleFunc("/get-one/{id}", servidor.GetOneToodo).Methods(http.MethodGet) //rota com parametro
	router.HandleFunc("/update/{id}", servidor.UpdateToodo).Methods(http.MethodPut)
	router.HandleFunc("/delete/{id}", servidor.DeleteToodo).Methods(http.MethodDelete)

	fmt.Println("Server is running on port 3030")
	log.Fatal(http.ListenAndServe(":3030", router))
}
