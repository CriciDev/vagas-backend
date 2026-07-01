## Purpose

Define request id propagation behavior for API requests.

## Requirements

### Requirement: Request id propagation
The system SHALL associate every API request with a request id and SHALL return that id in the `X-Request-ID` response header.

The system SHALL reuse the inbound `X-Request-ID` header when it is a non-empty value of at most 128 characters using only the characters `A-Z`, `a-z`, `0-9`, `.`, `_`, or `-`.

The system SHALL generate a new request id when the inbound header is missing or invalid.

The system SHALL make the request id available on the request context so handlers and logs can use the same identifier.

#### Scenario: Request without a request id header
- **WHEN** a client sends a request without an `X-Request-ID` header
- **THEN** the system generates a request id and returns it in the `X-Request-ID` response header

#### Scenario: Request with a valid request id header
- **WHEN** a client sends a request with a valid `X-Request-ID` header
- **THEN** the system reuses that value and returns the same id in the `X-Request-ID` response header

#### Scenario: Request with an invalid request id header
- **WHEN** a client sends a request with an empty, oversized, or malformed `X-Request-ID` header
- **THEN** the system ignores it, generates a new request id, and returns the generated id in the `X-Request-ID` response header
