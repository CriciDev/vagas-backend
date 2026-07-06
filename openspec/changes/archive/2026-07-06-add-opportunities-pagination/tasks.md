## 1. Model

- [x] 1.1 Add `DefaultPageSize` (20) and `MaxPageSize` (100) constants.
- [x] 1.2 Add a `Pagination` type with `Page` and `PageSize` plus `Offset`/`Limit` helpers.
- [x] 1.3 Add a `PageMeta` type (`page`, `page_size`, `total`, `total_pages`) and an `OpportunityPage` envelope (`data`, `meta`).

## 2. Usecase

- [x] 2.1 Add `NewPagination(page, pageSize)` clamping to the valid range.
- [x] 2.2 Change `List` to accept pagination and return an `OpportunityPage` with computed `total_pages`.

## 3. Repository

- [x] 3.1 Change the `List` interface and Postgres implementation to accept pagination and return the page plus the total count.
- [x] 3.2 Add a `count(*)` query with the same filters and apply `LIMIT`/`OFFSET` to the data query.

## 4. Controller

- [x] 4.1 Read `page` and `page_size` query params and clamp them via `NewPagination`.
- [x] 4.2 Return the `OpportunityPage` envelope.

## 5. Tests

- [x] 5.1 Test default pagination returns the first page with correct metadata.
- [x] 5.2 Test custom `page`/`page_size` returns the right slice.
- [x] 5.3 Test invalid values are clamped and still return 200.
- [x] 5.4 Update existing listing tests to the `{ data, meta }` envelope.

## 6. Verification

- [x] 6.1 Run `make check` and confirm fmt, vet, test, and build pass.
- [x] 6.2 Reference the OpenSpec artifacts in the PR for issue #31.
