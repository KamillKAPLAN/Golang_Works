package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Person struct {
	ID        string
	FirstName string
	LastName  string
	Address   *Address
}
type Address struct {
	City  string
	State string
}

var people []Person

func GetPeople(w http.ResponseWriter, r *http.Request) {
	fmt.Print("List users\n")
	//Veri setimizle, tüm insan dilimlerinden getirelim:
	json.NewEncoder(w).Encode(people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Get user\n")
	// Tek bir veri al
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			fmt.Print(item)
			// mySlice := []Person{item}
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{"4", "kamil", "kaplan", &Address{}})

}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Create  user\n")
	// yeni bir öğe oluştur
	params := mux.Vars(r)
	var person Person
	json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Bir öğeyi sil
	params := mux.Vars(r)
	fmt.Printf("Delete user: %s\n", params["id"])
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}

func main() {
	fmt.Print("Server starts...\n")
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", FirstName: "John", LastName: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", FirstName: "Koko", LastName: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: "3", FirstName: "Francis", LastName: "Sunday"})

	router.HandleFunc("/people", GetPeople).Methods("GET")            // Telefon rehberi belgesindeki herkesi al
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")       // Tek bir kişi al
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")   // Yeni bir kişi oluştur
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE") // Bir kişi sil

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("Create  user\n")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// Insert the middleware
	handler := c.Handler(router)

	// handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":3000", handler))
}
