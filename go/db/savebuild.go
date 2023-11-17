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
	Stats  []Stat  `json:"stats"`
	Skills []Skill `json:"skills"`
	Feats  []Feat  `json:"feats"`
}

func saveBuildData(db *sql.DB) error {

}

func SaveBuild(build SaveData) {
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
	err = saveBuildData(dbConn)
}
