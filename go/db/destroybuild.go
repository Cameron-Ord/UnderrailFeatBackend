package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func DestroyBuild(delete_builds DeleteData) error {
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
	stmt, err := dbConn.Prepare("DELETE FROM builds WHERE id = ? AND client_id = ?")
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(delete_builds.Build_ID, delete_builds.Client_ID)
	if err != nil {
		fmt.Println("Error executing statement:", err)
		return err
	}

	rows_affected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected:", err)
		return err
	}

	fmt.Println("Rows affected:", rows_affected)
	if rows_affected == 0 {
		fmt.Println("No rows affected")
		return errors.New("build not found")
	}
	return nil
}
