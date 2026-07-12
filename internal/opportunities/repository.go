package opportunities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type (
	Repository interface {
		Create(ctx context.Context, opportunity Opportunity) (Opportunity, error)
		List(ctx context.Context, filters ListFilters, pagination Pagination) ([]Opportunity, int, error)
		FindByID(ctx context.Context, id int64) (Opportunity, error)
		FindPublishedByID(ctx context.Context, id int64) (Opportunity, error)
		Update(ctx context.Context, opportunity Opportunity) (Opportunity, error)
		Delete(ctx context.Context, id int64) error
	}

	PostgresRepository struct {
		db *sql.DB
	}
)

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (repo *PostgresRepository) Create(ctx context.Context, opportunity Opportunity) (Opportunity, error) {
	row := repo.db.QueryRowContext(ctx, `
		INSERT INTO opportunities (
			title, description, organization_name, organization_url, type, work_mode,
			location, salary_range, seniority, skills, contact_email, contact_url, expires_at, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, title, description, organization_name, organization_url, type, work_mode,
			location, salary_range, seniority, skills, contact_email, contact_url, expires_at, status, created_at, updated_at
	`, opportunity.Title, opportunity.Description, opportunity.OrganizationName, nullString(opportunity.OrganizationURL), opportunity.Type,
		opportunity.WorkMode, nullString(opportunity.Location), nullString(opportunity.SalaryRange), nullString(opportunity.Seniority),
		pq.Array(opportunity.Skills), nullString(opportunity.ContactEmail), nullString(opportunity.ContactURL), opportunity.ExpiresAt, opportunity.Status)

	return scanOpportunity(row)
}

func (repo *PostgresRepository) List(ctx context.Context, filters ListFilters, pagination Pagination) ([]Opportunity, int, error) {
	clauses := []string{"status = 'published'"}
	args := []any{}

	if filters.Type != "" {
		args = append(args, filters.Type)
		clauses = append(clauses, fmt.Sprintf("type = $%d", len(args)))
	}
	if filters.WorkMode != "" {
		args = append(args, filters.WorkMode)
		clauses = append(clauses, fmt.Sprintf("work_mode = $%d", len(args)))
	}
	if filters.Location != "" {
		args = append(args, "%"+filters.Location+"%")
		clauses = append(clauses, fmt.Sprintf("location ILIKE $%d", len(args)))
	}

	where := strings.Join(clauses, " AND ")

	var total int
	if err := repo.db.QueryRowContext(ctx, `SELECT count(*) FROM opportunities WHERE `+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, pagination.Limit())
	limitPlaceholder := len(args)
	args = append(args, pagination.Offset())
	offsetPlaceholder := len(args)

	query := `
		SELECT id, title, description, organization_name, organization_url, type, work_mode,
			location, salary_range, seniority, skills, contact_email, contact_url, expires_at, status, created_at, updated_at
		FROM opportunities
		WHERE ` + where + `
		ORDER BY created_at DESC, id DESC
		LIMIT ` + fmt.Sprintf("$%d", limitPlaceholder) + ` OFFSET ` + fmt.Sprintf("$%d", offsetPlaceholder)

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	opportunities := []Opportunity{}
	for rows.Next() {
		opportunity, err := scanOpportunity(rows)
		if err != nil {
			return nil, 0, err
		}
		opportunities = append(opportunities, opportunity)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return opportunities, total, nil
}

func (repo *PostgresRepository) FindByID(ctx context.Context, id int64) (Opportunity, error) {
	return repo.findByID(ctx, id, "id = $1")
}

func (repo *PostgresRepository) FindPublishedByID(ctx context.Context, id int64) (Opportunity, error) {
	return repo.findByID(ctx, id, "id = $1 AND status = 'published'")
}

func (repo *PostgresRepository) findByID(ctx context.Context, id int64, condition string) (Opportunity, error) {
	row := repo.db.QueryRowContext(ctx, `
		SELECT id, title, description, organization_name, organization_url, type, work_mode,
			location, salary_range, seniority, skills, contact_email, contact_url, expires_at, status, created_at, updated_at
		FROM opportunities
		WHERE `+condition, id)

	opportunity, err := scanOpportunity(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Opportunity{}, ErrNotFound
	}
	if err != nil {
		return Opportunity{}, err
	}

	return opportunity, nil
}

func (repo *PostgresRepository) Update(ctx context.Context, opportunity Opportunity) (Opportunity, error) {
	row := repo.db.QueryRowContext(ctx, `
		UPDATE opportunities
		SET title = $2, description = $3, organization_name = $4, organization_url = $5, type = $6,
			work_mode = $7, location = $8, salary_range = $9, seniority = $10, skills = $11,
			contact_email = $12, contact_url = $13, expires_at = $14, status = $15, updated_at = now()
		WHERE id = $1
		RETURNING id, title, description, organization_name, organization_url, type, work_mode,
			location, salary_range, seniority, skills, contact_email, contact_url, expires_at, status, created_at, updated_at
	`, opportunity.ID, opportunity.Title, opportunity.Description, opportunity.OrganizationName, nullString(opportunity.OrganizationURL), opportunity.Type,
		opportunity.WorkMode, nullString(opportunity.Location), nullString(opportunity.SalaryRange), nullString(opportunity.Seniority),
		pq.Array(opportunity.Skills), nullString(opportunity.ContactEmail), nullString(opportunity.ContactURL), opportunity.ExpiresAt, opportunity.Status)

	updated, err := scanOpportunity(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Opportunity{}, ErrNotFound
	}
	if err != nil {
		return Opportunity{}, err
	}

	return updated, nil
}

func (repo *PostgresRepository) Delete(ctx context.Context, id int64) error {
	result, err := repo.db.ExecContext(ctx, `DELETE FROM opportunities WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

type (
	scanner interface {
		Scan(dest ...any) error
	}
)

func scanOpportunity(row scanner) (Opportunity, error) {
	var opportunity Opportunity
	var organizationURL sql.NullString
	var location sql.NullString
	var salaryRange sql.NullString
	var seniority sql.NullString
	var skills pq.StringArray
	var contactEmail sql.NullString
	var contactURL sql.NullString
	var expiresAt sql.NullTime

	err := row.Scan(
		&opportunity.ID,
		&opportunity.Title,
		&opportunity.Description,
		&opportunity.OrganizationName,
		&organizationURL,
		&opportunity.Type,
		&opportunity.WorkMode,
		&location,
		&salaryRange,
		&seniority,
		&skills,
		&contactEmail,
		&contactURL,
		&expiresAt,
		&opportunity.Status,
		&opportunity.CreatedAt,
		&opportunity.UpdatedAt,
	)
	if err != nil {
		return Opportunity{}, err
	}

	opportunity.OrganizationURL = organizationURL.String
	opportunity.Location = location.String
	opportunity.SalaryRange = salaryRange.String
	opportunity.Seniority = seniority.String
	opportunity.Skills = []string(skills)
	opportunity.ContactEmail = contactEmail.String
	opportunity.ContactURL = contactURL.String
	if expiresAt.Valid {
		expiresAtValue := expiresAt.Time
		opportunity.ExpiresAt = &expiresAtValue
	}

	return opportunity, nil
}

func nullString(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}
