# Delta Spec: TDMQ RabbitMQ VIP Instance

This delta specification modifies the existing `tdmq-rabbitmq-vip-instance` spec to enable dynamic field updates.

## REMOVED Requirements

### Requirement: Public Access Field Immutability
**Reason**: The `enable_public_access` and `band_width` fields are now supported for dynamic updates by Tencent Cloud API. The previous immutability constraint was too restrictive and prevented users from toggling public access or adjusting bandwidth without recreating the instance.

**Migration**: Users who previously had to delete and recreate instances to change public access settings can now update these fields directly via Terraform. Existing instances with public access disabled can be enabled in-place, and bandwidth can be adjusted dynamically.

## MODIFIED Requirements

### Requirement: Public Network Access Configuration
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support public network access configuration through `enable_public_access` and `band_width` fields, and these fields SHALL be updatable after instance creation.

#### Scenario: User creates instance with public access enabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user sets `enable_public_access = true` and `band_width = 100`
- **THEN** the instance is created with public network access enabled
- **AND** the public network bandwidth is set to 100 Mbps
- **AND** the fields are visible in Terraform state after creation
- **AND** the `public_access_endpoint` computed field shows the public access endpoint

#### Scenario: User creates instance without public access (default)
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user does not specify `enable_public_access` field
- **THEN** the instance is created without public network access (default behavior)
- **AND** the `enable_public_access` field in state shows `false`
- **AND** the `band_width` field reflects API default value if present

#### Scenario: User creates instance with bandwidth but public access disabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user sets `band_width = 100` but `enable_public_access = false`
- **THEN** the instance is created without public network access
- **AND** the bandwidth configuration is sent to API but has no practical effect
- **AND** no error is raised (API-level validation)

#### Scenario: User enables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** the user changes `enable_public_access` to `true` in Terraform configuration
- **THEN** `terraform plan` shows the public access enable change
- **AND** `terraform apply` successfully enables public network access via API
- **AND** the instance becomes accessible from public network
- **AND** the Terraform state reflects the enabled status
- **AND** the `public_access_endpoint` computed field is populated

#### Scenario: User disables public access on existing instance
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true`
- **WHEN** the user changes `enable_public_access` to `false` in Terraform configuration
- **THEN** `terraform plan` shows the public access disable change
- **AND** `terraform apply` successfully disables public network access via API
- **AND** the instance is no longer accessible from public network
- **AND** the Terraform state reflects the disabled status
- **AND** the `public_access_endpoint` computed field becomes empty

#### Scenario: User updates bandwidth on existing instance with public access enabled
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = true` and `band_width = 100`
- **WHEN** the user changes `band_width` to `200` in Terraform configuration
- **THEN** `terraform plan` shows the bandwidth change
- **AND** `terraform apply` successfully updates the bandwidth via API
- **AND** the instance public network bandwidth is updated
- **AND** the Terraform state reflects the updated bandwidth

#### Scenario: User attempts to update bandwidth when public access is disabled
- **GIVEN** an existing RabbitMQ VIP instance with `enable_public_access = false`
- **WHEN** the user changes `band_width` in Terraform configuration
- **THEN** `terraform plan` shows the bandwidth change
- **AND** `terraform apply` attempts to update the bandwidth via API
- **AND** the behavior depends on API validation (may succeed but have no effect until public access is enabled)
- **AND** the Terraform state reflects the updated bandwidth value

### Requirement: Schema Definition for Public Access
The resource schema SHALL define public access fields with appropriate attributes, and these fields SHALL be updatable.

#### Scenario: enable_public_access schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_public_access` field
- **THEN** the field type is `schema.TypeBool`
- **AND** the field is marked as `Optional: true`
- **AND** the field has description: "Whether to enable public network access. Default is false. This field can be updated after creation."

#### Scenario: band_width schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `band_width` field
- **THEN** the field type is `schema.TypeInt`
- **AND** the field is marked as `Optional: true` and `Computed: true`
- **AND** the field has description: "Public network bandwidth in Mbps. This field can be updated after creation."

#### Scenario: Schema validation
- **GIVEN** a user provides `band_width` in configuration
- **WHEN** Terraform validates the configuration
- **THEN** negative values are rejected by Terraform type validation
- **AND** API-level validation (bandwidth limits, quota) is enforced by Tencent Cloud
- **AND** invalid combinations (e.g., bandwidth without public access) are accepted but may have no practical effect

### Requirement: Documentation Completeness for Public Access
The resource documentation SHALL clearly describe public access configuration and update capabilities.

#### Scenario: Field documentation
- **GIVEN** the resource documentation file `tdmq_rabbitmq_vip_instance.html.markdown`
- **WHEN** reviewing the arguments reference section
- **THEN** `enable_public_access` is documented with type, default value, and update capability
- **AND** `band_width` is documented with type, unit (Mbps), and update capability
- **AND** both fields are marked as updatable (can be changed after creation)
- **AND** documentation includes note about instance status requirements for updates

#### Scenario: Usage example
- **GIVEN** the resource documentation file
- **WHEN** reviewing the example usage section
- **THEN** an example shows creating an instance with public access enabled
- **AND** the example includes both `enable_public_access = true` and `band_width = 100`
- **AND** an example shows updating public access on an existing instance
- **AND** the example demonstrates the relationship between the two fields
