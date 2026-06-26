## 1. Product Contract

- [ ] 1.1 Review the proposed `opportunity` required fields: `title`, `description`, `organizationName`, `type`, `workMode`, and contact channel.
- [ ] 1.2 Review the proposed optional fields: `organizationUrl`, `location`, `salaryRange`, `seniority`, `skills`, `expiresAt`, and `status`.
- [ ] 1.3 Validate the initial enum values for `type`, `workMode`, and `status` with the community.
- [ ] 1.4 Confirm that public reads should expose only `published opportunities` in the MVP.

## 2. Permissions And Routes

- [ ] 2.1 Confirm that anonymous users can list and view published opportunities.
- [ ] 2.2 Confirm that create, update, and delete operations remain restricted to role `admin` for the first implementation.
- [ ] 2.3 Confirm future endpoint naming around `/api/opportunities` before code implementation starts.

## 3. Open Decisions

- [ ] 3.1 Discuss whether each opportunity must be linked to a future `organization` entity.
- [ ] 3.2 Discuss whether all proposed opportunity types belong in the MVP or whether the first release should start smaller.
- [ ] 3.3 Discuss whether authenticated organizations should be able to publish opportunities after the admin-only MVP.
- [ ] 3.4 Discuss whether `expiresAt` or automatic expiration should be required.

## 4. Documentation Verification

- [ ] 4.1 Ensure `proposal.md`, `design.md`, and `specs/opportunity-management/spec.md` clearly state that no code implementation is part of this change.
- [ ] 4.2 Run `openspec status --change "define-opportunities"` and verify the change is ready for review.
- [ ] 4.3 Reference the generated OpenSpec artifacts in issue #12 or in the related PR for community feedback.
