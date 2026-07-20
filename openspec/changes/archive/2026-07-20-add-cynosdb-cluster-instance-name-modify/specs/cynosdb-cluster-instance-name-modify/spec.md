## ADDED Requirements

### Requirement: instance_name supports in-place update for cluster resources

The `tencentcloud_cynosdb_cluster` and `tencentcloud_cynosdb_cluster_v2` resources SHALL allow the `instance_name` field to be updated in-place. When `instance_name` changes, the provider SHALL call the `ModifyInstanceName` API to update the read-write instance name on the cloud, then refresh state.

#### Scenario: Update instance_name triggers ModifyInstanceName API

- **WHEN** the user updates `instance_name` in the Terraform configuration for an existing `tencentcloud_cynosdb_cluster` or `tencentcloud_cynosdb_cluster_v2`
- **THEN** the provider SHALL call `ModifyInstanceName` with the read-write instance ID and the new instance name, and SHALL NOT recreate the resource

#### Scenario: instance_name is Optional and Computed

- **WHEN** the schema for `instance_name` is inspected
- **THEN** the field SHALL be `Optional`, `Computed`, and `String`, and SHALL NOT have `ForceNew` set

#### Scenario: Read reflects updated instance_name

- **WHEN** `ModifyInstanceName` succeeds and a subsequent Read is performed
- **THEN** the provider SHALL set `instance_name` in state to the value returned by the `DescribeInstanceById` API for the read-write instance

#### Scenario: ModifyInstanceName failure is retried

- **WHEN** the `ModifyInstanceName` API returns an error
- **THEN** the provider SHALL wrap the error with `tccommon.RetryError` and retry within `WriteRetryTimeout`
