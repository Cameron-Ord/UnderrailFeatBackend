package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func SaveBuild(build SaveData) error {
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
	err = saveBuildTitle(dbConn, &build)
	if err != nil {
		fmt.Println("Error saving build title: ", err)
		return err
	}
	err = saveBuildStats(dbConn, &build)
	if err != nil {
		fmt.Println("Error saving build stats: ", err)
		return err
	}
	err = saveBuildSkills(dbConn, &build)
	if err != nil {
		fmt.Println("Error saving build skills: ", err)
		return err
	}

	err = saveBuildFeats(dbConn, &build)
	if err != nil {
		fmt.Println("Error saving build feats: ", err)
		return err
	}

	fmt.Println("Finished..")
	return nil
}

func saveBuildTitle(db *sql.DB, build *SaveData) error {
	fmt.Println("Commiting build title to DB..")
	rows, err := db.Query("CALL insert_build(?,?,?)", build.Title, build.Client_ID, build.Session_Token)
	if err != nil {
		return err
	}
	var build_id uint
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&build_id)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func saveBuildStats(db *sql.DB, build *SaveData) error {
	fmt.Println("Commiting build stats to DB..")
	var query string
	var build_id uint
	query = "SELECT id FROM builds WHERE client_id = ? AND title = ?"
	rows, err := db.Query(query, build.Client_ID, build.Title)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&build_id)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	query = "CALL save_stats(?,?,?,?,?)"
	for i := 0; i < len(build.Stats); i++ {
		stat := build.Stats[i]
		_, err = db.Exec(query, build_id, stat.StatName, stat.StatValue, build.Session_Token, build.Client_ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveBuildSkills(db *sql.DB, build *SaveData) error {
	fmt.Println("Commiting build skills to DB..")
	var query string
	var build_id uint
	query = "SELECT id FROM builds WHERE client_id = ? AND title = ?"
	rows, err := db.Query(query, build.Client_ID, build.Title)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&build_id)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	query = "CALL save_skills(?,?,?,?,?)"
	for i := 0; i < len(build.Skills); i++ {
		skill := build.Skills[i]
		_, err = db.Exec(query, build_id, skill.SkillName, skill.SkillValue, build.Session_Token, build.Client_ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveBuildFeats(db *sql.DB, build *SaveData) error {
	fmt.Println("Commiting build feats to DB..")
	var query string
	var build_id uint
	query = "SELECT id FROM builds WHERE client_id = ? AND title = ?"
	rows, err := db.Query(query, build.Client_ID, build.Title)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&build_id)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	query = "CALL save_feats(?,?,?,?)"
	for i := 0; i < len(build.Feats); i++ {
		feat := build.Feats[i]
		_, err := db.Exec(query, build_id, build.Session_Token, build.Client_ID, feat.FeatName)
		if err != nil {
			return err
		}
	}
	return nil
}
