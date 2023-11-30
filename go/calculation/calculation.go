// part of package calculation
package calculation

import (
	"encoding/json"
	"fmt"
	"strconv"

	"log"
	"os"
)

// directing the unmarshalling and allocation of the data sent from the frontend

type Stat struct {
	StatName  string `json:"statName"`
	StatValue int    `json:"statValue"`
}

type Skill struct {
	SkillName  string `json:"skillName"`
	SkillValue int    `json:"skillValue"`
}

type RequestData struct {
	Stats  []Stat  `json:"stats"`
	Skills []Skill `json:"skills"`
}

/*
START OF FILE:
	This function is being called from main.go
	Here we are simply declaring 2 variables, an error and a slice
	After that we call the unloadJson function which then reads from the FeatsData.json to
	create a map out of the feats containing json data.
	After the map has been populated we begin running each subsequent function
	Once each function has been executed the allAllocatedFeats slice gets marshalled into JSON,
	which we then return the json in type []byte.
*/

func PrepareData(data RequestData) ([]byte, error) {
	var allAllocatedFeats []string
	var err error
	var feats []map[string]string
	fmt.Println(" ")
	fmt.Println("////////////////////////")
	fmt.Println("Starting calculation..")
	fmt.Println("////////////////////////")
	fmt.Println(" ")

	feats = unloadJson()

	err = checkSkill(feats, data, &allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	err = checkStat(feats, data, &allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	err = checkSkillStat(feats, data, &allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	jsonData, err := marshalData(allAllocatedFeats)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
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

/*
	Going from here we get into the main functions that append data to the allAllocatedFeats pointer
	Essentially these are just performing a linear search looping through the unordered feats map,
	performing moderately complex checks for each iteration of Feats. The Feats map isn't very large;
	only a few hundred or so indexes.

	There are 3 functions that do a similar operation as the one immediately below.
	I split them up as there are a ton of checks for each loop, I wanted my program to be accurate
	in it's assessment, and be somewhat easy to debug/figure out what was going on. Later on I may
	condense this into a single loop.

	Each function has a specific thing it's looking for and it's looking for that one specific case only.
	As for below, it's matching feat requirements that only require skills, against the skills that the person
	using the application has sent. If there are stat requirements in any form, we delegate that to the last
	function that checks for feats that require both stats and skills only.
*/

func convertToInt(givenNumStr string) (int, error) {
	//string to int conversion
	NumToInt, err := strconv.Atoi(givenNumStr)
	return NumToInt, err
}

func performCheck(metPtr *[]string, failedPtr *[]string, requirement string, name string, given_value int) error {
	required_value, err := convertToInt(requirement)
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

func checkSkill(feats []map[string]string, data RequestData, allAllocatedFeats *[]string) error {
	fmt.Println(" ")
	fmt.Println("////////////////////////")
	fmt.Println("Checking skills...")
	fmt.Println("////////////////////////")
	fmt.Println(" ")

	for i := 0; i < len(feats); i++ {
		feat := feats[i]
		var hasStatRequirement bool = false
		var noFails bool = true
		var skillsMet = []string{}
		var skillsFailed = []string{}
		for j := 0; j < len(data.Skills); j++ {
			skill := data.Skills[j]
			var skillName string = skill.SkillName
			var skillValue int = skill.SkillValue
			if requirement, ok := feat[skillName]; ok {
				err := performCheck(&skillsMet, &skillsFailed, requirement, skillName, skillValue)
				if err != nil {
					return err
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
				if hardreq, ok := feat["HardRequirement"]; ok {
					for _, skillmet := range skillsMet {
						if skillmet != hardreq {
							fmt.Println(skillmet, "->", "{Exists}", "Continuing..")
							noFails = true
							break
						} else {
							fmt.Println("Hardreq only->", feat["HardRequirement"], "{SKIPPED}")
							noFails = false
						}
					}
					//skips the block entirely if ok is false
				} else {
					noFails = true
				}
				if noFails && len(skillsFailed) > 0 {
					for _, failed := range skillsFailed {
						fmt.Println("SKILL = {FAILED}:", "({"+failed, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[failed], failed)
						//if ok is true here, it will then assign and check the NeedsAllSkills value
						if skillrequire, ok := feat["NeedsAllSkills"]; ok {
							//if the value is a string representing "false"
							if skillrequire == "false" {
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

func checkStat(feats []map[string]string, data RequestData, allAllocatedFeats *[]string) error {
	fmt.Println(" ")
	fmt.Println("////////////////////////")
	fmt.Println("Checking stats...")
	fmt.Println("////////////////////////")
	fmt.Println(" ")
	for i := 0; i < len(feats); i++ {
		feat := feats[i]
		var hasSkillRequirement bool = false
		var noFails bool = true
		var statsMet = []string{}
		var statsFailed = []string{}
		for j := 0; j < len(data.Stats); j++ {
			stat := data.Stats[j]
			var statName string = stat.StatName
			var statValue int = stat.StatValue
			if requirement, ok := feat[statName]; ok {
				err := performCheck(&statsMet, &statsFailed, requirement, statName, statValue)
				if err != nil {
					return err
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
			if len(statsFailed) > 0 && !hasSkillRequirement {
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
			//finally, if there are no fails, we append the feat to the allAllocatedFeats slice using a ptr.
			if noFails && !hasSkillRequirement {
				*allAllocatedFeats = append(*allAllocatedFeats, feat["Feat"])
				fmt.Println("STAT = {MET}:", feat["Feat"], "Appended at iteration:", "->", i, "noFails:", noFails, "hasSkillRequirement:", hasSkillRequirement)
			}
		}
	}
	//returning nil, as the function reached the end without errors.
	return nil
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
			var statValue int = stat.StatValue
			if requirement, ok := feat[statName]; ok {
				err := performCheck(&statsMet, &statsFailed, requirement, statName, statValue)
				if err != nil {
					return err
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
			//if there are stats failed, performs a basic check to determine whether every stat is a required
			if len(statsFailed) > 0 {
				for _, statfail := range statsFailed {
					fmt.Println("SKILL/STAT = {FAILED}:", "({"+statfail, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[statfail], statfail)
					if statrequire, ok := feat["NeedsAllStats"]; ok {
						if statrequire == "true" {
							fmt.Println(feat["Feat"], "{FAILED AT STATREQUIRE}:", "Cause ->", statfail)
							noFails = false
							break
						} else if statrequire == "false" {
							noFails = true
						}
					}
				}
			}
			if noFails && hasSkillRequirement {
				for l := 0; l < len(data.Skills); l++ {
					skill := data.Skills[l]
					var skillName string = skill.SkillName
					var skillValue int = skill.SkillValue
					if requirement, ok := feat[skillName]; ok {
						err := performCheck(&skillsMet, &skillsFailed, requirement, skillName, skillValue)
						if err != nil {
							return err
						}
					}
				}
				//if the skillsMet slice is populated
				if len(skillsMet) > 0 {
					if hardreq, ok := feat["HardRequirement"]; ok {
						for _, skillmet := range skillsMet {
							if skillmet != hardreq {
								fmt.Println("{Breaking at hardreq check}:", "->", feat["Feat"])
								noFails = true
								break
							} else {
								fmt.Println("Hardreq only->", feat["HardRequirement"], "{SKIPPED}")
								noFails = false
							}
						}
					} else {
						noFails = true
					}
					if noFails && len(skillsFailed) > 0 {
						for _, failed := range skillsFailed {
							fmt.Println("SKILL = {FAILED}:", "({"+failed, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[failed], failed)
							if skillrequire, ok := feat["NeedsAllSkills"]; ok {
								if skillrequire == "false" {
									if requirement, ok := feat["HardRequirement"]; ok && requirement == failed {
										noFails = false
										fmt.Println("{Failed break}:", "->", feat["Feat"])
										break
									} else {
										noFails = true
									}
									//if all skills are required, then we've failed pretty quick
								} else if skillrequire == "true" {
									noFails = false

									fmt.Println("{Failed break}:", "->", feat["Feat"])

									//loop breaks as a point of failure is found for this iteration of Feats
									break
								}
							}
						}

					}
					//if noFails is true
					if noFails {
						//we append the feat to the corresponding slice.
						*allAllocatedFeats = append(*allAllocatedFeats, feat["Feat"])
						fmt.Println("SKILL/STAT = {MET}", feat["Feat"], "Appended at iteration:", "->", i, "noFails:", noFails)
					}
				}
			}
		}
	}
	return nil
}
