package main

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Note struct {
	Id             int    `json:""`
	Title          string `json:""`
	PositivePoints string `json:""`
	NegativePoints string `json:""`
	KeyWords       []string
}

func (n *Note) InsertNote() (int64,error) {
	db, err := connectDB()
	if err != nil {
		log.Panic(err)
		return -1, err
	}
	defer db.Close()
	query := "INSERT INTO note(Title, PositivePoints, NegativePoints) VALUES(?, ?, ?)"
	result , err := db.Exec(query, n.Title, n.PositivePoints, n.NegativePoints)
	if err != nil {
		log.Panic(err)
		return -1, err
	}
	noteID, err := result.LastInsertId()
	if err != nil {
		log.Panic(err)
		return 0, err
	}
	
	err = n.InsertKeyword(noteID)
	if err != nil {
		return noteID, err
	}
	return noteID,nil
}

//Cette fonction est appelée à la fin de l'insertion d'une note.
func (n *Note) InsertKeyword(noteID int64) error {
	db, err := connectDB()
	if err != nil {
		log.Panic(err)
		return err
	}
	defer db.Close()
	query := "INSERT INTO NoteKeywords(NoteId, Keyword) VALUES(?, ?)"
	for _,keyword := range n.KeyWords{
		_, err = db.Exec(query, noteID, keyword)
		if err != nil {
			log.Panic(err)
			return err
		}
	}
	return nil
}


func GetAllNotes()[]Note{
	db, err := connectDB()
	if err != nil{
		log.Panic(err)
	}
	defer db.Close()
	query := "SELECT * FROM Note"
	var myNotes []Note
	rows, err := db.Query(query)
	defer rows.Close()
	fmt.Println("Mes rows sont:",rows)
	for rows.Next(){
		var note Note

		if err := rows.Scan(&note.Id, &note.Title, &note.PositivePoints,&note.NegativePoints); err != nil{
			log.Fatal(err)
		}

		note.KeyWords, err = getKeywordsForNote(db, note.Id)
		if err != nil {
			log.Panic(err)
		}
		myNotes = append(myNotes, note)
	}
	return myNotes
}

//Cette fontion est appelée directement
func getKeywordsForNote(db *sql.DB, noteID int) ([]string, error) {
	query := "SELECT Keyword FROM NoteKeywords WHERE NoteId = ?"
	rows, err := db.Query(query, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []string
	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			return nil, err
		}
		keywords = append(keywords, keyword)
	}
	return keywords, rows.Err()
}