## ADDED Requirements

### Requirement: Increase Plan Quota Action
The system SHALL provide a Terraform framework Action `tencentcloud_teo_increase_plan_quota` (terraform-plugin-framework `action.Action`) that calls the TEO `IncreasePlanQuota` API to purchase additional plan quotas. The action is a one-time operation that does not persist any state.

#### Scenario: Successful quota increase
- **WHEN** a user invokes the action with valid `plan_id`, `quota_type`, and `quota_number`
- **THEN** the provider calls the `IncreasePlanQuota` API and reports the `deal_name` (order number) back to the user via a warning diagnostic

#### Scenario: Missing required parameter
- **WHEN** a user invokes the action without `plan_id`, `quota_type`, or `quota_number`
- **THEN** the provider returns an error diagnostic before calling the API

#### Scenario: API returns an error
- **WHEN** the `IncreasePlanQuota` API call fails (e.g., insufficient balance, invalid quota type)
- **THEN** the provider returns the API error to the user as an error diagnostic, with retry on retriable errors

#### Scenario: API returns nil response
- **WHEN** the `IncreasePlanQuota` API call returns a nil response (or nil `DealName`)
- **THEN** the provider returns a non-retryable error indicating the response was empty

#### Scenario: Action does not persist state
- **WHEN** the action invocation completes (success or failure)
- **THEN** no resource state is written to tfstate; the action is a one-time operation

#### Scenario: Action is registered via framework registry
- **WHEN** the provider binary starts
- **THEN** the `tencentcloud_teo_increase_plan_quota` action is registered through `tencentcloud/framework/registry.go` (`actionFactories`) and exposed via `Provider.Actions()`, not through the SDKv2 `ResourcesMap`