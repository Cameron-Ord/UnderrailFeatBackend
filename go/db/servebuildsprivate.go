package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func GetUserBuilds(user_session_data User_Session_Data) ([]byte, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBHost, DBPort, DBName)
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer dbConn.Close()

	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Connected to the database!")
	var UserBuildIdSlice []User_Build_IDs
	UserBuildIdSlice, err = getUserBuildIds(dbConn, user_session_data)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(UserBuildIdSlice); i++ {
		fmt.Println(UserBuildIdSlice[i].Title)
	}
	return nil, nil
}
func getUserBuildIds(db *sql.DB, user_session_data User_Session_Data) ([]User_Build_IDs, error) {
	var query string
	user_builds := []User_Build_IDs{}
	query = "CALL select_user_build_ids(?,?)"
	rows, err := db.Query(query, user_session_data.Client_ID_Value, user_session_data.Client_Session_Token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id_holder uint
		var title_holder string
		err := rows.Scan(&id_holder, &title_holder)
		if err != nil {
			return nil, err
		}
		currentBuild := User_Build_IDs{
			Build_ID: id_holder,
			Title:    title_holder,
		}
		user_builds = append(user_builds, currentBuild)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return user_builds, nil
}
