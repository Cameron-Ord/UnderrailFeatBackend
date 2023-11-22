package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Feat_Info struct {
	Name     string
	Build_ID uint
}

type Skill_Info struct {
	Name     string
	Value    int
	Build_ID uint
}
type Stat_Info struct {
	Name     string
	Value    int
	Build_ID uint
}
type Build_ID_Titles struct {
	Build_ID uint
	Title    string
}

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
	type All_Returned_Data struct {
		Returned_Feat_Slice  [][]Feat_Info
		Returned_Stat_Slice  [][]Stat_Info
		Returned_Skill_Slice [][]Skill_Info
		Returned_Id_Slice    []Build_ID_Titles
	}
	var returnedData All_Returned_Data
	returnedData.Returned_Id_Slice, err = getBuildIds(dbConn)
	if err != nil {
		return nil, err
	}
	returnedData.Returned_Skill_Slice, err = retrieve_skills(dbConn, returnedData.Returned_Id_Slice)
	if err != nil {
		return nil, err
	}
	returnedData.Returned_Stat_Slice, err = retrieve_stats(dbConn, returnedData.Returned_Id_Slice)
	if err != nil {
		return nil, err
	}

	returnedData.Returned_Feat_Slice, err = retrieve_feats(dbConn, returnedData.Returned_Id_Slice)
	if err != nil {
		return nil, err
	}
	jsonified_data, err := json.Marshal(returnedData)
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
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id_holder uint
		var title_holder string
		fmt.Println("row iteration")
		err := rows.Scan(&id_holder, &title_holder)
		if err != nil {
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
		return nil, err
	}

	return builds, nil
}

func retrieve_skills(db *sql.DB, id_slice []Build_ID_Titles) ([][]Skill_Info, error) {

	var tha_big_slice [][]Skill_Info

	for i := 0; i < len(id_slice); i++ {
		all_skill_info := []Skill_Info{}
		id_val := id_slice[i]
		rows, err := db.Query("CALL get_all_skills(?)", id_val.Build_ID)
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
				Build_ID: id_val.Build_ID,
			}
			all_skill_info = append(all_skill_info, currentSkill)
		}
		tha_big_slice = append(tha_big_slice, all_skill_info)
	}
	return tha_big_slice, nil
}

func retrieve_stats(db *sql.DB, id_slice []Build_ID_Titles) ([][]Stat_Info, error) {
	var tha_big_slice [][]Stat_Info

	for i := 0; i < len(id_slice); i++ {
		all_stat_info := []Stat_Info{}
		id_val := id_slice[i]
		rows, err := db.Query("CALL get_all_stats(?)", id_val.Build_ID)
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
				Build_ID: id_val.Build_ID,
			}
			all_stat_info = append(all_stat_info, currentStat)
		}
		tha_big_slice = append(tha_big_slice, all_stat_info)
	}
	return tha_big_slice, nil
}

func retrieve_feats(db *sql.DB, id_slice []Build_ID_Titles) ([][]Feat_Info, error) {

	var tha_big_slice [][]Feat_Info

	for i := 0; i < len(id_slice); i++ {
		all_feat_info := []Feat_Info{}
		id_val := id_slice[i]
		rows, err := db.Query("CALL get_all_feats(?)", id_val.Build_ID)
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
				Build_ID: id_val.Build_ID,
			}
			all_feat_info = append(all_feat_info, currentFeat)
		}
		tha_big_slice = append(tha_big_slice, all_feat_info)
	}
	return tha_big_slice, nil
}
