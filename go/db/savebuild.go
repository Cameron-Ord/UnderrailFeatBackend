package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Stat struct {
	StatName  string `json:"statName"`
	StatValue string `json:"statValue"`
}

type Skill struct {
	SkillName  string `json:"skillName"`
	SkillValue string `json:"skillValue"`
}

type Feat struct {
	FeatName string `json:"Feat"`
}

type SaveData struct {
	Stats         []Stat  `json:"stats"`
	Skills        []Skill `json:"skills"`
	Feats         []Feat  `json:"feats"`
	Title         string  `json:"title"`
	Client_ID     string  `json:"client_id"`
	Session_Token string  `json:"session_token"`
}

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
	fmt.Println(build.Client_ID)
	err = saveBuildTitle(dbConn, &build)
	if err != nil {
		return err
	}
	err = saveBuildStats(dbConn, &build)
	if err != nil {
		return err
	}
	return nil
}

func saveBuildTitle(db *sql.DB, build *SaveData) error {

	fmt.Println("Commiting build to DB..")
	rows, err := db.Query("CALL insert_build(?,?,?)", build.Title, build.Client_ID, build.Session_Token)
	if err != nil {
		return err
	}
	var build_id uint
	defer rows.Close()
	for rows.Next(){
		err := rows.Scan(&build_id)
		if err != nil {
			return err
		}	
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	fmt.Println("Inserted build ID: ",build_id)
	return nil
}


func saveBuildStats(db *sql.DB, build *SaveData) error {
	var query string
	var build_id uint
	query = "SELECT id FROM builds WHERE client_id = ? AND title = ?"
	rows, err := db.Query(query, build.Client_ID, build.Title)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next(){
		err := rows.Scan(&build_id)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err !=nil {
		return err
		}
	fmt.Println("Selected build ID: ",build_id)
	
	query = "CALL save_stats(?,?,?,?,?)"
	for i:=0; i<len(build.Stats); i++ {
		stat := build.Stats[i]
		_, err := db.Exec(query, build_id, stat.StatName, stat.StatValue, build.Session_Token, build.Client_ID)
		if err != nil {
			return err
		}
	}
	return nil
}
