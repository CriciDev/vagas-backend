package devs

import (
	"context"
	"net/mail"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (service *Service) Create(ctx context.Context, request SaveDeveloperRequest) (Developer, error) {
	developer, err := developerFromRequest(request)
	if err != nil {
		return Developer{}, err
	}
	return service.repo.Create(ctx, developer)
}

func (service *Service) List(ctx context.Context) ([]Developer, error) {
	return service.repo.List(ctx)
}

func (service *Service) FindByID(ctx context.Context, id int64) (Developer, error) {
	if id <= 0 {
		return Developer{}, ErrNotFound
	}
	return service.repo.FindByID(ctx, id)
}

func (service *Service) Update(ctx context.Context, id int64, request SaveDeveloperRequest) (Developer, error) {
	if id <= 0 {
		return Developer{}, ErrNotFound
	}

	developer, err := developerFromRequest(request)
	if err != nil {
		return Developer{}, err
	}
	developer.ID = id

	return service.repo.Update(ctx, developer)
}

func (service *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrNotFound
	}
	return service.repo.Delete(ctx, id)
}

func developerFromRequest(request SaveDeveloperRequest) (Developer, error) {
	name := strings.TrimSpace(request.Name)
	email := strings.TrimSpace(strings.ToLower(request.Email))
	bio := strings.TrimSpace(request.Bio)
	skills := normalizeSkills(request.Skills)

	if name == "" || email == "" || len(skills) == 0 {
		return Developer{}, ErrValidation
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return Developer{}, ErrValidation
	}

	return Developer{Name: name, Email: email, Skills: skills, Available: request.Available, Bio: bio}, nil
}

func normalizeSkills(values []string) []string {
	skills := []string{}
	seen := map[string]bool{}
	for _, value := range values {
		skill := strings.TrimSpace(value)
		key := strings.ToLower(skill)
		if skill == "" || seen[key] {
			continue
		}
		seen[key] = true
		skills = append(skills, skill)
	}
	return skills
}
