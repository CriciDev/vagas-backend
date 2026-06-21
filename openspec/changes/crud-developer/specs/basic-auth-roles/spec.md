## ADDED Requirements

### Requirement: Admin authentication
The system SHALL authenticate an admin user and issue a bearer token containing the user's role.

#### Scenario: Valid admin credentials are submitted
- **WHEN** a client submits valid admin credentials to the login endpoint
- **THEN** the system returns a bearer token that identifies the user role as admin

#### Scenario: Invalid credentials are submitted
- **WHEN** a client submits invalid credentials to the login endpoint
- **THEN** the system rejects the request with an unauthorized response

### Requirement: Bearer token authentication
The system SHALL authenticate protected requests using bearer tokens.

#### Scenario: Valid bearer token is provided
- **WHEN** a protected endpoint receives a valid bearer token
- **THEN** the system allows the request to continue with authenticated user context

#### Scenario: Missing bearer token is provided
- **WHEN** a protected endpoint receives no bearer token
- **THEN** the system rejects the request with an unauthorized response

#### Scenario: Invalid bearer token is provided
- **WHEN** a protected endpoint receives an invalid bearer token
- **THEN** the system rejects the request with an unauthorized response

### Requirement: Role based authorization
The system SHALL restrict protected write operations to users with the required role.

#### Scenario: Admin accesses admin-only route
- **WHEN** an authenticated user with role admin accesses an admin-only route
- **THEN** the system allows the request to continue

#### Scenario: Non-admin accesses admin-only route
- **WHEN** an authenticated user without role admin accesses an admin-only route
- **THEN** the system rejects the request with a forbidden response
