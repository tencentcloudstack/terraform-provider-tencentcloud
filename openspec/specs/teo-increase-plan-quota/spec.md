## ADDED Requirements

### Requirement: Increase TEO plan quota

The resource `tencentcloud_teo_increase_plan_quota` SHALL call the `IncreasePlanQuota` API to purchase additional quota for a TEO plan. As a one-time operation resource, it MUST auto-generate its unique ID via `helper.BuildToken()`. The resource SHALL accept `plan_id`, `quota_type`, and `quota_number` as required input parameters.

#### Scenario: Increase plan quota successfully

- **WHEN** the user applies the resource with valid `plan_id`, `quota_type`, and `quota_number`
- **THEN** the resource calls `IncreasePlanQuota` wrapped in a retry mechanism
- **AND** sets an auto-generated unique resource ID
- **AND** exposes the returned `deal_name` as a computed attribute

### Requirement: Return deal name

The resource SHALL capture the `DealName` from the `IncreasePlanQuota` API response and expose it as a computed attribute `deal_name`.

#### Scenario: Read deal name after creation

- **WHEN** the resource is created successfully
- **THEN** the `deal_name` attribute is populated from the API response in a nil-safe manner

### Requirement: Delete is a no-op

Because quota purchases cannot be reverted through the API, the delete lifecycle SHALL only remove the resource from Terraform state without calling any API.

#### Scenario: Destroy the operation resource

- **WHEN** the resource is destroyed
- **THEN** no API is called and the resource is removed from state