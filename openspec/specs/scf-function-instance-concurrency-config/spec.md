# scf-function-instance-concurrency-config Specification

## Purpose
TBD - created by archiving change add-scf-function-instance-concurrency-config. Update Purpose after archive.
## Requirements
### Requirement: Instance Concurrency Config parameter
The `tencentcloud_scf_function` resource SHALL support an optional `instance_concurrency_config` parameter that allows users to configure single-instance multi-concurrency behavior for Web functions.

#### Scenario: Create function with instance concurrency config
- **WHEN** a user creates a `tencentcloud_scf_function` resource with `instance_concurrency_config` specified
- **THEN** the `InstanceConcurrencyConfig` field SHALL be passed to the `CreateFunction` API call
- **AND** the function SHALL be created with the specified concurrency configuration

#### Scenario: Create function without instance concurrency config
- **WHEN** a user creates a `tencentcloud_scf_function` resource without `instance_concurrency_config`
- **THEN** the `InstanceConcurrencyConfig` field SHALL NOT be set in the `CreateFunction` API call
- **AND** the function SHALL be created with default concurrency behavior

#### Scenario: Read function with instance concurrency config
- **WHEN** Terraform reads the state of a `tencentcloud_scf_function` resource that has `InstanceConcurrencyConfig` set in the API response
- **THEN** the `instance_concurrency_config` attribute SHALL be populated with the API-returned values
- **AND** all sub-fields (`dynamic_enabled`, `max_concurrency`, `instance_isolation_enabled`, `type`, `mix_node_config`, `session_config`) SHALL be set if present in the API response

#### Scenario: Read function without instance concurrency config
- **WHEN** Terraform reads the state of a `tencentcloud_scf_function` resource where the API returns `InstanceConcurrencyConfig` as nil
- **THEN** the `instance_concurrency_config` attribute SHALL remain unset
- **AND** no error SHALL be returned

#### Scenario: Update function instance concurrency config
- **WHEN** a user modifies the `instance_concurrency_config` parameter on an existing `tencentcloud_scf_function` resource
- **THEN** the `InstanceConcurrencyConfig` field SHALL be passed to the `UpdateFunctionConfiguration` API call
- **AND** the function's concurrency configuration SHALL be updated

#### Scenario: Remove instance concurrency config
- **WHEN** a user removes the `instance_concurrency_config` parameter from the Terraform configuration
- **THEN** the `UpdateFunctionConfiguration` call SHALL be made without `InstanceConcurrencyConfig` (allowing the API to reset to defaults)

