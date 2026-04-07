# Spec: TDMQ RabbitMQ VIP Instance (Delta)

This document describes changes to the `tdmq-rabbitmq-vip-instance` specification to support enhanced update capabilities.

## MODIFIED Requirements

### Requirement: Public Access Field Immutability
The `enable_public_access` and `band_width` fields SHALL remain immutable after instance creation.

#### Scenario: User attempts to enable public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** the user changes `enable_public_access` to `true` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `enable_public_access` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to modify bandwidth on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 100`
- **WHEN** the user changes `band_width` to `200` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `band_width` cannot be changed"
- **AND** the instance is not modified

#### Scenario: Recreating instance with different public access settings
- **GIVEN** an existing instance with public access configuration
- **WHEN** the user changes immutable fields (`enable_public_access` or `band_width`)
- **THEN** Terraform requires manual resource destruction and recreation
- **AND** users must use `terraform taint` or manual delete/create workflow

### Requirement: Delete Protection Field Mutability
The `enable_deletion_protection` field SHALL be mutable after instance creation.

#### Scenario: User enables deletion protection on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_deletion_protection = false` or without the field set
- **WHEN** the user sets `enable_deletion_protection` to `true` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with `EnableDeletionProtection=true`
- **AND** the deletion protection is enabled on the instance
- **AND** the Terraform state reflects the updated value

#### Scenario: User disables deletion protection on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_deletion_protection = true`
- **WHEN** the user sets `enable_deletion_protection` to `false` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with `EnableDeletionProtection=false`
- **AND** the deletion protection is disabled on the instance
- **AND** the Terraform state reflects the updated value

### Requirement: Instance Remark Field Support
The resource SHALL support a `remark` field for instance description/configuration notes.

#### Scenario: User sets remark during instance creation
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user includes `remark = "Production RabbitMQ cluster"` in the configuration
- **THEN** the instance is created with the specified remark
- **AND** the remark field is visible in Terraform state after creation
- **AND** the remark is stored as instance metadata in Tencent Cloud

#### Scenario: User modifies remark on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `remark = "Initial setup"`
- **WHEN** the user changes the remark to `remark = "Production environment"`
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with the new `Remark` value
- **AND** the instance remark is updated
- **AND** the Terraform state reflects the updated value

#### Scenario: User clears remark on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `remark = "Old description"`
- **WHEN** the user removes the `remark` field or sets it to empty string
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with `Remark` set to empty string or null
- **AND** the instance remark is cleared
- **AND** the Terraform state reflects the change

#### Scenario: Read remark from existing instance
- **GIVEN** a RabbitMQ VIP instance exists with a remark in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** the `remark` field in state is populated with the current remark from the cloud
- **AND** nil values are handled gracefully (field not set in state)

### Requirement: Risk Warning Field Support
The resource SHALL support an `enable_risk_warning` field for cluster risk warning configuration.

#### Scenario: User enables risk warning on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_risk_warning = false` or without the field set
- **WHEN** the user sets `enable_risk_warning` to `true` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with `EnableRiskWarning=true`
- **AND** the risk warning is enabled on the instance
- **AND** the Terraform state reflects the updated value

