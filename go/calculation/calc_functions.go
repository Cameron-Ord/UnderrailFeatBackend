package calculation

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Check_For_Type(Feat_Type string, type_substr string) bool {
	trimmed_substr := strings.TrimSpace(type_substr)
	contains_substr := strings.Contains(Feat_Type, trimmed_substr)
	return contains_substr
}

func Marshal_Data(sentSlice []string) ([]byte, error) {
	//assigning jsondata from the sentslice(the amalgamation of the results from all the checks)
	jsonData, err := json.Marshal(sentSlice)
	if err != nil {
		return nil, err
	}
	//returning the data
	return jsonData, nil
}

func Unload_Json() []map[string]string {
	//reading the json file and unmarshalling it into the set data type
	jsonFile, err := os.ReadFile("jsonData/FeatsData.json")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	var feats []map[string]string
	err = json.Unmarshal(jsonFile, &feats)
	if err != nil {
		log.Fatal("Error decoding JSON data:", err)
	}
	//returning the feats map
	return feats
}

func Int_Conversion(givenNumStr string) (int, error) {
	//string to int conversion
	NumToInt, err := strconv.Atoi(givenNumStr)
	return NumToInt, err
}

func Str_Conversion(givenNumInt int) string {
	NumToStr := strconv.Itoa(givenNumInt)
	return NumToStr
}

func Perform_Check(metPtr *[]string, failedPtr *[]string, requirement string, name string, given_value int) error {
	required_value, err := Int_Conversion(requirement)
	if err != nil {
		fmt.Println("Error during str conversion", err)
		return err
	}
	if given_value >= required_value {
		fmt.Println("SUCCESS: ", name, "-> REASON: ", given_value, " > ", required_value)
		*metPtr = append(*metPtr, name)
	} else if given_value < required_value {
		fmt.Println("FAILED: ", name, "-> REASON: ", given_value, " < ", required_value)
		*failedPtr = append(*failedPtr, name)
	}
	return nil
}

func Get_Requirements(Feat map[string]string, data RequestData, stats_req *[]Stats_Tracker, skills_req *[]Skills_Tracker) {
	for k := 0; k < len(data.Stats); k++ {
		Stat := data.Stats[k]
		if feat_value, ok := Feat[Stat.StatName]; ok {
			stats_tracker := Stats_Tracker{
				Stat_Name:  Stat.StatName,
				Stat_Value: Stat.StatValue,
				Feat_Value: feat_value,
			}
			*stats_req = append(*stats_req, stats_tracker)
		}
	}

	for j := 0; j < len(data.Skills); j++ {
		Skill := data.Skills[j]
		if feat_value, ok := Feat[Skill.SkillName]; ok {
			skills_tracker := Skills_Tracker{
				Skill_Name:  Skill.SkillName,
				Skill_Value: Skill.SkillValue,
				Feat_Value:  feat_value,
			}
			*skills_req = append(*skills_req, skills_tracker)
		}
	}

	fmt.Println("Requirements retrieved..")
}

