# cynosdb-readonly-instance-name-modify Specification

## Purpose
TBD - created by archiving change add-cynosdb-readonly-instance-name-modify. Update Purpose after archive.
## Requirements
### Requirement: instance_name supports in-place update

The `tencentcloud_cynosdb_readonly_instance` resource SHALL allow the `instance_name` field to be updated in-place without recreating the resource. When `instance_name` changes, the provider SHALL call the `ModifyInstanceName` API to update the instance name on the cloud, then refresh state.

#### Scenario: Update instance_name triggers ModifyInstanceName API

- **WHEN** the user updates `instance_name` in the Terraform configuration for an existing `tencentcloud_cynosdb_readonly_instance`
- **THEN** the provider SHALL call `ModifyInstanceName` with the instance ID and the new instance name, and SHALL NOT recreate the resource

#### Scenario: instance_name is no longer ForceNew

- **WHEN** the schema for `instance_name` is inspected
- **THEN** the field SHALL be `Required` and `String`, and SHALL NOT have `ForceNew` set

#### Scenario: Read reflects updated instance_name

- **WHEN** `ModifyInstanceName` succeeds and a subsequent Read is performed
- **THEN** the provider SHALL set `instance_name` in state to the value returned by the `DescribeInstanceById` API

#### Scenario: ModifyInstanceName failure is retried

- **WHEN** the `ModifyInstanceName` API returns an error
- **THEN** the provider SHALL wrap the error with `tccommon.RetryError` and retry within `WriteRetryTimeout`

