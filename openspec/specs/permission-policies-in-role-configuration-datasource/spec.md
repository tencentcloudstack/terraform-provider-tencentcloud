## ADDED Requirements

### Requirement: Query permission policies in role configuration

The system SHALL provide a Terraform data source `tencentcloud_organization_permission_policies_in_role_configuration` that queries the list of permission policies attached to a specified role configuration via the `ListPermissionPoliciesInRoleConfiguration` API.

#### Scenario: Successful query with required parameters only

- **WHEN** user specifies `zone_id` and `role_configuration_id` in the data source configuration
- **THEN** the system calls `ListPermissionPoliciesInRoleConfiguration` API with those parameters and returns `total_counts` and `role_policies` list containing all attached policies

#### Scenario: Query filtered by policy type

- **WHEN** user specifies `zone_id`, `role_configuration_id`, and `role_policy_type` set to "System"
- **THEN** the system calls the API with `RolePolicyType` set to "System" and returns only system policies in the `role_policies` list

#### Scenario: Query filtered by policy name

- **WHEN** user specifies `zone_id`, `role_configuration_id`, and `filter` set to a keyword string
- **THEN** the system calls the API with `Filter` set to that keyword and returns only policies whose names match the filter

#### Scenario: API call failure with retry

- **WHEN** the API call fails with a retryable error
- **THEN** the system retries the call within `tccommon.ReadRetryTimeout` and returns an error only if all retries are exhausted

### Requirement: Role policies output structure

The `role_policies` output attribute SHALL be a list of objects, each containing the following fields mapped from the `RolePolicie` SDK struct.

#### Scenario: Complete policy object returned

- **WHEN** the API returns a policy entry in `RolePolicies`
- **THEN** the data source exposes it as an object with fields: `role_policy_id` (int), `role_policy_name` (string), `role_policy_type` (string), `role_policy_document` (string), `add_time` (string)

#### Scenario: Nil fields handled gracefully

- **WHEN** the API returns a policy entry with some nil fields (e.g., `RolePolicyDocument` is nil for system policies)
- **THEN** the data source omits setting those fields rather than setting them to empty values

### Requirement: Data source registration

The data source SHALL be registered in `tencentcloud/provider.go` under the organization data sources section so that it is available to users as `tencentcloud_organization_permission_policies_in_role_configuration`.

#### Scenario: Provider includes the data source

- **WHEN** a user references `data.tencentcloud_organization_permission_policies_in_role_configuration` in their Terraform configuration
- **THEN** the provider recognizes it as a valid data source and invokes its Read function

### Requirement: Result output file support

The data source SHALL support an optional `result_output_file` parameter that writes the query results to a local file.

#### Scenario: Output file specified

- **WHEN** user specifies `result_output_file` with a file path
- **THEN** the system writes the role_policies results to the specified file after a successful query
