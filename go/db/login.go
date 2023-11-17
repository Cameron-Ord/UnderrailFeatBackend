package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var client_id int
var token_db string

func comparePWHash(hashed_password string, loginQuery LoginData) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(loginQuery.Password))
	fmt.Println("Password -> Match")
	if err != nil {
		return err
	}
	return nil
}

func checkPassword(db *sql.DB, loginQuery LoginData) error {
	rows, err := db.Query("CALL get_hpw(?)", loginQuery.Username)
	var hashed_password string
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var hashedPW string
		err := rows.Scan(&hashedPW)
		if err != nil {
			return err
		}
		hashed_password = hashedPW
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	err = comparePWHash(hashed_password, loginQuery)
	if err != nil {
		return err
	}

	err = hashLoginPassword(&loginQuery)
	if err != nil {
		return err
	}
	err = commitLogin(db, loginQuery.Username, loginQuery.Password)

	return nil
}

func generateToken(length int) (string, error) {

	hexlength := length / 2
	randomBytes := make([]byte, hexlength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(randomBytes)
	return token, nil
}

func commitLogin(db *sql.DB, username string, password string) error {

	tokenLength := 16
	token, err := generateToken(tokenLength)
	if err != nil {
		return err
	}
	rows, err := db.Query("CALL client_login(?,?,?)", username, password, token)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&client_id, &token_db)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	fmt.Println(client_id)
	fmt.Println(token_db)
	return nil
}

func hashLoginPassword(loginQuery *LoginData) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginQuery.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	loginQuery.Password = string(hashedPassword)
	return nil
}

func ConnectForLogin(loginQuery LoginData) {
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
	err = checkPassword(dbConn, loginQuery)
}
