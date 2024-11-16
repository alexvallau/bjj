package main

import (
	"fmt"
	"net/http"
	"log"
)

func main() {



	http.HandleFunc("/getallnotes",GetAllNotesHandler)
	http.HandleFunc("/addnote", AddNoteHandler)
	http.HandleFunc("/deletenote", DeleteNoteHandler)
	http.HandleFunc("/createUser", CreateUserUserHandler)
	http.HandleFunc("/login", loginHandler)
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
