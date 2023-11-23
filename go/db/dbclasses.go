package db

// ---------------------------------- //
// SIGNUP CLASSES
type SignupData struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
// PRIVATE SERVEBUILDS CLASSES
type User_Session_Data struct {
	Client_Session_Token string `json:"session_token"`
	Client_ID_Value      uint   `json:"client_id"`
}

type User_Build_IDs struct {
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

// ---------------------------------- //
// LOGIN CLASSES
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionDataStruct struct {
	Client_ID     uint
	Session_Token string
}
