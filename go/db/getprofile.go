package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func GetProfileInfo(user_session_data User_Session_Data) ([]byte, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBHost, DBPort, DBName)
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")
	user_info, err := fetchUsername(dbConn, user_session_data)
	if err != nil {
		return nil, err
	}
	jsonified_username, err := json.Marshal(user_info)
	if err != nil {
		return nil, err
	}
	return jsonified_username, nil
}

func fetchUsername(db *sql.DB, seshdata User_Session_Data) (Basic_User_Info, error) {

	currentUser := Basic_User_Info{}
	rows, err := db.Query("SELECT username, created_at FROM client WHERE id = ?", seshdata.Client_ID_Value)
	if err != nil {
		return currentUser, err
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		var created_at string
		err := rows.Scan(&username, &created_at)
		if err != nil {
			return currentUser, err
		}
		currentUser.Username = username
		currentUser.Created_at = created_at
	}
	err = rows.Err()
	if err != nil {
		return currentUser, err
	}

	return currentUser, nil
}
