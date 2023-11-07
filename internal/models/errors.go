package models

import "errors"

var (
	ErrClientAlreadyExists = errors.New("клиент с таким логином уже существует")
	ErrUserNotFound        = errors.New("клиента с таким логином не существует")
)
