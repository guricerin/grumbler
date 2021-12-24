package model

import "time"

type Grumble struct {
	Pk        string
	Content   string
	UserId    string
	CreatedAt time.Time
}
