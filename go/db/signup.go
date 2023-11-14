package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type SignupData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ConnectForSignup(signupQuery SignupData) {
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
	fmt.Println()
}