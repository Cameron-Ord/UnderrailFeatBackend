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
	fmt.Println("////////////////////////")
	fmt.Println("Checking skills and stats...")
	fmt.Println("////////////////////////")
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
				} else {
					hasSkillRequirement = false
				}
			}

			if len(statsFailed) > 0 {
				for _, statfail := range statsFailed {
					fmt.Println("-----------------------------")
					fmt.Println("SKILL/STAT = {FAILED}:", "({"+statfail, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[statfail], statfail)
					fmt.Println("-----------------------------")

					if statrequire, ok := feat["NeedsAllStats"]; ok {
						if statrequire == "true" {
							fmt.Println(" ")
							fmt.Println(feat["Feat"], "{FAILED AT STATREQUIRE}:", "Cause ->", statfail)
							fmt.Println(" ")
							noFails = false
							break
						} else if statrequire == "false" {
							noFails = true
						}
					}
				}
			}

			if hasSkillRequirement && noFails {
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
					if hardreq, ok := feat["HardRequirement"]; ok {
						for _, skillmet := range skillsMet {
							if skillmet != hardreq {
								fmt.Println(" ")
								fmt.Println("{Breaking at hardreq check}:", "->", feat["Feat"])
								fmt.Println(" ")
								noFails = true
								break
							} else {
								noFails = false
							}
						}
					} else {
						noFails = true
					}

					if noFails && len(skillsFailed) > 0 {
						for _, failed := range skillsFailed {
							fmt.Println("-----------------------------")
							fmt.Println("SKILL = {FAILED}:", "({"+failed, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[failed], failed)
							fmt.Println("-----------------------------")
							if skillrequire, ok := feat["NeedsAllSkills"]; ok {
								if skillrequire == "false" {
									if requirement, ok := feat["HardRequirement"]; ok && requirement == failed {
										noFails = false
										fmt.Println(" ")
										fmt.Println("{Failed break}:", "->", feat["Feat"])
										fmt.Println(" ")
										break
									} else {
										noFails = true
									}
								} else if skillrequire == "true" {
									noFails = false
									fmt.Println(" ")
									fmt.Println("{Failed break}:", "->", feat["Feat"])
									fmt.Println(" ")
									break
								}
							}
						}

					}
					if noFails {
						StatSkillFeats = append(StatSkillFeats, feat["Feat"])
						fmt.Println("-----------------------------")
						fmt.Println("SKILL/STAT = {MET}", feat["Feat"], "Appended at iteration:", "->", i, "noFails:", noFails)
						fmt.Println("-----------------------------")
					}
				}
			}
		}
	}
	return StatSkillFeats, nil
}

