# mongodb-instance-cpu Specification

## Purpose
TBD - created by archiving change add-mongodb-instance-cpu-param. Update Purpose after archive.
## Requirements
### Requirement: MongoDB Instance CPU Parameter
The `tencentcloud_mongodb_instance` resource SHALL provide a `cpu` schema field to manage the instance CPU core count, allowing users to change the CPU specification of a running MongoDB instance through in-place update via the `ModifyDBInstanceSpec` API.

#### Scenario: Schema defines optional and computed cpu field
- **GIVEN** the `tencentcloud_mongodb_instance` resource schema definition
- **THEN** `cpu` SHALL be a field of type `schema.TypeInt`
- **AND** `cpu` SHALL have `Optional: true` and `Computed: true`
- **AND** `cpu` SHALL NOT have `ForceNew` set (changes are applied via in-place update)
- **AND** the field description SHALL document that the unit is C and that the supported CPU specifications can be obtained via the `DescribeSpecInfo` API

#### Scenario: Update triggers ModifyDBInstanceSpec when cpu changes
- **GIVEN** an existing `tencentcloud_mongodb_instance` resource
- **WHEN** the user updates only the `cpu` field (without changing `memory`, `volume`, or `node_num`)
- **THEN** the system SHALL detect `d.HasChange("cpu")` and enter the spec-modification branch
- **AND** the system SHALL read the current `cpu` value from the schema
- **AND** the system SHALL pass `cpu` into the params map of `MongodbService.UpgradeInstance`
- **AND** the system SHALL call the `ModifyDBInstanceSpec` API with the `Cpu` request field set to the new CPU value
- **AND** the system SHALL poll the returned deal id via `DescribeDBInstanceDeal` until the deal succeeds (when `in_maintenance == 0`)
- **AND** the resource state SHALL be updated after a successful modification

#### Scenario: Update triggers ModifyDBInstanceSpec when cpu changes together with other spec params
- **GIVEN** an existing `tencentcloud_mongodb_instance` resource
- **WHEN** the user updates `cpu` together with `memory`, `volume`, or `node_num`
- **THEN** the system SHALL include the `cpu` value alongside the other changed spec parameters in a single `ModifyDBInstanceSpec` call
- **AND** the `Cpu` request field SHALL be set from the user-specified `cpu` value

#### Scenario: cpu omitted in update keeps existing behavior
- **GIVEN** an existing `tencentcloud_mongodb_instance` resource where `cpu` was not configured
- **WHEN** the user updates fields unrelated to spec modification (e.g., tags, project_id)
- **THEN** the system SHALL NOT call `ModifyDBInstanceSpec`
- **AND** the `cpu` field SHALL remain unset in the schema (populated only via Read)

#### Scenario: Service layer maps cpu to ModifyDBInstanceSpecRequest.Cpu
- **GIVEN** the `MongodbService.UpgradeInstance` service method
- **WHEN** `params["cpu"]` is present with an integer value
- **THEN** the system SHALL set `request.Cpu` to a `*int64` pointer constructed from the integer value
- **AND** the type conversion SHALL use the `int64` helper matching the `ModifyDBInstanceSpecRequest.Cpu` field type
- **WHEN** `params["cpu"]` is not present
- **THEN** the system SHALL leave `request.Cpu` unset (nil) so the cloud API defaults to the current CPU size

### Requirement: MongoDB Instance CPU Read Backfill
The resource Read operation SHALL populate the `cpu` field from the cloud API response so that Terraform state reflects the actual CPU core count of the instance.

#### Scenario: Read populates cpu from DescribeDBInstances
- **GIVEN** an existing `tencentcloud_mongodb_instance` resource
- **WHEN** the Read operation queries the instance via `DescribeDBInstances` (through `DescribeInstanceById`)
- **AND** the returned `InstanceDetail.CpuNum` is not nil
- **THEN** the system SHALL set the `cpu` field in state to `int(*instance.CpuNum)`

#### Scenario: Read handles nil CpuNum gracefully
- **GIVEN** an existing `tencentcloud_mongodb_instance` resource
- **WHEN** the Read operation queries the instance
- **AND** the returned `InstanceDetail.CpuNum` is nil
- **THEN** the system SHALL NOT call `d.Set("cpu", ...)` (skip setting the field)
- **AND** the Read operation SHALL NOT fail solely because `CpuNum` is nil

#### Scenario: CpuNum not added to mandatory nil check
- **GIVEN** the Read operation's `CheckNil` validation list
- **THEN** `CpuNum` SHALL NOT be added to the mandatory nil-check map
- **AND** existing instances that do not return `CpuNum` SHALL still be readable without error

### Requirement: CPU Parameter Testing
The resource SHALL have unit tests using gomonkey mocks (not Terraform test suites) covering the new `cpu` parameter behavior, runnable via `go test` with `-gcflags=all=-l`.

#### Scenario: Unit test for cpu change triggering spec modification
- **GIVEN** the `tencentcloud_mongodb_instance` resource update function
- **WHEN** a unit test simulates a `cpu` change in the schema
- **THEN** the test SHALL mock the `ModifyDBInstanceSpec` and `DescribeDBInstanceDeal` cloud APIs using gomonkey
- **AND** verify that the update flow calls `ModifyDBInstanceSpec` with the expected `Cpu` value
- **AND** verify the update completes without error

#### Scenario: Unit test for cpu read backfill
- **GIVEN** the `tencentcloud_mongodb_instance` resource read function
- **WHEN** a unit test mocks `DescribeDBInstances` to return an instance with `CpuNum` set
- **THEN** the test SHALL verify that `d.Set("cpu", ...)` is called with the correct CPU value
- **AND** verify the read completes without error

#### Scenario: Unit test for nil CpuNum read
- **GIVEN** the `tencentcloud_mongodb_instance` resource read function
- **WHEN** a unit test mocks `DescribeDBInstances` to return an instance with `CpuNum` nil
- **THEN** the test SHALL verify the read operation does not fail
- **AND** verify the `cpu` field is not set

### Requirement: CPU Parameter Documentation
The resource documentation SHALL be updated to describe the new `cpu` parameter, generated via the `make doc` command during the finalize phase.

#### Scenario: cpu field documented
- **GIVEN** the resource documentation file `resource_tc_mongodb_instance.md`
- **THEN** the documentation SHALL include the `cpu` parameter description
- **AND** the description SHALL state the unit is C and that supported specifications can be obtained via the `DescribeSpecInfo` API
- **AND** the documentation SHALL NOT include manual `Argument Reference` or `Attribute Reference` sections (auto-generated)

#### Scenario: Example usage includes cpu
- **GIVEN** the resource documentation
- **THEN** the example usage SHALL demonstrate specifying the `cpu` parameter in an update scenario

