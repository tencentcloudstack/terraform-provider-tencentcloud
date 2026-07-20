## ADDED Requirements

### Requirement: Purchase a plan for a TEO zone

The resource `tencentcloud_teo_plan_for_zone` SHALL purchase a plan for a zone that has not yet bound a plan via the `CreatePlanForZone` API. As a one-time operation resource, it MUST auto-generate its unique ID. The `CreatePlanForZone` API accepts `ZoneId` and `PlanType` inputs and returns `ResourceNames` and `DealNames`.

#### Scenario: Purchase a plan for a zone

- **WHEN** the user applies the resource with `zone_id` and `plan_type` set
- **THEN** the resource calls `CreatePlanForZone` wrapped in a retry mechanism with the provided `ZoneId` and `PlanType`
- **AND** sets an auto-generated unique resource ID
- **AND** populates the computed `resource_names` and `deal_names` attributes from the response in a nil-safe manner
- **AND** then calls Read

#### Scenario: Plan purchase with empty response

- **WHEN** the `CreatePlanForZone` response is nil or `response.Response` is nil
- **THEN** the resource SHALL return a NonRetryableError to avoid persisting an empty state

### Requirement: Resource schema definition

The system SHALL define the following schema for `tencentcloud_teo_plan_for_zone`:

Required parameters (ForceNew):
- `zone_id` (TypeString): Zone ID
- `plan_type` (TypeString): Plan type to purchase

Computed parameters:
- `resource_names` (TypeList, element TypeString): List of purchased resource names returned by `CreatePlanForZone`
- `deal_names` (TypeList, element TypeString): List of purchased order/deal names returned by `CreatePlanForZone`

#### Scenario: Schema validates required fields

- **WHEN** user omits `zone_id` or `plan_type` in the resource configuration
- **THEN** Terraform SHALL produce a validation error indicating the missing required field

#### Scenario: All input parameters are ForceNew

- **WHEN** user changes `zone_id` or `plan_type` after creation
- **THEN** Terraform SHALL force resource recreation

### Requirement: Empty Read/Update/Delete methods

The system SHALL implement empty Read, Update, and Delete methods for `tencentcloud_teo_plan_for_zone` resource, as this is a RESOURCE_KIND_OPERATION type. There is no query API for the purchase result, so Read is a no-op.

#### Scenario: Read method returns nil

- **WHEN** Terraform calls the Read method
- **THEN** the system SHALL return nil without making any API calls

#### Scenario: Delete method returns nil

- **WHEN** Terraform calls the Delete method
- **THEN** the system SHALL return nil without making any API calls, and the purchased plan is not reverted

### Requirement: Resource registration in provider

The system SHALL register `tencentcloud_teo_plan_for_zone` resource in `tencentcloud/provider.go` with the resource name "tencentcloud_teo_plan_for_zone" and add the corresponding entry in `tencentcloud/provider.md`.

#### Scenario: Resource available in provider

- **WHEN** user references `tencentcloud_teo_plan_for_zone` in Terraform configuration
- **THEN** the provider SHALL recognize and process the resource

### Requirement: Unit tests with gomonkey mock

The system SHALL include unit tests in `resource_tc_teo_plan_for_zone_operation_test.go` using gomonkey to mock the TEO API client calls. Tests SHALL cover the Create operation with a successful response.

#### Scenario: Test successful plan purchase

- **WHEN** unit test runs with mocked `CreatePlanForZone` returning `ResourceNames` and `DealNames`
- **THEN** the test SHALL verify the resource ID is set and computed fields are populated

### Requirement: Resource documentation

The system SHALL include a `resource_tc_teo_plan_for_zone_operation.md` documentation file with a description, and example usage. No Import section is needed for an OPERATION type resource.

#### Scenario: Documentation file exists

- **WHEN** the resource is added
- **THEN** a corresponding .md file SHALL exist in the teo service directory with proper format following the gendoc README guidelines
