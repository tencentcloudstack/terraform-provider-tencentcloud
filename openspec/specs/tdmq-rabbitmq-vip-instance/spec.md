# Spec: TDMQ RabbitMQ VIP Instance

This specification defines the requirements for the `tencentcloud_tdmq_rabbitmq_vip_instance` Terraform resource.

## Requirements

### Requirement: Resource Tags Field Support
The `tencentcloud_tdmq_rabbitmq_vip_instance` resource SHALL support a `resource_tags` field for managing instance-level resource tags.

#### Scenario: User creates instance with tags
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user includes resource_tags blocks with `tag_key` and `tag_value` pairs in the configuration
- **THEN** the instance is created with the specified tags
- **AND** the tags are visible in Terraform state after creation as a list of tag objects
- **AND** the tags are applied to the cloud resource as resource tags

#### Scenario: User creates instance without tags
- **GIVEN** a user defines a RabbitMQ VIP instance configuration
- **WHEN** the user does not include the `resource_tags` field
- **THEN** the instance is created successfully without any tags
- **AND** the `resource_tags` field is not present in Terraform state

#### Scenario: User reads instance tags into state
- **GIVEN** a RabbitMQ VIP instance exists with tags in Tencent Cloud
- **WHEN** Terraform performs a refresh operation
- **THEN** the `resource_tags` field in state is populated with the current tags from the cloud
- **AND** tags with nil keys or values are ignored gracefully

### Requirement: Tag Update Support
The resource SHALL allow users to update tags through Terraform apply operations.

#### Scenario: User adds new tags to existing instance
- **GIVEN** an existing RabbitMQ VIP instance without tags in Terraform configuration
- **WHEN** the user adds resource_tags blocks (e.g., `tag_key = "cost-center", tag_value = "123"`) to the configuration
- **THEN** `terraform plan` shows the tag addition
- **AND** `terraform apply` successfully adds the tags to the instance
- **AND** the Terraform state reflects the updated tags as a list

#### Scenario: User modifies existing tags
- **GIVEN** an existing instance with resource_tags blocks for `env = "dev"`
- **WHEN** the user changes the tag value to `env = "prod"`
- **THEN** `terraform plan` shows the tag value change
- **AND** `terraform apply` replaces all tags with the new configuration
- **AND** the Terraform state reflects the updated tags

#### Scenario: User removes all tags
- **GIVEN** an existing instance with multiple resource_tags blocks
- **WHEN** the user removes all resource_tags blocks from the configuration
- **THEN** `terraform plan` shows tag removal
- **AND** `terraform apply` removes all tags from the instance using RemoveAllTags flag
- **AND** the Terraform state shows an empty tags list

#### Scenario: User removes resource_tags field
- **GIVEN** an existing instance with resource_tags blocks
- **WHEN** the user removes the `resource_tags` field from the configuration
- **THEN** `terraform plan` shows no changes to tags
- **AND** tags on the cloud resource remain unchanged

### Requirement: Tag API Integration
The resource SHALL correctly map tags between Terraform and Tencent Cloud APIs.

#### Scenario: Create API integration
- **GIVEN** a user creates an instance with resource_tags blocks (e.g., `tag_key = "key1", tag_value = "value1"`)
- **WHEN** the Create operation calls `CreateRabbitMQVipInstance`
- **THEN** the request includes `ResourceTags: [{"TagKey": "key1", "TagValue": "value1"}]`
- **AND** the API response confirms successful tag creation

#### Scenario: Read API integration
- **GIVEN** a RabbitMQ VIP instance exists with tags in Tencent Cloud
- **WHEN** the Read operation calls `DescribeRabbitMQVipInstances`
- **THEN** tags are extracted from `response.ClusterInfo.Tags` array
- **AND** each SDK Tag object with TagKey and TagValue is converted to a list entry with tag_key and tag_value fields
- **AND** the result is set in Terraform state as `resource_tags` list

#### Scenario: Update API integration
- **GIVEN** a user modifies tags in an existing instance configuration
- **WHEN** the Update operation detects `resource_tags` changes via `d.HasChange()`
- **THEN** the request to `ModifyRabbitMQVipInstance` includes `Tags` array with all current tags
- **AND** the API call succeeds and tags are updated on the cloud resource

### Requirement: Data Type Handling
The resource SHALL handle tag data types and formats correctly.

#### Scenario: Tag conversion from Terraform to SDK
- **GIVEN** Terraform configuration with resource_tags blocks
- **WHEN** the provider converts tags for API call
- **THEN** the list of tag objects is converted to an array of `tdmq.Tag` structs
- **AND** each block becomes `{TagKey: helper.String(tag_key), TagValue: helper.String(tag_value)}`

