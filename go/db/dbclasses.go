package db

// ---------------------------------- //
// SIGNUP CLASSES
type SignupData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ---------------------------------- //
// LOGIN CLASSES
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ---------------------------------- //
// SESSION DATA
type User_Session_Data struct {
	Client_Session_Token string
	Client_ID_Value      uint
}

type Basic_User_Info struct {
	Username   string
	Created_at string
}

// ---------------------------------- //
// SERVEBUILDS CLASSES
type Feat_Info struct {
	Name     string
	Build_ID uint
}
type Skill_Info struct {
	Name     string
	Value    int
	Build_ID uint
}
type Stat_Info struct {
	Name     string
	Value    int
	Build_ID uint
}
type Build_ID_Titles struct {
	Build_ID uint
	Title    string
}

// ---------------------------------- //
// SAVEBUILD CLASSES
type Stat struct {
	StatName  string `json:"statName"`
	StatValue int    `json:"statValue"`
}

type Skill struct {
	SkillName  string `json:"skillName"`
	SkillValue int    `json:"skillValue"`
}

type Feat struct {
	FeatName string `json:"Feat"`
}

type SaveData struct {
	Stats         []Stat  `json:"stats"`
	Skills        []Skill `json:"skills"`
	Feats         []Feat  `json:"feats"`
	Title         string  `json:"title"`
	Client_ID     uint    `json:"client_id"`
	Session_Token string  `json:"session_token"`
}

// DELETE

type DeleteData struct {
	Client_ID uint `json:"client_id"`
	Build_ID  uint `json:"build_id"`
}
