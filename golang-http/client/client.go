package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var SERVER_ENDPOINT string = "http://localhost"
var SERVER_PORT string = ":3000"

type ClientApp struct {
	ClientApp http.Client
}

type User struct {
	Name string `json:"name"`
}

func CreateNewClientApp(timeout int) ClientApp {
	return ClientApp{http.Client{Timeout: time.Duration(timeout) * time.Second}}
}

func (c *ClientApp) getAllUsers() {
	resp, err := c.ClientApp.Get(SERVER_ENDPOINT + SERVER_PORT + "/users")
	if err != nil {
		fmt.Println("Error in GET", err)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	// fmt.Printf("Body : %v", body)

	var user map[string]User
	e := json.Unmarshal(body, &user)
	if e != nil {
		fmt.Println("Error decoding", e)
		return
	}
	fmt.Printf("Body : %+v", user)
}

func (c *ClientApp) addUser(name string) {
	newUser := User{Name: name}
	jsondata, err := json.Marshal(newUser)
	reqbody := bytes.NewBuffer(jsondata)
	// if err != nil {
	// 	fmt.Printf("Error marshalling :", err)
	// 	return
	// }

	resp, err := c.ClientApp.Post(SERVER_ENDPOINT+SERVER_PORT+"/users", "application/json", reqbody)

	if err != nil {
		fmt.Println("Error adding user :", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.Header)
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("Body : %v", body)

}

func main() {
	fmt.Println("http client implementation")

	c := CreateNewClientApp(5)
	c.addUser("Advith")
	c.addUser("Aezack")
	c.addUser("Normal")

	c.getAllUsers()

}
