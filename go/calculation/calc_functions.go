package calculation

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

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
		return err
	}

	if given_value >= required_value {
		*metPtr = append(*metPtr, name)
	} else if given_value < required_value {
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
		Feat := Feats[i]
		var stats_required = []Stats_Tracker{}
		var skills_required = []Skills_Tracker{}
		Get_Requirements(Feat, data, &stats_required, &skills_required)
		if len(stats_required) > 0 && len(skills_required) > 0 {
			var stats_met = []string{}
			var stats_failed = []string{}
			var skills_met = []string{}
			var skills_failed = []string{}
			var no_fails bool = true
			err := Check_All(
				stats_required, skills_required, &stats_met,
				&stats_failed, &skills_met, &skills_failed,
			)
			if err != nil {
				return err
			}

			if len(stats_met) > 0 {
				if len(stats_failed) > 0 {
					for l := 0; l < len(stats_failed); l++ {
						if hard_req, ok := Feat["NeedsAllStats"]; ok {
							if hard_req == "true" {
								no_fails = false
								break
							} else if hard_req == "false" {
								no_fails = true
							}
						}
					}
				}
				if len(skills_met) > 0 {
					if hard_req, ok := Feat["HardRequirement"]; ok {
						for g := 0; g < len(skills_met); g++ {
							skill := skills_met[g]
							if skill != hard_req {
								no_fails = true
								break
							} else {
								no_fails = false
							}
						}
					} else {
						no_fails = true
					}

					if no_fails && len(skills_failed) > 0 {
						for h := 0; h < len(skills_failed); h++ {
							if skill_required, ok := Feat["NeedsAllSkills"]; ok {
								if skill_required == "false" {
									if hard_requirement, ok := Feat["HardRequirement"]; ok && hard_requirement == skills_failed[h] {
										no_fails = false
										break
									} else {
										no_fails = true
									}
								} else if skill_required == "true" {
									no_fails = false
									break
								}
							}
						}
					}
				}

				if len(skills_met) == 0 {
					no_fails = false
				}

				if no_fails {
					*allAllocatedFeats = append(*allAllocatedFeats, Feat["Feat"])
				}
			}
		} else if len(stats_required) > 0 && len(skills_required) == 0 {
			var stats_met = []string{}
			var stats_failed = []string{}
			var no_fails bool = true
			err := Check_Stats(stats_required, &stats_met, &stats_failed)
			if err != nil {
				return err
			}
			if len(stats_met) == 0 {
				no_fails = false
			}

			if len(stats_met) > 0 {
				if len(stats_failed) > 0 {
					for l := 0; l < len(stats_failed); l++ {
						if hard_req, ok := Feat["NeedsAllStats"]; ok {
							if hard_req == "true" {
								no_fails = false
								break
							} else if hard_req == "false" {
								no_fails = true
							}
						}
					}
				}
				if no_fails {
					*allAllocatedFeats = append(*allAllocatedFeats, Feat["Feat"])
				}
			}

		} else if len(skills_required) > 0 && len(stats_required) == 0 {
			var skills_met = []string{}
			var skills_failed = []string{}
			var no_fails bool = true
			err := Check_Skills(skills_required, &skills_met, &skills_failed)
			if err != nil {
				return err
			}
			if len(skills_met) == 0 {
				no_fails = false
			}

			if len(skills_met) > 0 {
				if hard_req, ok := Feat["HardRequirement"]; ok {
					for g := 0; g < len(skills_met); g++ {
						skill := skills_met[g]
						if skill != hard_req {
							no_fails = true
							break
						} else {
							no_fails = false
						}
					}
				} else {
					no_fails = true
				}

				if no_fails && len(skills_failed) > 0 {
					for h := 0; h < len(skills_failed); h++ {
						if skill_required, ok := Feat["NeedsAllSkills"]; ok {
							if skill_required == "false" {
								if hard_requirement, ok := Feat["HardRequirement"]; ok && hard_requirement == skills_failed[h] {
									no_fails = false
									break
								} else {
									no_fails = true
								}
							} else if skill_required == "true" {
								no_fails = false
								break
							}
						}
					}
				}
				if no_fails {
					*allAllocatedFeats = append(*allAllocatedFeats, Feat["Feat"])
				}
			}
		}
	}
	return nil
}
