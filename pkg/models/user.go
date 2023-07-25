package models

import (
	"database/sql"
	"encoding/json"
)

func NewUser(number Number, id uint32) User {
	return User{
		ID:     id,
		Number: number,
	}
}

func (n Number) Marshall() ([]byte, error) {
	jsonBytes, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (codeRequest AuthNumberRequest) Marshall() ([]byte, error) {
	jsonBytes, err := json.Marshal(codeRequest)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (codeRequest AuthCodeRequest) Marshall() ([]byte, error) {
	jsonBytes, err := json.Marshal(codeRequest)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (request Auth2FARequest) Marshall() ([]byte, error) {
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (request AuthCheckClientRequest) Marshall() ([]byte, error) {
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

type Number string

type AuthNumberRequest struct {
	Id      uint32 `json:"id"`
	Command string `json:"command"`
	Number  Number `json:"phone"`
}

type AuthCodeRequest struct {
	UserID  uint32 `json:"id"`
	Command string `json:"command"`
	Number  Number `json:"phone"`
	Code    string `json:"phone_code"`
}

type Auth2FARequest struct {
	UserID   uint32
	Password string
}

type AuthCheckClientRequest struct {
	UserID  uint32 `json:"id"`
	Command string `json:"command"`
}

type User struct {
	ID       uint32         `json:"id"`
	Nickname sql.NullString `json:"nickname"`
	Number   Number         `json:"number"`
}

type Session struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}

type AuthTelethonFinishResponse struct {
	Command string  `json:"command"`
	Details Details `json:"details"`
	ID      int     `json:"id"`
}

type Details struct {
	Nickname string `json:"nickname"`
}

type AuthFinishResponse struct {
	Nickname string `json:"nickname"`
}
