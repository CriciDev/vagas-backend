## 1. Middleware

- [x] 1.1 Create package `internal/middleware` with a `RequestID` gin middleware.
- [x] 1.2 Reuse the inbound `X-Request-ID` header when it is non-empty, at most 128 chars, and matches the safe charset `A-Za-z0-9._-`.
- [x] 1.3 Generate a new id with `crypto/rand` when the inbound header is missing or invalid.
- [x] 1.4 Set the id on the response `X-Request-ID` header and on the request context.
- [x] 1.5 Expose a helper to read the id from a gin context for handlers and logs.

## 2. Router Integration

- [x] 2.1 Replace `gin.Default()` with `gin.New()` in `cmd/server/main.go`.
- [x] 2.2 Register middlewares in order: `RequestID`, access logger, recovery.
- [x] 2.3 Configure the access logger to include the request id from the context.

## 3. Tests

- [x] 3.1 Test that a request without the header receives a generated id in the response header.
- [x] 3.2 Test that a request with a valid `X-Request-ID` header has that id echoed back.
- [x] 3.3 Test that an invalid inbound id is replaced by a generated one.

## 4. Verification

- [x] 4.1 Run `make check` and confirm fmt, vet, test, and build pass.
- [x] 4.2 Reference the generated OpenSpec artifacts in the related PR for issue #43.
