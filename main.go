package main

import (
	// "encoding/json"
	"fmt"
	//"github.com/graphql-go/graphql"
	//	"github.com/graphql-go/handler"
	"log"
	//	"net/http"
	"api.reeych.com/models"

	"database/sql"
	_ "github.com/lib/pq"
	shuuid "github.com/lithammer/shortuuid"
	//"reflect"
)

func main() {

	connStr := "user=brent dbname=reeych_backend_dev sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db)

	rows, err := db.Query("SELECT uc.id, first_name, last_name, title, owner FROM (users JOIN cards_users_rel ON users.id = user_id) as uc JOIN cards ON uc.card_id = cards.id")

	defer rows.Close()

	for rows.Next() {
		var (
			userId    string
			firstName string
			lastName  string
			title     string
			owner     string
		)

		if err := rows.Scan(&userId, &firstName, &lastName, &title, &owner); err != nil {
			log.Fatal(err)
		}

		fmt.Println(userId, firstName, lastName, title, owner)
	}

	fmt.Println(shuuid.New())

	models.PrintMsg("Hello There Everyone!")
}
