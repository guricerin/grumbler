package model

type User struct {
	Pk       uint
	Id       string
	Name     string
	Password string
	Profile  string
}

type UserJson struct {
	Pk      uint   `json:"pk,number"`
	Id      string `json:"id,string"`
	Name    string `json:"name,string"`
	Profile string `json:"profile,string"`
}

func (u *User) ToJson() UserJson {
	return UserJson{
		Pk:      u.Pk,
		Id:      u.Id,
		Name:    u.Name,
		Profile: u.Profile,
	}
}
