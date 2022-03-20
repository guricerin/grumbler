package model

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

type User struct {
	Pk       uint64
	Id       string
	Name     string
	Password string
	Profile  string
}

func ValidateUserName(name string) error {
	l := utf8.RuneCountInString(name)
	if 0 < l && l < 33 {
		return nil
	} else {
		return errors.New("ユーザ名は1文字以上32文字以下の範囲で設定してください。")
	}
}

var idReg = regexp.MustCompile(`[a-zA-Z0-9\_]{1,32}`)

func ValidateUserId(id string) error {
	l := len(id)
	if !(0 < l && l < 33) {
		return errors.New("ユーザIDは1文字以上32文字以下の半角英数字とアンダーバーの組み合わせで設定してください。")
	}

	if idReg.MatchString(id) {
		return nil
	} else {
		return errors.New("ユーザIDは1文字以上32文字以下の半角英数字とアンダーバーの組み合わせで設定してください。")
	}
}

var passwordReg = regexp.MustCompile(`[a-zA-Z0-9]{8,127}`)

func ValidateUserPassword(plain string) error {
	l := len(plain)
	if !(7 < l && l < 128) {
		return errors.New("パスワードは8文字以上127文字以下の半角英数字で設定してください。")
	}

	if passwordReg.MatchString(plain) {
		return nil
	} else {
		return errors.New("パスワードは8文字以上127文字以下の半角英数字で設定してください。")
	}
}

func ValidateUserProfile(profile string) error {
	l := utf8.RuneCountInString(profile)
	if 0 <= l && l < 201 {
		return nil
	} else {
		return errors.New("プロフィールは0文字以上200文字以下の範囲で設定してください。")
	}
}
