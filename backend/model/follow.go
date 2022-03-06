package model

type Follow struct {
	Pk        int64
	SrcUserId string // フォロー元ユーザID
	DstUserId string // フォロー先ユーザID
}
