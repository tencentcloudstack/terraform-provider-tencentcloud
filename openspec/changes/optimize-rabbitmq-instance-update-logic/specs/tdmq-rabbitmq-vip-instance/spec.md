# Spec: TDMQ RabbitMQ VIP Instance

This is a delta spec for the `tdmq-rabbitmq-vip-instance` capability.

## MODIFIED Requirements

### Requirement: Public Access Field Mutability
The `enable_public_access` and `band_width` fields SHALL be mutable after instance creation, allowing users to update these fields through Terraform apply operations.

#### Scenario: User enables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** the user changes `enable_public_access` to `true` in the configuration
- **THEN** `terraform plan` shows the public access change
- **AND** `terraform apply` successfully enables public network access
- **AND** the Terraform state reflects the updated enable_public_access value
- **AND** the update operation calls the appropriate Tencent Cloud API

#### Scenario: User disables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true`
- **WHEN** the user changes `enable_public_access` to `false` in the configuration
- **THEN** `terraform plan` shows the public access change
- **AND** `terraform apply` successfully disables public network access
- **AND** the Terraform state reflects the updated enable_public_access value
- **AND** the update operation calls the appropriate Tencent Cloud API

#### Scenario: User increases bandwidth on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 100`
- **WHEN** the user changes `band_width` to `200` in the configuration
- **THEN** `terraform plan` shows the bandwidth change
- **AND** `terraform apply` successfully updates the instance bandwidth
- **AND** the Terraform state reflects the updated band_width value
- **AND** the update operation calls the appropriate Tencent Cloud API

#### Scenario: User decreases bandwidth on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `band_width = 200`
- **WHEN** the user changes `band_width` to `100` in the configuration
- **THEN** `terraform plan` shows the bandwidth change
- **AND** `terraform apply` successfully updates the instance bandwidth
- **AND** the Terraform state reflects the updated band_width value
- **AND** the update operation calls the appropriate Tencent Cloud API

#### Scenario: User toggles public access with bandwidth update
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false` and `band_width = 100`
- **WHEN** the user changes `enable_public_access` to `true` and `band_width` to `200` in the configuration
- **THEN** `terraform plan` shows both field changes
- **AND** `terraform apply` successfully enables public access and updates bandwidth
- **AND** the Terraform state reflects both updated values
- **AND** the update operation calls the appropriate Tencent Cloud API(s)

#### Scenario: User validates bandwidth constraints during update
- **GIVEN** a user provides a `band_width` value in the configuration for update
- **WHEN** Terraform validates the configuration
- **THEN** the value must be a positive integer
- **AND** the value must meet Tencent Cloud's minimum bandwidth requirement
- **AND** the value must not exceed the maximum bandwidth limit for the account
- **AND** invalid values are rejected with a clear error message

#### Scenario: Backward compatibility with existing resources
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the provider is upgraded to support public access field updates
- **THEN** `terraform plan` shows no changes for resources without field modifications
- **AND** existing resources continue to function normally
- **AND** the provider does not force any unintended updates
