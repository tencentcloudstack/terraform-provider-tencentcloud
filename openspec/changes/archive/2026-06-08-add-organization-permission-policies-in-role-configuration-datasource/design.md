## Context

The TencentCloud Organization service provides an Identity Center (CIC) feature that manages role configurations and their associated permission policies. The `ListPermissionPoliciesInRoleConfiguration` API allows querying all permission policies (system and custom) attached to a specific role configuration within a zone.

Currently, no Terraform data source exists to expose this information. Users who need to reference permission policies in their Terraform configurations must look them up manually.

The existing provider already has organization-related data sources in `tencentcloud/services/tco/` and uses the `organization/v20210331` SDK package.

## Goals / Non-Goals

**Goals:**
- Provide a read-only data source `tencentcloud_organization_permission_policies_in_role_configuration` that queries permission policies for a given role configuration.
- Follow existing patterns in the provider (retry logic, error handling, schema conventions).
- Support optional filtering by policy type and policy name.

**Non-Goals:**
- Pagination handling (the API does not have pagination parameters; it returns all results in one call).
- Write operations (this is a data source only).
- Modifying any existing resources or data sources.

## Decisions

1. **File location**: Place the data source in `tencentcloud/services/tco/` following existing organization data sources convention.
   - Rationale: All organization-related resources are already in this directory.

2. **No service layer function**: Call the SDK directly in the Read function with retry wrapper, similar to other simple data sources.
   - Rationale: The API is a single call with no pagination. Adding a service layer function adds unnecessary indirection.
   - Alternative considered: Adding a method to `service_tencentcloud_tco.go`. Rejected because the call is simple enough to inline with retry.

3. **Schema design**: Expose all request parameters as Optional inputs and all response fields as Computed outputs.
   - `zone_id` (Required): Space ID, needed to identify the zone.
   - `role_configuration_id` (Required): Role configuration ID, needed to identify which role configuration to query.
   - `role_policy_type` (Optional): Filter by policy type (System/Custom).
   - `filter` (Optional): Filter by policy name keyword.
   - `total_counts` (Computed): Total number of policies returned.
   - `role_policies` (Computed, List): List of policy objects with fields: role_policy_id, role_policy_name, role_policy_type, role_policy_document, add_time.
   - `result_output_file` (Optional): Standard output file parameter for saving results.

4. **ID strategy**: Use a composite ID of `zone_id#role_configuration_id` since the data source is uniquely identified by these two required parameters.

5. **Retry logic**: Use `tccommon.ReadRetryTimeout` with `resource.Retry` for the API call, wrapping errors with `tccommon.RetryError`.

## Risks / Trade-offs

- [Risk] API returns all results without pagination → If a role configuration has many policies, the response could be large. Mitigation: This is unlikely in practice as role configurations typically have a limited number of policies.
- [Risk] The `RolePolicie` struct name has a typo in the SDK (missing 'y') → Mitigation: Use the SDK type as-is; this is a vendor dependency we don't control.
