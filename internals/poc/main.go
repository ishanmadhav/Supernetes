package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

func main() {
	// Read the JSON file into a byte array
	jsonData, err := ioutil.ReadFile("example.json")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Unmarshal the JSON data into a Person struct
	var person Person
	err = json.Unmarshal(jsonData, &person)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the Person struct
	fmt.Println(person)
}
