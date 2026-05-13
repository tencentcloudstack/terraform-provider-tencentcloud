## ADDED Requirements

### Requirement: Resource schema definition for tencentcloud_teo_confirm_multi_path_gateway_origin_acl
The resource SHALL define a Terraform RESOURCE_KIND_CONFIG resource named `tencentcloud_teo_confirm_multi_path_gateway_origin_acl` with the following schema fields:
- `zone_id` (TypeString, Required, ForceNew): Zone ID identifying the TEO site
- `gateway_id` (TypeString, Required, ForceNew): Gateway ID identifying the multi-path gateway
- `origin_acl_version` (TypeInt, Optional): Origin ACL version number to confirm
- `multi_path_gateway_origin_acl_info` (TypeList, Computed, MaxItems: 1): ACL info output containing current and next origin ACL details

The resource SHALL support Import via `schema.ImportStatePassthroughContext`.

#### Scenario: Resource schema defines required and computed fields
- **WHEN** the resource schema is registered
- **THEN** `zone_id` and `gateway_id` are Required with ForceNew, `origin_acl_version` is Optional, and `multi_path_gateway_origin_acl_info` is Computed

#### Scenario: Resource supports import
- **WHEN** a user imports an existing resource with ID `zone_id#gateway_id`
- **THEN** Terraform reads the current state using the parsed zone_id and gateway_id

### Requirement: Composite resource ID using zone_id and gateway_id
The resource SHALL use `zone_id` and `gateway_id` joined by `tccommon.FILED_SP` as the composite resource ID.

#### Scenario: Create sets composite ID
- **WHEN** the resource is created with zone_id "zone-123" and gateway_id "gw-456"
- **THEN** the resource ID is set to "zone-123#gw-456" (using FILED_SP separator)

#### Scenario: Read parses composite ID
- **WHEN** the resource ID is "zone-123#gw-456"
- **THEN** zone_id is parsed as "zone-123" and gateway_id is parsed as "gw-456"

### Requirement: Read operation via DescribeMultiPathGatewayOriginACL
The resource SHALL implement the Read operation by calling `DescribeMultiPathGatewayOriginACL` with `ZoneId` and `GatewayId` from the resource state. The response's `MultiPathGatewayOriginACLInfo` SHALL be flattened into the `multi_path_gateway_origin_acl_info` computed output.

#### Scenario: Successful read returns ACL info
- **WHEN** Read is called and the API returns `MultiPathGatewayOriginACLInfo` with current and next ACL data
- **THEN** the `multi_path_gateway_origin_acl_info` output is populated with `multi_path_gateway_current_origin_acl` and `multi_path_gateway_next_origin_acl` blocks

#### Scenario: Resource not found during read
- **WHEN** Read is called and the API returns nil/empty data
- **THEN** the resource ID is cleared and a warning is logged

### Requirement: Read operation with retry
The resource SHALL use `helper.Retry()` with `tccommon.ReadRetryTimeout` when calling `DescribeMultiPathGatewayOriginACL`. If the API call fails, the error SHALL be wrapped using `tccommon.RetryError()` and returned.

#### Scenario: API call succeeds on first attempt
- **WHEN** the DescribeMultiPathGatewayOriginACL API call succeeds
- **THEN** the result is returned immediately without retry

#### Scenario: API call fails and retry succeeds
- **WHEN** the DescribeMultiPathGatewayOriginACL API call fails on the first attempt
- **THEN** the retry mechanism retries the call within the ReadRetryTimeout period

### Requirement: Create operation delegates to Update
The resource SHALL implement the Create operation by setting the composite ID from `zone_id` and `gateway_id`, then delegating to the Update function.

#### Scenario: Create with origin_acl_version specified
- **WHEN** the resource is created with zone_id, gateway_id, and origin_acl_version
- **THEN** the composite ID is set and the Update function is called to confirm the ACL version

#### Scenario: Create without origin_acl_version
- **WHEN** the resource is created with zone_id and gateway_id but no origin_acl_version
- **THEN** the composite ID is set and the Update function is called (which may skip the Confirm API call)

