package db

import "fmt"

type Stat struct {
	StatName  string `json:"statName"`
	StatValue string `json:"statValue"`
}

type Skill struct {
	SkillName  string `json:"skillName"`
	SkillValue string `json:"skillValue"`
}

type Feat struct {
	FeatName string `json:"Feat"`
}

type SaveData struct {
	Stats  []Stat  `json:"stats"`
	Skills []Skill `json:"skills"`
	Feats  []Feat  `json:"feats"`
}

func SaveBuild(build SaveData) {
	for i := 0; i < len(build.Feats); i++ {
		feat := build.Feats[i]
		fmt.Println(feat.FeatName)
	}
}
