package opportunities

import (
	"context"
	"net/mail"
	"net/url"
	"strings"
)

type (
	Service struct {
		repo Repository
	}
)

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (service *Service) Create(ctx context.Context, request SaveOpportunityRequest) (Opportunity, error) {
	opportunity, err := opportunityFromRequest(request, StatusPublished)
	if err != nil {
		return Opportunity{}, err
	}
	return service.repo.Create(ctx, opportunity)
}

func (service *Service) List(ctx context.Context, filters ListFilters, pagination Pagination) (OpportunityPage, error) {
	filters = normalizeFilters(filters)
	if filters.Type != "" && !validOpportunityType(filters.Type) {
		return OpportunityPage{}, ErrValidation
	}
	if filters.WorkMode != "" && !validWorkMode(filters.WorkMode) {
		return OpportunityPage{}, ErrValidation
	}

	items, total, err := service.repo.List(ctx, filters, pagination)
	if err != nil {
		return OpportunityPage{}, err
	}

	return OpportunityPage{
		Data: items,
		Meta: PageMeta{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      total,
			TotalPages: totalPages(total, pagination.PageSize),
		},
	}, nil
}

func NewPagination(page, pageSize int) Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Pagination{Page: page, PageSize: pageSize}
}

func (pagination Pagination) Offset() int {
	return (pagination.Page - 1) * pagination.PageSize
}

func (pagination Pagination) Limit() int {
	return pagination.PageSize
}

func totalPages(total, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return (total + pageSize - 1) / pageSize
}

func (service *Service) FindByID(ctx context.Context, id int64, includeUnpublished bool) (Opportunity, error) {
	if id <= 0 {
		return Opportunity{}, ErrNotFound
	}
	if includeUnpublished {
		return service.repo.FindByID(ctx, id)
	}
	return service.repo.FindPublishedByID(ctx, id)
}

func (service *Service) Update(ctx context.Context, id int64, request SaveOpportunityRequest) (Opportunity, error) {
	if id <= 0 {
		return Opportunity{}, ErrNotFound
	}

	existing, err := service.repo.FindByID(ctx, id)
	if err != nil {
		return Opportunity{}, err
	}

	opportunity, err := opportunityFromRequest(request, existing.Status)
	if err != nil {
		return Opportunity{}, err
	}
	opportunity.ID = id

	return service.repo.Update(ctx, opportunity)
}

func (service *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrNotFound
	}
	return service.repo.Delete(ctx, id)
}

func opportunityFromRequest(request SaveOpportunityRequest, defaultStatus string) (Opportunity, error) {
	title := strings.TrimSpace(request.Title)
	description := strings.TrimSpace(request.Description)
	organizationName := strings.TrimSpace(request.OrganizationName)
	organizationURL := strings.TrimSpace(request.OrganizationURL)
	opportunityType := strings.TrimSpace(strings.ToLower(request.Type))
	workMode := strings.TrimSpace(strings.ToLower(request.WorkMode))
	location := strings.TrimSpace(request.Location)
	salaryRange := strings.TrimSpace(request.SalaryRange)
	seniority := strings.TrimSpace(request.Seniority)
	contactEmail := strings.TrimSpace(strings.ToLower(request.ContactEmail))
	contactURL := strings.TrimSpace(request.ContactURL)
	status := strings.TrimSpace(strings.ToLower(request.Status))
	skills := normalizeSkills(request.Skills)

	if status == "" {
		status = defaultStatus
	}

	if title == "" || description == "" || organizationName == "" || opportunityType == "" || workMode == "" {
		return Opportunity{}, ErrValidation
	}
	if contactEmail == "" && contactURL == "" {
		return Opportunity{}, ErrValidation
	}
	if !validOpportunityType(opportunityType) || !validWorkMode(workMode) || !validStatus(status) {
		return Opportunity{}, ErrValidation
	}
	if contactEmail != "" {
		if _, err := mail.ParseAddress(contactEmail); err != nil {
			return Opportunity{}, ErrValidation
		}
	}
	if contactURL != "" && !validHTTPURL(contactURL) {
		return Opportunity{}, ErrValidation
	}
	if organizationURL != "" && !validHTTPURL(organizationURL) {
		return Opportunity{}, ErrValidation
	}

	return Opportunity{
		Title:            title,
		Description:      description,
		OrganizationName: organizationName,
		OrganizationURL:  organizationURL,
		Type:             opportunityType,
		WorkMode:         workMode,
		Location:         location,
		SalaryRange:      salaryRange,
		Seniority:        seniority,
		Skills:           skills,
		ContactEmail:     contactEmail,
		ContactURL:       contactURL,
		ExpiresAt:        request.ExpiresAt,
		Status:           status,
	}, nil
}

func normalizeFilters(filters ListFilters) ListFilters {
	return ListFilters{
		Type:     strings.TrimSpace(strings.ToLower(filters.Type)),
		WorkMode: strings.TrimSpace(strings.ToLower(filters.WorkMode)),
		Location: strings.TrimSpace(filters.Location),
	}
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

func validOpportunityType(value string) bool {
	switch value {
	case TypeFullTime, TypePartTime, TypeContract, TypeFreelance, TypeVolunteer, TypeProject, TypeMentorship:
		return true
	default:
		return false
	}
}

func validWorkMode(value string) bool {
	switch value {
	case WorkModeRemote, WorkModeHybrid, WorkModeOnSite:
		return true
	default:
		return false
	}
}

func validStatus(value string) bool {
	switch value {
	case StatusDraft, StatusPublished, StatusClosed, StatusArchived:
		return true
	default:
		return false
	}
}

func validHTTPURL(value string) bool {
	parsed, err := url.ParseRequestURI(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}
