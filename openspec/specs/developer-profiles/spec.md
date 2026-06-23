## Purpose

Define the first public developer profile CRUD behavior for the community jobs platform.

## Requirements

### Requirement: Developer profile creation
The system SHALL allow an authorized admin to create a developer profile with at least name, email, skills, and availability status.

#### Scenario: Admin creates developer profile
- **WHEN** an authenticated admin sends a valid create developer request
- **THEN** the system stores the developer profile and returns it with a generated identifier

#### Scenario: Anonymous user attempts profile creation
- **WHEN** a request without valid authentication attempts to create a developer profile
- **THEN** the system rejects the request with an unauthorized response

### Requirement: Developer profile listing
The system SHALL allow any client to list developer profiles with their public fields.

#### Scenario: Public client lists developers
- **WHEN** any client requests the developers collection
- **THEN** the system returns the available developer profiles without requiring authentication

### Requirement: Developer profile lookup
The system SHALL allow any client to retrieve a developer profile by identifier.

#### Scenario: Existing developer is requested
- **WHEN** any client requests a developer profile by an existing identifier
- **THEN** the system returns that developer profile

#### Scenario: Missing developer is requested
- **WHEN** any client requests a developer profile by an unknown identifier
- **THEN** the system returns a not found response

### Requirement: Developer profile update
The system SHALL allow an authorized admin to update a developer profile's editable fields.

#### Scenario: Admin updates developer profile
- **WHEN** an authenticated admin sends a valid update request for an existing developer profile
- **THEN** the system persists the changes and returns the updated profile

#### Scenario: Unauthorized update attempt
- **WHEN** a request without an admin role attempts to update a developer profile
- **THEN** the system rejects the request with an unauthorized or forbidden response

### Requirement: Developer profile deletion
The system SHALL allow an authorized admin to remove a developer profile.

#### Scenario: Admin deletes developer profile
- **WHEN** an authenticated admin requests deletion of an existing developer profile
- **THEN** the system removes the profile and returns a successful empty response

#### Scenario: Unauthorized deletion attempt
- **WHEN** a request without an admin role attempts to delete a developer profile
- **THEN** the system rejects the request with an unauthorized or forbidden response

### Requirement: Developer profile validation
The system SHALL validate developer profile input before persisting changes.

#### Scenario: Invalid developer payload is submitted
- **WHEN** a create or update request omits required developer fields or uses invalid values
- **THEN** the system rejects the request with a validation error response
