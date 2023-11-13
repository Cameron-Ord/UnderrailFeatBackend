package calculation

import (
	"encoding/json"
	"fmt"
	"strconv"

	"log"
	"os"
)

type Stat struct {
	StatName  string `json:"statName"`
	StatValue string `json:"statValue"`
}

type Skill struct {
	SkillName  string `json:"skillName"`
	SkillValue string `json:"skillValue"`
}

type RequestData struct {
	Stats  []Stat  `json:"stats"`
	Skills []Skill `json:"skills"`
}

func checkSkillStat(feats []map[string]string, data RequestData) ([]string, error) {

	var StatSkillFeats = []string{}

	for i := 0; i < len(feats); i++ {
		feat := feats[i]
		var hasSkillRequirement = false
		var noFails bool = true
		var statsMet = []string{}
		var statsFailed = []string{}
		var skillsMet = []string{}
		var skillsFailed = []string{}
		for j := 0; j < len(data.Stats); j++ {
			stat := data.Stats[j]
			var statName string = stat.StatName
			var statValue string = stat.StatValue
			if requirement, ok := feat[statName]; ok {
				var featRequirement string = requirement
				statRequirement, err := convertToInt(featRequirement)
				if err != nil {
					return nil, err
				}
				givenStatValue, err := convertToInt(statValue)
				if err != nil {
					return nil, err
				}
				if givenStatValue >= statRequirement {
					statsMet = append(statsMet, statName)
				} else if givenStatValue < statRequirement {
					statsFailed = append(statsFailed, statName)
				}
			}
		}
		if len(statsMet) > 0 {
			for k := 0; k < len(data.Skills); k++ {
				skill := data.Skills[k]
				var skillName string = skill.SkillName
				if _, ok := feat[skillName]; ok {
					hasSkillRequirement = true
					break
				}
			}
			if hasSkillRequirement {
				for l := 0; l < len(data.Skills); l++ {
					skill := data.Skills[l]
					var skillName string = skill.SkillName
					var skillValue string = skill.SkillValue
					if requirement, ok := feat[skillName]; ok {
						var featRequirement string = requirement
						skillRequirement, err := convertToInt(featRequirement)
						if err != nil {
							return nil, err
						}
						givenSkillValue, err := convertToInt(skillValue)
						if err != nil {
							return nil, err
						}
						if givenSkillValue >= skillRequirement {
							skillsMet = append(skillsMet, skillName)
						} else if givenSkillValue < skillRequirement {
							skillsFailed = append(skillsFailed, skillName)
						}
					}
				}
			}

			if len(statsFailed) > 0 {
				for _, failed := range statsFailed {
					if _, ok := feat[failed]; ok {
						noFails = false
					}
				}
			}
			if noFails == true {
				if len(skillsMet) > 0 {
					StatSkillFeats = append(StatSkillFeats, feat["Feat"])
				}
			}
		}
	}

	return StatSkillFeats, nil
}

func checkStat(feats []map[string]string, data RequestData) ([]string, error) {

	var StatFeats = []string{}

	for i := 0; i < len(feats); i++ {
		feat := feats[i]
		var hasSkillRequirement bool = false
		var noFails bool = true
		var statsMet = []string{}
		var statsFailed = []string{}
		for j := 0; j < len(data.Stats); j++ {
			stat := data.Stats[j]
			var statName string = stat.StatName
			var statValue string = stat.StatValue
			if requirement, ok := feat[statName]; ok {
				var featRequirement string = requirement
				statRequirement, err := convertToInt(featRequirement)
				if err != nil {

					return nil, err
				}
				givenStatValue, err := convertToInt(statValue)
				if err != nil {

					return nil, err
				}
				if givenStatValue >= statRequirement {
					statsMet = append(statsMet, statName)
				} else if givenStatValue < statRequirement {
					statsFailed = append(statsFailed, statName)
				}
			}
		}

		if len(statsMet) > 0 {
			for k := 0; k < len(data.Skills); k++ {
				skill := data.Skills[k]
				var skillName string = skill.SkillName
				if _, ok := feat[skillName]; ok {
					hasSkillRequirement = true
					break
				} else {
					hasSkillRequirement = false
				}
			}
			if len(statsFailed) > 0 {
				// this is kind of unnecessary as this slice being populated alone is enough to make this fail, but hey, why not
				for _, failed := range statsFailed {
					if _, ok := feat[failed]; ok {
						noFails = false
					}
				}
			}
			if noFails == true {
				if !hasSkillRequirement {
					StatFeats = append(StatFeats, feat["Feat"])
				}
			}
		}
	}
	return StatFeats, nil
}

func convertToInt(givenNumStr string) (int, error) {
	NumToInt, err := strconv.Atoi(givenNumStr)
	return NumToInt, err
}

func checkSkill(feats []map[string]string, data RequestData) ([]string, error) {

	var SkillFeats = []string{}

	for i := 0; i < len(feats); i++ {
		feat := feats[i]
		var hasStatRequirement bool = false
		var skillsMet = []string{}
		var skillsFailed = []string{}

		for j := 0; j < len(data.Skills); j++ {
			skill := data.Skills[j]
			var skillName string = skill.SkillName
			var skillValue string = skill.SkillValue
			if requirement, ok := feat[skillName]; ok {
				var featRequirement string = requirement
				skillRequirement, err := convertToInt(featRequirement)
				if err != nil {
					return nil, err
				}
				givenSkillValue, err := convertToInt(skillValue)
				if err != nil {
					return nil, err
				}
				if givenSkillValue >= skillRequirement {
					skillsMet = append(skillsMet, skillName)
				} else if givenSkillValue < skillRequirement {
					skillsFailed = append(skillsFailed, skillName)
				}
			}
		}
		if len(skillsMet) > 0 {
			for k := 0; k < len(data.Stats); k++ {
				stat := data.Stats[k]
				var statName string = stat.StatName
				if _, ok := feat[statName]; ok {
					hasStatRequirement = true
					break
				} else {
					hasStatRequirement = false
				}
			}
			if len(skillsFailed) > 0 {
				for _, failed := range skillsFailed {
					fmt.Println("failed: ", failed)
				}
			}
			if !hasStatRequirement {
				SkillFeats = append(SkillFeats, feat["Feat"])
			}
		}

	}
	return SkillFeats, nil
}

func unloadJson() []map[string]string {
	jsonFile, err := os.ReadFile("jsonData/FeatsData.json")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	var feats []map[string]string
	err = json.Unmarshal(jsonFile, &feats)
	if err != nil {
		log.Fatal("Error decoding JSON data:", err)
	}
	return feats
}

func marshalData(sentSlice []string) ([]byte, error) {
	jsonData, err := json.Marshal(sentSlice)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func PrepareData(data RequestData) ([]byte, error) {
	var feats []map[string]string
	feats = unloadJson()
	SkillFeats, err := checkSkill(feats, data)
	if err != nil {
		return nil, err
	}

	StatFeats, err := checkStat(feats, data)
	if err != nil {
		return nil, err
	}

	StatSkillFeats, err := checkSkillStat(feats, data)
	if err != nil {
		return nil, err
	}

	var allAllocatedFeats []string
	allAllocatedFeats = append(StatFeats, SkillFeats...)
	allAllocatedFeats = append(allAllocatedFeats, StatSkillFeats...)
	jsonData, err := marshalData(allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
