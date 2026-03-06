# Spec Delta: TDMQ RabbitMQ VIP Instance - Public Network Access

## ADDED Requirements

### Requirement: Public Network Access Configuration
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support public network access configuration through `enable_public_access` and `band_width` fields.

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
- **AND** the `band_width` field has no effect

#### Scenario: User creates instance with bandwidth but public access disabled
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user sets `band_width = 100` but `enable_public_access = false`
- **THEN** the instance is created without public network access
- **AND** the bandwidth configuration is sent to API but has no practical effect
- **AND** no error is raised (API-level validation)

### Requirement: Public Access Field Immutability
The `enable_public_access` and `band_width` fields SHALL be immutable after instance creation.

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

### Requirement: Public Access State Reading
The resource SHALL correctly read public access configuration from Tencent Cloud API responses.

#### Scenario: Read bandwidth from API response
- **GIVEN** a RabbitMQ VIP instance exists with public access enabled
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** the bandwidth value is extracted from `response.ClusterSpecInfo.PublicNetworkTps`
- **AND** the value is set in Terraform state as `band_width` field
- **AND** nil values are handled gracefully (not set in state)

#### Scenario: Read public access status from API response
- **GIVEN** a RabbitMQ VIP instance exists in Tencent Cloud
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** the public access status is extracted from `response.ClusterNetInfo.PublicDataStreamStatus`
- **AND** the string value "ON" is converted to boolean `true`
- **AND** the string value "OFF" is converted to boolean `false`
- **AND** the boolean value is set in Terraform state as `enable_public_access`

#### Scenario: Handle missing public access data in API response
- **GIVEN** a RabbitMQ VIP instance without public access configuration
- **WHEN** the Read operation receives `ClusterNetInfo.PublicDataStreamStatus` as nil or empty
- **THEN** the `enable_public_access` field defaults to `false` or is not set
- **AND** no error is raised during state refresh

### Requirement: Create API Integration
The resource SHALL correctly pass public access fields to the `CreateRabbitMQVipInstance` API.

#### Scenario: Send bandwidth to Create API
- **GIVEN** a user creates an instance with `band_width = 200`
- **WHEN** the Create operation calls `CreateRabbitMQVipInstance`
- **THEN** the request includes `Bandwidth: helper.IntInt64(200)`
- **AND** the API accepts the parameter

#### Scenario: Send public access flag to Create API
- **GIVEN** a user creates an instance with `enable_public_access = true`
- **WHEN** the Create operation calls `CreateRabbitMQVipInstance`
- **THEN** the request includes `EnablePublicAccess: helper.Bool(true)`
- **AND** the API enables public network access

#### Scenario: Omit public access fields when not specified
- **GIVEN** a user creates an instance without `enable_public_access` or `band_width`
- **WHEN** the Create operation calls `CreateRabbitMQVipInstance`
- **THEN** the request does not include `Bandwidth` or `EnablePublicAccess` fields
- **AND** the API applies default values (no public access)

### Requirement: Schema Definition for Public Access
The resource schema SHALL define public access fields with appropriate attributes.

#### Scenario: enable_public_access schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `enable_public_access` field
- **THEN** the field type is `schema.TypeBool`
- **AND** the field is marked as `Optional: true`
- **AND** the field has description: "Whether to enable public network access. Default is false."

#### Scenario: band_width schema properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `band_width` field
- **THEN** the field type is `schema.TypeInt`
- **AND** the field is marked as `Optional: true`
- **AND** the field has description: "Public network bandwidth in Mbps. Only takes effect when enable_public_access is true."

#### Scenario: Schema validation
- **GIVEN** a user provides `band_width` in configuration
- **WHEN** Terraform validates the configuration
- **THEN** negative values are rejected by Terraform type validation
- **AND** API-level validation (bandwidth limits, quota) is enforced by Tencent Cloud
- **AND** invalid combinations (e.g., bandwidth without public access) are accepted but have no effect

### Requirement: Documentation Completeness
The resource documentation SHALL clearly describe public access configuration.

#### Scenario: Field documentation
- **GIVEN** the resource documentation file `tdmq_rabbitmq_vip_instance.html.markdown`
- **WHEN** reviewing the arguments reference section
- **THEN** `enable_public_access` is documented with type, default value, and immutability note
- **AND** `band_width` is documented with type, unit (Mbps), and dependency on `enable_public_access`
- **AND** both fields are marked as immutable (cannot be changed after creation)

#### Scenario: Usage example
- **GIVEN** the resource documentation file
- **WHEN** reviewing the example usage section
- **THEN** an example shows creating an instance with public access enabled
- **AND** the example includes both `enable_public_access = true` and `band_width = 100`
- **AND** the example demonstrates the relationship between the two fields

### Requirement: Error Handling for Public Access
The resource SHALL handle public access errors gracefully.

#### Scenario: API error during creation with public access
- **GIVEN** a user creates an instance with `enable_public_access = true` but insufficient quota
- **WHEN** the `CreateRabbitMQVipInstance` API returns a quota error
- **THEN** the error is propagated to the user with context
- **AND** the instance creation is not completed
- **AND** no partial state is written

#### Scenario: Nil value handling during read
- **GIVEN** API response contains nil `PublicNetworkTps` or `PublicDataStreamStatus`
- **WHEN** the Read operation processes the response
- **THEN** nil values are handled without panicking
- **AND** fields are not set in state (or set to default values)
- **AND** no error is returned

### Requirement: Backward Compatibility with Public Access
The public access fields SHALL be backward compatible with existing resources.

#### Scenario: Existing instance without public access fields
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the provider is upgraded to include public access support
- **THEN** `terraform plan` shows no changes for resources without these fields
- **AND** existing resources continue to function normally
- **AND** state refresh correctly populates the new fields from API response

#### Scenario: First-time public access addition requires recreation
- **GIVEN** an existing instance without public access
- **WHEN** the user adds `enable_public_access = true` to configuration
- **THEN** Terraform detects this as a change to an immutable field
- **AND** the user must manually destroy and recreate the resource
- **AND** documentation clearly explains this limitation
