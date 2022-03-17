package model

import (
	"errors"
	"sort"
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

type ReplyInfoForGrumbleRes struct {
	DstGrumblePk string
	DstUserId    string
	RepliedCount int
}

type RegrumbleInfoForGrumbleRes struct {
	RegrumbledCount          int
	IsRegrumble              bool
	CreatedAt                time.Time
	ByUserId                 string
	IsRegrumbledBySigninUser bool
}

type GrumbleRes struct {
	Pk                       string
	Content                  string
	UserId                   string
	CreatedAt                time.Time
	UserName                 string
	Reply                    ReplyInfoForGrumbleRes
	Regrumble                RegrumbleInfoForGrumbleRes
	BookmarkedCount          int
	IsBookmarkedBySigninUser bool
}

// 最新日時順
// リグランブルの場合は、元のグランブルが投稿された日時ではなくリグランブルされた日時を元にソート
func SortGrumblesForNewest(grumbles []GrumbleRes) {
	sort.Slice(grumbles, func(i, j int) bool {
		var createdAtI, createdAtJ time.Time
		if grumbles[i].Regrumble.IsRegrumble && grumbles[j].Regrumble.IsRegrumble {
			createdAtI = grumbles[i].Regrumble.CreatedAt
			createdAtJ = grumbles[j].Regrumble.CreatedAt
		} else if grumbles[i].Regrumble.IsRegrumble {
			createdAtI = grumbles[i].Regrumble.CreatedAt
			createdAtJ = grumbles[j].CreatedAt
		} else if grumbles[j].Regrumble.IsRegrumble {
			createdAtI = grumbles[i].CreatedAt
			createdAtJ = grumbles[j].Regrumble.CreatedAt
		} else {
			createdAtI = grumbles[i].CreatedAt
			createdAtJ = grumbles[j].CreatedAt
		}
		return createdAtI.After(createdAtJ)
	})
}
