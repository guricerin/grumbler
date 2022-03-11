package model

import (
	"errors"
	"time"
)

type Grumble struct {
	Pk        string
	Content   string
	UserId    string
	CreatedAt time.Time
}

func ValidateGrumble(text string) error {
	l := len(text)
	if 0 < l && l < 301 {
		return nil
	} else {
		return errors.New("グランブルの文字数制限は1文字以上300文字以下です。")
	}
}

type GrumbleRes struct {
	Pk                       string
	Content                  string
	UserId                   string
	CreatedAt                time.Time
	UserName                 string
	RepliedCount             int
	BookmarkedCount          int
	IsBookmarkedBySigninUser bool
}
