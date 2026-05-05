package devs

import (
	"context"
	"log"
	"net/http"

	"github.com/CriciumaDevJobs/backend/handlers"
	"github.com/CriciumaDevJobs/backend/utils"
)

type DevUseCase struct {
	Repository *Queries
}

func NewDevUseCase(devRepository *Queries) *DevUseCase {
	usecase := DevUseCase{
		Repository: devRepository,
	}

	return &usecase
}

func (usecase *DevUseCase) CreateDev(ctx context.Context, dev *Dev) (*CreateDevRow, *handlers.ErrorResponse) {

	http_err := dev.validate()

	if http_err != nil {
		return nil, http_err
	}

	_, db_err := usecase.Repository.FindDevByEmail(ctx, dev.Email)

	if db_err == nil {
		return nil, handlers.ErrEmailAlreadyInUse
	}

	hashedPassword, http_err := utils.EncryptPassword(dev.Password)

	if http_err != nil {
		return nil, http_err
	}

	var devParams = CreateDevParams{
		Name:         dev.Name,
		Email:        dev.Email,
		Password:     hashedPassword,
		Skills:       dev.Skills,
		Bio:          dev.Bio,
		Availability: dev.Availability,
		Socials:      dev.Socials,
	}

	resp, db_err := usecase.Repository.CreateDev(ctx, devParams)

	if db_err != nil {
		log.Printf("ERRO: Falha no banco de dados ao salvar novo usuário! Message: %s", db_err.Error())
		return nil, handlers.NewError(http.StatusInternalServerError, "Erro Interno!")
	}

	return &resp, nil
}

func (usecase *DevUseCase) FindDevByEmail(ctx context.Context, email string) (*FindDevByEmailRow, *handlers.ErrorResponse) {
	dev, err := usecase.Repository.FindDevByEmail(ctx, email)

	if err != nil {
		return nil, handlers.ErrInvalidEmailOrPassword
	}

	return &dev, nil
}

func (usecase *DevUseCase) FindDevByID(ctx context.Context, id int32) (*FindDevByIDRow, *handlers.ErrorResponse) {
	dev, err := usecase.Repository.FindDevByID(ctx, id)

	if err != nil {
		return nil, handlers.ErrProfileNotFound
	}

	return &dev, nil
}

func (dev *Dev) validate() *handlers.ErrorResponse {
	if dev.Name == "" {
		return handlers.ErrNameNotEmpty
	}

	if dev.Email == "" {
		return handlers.ErrEmailNotEmpty
	}

	if dev.Bio == "" {
		return handlers.ErrBioNotEmpty
	}

	return nil
}
