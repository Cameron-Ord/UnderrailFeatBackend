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

func checkSkillStat(feats []map[string]string, data RequestData, allAllocatedFeats *[]string) error {
	fmt.Println(" ")
	fmt.Println("////////////////////////")
	fmt.Println("Checking skills and stats...")
	fmt.Println("////////////////////////")
	fmt.Println(" ")

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
					return err
				}
				givenStatValue, err := convertToInt(statValue)
				if err != nil {
					return err
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
							return err
						}
						givenSkillValue, err := convertToInt(skillValue)
						if err != nil {
							return err
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
						*allAllocatedFeats = append(*allAllocatedFeats, feat["Feat"])
						fmt.Println("-----------------------------")
						fmt.Println("SKILL/STAT = {MET}", feat["Feat"], "Appended at iteration:", "->", i, "noFails:", noFails)
						fmt.Println("-----------------------------")
					}
				}
			}
		}
	}
	return nil
}

func checkStat(feats []map[string]string, data RequestData, allAllocatedFeats *[]string) error {
	fmt.Println(" ")
	fmt.Println("////////////////////////")
	fmt.Println("Checking stats...")
	fmt.Println("////////////////////////")
	fmt.Println(" ")
	//looping using the length of the feats map
	for i := 0; i < len(feats); i++ {
		//assigning each index of the feat map to a variable
		feat := feats[i]
		//declaring variables
		var hasSkillRequirement bool = false
		var noFails bool = true
		var statsMet = []string{}
		var statsFailed = []string{}
		//looping through the length of the Stats slice (sent from frontend)
		for j := 0; j < len(data.Stats); j++ {
			//assigning each index of the Stats slice to a variable
			stat := data.Stats[j]
			//copying over the statName to its own variable
			var statName string = stat.StatName
			//doing the same, but with the value of the stat received front the frontend
			var statValue string = stat.StatValue
			//if ok = true, this index of feats contains the statName(data.Stats contains every stat so by looping it, this if check is being checked until every stat has been thrown at it)
			if requirement, ok := feat[statName]; ok {
				//assigning the value of the stat requirement listed in the feat variable
				var featRequirement string = requirement
				//performing a string to int conversion of the value
				statRequirement, err := convertToInt(featRequirement)
				if err != nil {

					return err
				}
				//performing a string to int conversion of the value assigned at the start of the loop
				givenStatValue, err := convertToInt(statValue)
				if err != nil {

					return err
				}
				//performing some simple checks
				if givenStatValue >= statRequirement {
					statsMet = append(statsMet, statName)
				} else if givenStatValue < statRequirement {
					statsFailed = append(statsFailed, statName)
				}
			}
		}
		/*if the statsMet slice is populated in this iteration
		start looping through the Skills slice to check for existing skill requirements
		for this index of feats. The exact same as the if check from before, really.
		*/

		if len(statsMet) > 0 {
			for k := 0; k < len(data.Skills); k++ {
				skill := data.Skills[k]
				var skillName string = skill.SkillName
				if _, ok := feat[skillName]; ok {
					//if the skillname of is found in the current index of feats, breaks the loop to prevent further checking as it is no longer necessary
					//hasSkillRequirement is set to true
					hasSkillRequirement = true
					break
				} else {
					hasSkillRequirement = false
				}
			}
			// if the statsFailed slice is populated and there is no skill requirement
			if len(statsFailed) > 0 && !hasSkillRequirement {
				//range looping the statsFailed slice
				for _, failed := range statsFailed {
					fmt.Println("STAT = {FAILED}:", "({"+failed, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[failed], failed)
					//checking if every stat is required for the skill to be obtainable
					if required, ok := feat["NeedsAllStats"]; ok {
						if required == "true" {
							//if every skill listed in the feat is required, breaks and sets noFails to false
							noFails = false
							break
						} else if required == "false" {
							//otherwise, we continue on
							noFails = true
						}
					}
				}
			}
			//finally, if there are no fails, we append the feat to the allAllocatedFeats slice using a ptr.
			if noFails == true {
				*allAllocatedFeats = append(*allAllocatedFeats, feat["Feat"])
				fmt.Println("STAT = {MET}:", feat["Feat"], "Appended at iteration:", "->", i, "noFails:", noFails, "hasSkillRequirement:", hasSkillRequirement)
			}
		}
	}
	//returning nil, as the function reached the end without errors.
	return nil
}

func convertToInt(givenNumStr string) (int, error) {
	NumToInt, err := strconv.Atoi(givenNumStr)
	return NumToInt, err
}

func checkSkill(feats []map[string]string, data RequestData, allAllocatedFeats *[]string) error {
	fmt.Println(" ")
	fmt.Println("////////////////////////")
	fmt.Println("Checking skills...")
	fmt.Println("////////////////////////")
	fmt.Println(" ")

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
					return err
				}
				//runs the conversion function and assigns returned values
				givenSkillValue, err := convertToInt(skillValue)
				if err != nil {
					return err
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
					*allAllocatedFeats = append(*allAllocatedFeats, feat["Feat"])
					fmt.Println("SKILL = {MET}", feat["Feat"], "Appended at iteration:", "->", i, "hasStatRequirement:", hasStatRequirement)
				}
			}
		}
	}
	//returning nil at the end of the function, as no error was raised
	return nil
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
	var allAllocatedFeats []string
	var err error
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

	err = checkSkill(feats, data, &allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	//checking feats that require only stats
	err = checkStat(feats, data, &allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	//checking feats that require stats and skills only
	err = checkSkillStat(feats, data, &allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	//gathering all the feats that got assigned from each check into a single slice
	//marshalling the data to JSON format using a return function
	jsonData, err := marshalData(allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	//returning the json data
	return jsonData, nil
}
