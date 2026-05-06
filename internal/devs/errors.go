package devs

import (
	"errors"
)

var (
	ErrNameNotEmpty           = errors.New("Campo Nome deve ser preenchido!")
	ErrEmailNotEmpty          = errors.New("Campo Email deve ser preenchido!")
	ErrPasswordNotEmpty       = errors.New("Campo Senha deve ser preenchido")
	ErrEmailAddressNotValid   = errors.New("Por favor, Insira um email válido")
	ErrEmailAlreadyInUse      = errors.New("Email já está em uso")
	ErrBioNotEmpty            = errors.New("Campo Biografia deve ser preenchido")
	ErrPasswordToLong         = errors.New("O Campo de Senha não pode ultrapassar 72 caracteres")
	ErrInvalidEmailOrPassword = errors.New("Email ou Senha inválidos")
	ErrProfileNotFound        = errors.New("Houve algum erro ao encontrar seu Perfil!")
	ErrJsonNotExpected        = errors.New("JSON Enviado não segue a estrutura esperada!")
)
