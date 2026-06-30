package opportunities

import (
	"errors"
	"time"
)

const (
	TypeFullTime   = "full_time"
	TypePartTime   = "part_time"
	TypeContract   = "contract"
	TypeFreelance  = "freelance"
	TypeVolunteer  = "volunteer"
	TypeProject    = "project"
	TypeMentorship = "mentorship"

	WorkModeRemote = "remote"
	WorkModeHybrid = "hybrid"
	WorkModeOnSite = "on_site"

	StatusDraft     = "draft"
	StatusPublished = "published"
	StatusClosed    = "closed"
	StatusArchived  = "archived"
)

var (
	ErrNotFound   = errors.New("opportunity not found")
	ErrValidation = errors.New("opportunity validation failed")
)

type (
	Opportunity struct {
		ID               int64      `json:"id"`
		Title            string     `json:"title"`
		Description      string     `json:"description"`
		OrganizationName string     `json:"organization_name"`
		OrganizationURL  string     `json:"organization_url"`
		Type             string     `json:"type"`
		WorkMode         string     `json:"work_mode"`
		Location         string     `json:"location"`
		SalaryRange      string     `json:"salary_range"`
		Seniority        string     `json:"seniority"`
		Skills           []string   `json:"skills"`
		ContactEmail     string     `json:"contact_email"`
		ContactURL       string     `json:"contact_url"`
		ExpiresAt        *time.Time `json:"expires_at,omitempty"`
		Status           string     `json:"status"`
		CreatedAt        time.Time  `json:"created_at"`
		UpdatedAt        time.Time  `json:"updated_at"`
	}

	SaveOpportunityRequest struct {
		Title            string     `json:"title"`
		Description      string     `json:"description"`
		OrganizationName string     `json:"organization_name"`
		OrganizationURL  string     `json:"organization_url"`
		Type             string     `json:"type"`
		WorkMode         string     `json:"work_mode"`
		Location         string     `json:"location"`
		SalaryRange      string     `json:"salary_range"`
		Seniority        string     `json:"seniority"`
		Skills           []string   `json:"skills"`
		ContactEmail     string     `json:"contact_email"`
		ContactURL       string     `json:"contact_url"`
		ExpiresAt        *time.Time `json:"expires_at"`
		Status           string     `json:"status"`
	}

	ListFilters struct {
		Type     string
		WorkMode string
		Location string
	}
)
