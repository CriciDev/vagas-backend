## Purpose

Define the opportunity management behavior for publishing, retrieving, and administering MVP opportunities.
## Requirements
### Requirement: Opportunity model
The system SHALL define an `opportunity` contract for the MVP with required and optional fields.

Required fields SHALL be `title`, `description`, `organization_name`, `type`, `work_mode`, and at least one contact channel through `contact_email` or `contact_url`.

Optional fields MAY include `organization_url`, `location`, `salary_range`, `seniority`, `skills`, `expires_at`, and `status`.

The `type` field SHALL use one of `full_time`, `part_time`, `contract`, `freelance`, `volunteer`, `project`, or `mentorship`.

The `work_mode` field SHALL use one of `remote`, `hybrid`, or `on_site`.

The `status` field SHALL use one of `draft`, `published`, `closed`, or `archived`, and SHALL default to `published` when omitted by an authorized create request.

#### Scenario: Valid opportunity payload is submitted
- **WHEN** an authorized admin submits an opportunity with all required fields and a valid contact channel
- **THEN** the system accepts the payload as a valid opportunity

#### Scenario: Required opportunity fields are missing
- **WHEN** an authorized admin submits an opportunity without title, description, organization_name, type, work_mode, or contact channel
- **THEN** the system rejects the payload with a validation error response

#### Scenario: Invalid opportunity enum is submitted
- **WHEN** an authorized admin submits an opportunity with an unsupported type, work_mode, or status value
- **THEN** the system rejects the payload with a validation error response

### Requirement: Opportunity creation
The system SHALL allow an authorized admin to create an opportunity.

#### Scenario: Admin creates opportunity
- **WHEN** an authenticated admin sends a valid create opportunity request
- **THEN** the system stores the opportunity and returns it with a generated identifier

#### Scenario: Anonymous user attempts opportunity creation
- **WHEN** a request without valid authentication attempts to create an opportunity
- **THEN** the system rejects the request with an unauthorized response

#### Scenario: Non-admin user attempts opportunity creation
- **WHEN** an authenticated user without role admin attempts to create an opportunity
- **THEN** the system rejects the request with a forbidden response

### Requirement: Opportunity listing
The system SHALL allow any client to list published opportunities with their public fields.

The listing MAY support simple filters for `type`, `work_mode`, and `location`.

The system SHALL paginate the listing using `page` and `page_size` query parameters, defaulting to page `1` and page size `20`, with a maximum page size of `100`.

The system SHALL clamp invalid pagination values to the valid range instead of returning an error, so that `page` below `1` becomes `1`, `page_size` below `1` becomes the default, `page_size` above the maximum becomes the maximum, and non-numeric values fall back to the defaults.

The listing response SHALL be an object with a `data` array of opportunities and a `meta` object containing `page`, `page_size`, `total`, and `total_pages`.

#### Scenario: Public client lists opportunities
- **WHEN** any client requests the opportunities collection
- **THEN** the system returns published opportunities without requiring authentication

#### Scenario: Public client filters opportunities by type
- **WHEN** any client requests the opportunities collection with a valid type filter
- **THEN** the system returns published opportunities matching that type

#### Scenario: Public client lists opportunities with unpublished records present
- **WHEN** draft, closed, or archived opportunities exist
- **THEN** the system excludes them from the public listing

#### Scenario: Public client requests a specific page
- **WHEN** any client requests the opportunities collection with `page` and `page_size`
- **THEN** the system returns only that page of published opportunities and a `meta` object with `page`, `page_size`, `total`, and `total_pages`

#### Scenario: Public client omits pagination parameters
- **WHEN** any client requests the opportunities collection without pagination parameters
- **THEN** the system returns the first page using the default page size and includes pagination metadata

#### Scenario: Public client sends invalid pagination values
- **WHEN** any client requests the opportunities collection with invalid `page` or `page_size` values
- **THEN** the system clamps them to the valid range and still returns a successful paginated response

### Requirement: Opportunity lookup
The system SHALL allow any client to retrieve a published opportunity by identifier and SHALL allow an authenticated admin to retrieve an opportunity in any status.

#### Scenario: Existing published opportunity is requested
- **WHEN** any client requests a published opportunity by an existing identifier
- **THEN** the system returns that opportunity

#### Scenario: Admin requests unpublished opportunity
- **WHEN** an authenticated admin requests a draft, closed, or archived opportunity by an existing identifier
- **THEN** the system returns that opportunity

#### Scenario: Missing opportunity is requested
- **WHEN** any client requests an opportunity by an unknown identifier
- **THEN** the system returns a not found response

#### Scenario: Unpublished opportunity is requested publicly
- **WHEN** any client requests a draft, closed, or archived opportunity without admin access
- **THEN** the system returns a not found response

### Requirement: Opportunity update
The system SHALL allow an authorized admin to update an opportunity's editable fields.

#### Scenario: Admin updates opportunity
- **WHEN** an authenticated admin sends a valid update request for an existing opportunity
- **THEN** the system persists the changes and returns the updated opportunity

#### Scenario: Unauthorized opportunity update attempt
- **WHEN** a request without an admin role attempts to update an opportunity
- **THEN** the system rejects the request with an unauthorized or forbidden response

#### Scenario: Invalid opportunity update payload is submitted
- **WHEN** an authenticated admin submits invalid opportunity fields during update
- **THEN** the system rejects the request with a validation error response

### Requirement: Opportunity deletion
The system SHALL allow an authorized admin to remove an opportunity.

#### Scenario: Admin deletes opportunity
- **WHEN** an authenticated admin requests deletion of an existing opportunity
- **THEN** the system removes the opportunity and returns a successful empty response

#### Scenario: Unauthorized opportunity deletion attempt
- **WHEN** a request without an admin role attempts to delete an opportunity
- **THEN** the system rejects the request with an unauthorized or forbidden response

#### Scenario: Missing opportunity deletion is requested
- **WHEN** an authenticated admin requests deletion of an unknown opportunity
- **THEN** the system returns a not found response

