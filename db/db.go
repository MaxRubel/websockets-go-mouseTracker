package db

import (
	"database/sql"
	"fmt"

	"github.com/MaxRubel/WebsocketsGo/models"
	_ "github.com/mattn/go-sqlite3"
)


func AddMessageToDb(newMessage models.Message) error {

    db, err := sql.Open("sqlite3", "./db/db.sqlite3")
    if err != nil {
        fmt.Println("error opening database:", err)
        return fmt.Errorf("unable to open database: %v", err)
    }
    defer db.Close()

	fmt.Println("new message: ", newMessage)

    query := `INSERT INTO Messages (name, message) VALUES (?, ?)`
    
    _, err = db.Exec(query, newMessage.Name, newMessage.Message)
    if err != nil {
        fmt.Println("error inserting message into db:", err)
        return fmt.Errorf("unable to insert new message into db: %v", err)
    } else {
        fmt.Println("Message inserted successfully")
    }

    return nil
}

func GetAllMessages() ([]models.Message, error) {
	db, err := sql.Open("sqlite3", "./db/db.sqlite3")

	if err != nil {
		return nil, fmt.Errorf("unable to open database: %v", err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM Messages")

	var messages []models.Message
	for rows.Next(){
		var message models.Message
		err := rows.Scan(&message.Id, &message.Name, &message.Message)

		if err != nil {
			return nil, fmt.Errorf("unable to scan and assign data correctly: %v", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}
