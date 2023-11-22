package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ServeBuilds() error {
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
	err = getBuildIds(dbConn)
	if err != nil {
		return err
	}
	return nil
}

func getBuildIds(db *sql.DB) error {

	var query string
	var ids uint
	var id_slice []uint
	slice_ptr := &id_slice
	query = "CALL select_all_build_ids()"
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		fmt.Println("row iteration")
		err := rows.Scan(&ids)
		if err != nil {
			return err
		}
		*slice_ptr = append(*slice_ptr, ids)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	fmt.Println(id_slice)
	return nil
}