#### Scenario: User disables risk warning on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_risk_warning = true`
- **WHEN** the user sets `enable_risk_warning` to `false` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` calls `ModifyRabbitMQVipInstance` API with `EnableRiskWarning=false`
- **AND** the risk warning is disabled on the instance
- **AND** the Terraform state reflects the updated value

#### Scenario: Read risk warning from existing instance
- **GIVEN** a RabbitMQ VIP instance exists with risk warning configuration in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** the `enable_risk_warning` field in state is populated with the current value from the cloud
- **AND** nil values are handled gracefully (field not set in state or defaults to false)

### Requirement: Immutable Infrastructure Parameters
The following fields SHALL remain immutable after instance creation as they represent infrastructure-level specifications: `zone_ids`, `vpc_id`, `subnet_id`, `node_spec`, `node_num`, `storage_size`, `cluster_version`, `auto_renew_flag`, `time_span`, `pay_mode`.

#### Scenario: User attempts to modify zone_ids
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes `zone_ids` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `zone_ids` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to modify vpc_id
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes `vpc_id` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `vpc_id` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to modify subnet_id
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes `subnet_id` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `subnet_id` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to modify node_spec
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes `node_spec` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `node_spec` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to modify node_num
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes `node_num` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `node_num` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to modify storage_size
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes `storage_size` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `storage_size` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to modify cluster_version
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes `cluster_version` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `cluster_version` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to modify auto_renew_flag
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes `auto_renew_flag` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `auto_renew_flag` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to modify time_span
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes `time_span` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `time_span` cannot be changed"
- **AND** the instance is not modified

#### Scenario: User attempts to modify pay_mode
- **GIVEN** an existing RabbitMQ VIP instance
- **WHEN** the user changes `pay_mode` in the configuration
- **THEN** `terraform plan` detects the change
- **AND** `terraform apply` fails with error message "argument `pay_mode` cannot be changed"
- **AND** the instance is not modified

### Requirement: Schema Definition for New Fields
The resource schema SHALL define new fields with appropriate attributes.

#### Scenario: remark schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `remark` field
- **THEN** the field type is `schema.TypeString`
- **AND** the field is marked as `Optional: true`
- **AND** the field has description: "Instance remark/description"
- **AND** the field is not marked as `Computed`

#### Scenario: enable_risk_warning schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_risk_warning` field
- **THEN** the field type is `schema.TypeBool`
- **AND** the field is marked as `Optional: true`
- **AND** the field has description: "Whether to enable cluster risk warning"

#### Scenario: enable_deletion_protection schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_deletion_protection` field
- **THEN** the field type is `schema.TypeBool`
- **AND** the field is marked as `Optional: true`
- **AND** the field has description: "Whether to enable deletion protection"
- **AND** the field is not marked as `Computed` (to support updates)

## ADDED Requirements

### Requirement: Update API Integration for New Parameters
The resource SHALL correctly pass new parameters to the `ModifyRabbitMQVipInstance` API.

#### Scenario: Send remark to Modify API
- **GIVEN** a user modifies the `remark` parameter on an existing instance
- **WHEN** the Update operation calls `ModifyRabbitMQVipInstance`
- **THEN** the request includes `Remark: helper.String("new remark")`
- **AND** the API accepts the parameter and updates the instance

#### Scenario: Send deletion protection to Modify API
- **GIVEN** a user modifies the `enable_deletion_protection` parameter on an existing instance
- **WHEN** the Update operation calls `ModifyRabbitMQVipInstance`
- **THEN** the request includes `EnableDeletionProtection: helper.Bool(true)`
- **AND** the API accepts the parameter and updates the instance

#### Scenario: Send risk warning to Modify API
- **GIVEN** a user modifies the `enable_risk_warning` parameter on an existing instance
- **WHEN** the Update operation calls `ModifyRabbitMQVipInstance`
- **THEN** the request includes `EnableRiskWarning: helper.Bool(true)`
- **AND** the API accepts the parameter and updates the instance

### Requirement: Backward Compatibility with New Fields
The new fields SHALL be backward compatible with existing resources.

#### Scenario: Existing instance without new fields
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the provider is upgraded to include new field support
- **THEN** `terraform plan` shows no changes for resources without these fields
- **AND** existing resources continue to function normally
- **AND** state refresh correctly populates the new fields from API response

#### Scenario: First-time new field addition to existing resource
- **GIVEN** an existing instance without `remark`, `enable_risk_warning`, or `enable_deletion_protection` in configuration
- **WHEN** the user adds any of these fields for the first time
- **THEN** Terraform treats this as an update (not recreation)
- **AND** only the `ModifyRabbitMQVipInstance` API is called
- **AND** no other instance properties are affected
- **AND** the specific parameter is updated based on the user's configuration
