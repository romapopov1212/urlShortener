package repository

import "errors"

var (
	ErrAliasAlreadyExist = errors.New("alias already exist")
	ErrNotFound          = errors.New("not found exist")
)

type Url struct {
	Url   string
	Alias string
}
