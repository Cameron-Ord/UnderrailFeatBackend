package db

import (
	"database/sql"
	"encoding/json"
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
	var UserBuildIdSlice []Build_ID_Titles
	UserBuildIdSlice, err = getUserBuildIds(dbConn, user_session_data)
	if err != nil {
		fmt.Println("Failed at retrieving IDS: ", err)
		return nil, err
	}
	type MiddleMan struct {
		Skill_Slice []Skill_Info
		Stat_Slice  []Stat_Info
		Feat_Slice  []Feat_Info
		Build_Title string
		Build_ID    uint
	}
	AllReturnedData := []MiddleMan{}
	for i := 0; i < len(UserBuildIdSlice); i++ {
		id_val := UserBuildIdSlice[i]
		ReturnedSkillSlice, err := retrieve_user_skills(dbConn, id_val.Build_ID, user_session_data)
		if err != nil {
			fmt.Println("Failed at retrieving skills: ", err)
			return nil, err
		}
		ReturnedStatSlice, err := retrieve_user_stats(dbConn, id_val.Build_ID, user_session_data)
		if err != nil {
			fmt.Println("Failed at retrieving stats: ", err)
			return nil, err
		}
		ReturnedFeatSlice, err := retrieve_user_feats(dbConn, id_val.Build_ID, user_session_data)
		if err != nil {
			fmt.Println("Failed at retrieving feats: ", err)
			return nil, err
		}
		middleMan := MiddleMan{
			Skill_Slice: ReturnedSkillSlice,
			Stat_Slice:  ReturnedStatSlice,
			Feat_Slice:  ReturnedFeatSlice,
			Build_Title: id_val.Title,
			Build_ID:    id_val.Build_ID,
		}
		AllReturnedData = append(AllReturnedData, middleMan)
	}
	jsonified_data, err := json.Marshal(AllReturnedData)
	if err != nil {
		return nil, err
	}
	return jsonified_data, nil
}
func getUserBuildIds(db *sql.DB, user_session_data User_Session_Data) ([]Build_ID_Titles, error) {
	var query string
	user_builds := []Build_ID_Titles{}
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
		currentBuild := Build_ID_Titles{
			Build_ID: id_holder,
			Title:    title_holder,
		}
		user_builds = append(user_builds, currentBuild)
		fmt.Printf("Build ID: %d, Build Title: %s\n", currentBuild.Build_ID, currentBuild.Title)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return user_builds, nil
}

func retrieve_user_skills(db *sql.DB, id_val uint, user_session_data User_Session_Data) ([]Skill_Info, error) {
	all_skill_info := []Skill_Info{}
	rows, err := db.Query("CALL get_client_skills(?,?,?)", id_val, user_session_data.Client_ID_Value, user_session_data.Client_Session_Token)
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
func retrieve_user_stats(db *sql.DB, id_val uint, user_session_data User_Session_Data) ([]Stat_Info, error) {
	all_stat_info := []Stat_Info{}
	rows, err := db.Query("CALL get_client_stats(?,?,?)", id_val, user_session_data.Client_ID_Value, user_session_data.Client_Session_Token)
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
		currentStat := Stat_Info{
			Name:     skill_name_holder,
			Value:    skill_value_holder,
			Build_ID: id_val,
		}
		all_stat_info = append(all_stat_info, currentStat)
	}
	return all_stat_info, nil
}
func retrieve_user_feats(db *sql.DB, id_val uint, user_session_data User_Session_Data) ([]Feat_Info, error) {

	all_feat_info := []Feat_Info{}
	rows, err := db.Query("CALL get_client_feats(?,?,?)", id_val, user_session_data.Client_ID_Value, user_session_data.Client_Session_Token)
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
