package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ServeBuilds() ([]byte, error) {
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

	var Returned_Id_Slice []Build_ID_Titles
	Returned_Id_Slice, err = getBuildIds(dbConn)
	if err != nil {
		fmt.Println("Failed at getBuildIds")
		return nil, err
	}

	fmt.Println(DBUsername)
	type MiddleMan struct {
		Skill_Slice []Skill_Info
		Stat_Slice  []Stat_Info
		Feat_Slice  []Feat_Info
		Build_Title string
		Build_ID    uint
	}
	AllReturnedData := []MiddleMan{}
	for i := 0; i < len(Returned_Id_Slice); i++ {
		id_val := Returned_Id_Slice[i]
		returned_Skill_Slice, err := retrieve_skills(dbConn, id_val.Build_ID)
		if err != nil {
			return nil, err
		}
		returned_Stat_Slice, err := retrieve_stats(dbConn, id_val.Build_ID)
		if err != nil {
			return nil, err
		}
		returned_Feat_Slice, err := retrieve_feats(dbConn, id_val.Build_ID)
		if err != nil {
			return nil, err
		}
		middleMan := MiddleMan{
			Skill_Slice: returned_Skill_Slice,
			Stat_Slice:  returned_Stat_Slice,
			Feat_Slice:  returned_Feat_Slice,
			Build_Title: id_val.Title,
			Build_ID:    id_val.Build_ID,
		}
		AllReturnedData = append(AllReturnedData, middleMan)
	}

	jsonified_data, err := json.Marshal(AllReturnedData)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return nil, err
	}
	return jsonified_data, nil
}

func getBuildIds(db *sql.DB) ([]Build_ID_Titles, error) {
	var query string
	builds := []Build_ID_Titles{}
	buildsPtr := &builds
	query = "CALL select_newest_build_ids()"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Failed at query")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id_holder uint
		var title_holder string
		err := rows.Scan(&id_holder, &title_holder)
		if err != nil {
			fmt.Println("Failed at scan")
			return nil, err
		}

		currentBuild := Build_ID_Titles{
			Build_ID: id_holder,
			Title:    title_holder,
		}

		*buildsPtr = append(*buildsPtr, currentBuild)

	}
	err = rows.Err()
	if err != nil {
		fmt.Println("Rows err")
		return nil, err
	}

	return builds, nil
}

func retrieve_skills(db *sql.DB, id_val uint) ([]Skill_Info, error) {
	all_skill_info := []Skill_Info{}
	rows, err := db.Query("CALL get_all_skills(?)", id_val)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var skill_name_holder string
		var skill_value_holder int
		err := rows.Scan(&skill_name_holder, &skill_value_holder)
		if err != nil {
			return nil, err
		}
		currentSkill := Skill_Info{
			Name:     skill_name_holder,
			Value:    skill_value_holder,
			Build_ID: id_val,
		}
		all_skill_info = append(all_skill_info, currentSkill)
	}
	return all_skill_info, nil
}

func retrieve_stats(db *sql.DB, id_val uint) ([]Stat_Info, error) {

	all_stat_info := []Stat_Info{}
	rows, err := db.Query("CALL get_all_stats(?)", id_val)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var stat_name_holder string
		var stat_value_holder int
		err := rows.Scan(&stat_name_holder, &stat_value_holder)
		if err != nil {
			return nil, err
		}
		currentStat := Stat_Info{
			Name:     stat_name_holder,
			Value:    stat_value_holder,
			Build_ID: id_val,
		}
		all_stat_info = append(all_stat_info, currentStat)
	}

	return all_stat_info, nil
}

func retrieve_feats(db *sql.DB, id_val uint) ([]Feat_Info, error) {

	all_feat_info := []Feat_Info{}
	rows, err := db.Query("CALL get_all_feats(?)", id_val)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var feat_name string
		err := rows.Scan(&feat_name)
		if err != nil {
			return nil, err
		}
		currentFeat := Feat_Info{
			Name:     feat_name,
			Build_ID: id_val,
		}
		all_feat_info = append(all_feat_info, currentFeat)
	}
	return all_feat_info, nil
}
