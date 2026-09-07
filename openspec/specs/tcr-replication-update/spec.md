## ADDED Requirements

### Requirement: Resource supports in-place update of replication rules
The `tencentcloud_tcr_replication` resource SHALL support an `Update` method that calls the TCR `ModifyReplication` API to modify an existing replication rule without destroying and recreating the resource.

#### Scenario: Update description in place
- **WHEN** a user changes the `description` field of an existing `tencentcloud_tcr_replication` resource
- **THEN** the provider SHALL call `ModifyReplication` with the new `Description` value and the existing `SourceRegistryId` and `RuleName` derived from the resource ID
- **AND** the resource SHALL NOT be destroyed and recreated

#### Scenario: Update rule destination namespace in place
- **WHEN** a user changes the `rule.dest_namespace` field of an existing `tencentcloud_tcr_replication` resource
- **THEN** the provider SHALL call `ModifyReplication` with the updated `Rule.DestNamespace` value
- **AND** the resource SHALL NOT be destroyed and recreated

#### Scenario: Update rule override flag in place
- **WHEN** a user changes the `rule.override` field of an existing `tencentcloud_tcr_replication` resource
- **THEN** the provider SHALL call `ModifyReplication` with the updated `Rule.Override` value
- **AND** the resource SHALL NOT be destroyed and recreated

#### Scenario: Update rule filters in place
- **WHEN** a user changes the `rule.filters` list of an existing `tencentcloud_tcr_replication` resource
- **THEN** the provider SHALL call `ModifyReplication` with the updated `Rule.Filters` values
- **AND** the resource SHALL NOT be destroyed and recreated

#### Scenario: Update rule deletion flag in place
- **WHEN** a user changes the `rule.deletion` field of an existing `tencentcloud_tcr_replication` resource
- **THEN** the provider SHALL call `ModifyReplication` with the updated `Rule.Deletion` value
- **AND** the resource SHALL NOT be destroyed and recreated

### Requirement: Non-updatable fields still trigger ForceNew
The `tencentcloud_tcr_replication` resource SHALL preserve `ForceNew: true` for fields that cannot be modified via the `ModifyReplication` API.

#### Scenario: Changing source_registry_id triggers destroy and recreate
- **WHEN** a user changes the `source_registry_id` field
- **THEN** the resource SHALL be destroyed and recreated because `ForceNew: true` is preserved

#### Scenario: Changing rule.name triggers destroy and recreate
- **WHEN** a user changes the `rule.name` field
- **THEN** the resource SHALL be destroyed and recreated because `rule.name` serves as the rule identifier and `ForceNew: true` is preserved

#### Scenario: Changing destination_region_id triggers destroy and recreate
- **WHEN** a user changes the `destination_region_id` field
- **THEN** the resource SHALL be destroyed and recreated because the API does not support in-place modification of this field

#### Scenario: Changing peer_replication_option triggers destroy and recreate
- **WHEN** a user changes any field within `peer_replication_option`
- **THEN** the resource SHALL be destroyed and recreated because the API does not support in-place modification of cross-account settings

### Requirement: Update method follows standard error handling and retry pattern
The `Update` method SHALL follow the same pattern as the existing `Create` and `Delete` methods: wrap the API call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`, use `tccommon.RetryError(e)` on failure, log with `logId`, and call `Read` at the end to refresh Terraform state.

#### Scenario: API call fails with retryable error
- **WHEN** the `ModifyReplication` API call returns a retryable error
- **THEN** the provider SHALL retry the call up to the write timeout
- **AND** the provider SHALL return the error if all retries are exhausted

#### Scenario: API call succeeds
- **WHEN** the `ModifyReplication` API call succeeds
- **THEN** the provider SHALL call `resourceTencentCloudTcrReplicationRead` to refresh the Terraform state
- **AND** the provider SHALL return nil (no error)

### Requirement: Backward compatibility is maintained
The `tencentcloud_tcr_replication` resource SHALL remain fully backward compatible with existing Terraform configurations and state files.

#### Scenario: Existing configuration with ForceNew fields
- **WHEN** a user applies an existing Terraform configuration that includes `ForceNew` fields (e.g., `source_registry_id`)
- **THEN** the behavior SHALL be identical to before the change

#### Scenario: Existing state file is valid after provider upgrade
- **WHEN** a user upgrades the provider and runs `terraform plan` on an existing `tencentcloud_tcr_replication` resource
- **THEN** no changes SHALL be detected (no diff) if the remote state matches the Terraform configuration