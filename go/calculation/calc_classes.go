package calculation

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

type Skills_Tracker struct {
	Skill_Name  string
	Skill_Value int
	Feat_Value  string
}

type Stats_Tracker struct {
	Stat_Name  string
	Stat_Value int
	Feat_Value string
}
