package model

import (
	"errors"
	"regexp"
)

type User struct {
	Pk       uint64
	Id       string
	Name     string
	Password string
	Profile  string
}

func ValidateUserName(name string) error {
	l := len(name)
	if 0 < l && l < 128 {
		return nil
	} else {
		return errors.New("ユーザ名は1文字以上127文字以下の範囲で設定してください。")
	}
}

var idReg = regexp.MustCompile(`[a-zA-Z0-9\_]{1,127}`)

func ValidateUserId(id string) error {
	if idReg.MatchString(id) {
		return nil
	} else {
		return errors.New("ユーザIDは1文字以上127文字以下の半角英数字とアンダーバーの組み合わせで設定してください。")
	}
}

var passwordReg = regexp.MustCompile(`[a-zA-Z0-9]{8,127}`)

func ValidateUserPassword(plain string) error {
	if passwordReg.MatchString(plain) {
		return nil
	} else {
		return errors.New("パスワードは8文字以上127文字以下の半角英数字で設定してください。")
	}
}
