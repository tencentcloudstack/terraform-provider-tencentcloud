## ADDED Requirements

### Requirement: Increase Plan Quota Operation
The system SHALL provide a Terraform resource `tencentcloud_teo_increase_plan_quota` that calls the TEO `IncreasePlanQuota` API to purchase additional plan quotas.

#### Scenario: Successful quota increase
- **WHEN** a user applies a configuration with valid `plan_id`, `quota_type`, and `quota_number`
- **THEN** the provider calls the `IncreasePlanQuota` API and returns the `deal_name` (order number) as a computed attribute

#### Scenario: Missing required parameter
- **WHEN** a user applies a configuration without `plan_id`, `quota_type`, or `quota_number`
- **THEN** Terraform returns a validation error before calling the API

#### Scenario: API returns an error
- **WHEN** the `IncreasePlanQuota` API call fails (e.g., insufficient balance, invalid quota type)
- **THEN** the provider returns the API error to the user with retry on retriable errors

#### Scenario: API returns nil response
- **WHEN** the `IncreasePlanQuota` API call returns a nil response
- **THEN** the provider returns a non-retryable error indicating the response was nil

#### Scenario: OPERATION resource read is no-op
- **WHEN** Terraform performs a read (refresh) on the `tencentcloud_teo_increase_plan_quota` resource
- **THEN** the read function returns nil without modifying state

#### Scenario: OPERATION resource delete is no-op
- **WHEN** Terraform performs a delete on the `tencentcloud_teo_increase_plan_quota` resource
- **THEN** the delete function returns nil without performing any API call