package handlers

import "errors"

var (
	ErrNameNotEmpty  = errors.New("Campo Nome deve ser preenchido!")
	ErrEmailNotEmpty = errors.New("Campo Email deve ser preenchido!")
)
