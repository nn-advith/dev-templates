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

var userCache = make(map[int]User) // simulating a database using a map
var cacheMutex sync.RWMutex

func rootHandler(
	w http.ResponseWriter,
	req *http.Request,
) {
	if req.URL.Path == "/" {
		// default / path handler
		fmt.Println(userCache)
		w.Write([]byte("Hello world. Wassup"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No. Not allowed."))
	}

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
	err := json.NewDecoder(req.Body).Decode(&user) // decode from request body and into user struct
	if err != nil {
		http.Error(w, "Error here", http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "Empty name is not allowed", http.StatusBadRequest)
		return
	}

	cacheMutex.Lock()
	userCache[len(userCache)+1] = user // write operation into db; so lock it
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func deleteUserFromId(
	w http.ResponseWriter,
	req *http.Request,
) {
	id, err := strconv.Atoi(req.PathValue("id")) // get query param id
	if err != nil {
		http.Error(w, "Error reading ID from path", http.StatusBadRequest)
		return
	}

	if _, ok := userCache[id]; !ok { //check existance in map
		http.Error(w, "User Not found", http.StatusNotFound)
		return
	}

	cacheMutex.Lock()
	delete(userCache, id)
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)

}

func unknownPath(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("This path is not defined"))
}

func hbrouting(
	w http.ResponseWriter,
	req *http.Request,
) {
	fmt.Fprintf(w, "Hit host based routing handlesr")
}

func hbroutingindex(
	w http.ResponseWriter,
	req *http.Request,
) {
	fmt.Fprintf(w, "Host-based-route server index page")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users", getAllUsers)
	mux.HandleFunc("GET /users/{id}", getUserFromId)
	mux.HandleFunc("DELETE /users/{id}", deleteUserFromId)
	mux.HandleFunc("samplehost.dev/", hbrouting)
	mux.HandleFunc("samplehost.dev/index", hbroutingindex)
	mux.HandleFunc("/", rootHandler)
	fmt.Println("Server is listening on :3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}

}
