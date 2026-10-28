## ADDED Requirements

### Requirement: BindZoneToPlan action
The system SHALL provide a Terraform framework action `tencentcloud_teo_bind_zone_to_plan` that calls the TEO `BindZoneToPlan` API to bind an unbound zone to an existing plan.

#### Scenario: Successful binding of zone to plan
- **WHEN** user invokes the `tencentcloud_teo_bind_zone_to_plan` action with valid `zone_id` and `plan_id`
- **THEN** the system SHALL call `BindZoneToPlan` API with the provided `zone_id` as `ZoneId` and `plan_id` as `PlanId`
- **AND** the system SHALL NOT persist any state (framework actions do not manage state)

#### Scenario: Invoke with retry on transient failure
- **WHEN** the `BindZoneToPlan` API call fails with a transient error
- **THEN** the system SHALL retry the API call using `tccommon.WriteRetryTimeout` and `resource.Retry`
- **AND** wrap non-retryable errors via `tccommon.RetryError`

#### Scenario: API returns zone already bound error
- **WHEN** the `BindZoneToPlan` API returns `InvalidParameter.ZoneHasBeenBound`
- **THEN** the system SHALL surface the error to the user via diagnostics

#### Scenario: Action injected with shared SDK client
- **WHEN** the framework provider configures the action
- **THEN** the system SHALL receive a `*sharedmeta.ProviderMeta` in `Configure` and store `meta.Client` for use in `Invoke`

### Requirement: Schema parameters
The action SHALL expose `zone_id` and `plan_id` parameters, both of type string and Required.

#### Scenario: Zone ID is required
- **WHEN** user invokes the action without specifying `zone_id`
- **THEN** the Terraform plan SHALL fail with a required field error

#### Scenario: Plan ID is required
- **WHEN** user invokes the action without specifying `plan_id`
- **THEN** the Terraform plan SHALL fail with a required field error

### Requirement: No Read/Update/Delete handlers
The action SHALL NOT define Read/Update/Delete handlers, as it is a one-shot operation that does not manage state.

#### Scenario: No Read handler
- **WHEN** the action is defined
- **THEN** there SHALL be no Read/Update/Delete lifecycle methods (action.Invoke is the only execution entry)

### Requirement: Framework registration
The action SHALL be registered in `tencentcloud/framework/registry.go` `actionFactories` and listed under the `Action` subsection of TEO in `provider.md`.

#### Scenario: Action available in provider
- **WHEN** the Terraform provider is initialized
- **THEN** the action `tencentcloud_teo_bind_zone_to_plan` SHALL be available for use

### Requirement: Action documentation
The system SHALL provide a markdown documentation file `action_tc_teo_bind_zone_to_plan.md` with a one-line description mentioning TEO and an example usage showing `zone_id` and `plan_id`.

#### Scenario: Documentation exists
- **WHEN** the action is created
- **THEN** a `.md` file `action_tc_teo_bind_zone_to_plan.md` SHALL exist with a one-line description mentioning TEO and an example usage with `zone_id` and `plan_id`

### Requirement: Unit tests
The system SHALL provide unit tests in `action_tc_teo_bind_zone_to_plan_test.go` using gomonkey to mock the cloud API, covering Invoke success, Invoke API error, Invoke missing parameters, Metadata, Schema, Configure, and Factory.

#### Scenario: Unit tests pass
- **WHEN** `go test` is run with `-gcflags=all=-l` on the test file
- **THEN** all test cases for Invoke success, API error, missing parameters, Metadata, Schema, Configure, and Factory SHALL pass

#### Scenario: Invoke success test case
- **WHEN** a test simulates a successful `BindZoneToPlan` call
- **THEN** the mocked `BindZoneToPlanWithContext` SHALL be invoked with the provided `ZoneId` and `PlanId`
- **AND** the system SHALL report no diagnostics