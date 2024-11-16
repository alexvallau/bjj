package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)


func GetAllNotesHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Only POST method allowed",http.StatusMethodNotAllowed)
		return
	}
	var notes []Note
	notes = GetAllNotes()
	json, err := json.Marshal(notes)
	if err != nil{
		log.Panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}


func AddNoteHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Only POST method allowed",http.StatusMethodNotAllowed)
		return
	}
	var nouvelleNote Note
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&nouvelleNote)
	if err != nil{
		log.Fatal(err)
		return
	}
	nouvelleNote.InsertNote()
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == ""{
		fmt.Fprintf(w,"The id is empty")
		return
	}
	idInt,err := strconv.Atoi(id)
	if err != nil{
		log.Panic(err)
		return
	}
	deleteNote(idInt)
}

//////////////////////////////////////////////////////////////////////////
//Utilisateurs
//////////////////////////////////////////////////////////////////////////
func CreateUserUserHandler(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost{
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}
	var user Utilisateur
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil{
		log.Fatal(err)
		return
	}
	CreateUser(user.Username, user.Password)
}


func loginHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var user Utilisateur
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil{
		log.Fatal(err)
		return
	}

	isLoggedIn, err := user.Login()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !isLoggedIn {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	if user.Id == -1 {
		return
	}
	jwtToken, err := user.GenerateJWT(user.Id)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	fmt.Println("The token is ", jwtToken)
	w.Header().Set("Authorization", jwtToken)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "User logged in successfully")

}