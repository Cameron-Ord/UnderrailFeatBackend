package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type SignupData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func hashPassword(signupQuery *SignupData) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupQuery.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	signupQuery.Password = string(hashedPassword)
	return nil
}

func commitSignup(db *sql.DB, username string, password string) error {

	_, err := db.Exec("CALL client_signup(?,?)", username, password)
	if err != nil {
		return err
	}

	return nil
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
	err = hashPassword(&signupQuery)
	err = commitSignup(dbConn, signupQuery.Username, signupQuery.Password)
	fmt.Println(signupQuery.Password)
}
