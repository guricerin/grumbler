package model

type Session struct {
	Pk     uint64
	Token  string
	UserPk uint64
}
