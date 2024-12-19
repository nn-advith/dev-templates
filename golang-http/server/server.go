package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	Name string `json:"name"`
}

var userCache = make(map[int]User)
var cacheMutex sync.RWMutex

func rootHandler(
	w http.ResponseWriter,
	req *http.Request,
) {
	fmt.Println(userCache)
	fmt.Fprintf(w, "Hello World")

}

func getAllUsers(
	w http.ResponseWriter,
	req *http.Request,
) {
	jsondata, err := json.Marshal(userCache)
	if err != nil {
		http.Error(w, "Error marshalling user cacje", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsondata)

}

func getUserFromId(
	w http.ResponseWriter,
	req *http.Request,
) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "Error reading ID from path", http.StatusBadRequest)
		return
	}

	cacheMutex.RLock()
	user, ok := userCache[id]
	cacheMutex.RUnlock()

	if !ok {
		http.Error(w, "User Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "error marshalling", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func createUser(
	w http.ResponseWriter,
	req *http.Request,
) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error here", http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "Empty name is not allowed", http.StatusBadRequest)
		return
	}

	cacheMutex.Lock()
	userCache[len(userCache)+1] = user
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func deleteUserFromId(
	w http.ResponseWriter,
	req *http.Request,
) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "Error reading ID from path", http.StatusBadRequest)
		return
	}

	if _, ok := userCache[id]; !ok {
		http.Error(w, "User Not found", http.StatusNotFound)
		return
	}

	cacheMutex.Lock()
	delete(userCache, id)
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)

}

func hbrouting(
	w http.ResponseWriter,
	req *http.Request,
) {
	fmt.Fprintf(w, "Hit host based routing handler")
}

func hbroutingindex(
	w http.ResponseWriter,
	req *http.Request,
) {
	fmt.Fprintf(w, "Host-based-route server index page")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users", getAllUsers)
	mux.HandleFunc("GET /users/{id}", getUserFromId)
	mux.HandleFunc("DELETE /users/{id}", deleteUserFromId)

	mux.HandleFunc("samplehost.dev/", hbrouting)
	mux.HandleFunc("samplehost.dev/index", hbroutingindex)

	fmt.Println("Server is listening on :3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}

}