func checkStat(feats []map[string]string, data RequestData) ([]string, error) {
	fmt.Println(" ")
	fmt.Println("////////////////////////")
	fmt.Println("Checking stats...")
	fmt.Println("////////////////////////")
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
	fmt.Println("////////////////////////")
	fmt.Println("Checking skills...")
	fmt.Println("////////////////////////")
	fmt.Println(" ")
	var SkillFeats = []string{}
	//looping for the length of the feats map
	for i := 0; i < len(feats); i++ {
		//assigning individual feat by index
		feat := feats[i]
		//declaring some variables
		var hasStatRequirement bool = false
		var noFails bool = true
		var skillsMet = []string{}
		var skillsFailed = []string{}
		//looping for the length of the skills sent from frontend
		for j := 0; j < len(data.Skills); j++ {
			//assigning individual skill from index in order to access the skillName and skillValue
			skill := data.Skills[j]
			var skillName string = skill.SkillName
			var skillValue string = skill.SkillValue
			//if the feat contains the current skillname(from index), sets the value from feat to the requirement variable if ok gets set to true
			if requirement, ok := feat[skillName]; ok {
				var featRequirement string = requirement
				//runs the conversion function and assigns returned values
				skillRequirement, err := convertToInt(featRequirement)
				if err != nil {
					return nil, err
				}
				//runs the conversion function and assigns returned values
				givenSkillValue, err := convertToInt(skillValue)
				if err != nil {
					return nil, err
				}
				//performs the checks, assigns to a respective slice depending on the outcome
				if givenSkillValue >= skillRequirement {
					skillsMet = append(skillsMet, skillName)
				} else if givenSkillValue < skillRequirement {
					skillsFailed = append(skillsFailed, skillName)
				}
			}
		}
		//if there are skills met
		if len(skillsMet) > 0 {
			//looping through stats sent from frontend
			for k := 0; k < len(data.Stats); k++ {
				//assiging variable by index
				stat := data.Stats[k]
				var statName string = stat.StatName
				//if the feat contains any of the stats, it will break the loop and set the hasStatRequirement variable to true
				if _, ok := feat[statName]; ok {
					hasStatRequirement = true
					break
				} else {
					hasStatRequirement = false
				}
			}
			//if the hasStatRequirement bool is false, continues
			if !hasStatRequirement {
				/*
					using a range loop to go through skillsMet
					the purpose of this is to make sure that in the case a feat has a hardrequirement,
					it will ensure that the other skills that are associated with it didn't get left out.
					Otherwise, with the following code, since the hardrequirement is literally required and all other skills aren't(You just need atleast one),
					it will return the feat despite none of the other skills associated with the feat being present
					So this is just a simple check in that case, if ok gets assigned false here it will just skip the block entirely
				*/
				if hardreq, ok := feat["HardRequirement"]; ok {
					for _, skillmet := range skillsMet {
						if skillmet != hardreq {
							fmt.Println(skillmet, "->", "{Exists}", "Continuing..")
							noFails = true
							break
						} else {
							noFails = false
						}
					}
					//skips the block entirely if ok is false
				} else {
					fmt.Println("Hardreq -> {SKIPPED}")
					noFails = true
				}
				//if noFails is still true and the skillsFailed slice is populated..
				if noFails && len(skillsFailed) > 0 {
					for _, failed := range skillsFailed {
						fmt.Println("SKILL = {FAILED}:", "({"+failed, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[failed], failed)
						//if ok is true here, it will then assign and check the NeedsAllSkills value
						if skillrequire, ok := feat["NeedsAllSkills"]; ok {
							//if the value is a string representing "false"
							if skillrequire == "false" {
								/*
									if ok is true, assigns the hardrequirement value to requirement then checks if the string equals the failed string.
									if the requirement is equal to the failed skill, it will break the loop and set the noFails bool to false, otherwise, stays true.
									the true statements are a little reduntant but its clear what it does atleast.
								*/
								if requirement, ok := feat["HardRequirement"]; ok && requirement == failed {
									noFails = false
									fmt.Println("failed break")
									break
								} else {
									noFails = true
								}
								//if the feat does require all skills, then it's a flatout fail for this one chief
							} else if skillrequire == "true" {
								noFails = false
								fmt.Println("failed break")
								break
							}
						}
					}
				}
				//finally, if there are no fails for this iteration of Feats, it will append the feat
				if noFails {
					SkillFeats = append(SkillFeats, feat["Feat"])
					fmt.Println("SKILL = {MET}", feat["Feat"], "Appended at iteration:", "->", i, "hasStatRequirement:", hasStatRequirement)
				}
			}
		}
	}
	//returning the populated slice
	return SkillFeats, nil
}

func unloadJson() []map[string]string {
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

func marshalData(sentSlice []string) ([]byte, error) {
	//assigning jsondata from the sentslice(the amalgamation of the results from all the checks)
	jsonData, err := json.Marshal(sentSlice)
	if err != nil {
		return nil, err
	}
	//returning the data
	return jsonData, nil
}

func PrepareData(data RequestData) ([]byte, error) {

	// init the calculation
	fmt.Println(" ")
	fmt.Println("////////////////////////")
	fmt.Println("Starting calculation..")
	fmt.Println("////////////////////////")
	fmt.Println(" ")
	//declare the variable to hold the featsdata json file
	var feats []map[string]string
	//assigning the data with a return function
	feats = unloadJson()
	//checking feats that require only skills
	SkillFeats, err := checkSkill(feats, data)
	if err != nil {
		return nil, err
	}
	//checking feats that require only stats
	StatFeats, err := checkStat(feats, data)
	if err != nil {
		return nil, err
	}
	//checking feats that require stats and skills only
	StatSkillFeats, err := checkSkillStat(feats, data)
	if err != nil {
		return nil, err
	}
	//gathering all the feats that got assigned from each check into a single slice
	var allAllocatedFeats []string
	allAllocatedFeats = append(StatFeats, SkillFeats...)
	allAllocatedFeats = append(allAllocatedFeats, StatSkillFeats...)
	//marshalling the data to JSON format using a return function
	jsonData, err := marshalData(allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	//returning the json data
	return jsonData, nil
}
