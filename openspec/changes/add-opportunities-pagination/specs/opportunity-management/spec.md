## MODIFIED Requirements

### Requirement: Opportunity listing
The system SHALL allow any client to list published opportunities with their public fields.

The listing MAY support simple filters for `type`, `workMode`, and `location`.

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
