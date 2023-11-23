package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func ConnectForLogin(loginQuery LoginData) ([]byte, error) {
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
	jsonified_struct, err := checkPassword(dbConn, loginQuery)
	if err != nil {
		return nil, err
	}

	return jsonified_struct, nil
}

func comparePWHash(hashed_password string, loginQuery LoginData) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(loginQuery.Password))
	if err != nil {
		return err
	}
	fmt.Println("Password -> Match")
	return nil
}

func checkPassword(db *sql.DB, loginQuery LoginData) ([]byte, error) {
	rows, err := db.Query("CALL get_hpw(?)", loginQuery.Username)
	var hashed_password string
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var hashedPW string
		err := rows.Scan(&hashedPW)
		if err != nil {
			return nil, err
		}
		hashed_password = hashedPW
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err, "rows error")
		return nil, err
	}

	err = comparePWHash(hashed_password, loginQuery)
	if err != nil {
		fmt.Println(err, "compare error")
		return nil, err
	}
	err = hashLoginPassword(&loginQuery)
	if err != nil {
		fmt.Println(err, "hashing error")
		return nil, err
	}
	client_id, token_db, err := commitLogin(db, loginQuery.Username, loginQuery.Password)
	if err != nil {
		return nil, err
	}

	structToMarshall := SessionDataStruct{
		Client_ID:     client_id,
		Session_Token: token_db,
	}

	session_json, errjson := marshall_session(structToMarshall)
	if errjson != nil {
		return nil, errjson
	}

	return session_json, nil
}
func marshall_session(structToMarshall SessionDataStruct) ([]byte, error) {
	json, err := json.Marshal(structToMarshall)
	if err != nil {
		return nil, err
	}
	return json, nil
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

func commitLogin(db *sql.DB, username string, password string) (uint, string, error) {
	fmt.Println("Commiting login..")
	tokenLength := 16
	token, err := generateToken(tokenLength)
	var client_id uint
	var token_db string
	if err != nil {
		return client_id, token_db, err
	}
	rows, err := db.Query("CALL client_login(?,?,?)", username, password, token)
	if err != nil {
		return client_id, token_db, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&client_id, &token_db)
		if err != nil {
			return client_id, token_db, err
		}
	}
	err = rows.Err()
	if err != nil {
		return client_id, token_db, err
	}
	return client_id, token_db, nil
}

func hashLoginPassword(loginQuery *LoginData) error {
	fmt.Println("hashing..")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(loginQuery.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	loginQuery.Password = string(hashedPassword)
	return nil
}
