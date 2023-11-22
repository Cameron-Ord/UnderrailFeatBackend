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

/*
The performCheck function uses the pointers of whatever met and failed slices are being sent to it. along with
the requirement and the name of the given skill it's checking along with the value of said skill. This serves
as a general function that each loop calls to when it reaches a point in it's code where it needs to perform
this check. All it does is check some numbers against each other, nothing too wild here.
*/
func convertToInt(givenNumStr string) (int, error) {
	//string to int conversion
	NumToInt, err := strconv.Atoi(givenNumStr)
	return NumToInt, err
}

func performCheck(metPtr *[]string, failedPtr *[]string, requirement string, name string, value string) error {
	skillRequirement, err := convertToInt(requirement)
	if err != nil {
		return err
	}

	givenSkillValue, err := convertToInt(value)
	if err != nil {
		return err
	}

	if givenSkillValue >= skillRequirement {
		*metPtr = append(*metPtr, name)
	} else if givenSkillValue < skillRequirement {
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
				err := performCheck(&skillsMet, &skillsFailed, requirement, skillName, skillValue)
				if err != nil {
					return err
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
							fmt.Println("Hardreq only->", feat["HardRequirement"], "{SKIPPED}")
							noFails = false
						}
					}
					//skips the block entirely if ok is false
				} else {
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
				err := performCheck(&statsMet, &statsFailed, requirement, statName, statValue)
				if err != nil {
					return err
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
							//otherwise, continue on
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
	//looping through the len of the feats map
	for i := 0; i < len(feats); i++ {
		//assigning each index to a variable
		feat := feats[i]
		//declaring variables at the beginning of the iteration
		var hasSkillRequirement = false
		var noFails bool = true
		var statsMet = []string{}
		var statsFailed = []string{}
		var skillsMet = []string{}
		var skillsFailed = []string{}
		//for loop for Stats
		for j := 0; j < len(data.Stats); j++ {
			//assigning the current index of Stats to a variable
			stat := data.Stats[j]
			//copying the StatName and StatValue over to their own variables
			var statName string = stat.StatName
			var statValue string = stat.StatValue
			//if this iteration of Feats(i) has the current iteration of Stats(j) ok gets set to true and the code gets executed
			if requirement, ok := feat[statName]; ok {
				err := performCheck(&statsMet, &statsFailed, requirement, statName, statValue)
				if err != nil {
					return err
				}
			}
		}

		//if the statsMet slice is populated for this iteration
		if len(statsMet) > 0 {
			//using a similar check to earlier to simply check whether any stats are found within the current index of Feats
			for k := 0; k < len(data.Skills); k++ {
				skill := data.Skills[k]
				var skillName string = skill.SkillName
				if _, ok := feat[skillName]; ok {
					//setting to true and breaking if a stat is found
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
			//if noFails is still true and there is a skill requirement
			if noFails && hasSkillRequirement {
				//doing the same as before, but for the skills
				for l := 0; l < len(data.Skills); l++ {
					skill := data.Skills[l]
					var skillName string = skill.SkillName
					var skillValue string = skill.SkillValue
					if requirement, ok := feat[skillName]; ok {
						err := performCheck(&skillsMet, &skillsFailed, requirement, skillName, skillValue)
						if err != nil {
							return err
						}
					}
				}
				//if the skillsMet slice is populated
				if len(skillsMet) > 0 {
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
					//if noFails is still true, and skillsFailed is populated
					if noFails && len(skillsFailed) > 0 {
						//looping through the skillsFailed slice
						for _, failed := range skillsFailed {
							fmt.Println("SKILL = {FAILED}:", "({"+failed, "->", "at iteration}):", i, "||", "{Feat}:", "({"+feat["Feat"], "->", "needs}):", feat[failed], failed)
							//checking whether all skills are required
							if skillrequire, ok := feat["NeedsAllSkills"]; ok {
								//if every skill isnt required to meet specifications,
								//then it begins checking for hardrequirements
								if skillrequire == "false" {
									//if the failed skill matches the hardrequirement, ok will be set to true and the block will execute
									if requirement, ok := feat["HardRequirement"]; ok && requirement == failed {
										//no fails gets set to false
										noFails = false

										fmt.Println("{Failed break}:", "->", feat["Feat"])

										//loop breaks as a point of failure is found for this iteration of Feats
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
