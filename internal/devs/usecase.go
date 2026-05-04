package devs

import (
	"context"

	"github.com/CriciumaDevJobs/backend/handlers"
	"github.com/CriciumaDevJobs/backend/internal/auth"
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

func (usecase *DevUseCase) CreateDev(ctx context.Context, dev *Dev) (CreateDevRow, error) {

	err := dev.validate()

	if err != nil {
		return CreateDevRow{}, err
	}

	_, err = usecase.Repository.FindDevByEmail(ctx, dev.Email)

	if err == nil {
		return CreateDevRow{}, handlers.ErrEmailAlreadyInUse
	}

	hashedPassword, err := auth.EncryptPassword(dev.Password)

	if err != nil {
		return CreateDevRow{}, err
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

	return usecase.Repository.CreateDev(ctx, devParams)
}

func (dev *Dev) validate() error {
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
