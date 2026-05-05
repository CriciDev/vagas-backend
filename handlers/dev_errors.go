package handlers

import (
	"net/http"
)

var (
	ErrNameNotEmpty           = NewError(http.StatusUnprocessableEntity, "Campo Nome deve ser preenchido!")
	ErrEmailNotEmpty          = NewError(http.StatusUnprocessableEntity, "Campo Email deve ser preenchido!")
	ErrEmailAlreadyInUse      = NewError(http.StatusUnprocessableEntity, "Email já está em uso")
	ErrBioNotEmpty            = NewError(http.StatusUnprocessableEntity, "Campo Biografia deve ser preenchido")
	ErrPasswordToLong         = NewError(http.StatusUnprocessableEntity, "O Campo de Senha não pode ultrapassar 72 caracteres")
	ErrInvalidEmailOrPassword = NewError(http.StatusUnauthorized, "Email ou Senha inválidos")
	ErrProfileNotFound        = NewError(http.StatusNotFound, "Houve algum erro ao encontrar seu Perfil!")
)
