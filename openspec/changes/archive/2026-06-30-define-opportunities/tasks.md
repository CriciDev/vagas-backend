## 1. Product Contract

- [x] 1.1 Review the proposed `opportunity` required fields: `title`, `description`, `organization_name`, `type`, `work_mode`, and contact channel.
- [x] 1.2 Review the proposed optional fields: `organization_url`, `location`, `salary_range`, `seniority`, `skills`, `expires_at`, and `status`.
- [x] 1.3 Validate the initial enum values for `type`, `work_mode`, and `status` with the community.
- [x] 1.4 Confirm that public reads should expose only `published opportunities` in the MVP.

## 2. Permissions And Routes

- [x] 2.1 Confirm that anonymous users can list and view published opportunities.
- [x] 2.2 Confirm that create, update, and delete operations remain restricted to role `admin` for the first implementation.
- [x] 2.3 Confirm future endpoint naming around `/api/opportunities` before code implementation starts.

## 3. Open Decisions

- [x] 3.1 Discuss whether each opportunity must be linked to a future `organization` entity.
- [x] 3.2 Discuss whether all proposed opportunity types belong in the MVP or whether the first release should start smaller.
- [x] 3.3 Discuss whether authenticated organizations should be able to publish opportunities after the admin-only MVP.
- [x] 3.4 Discuss whether `expires_at` or automatic expiration should be required.

## 4. Documentation Verification

- [x] 4.1 Ensure `proposal.md`, `design.md`, and `specs/opportunity-management/spec.md` clearly state that no code implementation is part of this change.
- [x] 4.2 Run `openspec status --change "define-opportunities"` and verify the change is ready for review.
- [x] 4.3 Reference the generated OpenSpec artifacts in issue #12 or in the related PR for community feedback.