### Requirement: Update operation via ConfirmMultiPathGatewayOriginACL
The resource SHALL implement the Update operation. When `origin_acl_version` is specified or changed, the resource SHALL call `ConfirmMultiPathGatewayOriginACL` with `ZoneId`, `GatewayId`, and `OriginACLVersion`. After the confirm call, the resource SHALL call Read to refresh the state.

#### Scenario: Update with new origin_acl_version
- **WHEN** the origin_acl_version field is changed from nil to a specific version number
- **THEN** ConfirmMultiPathGatewayOriginACL is called with the new version, and Read is called to refresh state

#### Scenario: Update with changed origin_acl_version
- **WHEN** the origin_acl_version field is changed from one version to another
- **THEN** ConfirmMultiPathGatewayOriginACL is called with the new version, and Read is called to refresh state

### Requirement: Delete operation is no-op
The resource SHALL implement the Delete operation as a no-op, since this is a CONFIG resource that cannot be truly deleted.

#### Scenario: Delete does nothing
- **WHEN** the resource is deleted from Terraform state
- **THEN** no API calls are made and the function returns nil

### Requirement: Output structure for multi_path_gateway_origin_acl_info
The `multi_path_gateway_origin_acl_info` output SHALL contain two nested blocks:
1. `multi_path_gateway_current_origin_acl` (TypeList, Computed, MaxItems: 1):
   - `entire_addresses` (TypeList, Computed, MaxItems: 1): Contains `ipv4` (TypeList, Computed) and `ipv6` (TypeList, Computed)
   - `version` (TypeInt, Computed): Current version number
   - `is_planed` (TypeString, Computed): Whether the update confirmation is completed
2. `multi_path_gateway_next_origin_acl` (TypeList, Computed, MaxItems: 1):
   - `version` (TypeInt, Computed): Next version number
   - `entire_addresses` (TypeList, Computed, MaxItems: 1): Contains `ipv4` and `ipv6`
   - `added_addresses` (TypeList, Computed, MaxItems: 1): Contains `ipv4` and `ipv6`
   - `removed_addresses` (TypeList, Computed, MaxItems: 1): Contains `ipv4` and `ipv6`
   - `no_change_addresses` (TypeList, Computed, MaxItems: 1): Contains `ipv4` and `ipv6`

#### Scenario: Full output structure is populated
- **WHEN** Read is called and both current and next ACL info exist
- **THEN** all nested fields are populated in the Terraform state

#### Scenario: Next ACL info is empty when no pending update
- **WHEN** Read is called and there is no pending next version
- **THEN** `multi_path_gateway_next_origin_acl` is empty/nil in the state

### Requirement: Resource registration in provider
The resource SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Resource appears in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_confirm_multi_path_gateway_origin_acl` is available as a resource type

### Requirement: Unit tests using gomonkey mock
The resource SHALL have unit tests in `resource_tc_teo_confirm_multi_path_gateway_origin_acl_config_test.go` using gomonkey to mock the cloud API calls. Tests SHALL be runnable with `go test -gcflags=all=-l`.

#### Scenario: Unit tests cover Read operation
- **WHEN** unit tests are run
- **THEN** the Read operation is tested with mocked DescribeMultiPathGatewayOriginACL response

#### Scenario: Unit tests cover Update operation
- **WHEN** unit tests are run
- **THEN** the Update/Confirm operation is tested with mocked ConfirmMultiPathGatewayOriginACL response

### Requirement: Documentation markdown file
The resource SHALL have a documentation file `resource_tc_teo_confirm_multi_path_gateway_origin_acl_config.md` following the format specified in `gendoc/README.md`, including a one-line description, Example Usage section, and Import section (since this is RESOURCE_KIND_CONFIG).

#### Scenario: Documentation file exists
- **WHEN** the resource is implemented
- **THEN** a .md file with description, example usage, and import section exists in the teo service directory