#### Scenario: Tag conversion from SDK to Terraform
- **GIVEN** API response contains `Tags: [{"TagKey": "key", "TagValue": "value"}]`
- **WHEN** the provider converts tags for state storage
- **THEN** the array is converted to a list of maps `[{"tag_key": "key", "tag_value": "value"}]`
- **AND** nil TagKey or TagValue entries are skipped

#### Scenario: Empty tags handling
- **GIVEN** Terraform configuration with no resource_tags blocks
- **WHEN** the provider processes the update
- **THEN** the `RemoveAllTags` flag is set to true in the API request
- **AND** all existing tags on the resource are removed

### Requirement: Backward Compatibility
The `resource_tags` field SHALL be backward compatible with existing resources.

#### Scenario: Existing instance without tags field
- **GIVEN** a RabbitMQ VIP instance managed by Terraform before this feature
- **WHEN** the provider is upgraded to include `resource_tags` support
- **THEN** `terraform plan` shows no changes for resources without the field
- **AND** existing resources continue to function normally

#### Scenario: First-time tag addition to existing resource
- **GIVEN** an existing instance without `resource_tags` in configuration
- **WHEN** the user adds `resource_tags` field for the first time
- **THEN** Terraform treats this as an update (not recreation)
- **AND** only the `ModifyRabbitMQVipInstance` API is called
- **AND** no other instance properties are affected

### Requirement: Schema Definition
The resource schema SHALL define the `resource_tags` field with appropriate attributes.

#### Scenario: Schema field properties
- **GIVEN** the resource schema definition
- **WHEN** examining the `resource_tags` field
- **THEN** the field type is `schema.TypeList`
- **AND** the element is a Resource with `tag_key` (TypeString) and `tag_value` (TypeString) fields
- **AND** both nested fields are marked as `Required: true`
- **AND** the field is marked as `Optional: true`
- **AND** the field has a clear description for documentation

#### Scenario: Field validation
- **GIVEN** a user provides `resource_tags` blocks in configuration
- **WHEN** Terraform validates the configuration
- **THEN** each block must contain both `tag_key` and `tag_value` string fields
- **AND** complex types (nested lists, maps) are rejected
- **AND** API-level validation (key/value length, format) is enforced by Tencent Cloud

### Requirement: Error Handling
The resource SHALL handle tag-related errors gracefully.

#### Scenario: API error during tag creation
- **GIVEN** a user creates an instance with invalid tag keys
- **WHEN** the `CreateRabbitMQVipInstance` API returns an error
- **THEN** the error is propagated to the user with context
- **AND** the instance creation is not completed
- **AND** no partial state is written

#### Scenario: API error during tag update
- **GIVEN** a user updates tags on an existing instance
- **WHEN** the `ModifyRabbitMQVipInstance` API fails
- **THEN** the error is logged and returned to the user
- **AND** Terraform state remains unchanged (previous tag values)
- **AND** subsequent `terraform apply` attempts the update again

#### Scenario: Nil tag handling during read
- **GIVEN** API response contains tags with nil TagKey or TagValue
- **WHEN** the Read operation processes the response
- **THEN** nil entries are skipped without panicking
- **AND** only valid tags are added to the state list
- **AND** no error is returned for nil tags

### Requirement: Code Formatting
All code changes SHALL be formatted using `go fmt`.

#### Scenario: File formatting after modification
- **GIVEN** code changes are made to the resource file
- **WHEN** all changes are complete
- **THEN** `go fmt` is executed on the file
- **AND** all Go code adheres to standard formatting rules
- **AND** no formatting warnings or errors exist

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
- **AND** the `band_width` field reflects API default value if present

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

### Requirement: Create API Integration for Public Access
The resource SHALL correctly pass public access fields to the `CreateRabbitMQVipInstance` API.

#### Scenario: Send bandwidth to Create API
- **GIVEN** a user creates an instance with `band_width = 200`
- **WHEN** the Create operation calls `CreateRabbitMQVipInstance`
- **THEN** the request includes `Bandwidth: helper.IntUint64(200)`
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
- **AND** the field is marked as `Optional: true` and `Computed: true`
- **AND** the field has description: "Public network bandwidth in Mbps."

#### Scenario: Schema validation
- **GIVEN** a user provides `band_width` in configuration
- **WHEN** Terraform validates the configuration
- **THEN** negative values are rejected by Terraform type validation
- **AND** API-level validation (bandwidth limits, quota) is enforced by Tencent Cloud
- **AND** invalid combinations (e.g., bandwidth without public access) are accepted but have no effect

### Requirement: Documentation Completeness for Public Access
The resource documentation SHALL clearly describe public access configuration.

#### Scenario: Field documentation
- **GIVEN** the resource documentation file `tdmq_rabbitmq_vip_instance.html.markdown`
- **WHEN** reviewing the arguments reference section
- **THEN** `enable_public_access` is documented with type, default value, and immutability note
- **AND** `band_width` is documented with type and unit (Mbps)
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
