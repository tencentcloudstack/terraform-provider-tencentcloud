# ckafka-instance-resource Specification

## Purpose
TBD - created by archiving change add-ckafka-instance-delete-protection. Update Purpose after archive.
## Requirements
### Requirement: CKafka Instance Resource Schema
The system SHALL define a Terraform resource `tencentcloud_ckafka_instance` (RESOURCE_KIND_GENERAL) covering the full CRUD lifecycle of a CKafka instance. The schema SHALL include a `delete_protection_enable` field with the following characteristics:
- `delete_protection_enable` (Optional, Computed, TypeInt): 实例删除保护开关，取值 `1`（开启）/ `0`（关闭）。未显式配置时由云端回填，不触发 plan diff。

#### Scenario: delete_protection_enable is optional and computed
- **WHEN** a user does not set `delete_protection_enable` in the Terraform configuration
- **THEN** the system SHALL populate `delete_protection_enable` from the `DescribeInstanceAttributes` API response (Computed behavior) without triggering a plan diff

#### Scenario: delete_protection_enable accepts user input
- **WHEN** a user sets `delete_protection_enable = 1` or `delete_protection_enable = 0` in the Terraform configuration
- **THEN** the schema SHALL accept the integer value and mark it as a candidate for create/update

### Requirement: CKafka Instance Create Operation
The system SHALL, after successfully creating a CKafka instance, call `ModifyInstanceAttributes` API to set instance attributes. When the user has explicitly configured `delete_protection_enable`, the system SHALL fill `request.DeleteProtectionEnable` with the configured int64 value (`1` or `0`) and include it in the same `ModifyInstanceAttributes` call that sets other post-creation attributes.

#### Scenario: Create with delete_protection_enable set to 1
- **WHEN** the user sets `delete_protection_enable = 1` in the Terraform configuration and creates the instance
- **THEN** the system SHALL call `ModifyInstanceAttributes` with `DeleteProtectionEnable = 1` after the instance is created

#### Scenario: Create with delete_protection_enable set to 0
- **WHEN** the user sets `delete_protection_enable = 0` in the Terraform configuration and creates the instance
- **THEN** the system SHALL call `ModifyInstanceAttributes` with `DeleteProtectionEnable = 0` (using GetOkExists so the explicit 0 value is not skipped)

#### Scenario: Create without delete_protection_enable configured
- **WHEN** the user does not set `delete_protection_enable` in the Terraform configuration
- **THEN** the system SHALL NOT fill `DeleteProtectionEnable` in the `ModifyInstanceAttributes` request, preserving the cloud-side default

### Requirement: CKafka Instance Update Operation
The system SHALL, when `delete_protection_enable` changes in the Terraform configuration, call `ModifyInstanceAttributes` API with `InstanceId` and the updated `DeleteProtectionEnable` value. This field SHALL be mutable in-place (not ForceNew, not in `immutableArgs`).

#### Scenario: Update delete_protection_enable from 0 to 1
- **WHEN** `delete_protection_enable` changes from `0` to `1` in the Terraform configuration
- **THEN** the system SHALL call `ModifyInstanceAttributes` with `DeleteProtectionEnable = 1`
- **AND** the update SHALL be performed in-place without resource recreation

#### Scenario: Update delete_protection_enable from 1 to 0
- **WHEN** `delete_protection_enable` changes from `1` to `0` in the Terraform configuration
- **THEN** the system SHALL call `ModifyInstanceAttributes` with `DeleteProtectionEnable = 0`

#### Scenario: delete_protection_enable not changed
- **WHEN** `delete_protection_enable` is not changed during an update
- **THEN** the system SHALL NOT include `DeleteProtectionEnable` in the `ModifyInstanceAttributes` request for this field

### Requirement: CKafka Instance Read Operation
The system SHALL read the `delete_protection_enable` value from the `DescribeInstanceAttributes` API response (`InstanceAttributesResponse.DeleteProtectionEnable`). Before calling `d.Set("delete_protection_enable", ...)`, the system SHALL check that the response field is not nil; if nil, the system SHALL skip the set operation.

#### Scenario: Read delete_protection_enable from DescribeInstanceAttributes
- **WHEN** the Read operation calls `DescribeInstanceAttributes` and the response `DeleteProtectionEnable` is not nil
- **THEN** the system SHALL set `delete_protection_enable` in Terraform state to the returned value

#### Scenario: Read with nil DeleteProtectionEnable
- **WHEN** the Read operation calls `DescribeInstanceAttributes` and the response `DeleteProtectionEnable` is nil
- **THEN** the system SHALL skip `d.Set("delete_protection_enable", ...)` to avoid a nil pointer issue

### Requirement: CKafka Instance Unit Tests
The system SHALL provide unit tests in `resource_tc_ckafka_instance_test.go` using gomonkey to mock cloud API calls, covering the `delete_protection_enable` handling in Create, Update, and Read operations.

#### Scenario: Unit tests pass
- **WHEN** `go test` is run with `-gcflags=all=-l` on the test file
- **THEN** all test cases for `delete_protection_enable` Create (set to 1 and 0), Update (change 0→1 and 1→0), and Read SHALL pass

#### Scenario: Create with delete_protection_enable test
- **WHEN** a test simulates creating an instance with `delete_protection_enable = 1`
- **THEN** the mocked `ModifyInstanceAttributes` SHALL be invoked with `DeleteProtectionEnable = 1`

#### Scenario: Update delete_protection_enable test
- **WHEN** a test simulates updating `delete_protection_enable` from `0` to `1`
- **THEN** the mocked `ModifyInstanceAttributes` SHALL be invoked with `DeleteProtectionEnable = 1`
- **AND** the test SHALL assert the API was called

### Requirement: CKafka Instance Resource Documentation
The system SHALL provide a markdown documentation file `resource_tc_ckafka_instance.md` with a one-line description mentioning CKafka, example usage, and import section. The example usage SHALL demonstrate the `delete_protection_enable` field.

#### Scenario: Documentation exists
- **WHEN** the resource is created
- **THEN** a `.md` file SHALL exist with a one-line description mentioning CKafka, example usage including `delete_protection_enable`, and import section showing the instance ID format

