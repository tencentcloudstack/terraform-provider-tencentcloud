## ADDED Requirements

### Requirement: BindZoneToPlan operation resource
The system SHALL provide a Terraform operation resource `tencentcloud_teo_bind_zone_to_plan` that calls the TEO `BindZoneToPlan` API to bind an unbound zone to an existing plan.

#### Scenario: Successful binding of zone to plan
- **WHEN** user creates a `tencentcloud_teo_bind_zone_to_plan` resource with valid `zone_id` and `plan_id`
- **THEN** the system SHALL call `BindZoneToPlan` API with the provided `zone_id` as `ZoneId` and `plan_id` as `PlanId`
- **AND** the resource ID SHALL be set to a generated token via `helper.BuildToken()`

#### Scenario: Create with retry on transient failure
- **WHEN** the `BindZoneToPlan` API call fails with a transient error
- **THEN** the system SHALL retry the API call using `tccommon.WriteRetryTimeout` and `resource.Retry`
- **AND** wrap non-retryable errors via `tccommon.RetryError`

#### Scenario: Set ID outside retry block
- **WHEN** the `BindZoneToPlan` API call succeeds
- **THEN** the system SHALL call `d.SetId(helper.BuildToken())` outside the retry block and after error handling

#### Scenario: API returns zone already bound error
- **WHEN** the `BindZoneToPlan` API returns `InvalidParameter.ZoneHasBeenBound`
- **THEN** the system SHALL surface the error to the user without clearing the resource ID

### Requirement: Schema parameters
The resource SHALL expose `zone_id` and `plan_id` parameters, both of type string, Required and ForceNew.

#### Scenario: Zone ID is required
- **WHEN** user creates the resource without specifying `zone_id`
- **THEN** the Terraform plan SHALL fail with a required field error

#### Scenario: Plan ID is required
- **WHEN** user creates the resource without specifying `plan_id`
- **THEN** the Terraform plan SHALL fail with a required field error

#### Scenario: Zone ID changes trigger recreation
- **WHEN** user changes the `zone_id` value in an existing resource configuration
- **THEN** the system SHALL destroy and recreate the resource

#### Scenario: Plan ID changes trigger recreation
- **WHEN** user changes the `plan_id` value in an existing resource configuration
- **THEN** the system SHALL destroy and recreate the resource

### Requirement: No-op Read handler
The Read handler SHALL be a no-op that returns nil without making any API calls.

#### Scenario: Read operation
- **WHEN** Terraform performs a refresh/read on the resource
- **THEN** the system SHALL return nil without calling any cloud API

### Requirement: No-op Delete handler
The Delete handler SHALL be a no-op that returns nil without making any API calls.

#### Scenario: Delete operation
- **WHEN** user destroys the resource
- **THEN** the system SHALL return nil without calling any cloud API

### Requirement: No Update handler and no Importer
The resource SHALL NOT define an Update handler and SHALL NOT support import, as it is a one-shot operation resource.

#### Scenario: No Update handler
- **WHEN** the resource schema is defined
- **THEN** the `Update` field SHALL be nil

#### Scenario: Not importable
- **WHEN** user attempts to import the resource
- **THEN** the import SHALL not be supported

### Requirement: Provider registration
The resource SHALL be registered in `provider.go` ResourcesMap and listed in `provider.md`.

#### Scenario: Resource available in provider
- **WHEN** the Terraform provider is initialized
- **THEN** the resource `tencentcloud_teo_bind_zone_to_plan` SHALL be available for use

### Requirement: Resource documentation
The system SHALL provide a markdown documentation file `resource_tc_teo_bind_zone_to_plan_operation.md` with a one-line description mentioning TEO and an example usage showing `zone_id` and `plan_id`.

#### Scenario: Documentation exists
- **WHEN** the resource is created
- **THEN** a `.md` file SHALL exist with a one-line description mentioning TEO and an example usage with `zone_id` and `plan_id`

### Requirement: Unit tests
The system SHALL provide unit tests in `resource_tc_teo_bind_zone_to_plan_operation_test.go` using gomonkey to mock the cloud API, covering Create success, API error, no-op Read, no-op Delete, and schema validation.

#### Scenario: Unit tests pass
- **WHEN** `go test` is run with `-gcflags=all=-l` on the test file
- **THEN** all test cases for Create success, API error, Read, Delete, and schema SHALL pass

#### Scenario: Create success test case
- **WHEN** a test simulates a successful `BindZoneToPlan` call
- **THEN** the mocked `BindZoneToPlanWithContext` SHALL be invoked with the provided `ZoneId` and `PlanId`
- **AND** the test SHALL assert the resource ID is non-empty
