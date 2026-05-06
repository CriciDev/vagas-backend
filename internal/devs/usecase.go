package devs

import (
	"context"
	"log"
	"net/mail"
	"unicode/utf8"

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

func (usecase *DevUseCase) CreateDev(ctx context.Context, dev *Dev) (*CreateDevRow, error) {

	err := dev.validate()

	if err != nil {
		return nil, err
	}

	row_count, err := usecase.DevExistsByEmail(ctx, dev.Email)

	if err != nil {
		return nil, err
	}

	if row_count > 0 {
		return nil, ErrEmailAlreadyInUse
	}

	hashedPassword, err := utils.EncryptPassword(dev.Password)

	if err != nil {
		return nil, err
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

	resp, err := usecase.Repository.CreateDev(ctx, devParams)

	if err != nil {
		log.Printf("ERRO: Falha no banco de dados ao salvar novo usuário! Message: %s", err.Error())
		return nil, err
	}

	return &resp, nil
}

func (usecase *DevUseCase) FindDevByEmail(ctx context.Context, email string) (*FindDevByEmailRow, error) {
	dev, err := usecase.Repository.FindDevByEmail(ctx, email)

	if err != nil {
		return nil, ErrInvalidEmailOrPassword
	}

	return &dev, nil
}

func (usecase *DevUseCase) FindDevByID(ctx context.Context, id int32) (*FindDevByIDRow, error) {
	dev, err := usecase.Repository.FindDevByID(ctx, id)

	if err != nil {
		return nil, ErrProfileNotFound
	}

	return &dev, nil
}

func (usecase *DevUseCase) DevExistsByEmail(ctx context.Context, email string) (int64, error) {
	row_count, err := usecase.Repository.EmailAlreadyRegistered(ctx, email)

	if err != nil {
		log.Printf("ERRO: Falha ao executar busca no banco de dados! Message: %s", err.Error())
		return 0, err
	}

	return row_count, nil
}

func (dev *Dev) validate() error {
	if dev.Name == "" {
		return ErrNameNotEmpty
	}

	if dev.Email == "" {
		return ErrEmailNotEmpty
	}

	_, err := mail.ParseAddress(dev.Email)

	if err != nil {
		return ErrEmailAddressNotValid
	}

	if dev.Password == "" {
		return ErrPasswordNotEmpty
	}

	//bcrypt encoder trunkate strings with more than 72 characters
	if utf8.RuneCountInString(dev.Password) > 72 {
		return ErrPasswordToLong
	}

	if dev.Bio == "" {
		return ErrBioNotEmpty
	}

	return nil
}
