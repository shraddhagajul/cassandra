package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Restful API using Go and Cassandra!")
}

func CreateStudent(w http.ResponseWriter, r *http.Request) {
	var newStudent Student
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "wrong data")
	}
	json.Unmarshal(reqBody, &newStudent)
	if err := Session.Query(
		"INSERT INTO students(id, firstname, lastname, age) VALUES(?, ?, ?, ?)",
		newStudent.Id, newStudent.FirstName, newStudent.LastName, newStudent.Age).Exec(); err != nil {
		fmt.Println("error while inserting")
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
	conv, _ := json.MarshalIndent(newStudent, "", "")
	fmt.Fprintf(w, "%s", string(conv))
}
