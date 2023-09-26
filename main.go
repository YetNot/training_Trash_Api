package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	port string = "8080"
	db   []Pizza
)

func init() {
	pizza1 := Pizza{
		ID:      1,
		Diametr: 22,
		Price:   500.50,
		Title:   "Pepperoni",
	}

	pizza2 := Pizza{
		ID:      2,
		Diametr: 25,
		Price:   400,
		Title:   "Mexico",
	}

	pizza3 := Pizza{
		ID:      3,
		Diametr: 30,
		Price:   901.34,
		Title:   "Domashniya",
	}

	db = append(db, pizza1, pizza2, pizza3)
}

type Pizza struct {
	ID      int     `json:"id"`
	Diametr int     `json:"diametr"`
	Price   float64 `json:"price"`
	Title   string  `json:"title"`
}

func FindPizzaById(id int) (Pizza, bool) {
	var pizza Pizza
	var found bool
	for _, p := range db {
		if p.ID == id {
			pizza = p
			found = true
			break
		}
	}
	return pizza, found
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func GetAllPizzas(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	log.Println("Get infos about pizzas in database")
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(db)
}

func GetPizzaById(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := ErrorMessage{Message: "do not use ID not supported int casting"}
		writer.WriteHeader(403)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	log.Println("Trying to send to client pizza with id #:", id)
	pizza, ok := FindPizzaById(id)
	if ok {
		writer.WriteHeader(200)
		json.NewEncoder(writer).Encode(pizza)
	} else {
		msg := ErrorMessage{Message: "pizza with that id does not exists in database"}
		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
	}
}

func main() {
	log.Println("Trying Api")
	router := mux.NewRouter()
	router.HandleFunc("/pizzas", GetAllPizzas).Methods("GET")
	router.HandleFunc("/pizza/{id}", GetPizzaById).Methods("GET")
	log.Println("Router configured succesfully!")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