func Check_All(stats_required []Stats_Tracker, skills_required []Skills_Tracker, stats_met *[]string, stats_failed *[]string, skills_met *[]string, skills_failed *[]string) error {
	for r := 0; r < len(stats_required); r++ {
		err := Perform_Check(stats_met, stats_failed, stats_required[r].Feat_Value, stats_required[r].Stat_Name, stats_required[r].Stat_Value)
		if err != nil {
			return err
		}
	}
	for k := 0; k < len(skills_required); k++ {
		err := Perform_Check(skills_met, skills_failed, skills_required[k].Feat_Value, skills_required[k].Skill_Name, skills_required[k].Skill_Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func Check_Stats(stats_required []Stats_Tracker, stats_met *[]string, stats_failed *[]string) error {
	for r := 0; r < len(stats_required); r++ {
		err := Perform_Check(stats_met, stats_failed, stats_required[r].Feat_Value, stats_required[r].Stat_Name, stats_required[r].Stat_Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func Check_Skills(skills_required []Skills_Tracker, skills_met *[]string, skills_failed *[]string) error {
	for k := 0; k < len(skills_required); k++ {
		err := Perform_Check(skills_met, skills_failed, skills_required[k].Feat_Value, skills_required[k].Skill_Name, skills_required[k].Skill_Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func Run_Calculation(Feats []map[string]string, data RequestData, allAllocatedFeats *[]string) error {

	for i := 0; i < len(Feats); i++ {
		fmt.Println("Beginning checks on: ", Feats[i]["Feat"], "...")
		fmt.Println("CHARACTER TYPE: ", data.Character_Type)
		fmt.Println()
		Feat := Feats[i]
		var stats_required = []Stats_Tracker{}
		var skills_required = []Skills_Tracker{}
		fmt.Println("Getting requirements..")
		Get_Requirements(Feat, data, &stats_required, &skills_required)
		if len(stats_required) > 0 && len(skills_required) > 0 {
			var stats_met = []string{}
			var stats_failed = []string{}
			var skills_met = []string{}
			var skills_failed = []string{}
			//failure tracking, if this is true even once it's failure
			var UTTERFAILURE bool = false

			err := Check_All(
				stats_required, skills_required, &stats_met,
				&stats_failed, &skills_met, &skills_failed,
			)
			if err != nil {
				fmt.Println("Errored: ", err)
				return err
			}

			if len(stats_met) == 0 {
				fmt.Println("Failed(stat+skill): ", Feat["Feat"], "-> REASON: (No stats met)")
				UTTERFAILURE = true
			}
			var stats_are_failed bool = false
			if len(stats_failed) > 0 {
				stats_are_failed = true
			}

			if stats_are_failed {
				var required []string
				for l := 0; l < len(stats_failed); l++ {
					if hard_req, ok := Feat["NeedsAllStats"]; ok {
						required = append(required, hard_req)
					}
				}

				for j := 0; j < len(required); j++ {
					if required[j] == "true" {
						fmt.Println("Failed(stat+skill): ", Feat["Feat"], "-> REASON: (Stat had hard requirement)")
						UTTERFAILURE = true
					}
				}
			}

			if len(skills_met) == 0 {
				fmt.Println("Failed(stat+skill): ", Feat["Feat"], "-> REASON: (No skills were met)")
				UTTERFAILURE = true
			}

			var skills_are_failed bool = false
			if len(skills_failed) > 0 {
				skills_are_failed = true
			}

			/*
				doing this because technically the skill_is_required == "false" isn't infallible.
				If I don't do this check, if the hardrequirement is within the skills_failed and all skills aren't (technically) required, it will bypass the check coming up
			*/
			if hard_req, ok := Feat["HardRequirement"]; ok && skills_are_failed {
				for g := 0; g < len(skills_failed); g++ {
					skill := skills_failed[g]
					if skill == hard_req {
						fmt.Println("Failed(stat+skill): ", Feat["Feat"], "-> REASON: (Skill had hard requirement: ", hard_req, ")")
						UTTERFAILURE = true
					}
				}
			}

			if skills_are_failed {
				for h := 0; h < len(skills_failed); h++ {
					var is_ok bool = false
					var skill_is_required string
					if skill_required, ok := Feat["NeedsAllSkills"]; ok {
						is_ok = true
						skill_is_required = skill_required
					}
					if !is_ok {
						fmt.Println("Failed(stat+skill): ", Feat["Feat"], "-> REASON: (Failed: Didn't have required JSON data)")
						UTTERFAILURE = true
					}
					var check_for_hard_requirement bool = false
					if skill_is_required == "false" {
						check_for_hard_requirement = true
					} else if skill_is_required == "true" {
						fmt.Println("Failed(stat+skill): ", Feat["Feat"], "-> REASON: (Failed skill was required(all were required))")
						UTTERFAILURE = true
					}

					if check_for_hard_requirement {
						if hard_requirement, ok := Feat["HardRequirement"]; ok && hard_requirement == skills_failed[h] {
							fmt.Println("Failed(stat+skill): ", Feat["Feat"], "-> REASON: (Failed on hard required skill)")
							UTTERFAILURE = true
						}
					}
				}
			}
			if !UTTERFAILURE {
				var matched_slice []bool
				if len(data.Character_Type) > 0 && Feat["Type"] != "" {
					for i := 0; i < len(data.Character_Type); i++ {
						var matches bool = Check_For_Type(Feat["Type"], data.Character_Type[i])
						if matches {
							matched_slice = append(matched_slice, matches)
						}
					}
					if len(matched_slice) > 0 {
						fmt.Println(Feat["Feat"], " ", "Checked for no failures on STAT+SKILL")
						*allAllocatedFeats = append(*allAllocatedFeats, Feat["Feat"])
					}
				} else {
					fmt.Println(Feat["Feat"], " ", "Checked for no failures on STAT+SKILL")
					*allAllocatedFeats = append(*allAllocatedFeats, Feat["Feat"])
				}
			}

		} else if len(stats_required) > 0 && len(skills_required) == 0 {
			var stats_met = []string{}
			var stats_failed = []string{}
			var UTTERFAILURE bool = false

			err := Check_Stats(stats_required, &stats_met, &stats_failed)
			if err != nil {
				fmt.Println("Errored: ", err)
				return err
			}
			if len(stats_met) == 0 {
				fmt.Println("Failed(stat): ", Feat["Feat"], "-> REASON: (No stats met)")
				UTTERFAILURE = true
			}
			var stats_are_failed bool = false
			if len(stats_failed) > 0 {
				stats_are_failed = true
			}

			if stats_are_failed {
				var required []string
				for l := 0; l < len(stats_failed); l++ {
					if hard_req, ok := Feat["NeedsAllStats"]; ok {
						required = append(required, hard_req)
					}
				}

				for j := 0; j < len(required); j++ {
					if required[j] == "true" {
						fmt.Println("Failed(stat): ", Feat["Feat"], "-> REASON: (Stat had hard requirement)")
						UTTERFAILURE = true
					}
				}
			}

			if !UTTERFAILURE {
				var matched_slice []bool
				if len(data.Character_Type) > 0 && Feat["Type"] != "" {
					for i := 0; i < len(data.Character_Type); i++ {
						var matches bool = Check_For_Type(Feat["Type"], data.Character_Type[i])
						if matches {
							matched_slice = append(matched_slice, matches)
						}
					}

					if len(matched_slice) > 0 {
						fmt.Println(Feat["Feat"], " ", "Checked for no failures on STAT")
						*allAllocatedFeats = append(*allAllocatedFeats, Feat["Feat"])
					}
				} else {
					fmt.Println(Feat["Feat"], " ", "Checked for no failures on STAT")
					*allAllocatedFeats = append(*allAllocatedFeats, Feat["Feat"])
				}
			}

		} else if len(skills_required) > 0 && len(stats_required) == 0 {
			var skills_met = []string{}
			var skills_failed = []string{}
			var UTTERFAILURE bool = false

			err := Check_Skills(skills_required, &skills_met, &skills_failed)
			if err != nil {
				return err
			}
			if len(skills_met) == 0 {
				fmt.Println("Failed(skill): ", Feat["Feat"], "-> REASON: (No skills were met)")
				UTTERFAILURE = true
			}
			var skills_are_failed bool = false
			if len(skills_failed) > 0 {
				skills_are_failed = true
			}

			/*
				doing this because technically the skill_is_required == "false" isn't infallible.
				If I don't do this check, if the hardrequirement is within the skills_failed and all skills aren't (technically) required, it will bypass the check coming up
			*/
			if hard_req, ok := Feat["HardRequirement"]; ok && skills_are_failed {
				for g := 0; g < len(skills_failed); g++ {
					skill := skills_failed[g]
					if skill == hard_req {
						fmt.Println("Failed(skill): ", Feat["Feat"], "-> REASON: (Skill had hard requirement: ", hard_req, ")")
						UTTERFAILURE = true
					}
				}
			}

			if skills_are_failed {
				for h := 0; h < len(skills_failed); h++ {
					var is_ok bool = false
					var skill_is_required string
					if skill_required, ok := Feat["NeedsAllSkills"]; ok {
						is_ok = true
						skill_is_required = skill_required
					}
					if !is_ok {
						fmt.Println("Failed(skill): ", Feat["Feat"], "-> REASON: (Failed: Didn't have required JSON data)")
						UTTERFAILURE = true
					}
					var check_for_hard_requirement bool = false
					if skill_is_required == "false" {
						check_for_hard_requirement = true
					} else if skill_is_required == "true" {
						fmt.Println("Failed(skill): ", Feat["Feat"], "-> REASON: (Failed skill was required(all were required))")
						UTTERFAILURE = true
					}

					if check_for_hard_requirement {
						if hard_requirement, ok := Feat["HardRequirement"]; ok && hard_requirement == skills_failed[h] {
							fmt.Println("Failed(skill): ", Feat["Feat"], "-> REASON: (Failed on hard required skill)")
							UTTERFAILURE = true
						}
					}
				}
			}

			if !UTTERFAILURE {
				var matched_slice []bool
				if len(data.Character_Type) > 0 && Feat["Type"] != "" {
					for i := 0; i < len(data.Character_Type); i++ {
						var matches bool = Check_For_Type(Feat["Type"], data.Character_Type[i])
						if matches {
							matched_slice = append(matched_slice, matches)
						}
					}
					if len(matched_slice) > 0 {
						fmt.Println(Feat["Feat"], " ", "Checked for no failures on SKILL")
						*allAllocatedFeats = append(*allAllocatedFeats, Feat["Feat"])
					}
				} else {
					fmt.Println(Feat["Feat"], " ", "Checked for no failures on SKILL")
					*allAllocatedFeats = append(*allAllocatedFeats, Feat["Feat"])
				}
			}
		}
		fmt.Println()
	}
	return nil
}
