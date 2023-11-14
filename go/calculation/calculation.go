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
	fmt.Println(" ")
	fmt.Println("----------------------")
	fmt.Println("Checking skills and stats...")
	fmt.Println("----------------------")
	fmt.Println(" ")
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

			if len(statsFailed) > 0 {
				for _, statfail := range statsFailed {
					fmt.Println("SKILL/STAT = {FAILED}:", "({"+statfail, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[statfail], statfail)
					if statrequire, ok := feat["NeedsAllStats"]; ok {
						if statrequire == "true" {
							noFails = false
							break
						} else if statrequire == "false" {
							noFails = true
						}
					}
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

				if len(skillsMet) > 0 {
					for _, skillmet := range skillsMet {
						if hardreq, ok := feat["HardRequirement"]; ok {
							if skillmet != hardreq {
								fmt.Println("breaking")
								break
							} else {
								noFails = false
							}
						} else {
							break
						}
					}
					if noFails && len(skillsFailed) > 0 {
						for _, failed := range skillsFailed {
							fmt.Println("SKILL = {FAILED}:", "({"+failed, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[failed], failed)
							if skillrequire, ok := feat["NeedsAllSkills"]; ok {
								if skillrequire == "false" {
									if requirement, ok := feat["HardRequirement"]; ok && requirement == failed {
										noFails = false
										fmt.Println("failed break")
										break
									}
								} else if skillrequire == "true" {
									noFails = false
									fmt.Println("failed break")
									break
								}
							}
						}
						if noFails {
							StatSkillFeats = append(StatSkillFeats, feat["Feat"])
							fmt.Println("SKILL/STAT = {MET}", feat["Feat"], "Appended at iteration:", "->", i, "noFails:", noFails)
						}
					}

				}
			}

		}
	}

	return StatSkillFeats, nil
}

func checkStat(feats []map[string]string, data RequestData) ([]string, error) {
	fmt.Println(" ")
	fmt.Println("----------------------")
	fmt.Println("Checking stats...")
	fmt.Println("----------------------")
	fmt.Println(" ")
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
					fmt.Println("STAT = {FAILED}:", "({"+failed, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[failed], failed)
					if required, ok := feat["NeedsAllStats"]; ok {
						if required == "true" {
							noFails = false
							break
						} else if required == "false" {
							noFails = true
						}
					}
				}
			}
			if noFails == true {
				if !hasSkillRequirement {
					StatFeats = append(StatFeats, feat["Feat"])
					fmt.Println("STAT = {MET}:", feat["Feat"], "Appended at iteration:", "->", i, "noFails:", noFails, "hasSkillRequirement:", hasSkillRequirement)
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
	fmt.Println(" ")
	fmt.Println("----------------------")
	fmt.Println("Checking skills...")
	fmt.Println("----------------------")
	fmt.Println(" ")
	var SkillFeats = []string{}

	for i := 0; i < len(feats); i++ {
		feat := feats[i]
		var hasStatRequirement bool = false
		var noFails bool = true
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
			if !hasStatRequirement {
				for _, skillmet := range skillsMet {
					if hardreq, ok := feat["HardRequirement"]; ok {
						if skillmet != hardreq {
							fmt.Println("breaking")
							break
						} else {
							noFails = false
						}
					} else {
						break
					}
				}

				if noFails && len(skillsFailed) > 0 {
					for _, failed := range skillsFailed {
						fmt.Println("SKILL = {FAILED}:", "({"+failed, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[failed], failed)
						if skillrequire, ok := feat["NeedsAllSkills"]; ok {
							if skillrequire == "false" {
								if requirement, ok := feat["HardRequirement"]; ok && requirement == failed {
									noFails = false
									fmt.Println("failed break")
									break
								}
							} else if skillrequire == "true" {
								noFails = false
								fmt.Println("failed break")
								break
							}
						}
					}
				}
				if noFails {
					SkillFeats = append(SkillFeats, feat["Feat"])
					fmt.Println("SKILL = {MET}", feat["Feat"], "Appended at iteration:", "->", i, "hasStatRequirement:", hasStatRequirement)
				}
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
	fmt.Println(" ")
	fmt.Println("----------------------")
	fmt.Println("Starting calculation..")
	fmt.Println("----------------------")
	fmt.Println(" ")

	var feats []map[string]string
	feats = unloadJson()
	SkillFeats, err := checkSkill(feats, data)
	if err != nil {
		return nil, err
	}
	fmt.Println(" ")
	fmt.Println("..Done")

	StatFeats, err := checkStat(feats, data)
	if err != nil {
		return nil, err
	}
	fmt.Println(" ")
	fmt.Println("..Done")
	StatSkillFeats, err := checkSkillStat(feats, data)
	if err != nil {
		return nil, err
	}
	fmt.Println(" ")
	fmt.Println("..Done")

	fmt.Println(" ")
	fmt.Println("----------------------")
	fmt.Println("Appending...")
	fmt.Println("----------------------")
	fmt.Println(" ")
	var allAllocatedFeats []string
	allAllocatedFeats = append(StatFeats, SkillFeats...)
	allAllocatedFeats = append(allAllocatedFeats, StatSkillFeats...)
	jsonData, err := marshalData(allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	fmt.Println(" ")
	fmt.Println("..Done")
	return jsonData, nil
}
