## ADDED Requirements

### Requirement: Resource Schema Definition
The resource `tencentcloud_teo_multi_path_gateway_line` SHALL define the following schema fields:
- `zone_id` (TypeString, Required, ForceNew): 站点 ID
- `gateway_id` (TypeString, Required, ForceNew): 多通道安全网关 ID
- `line_type` (TypeString, Required): 线路类型，取值为 direct/proxy/custom
- `line_address` (TypeString, Required): 线路地址，格式为 ip:port
- `proxy_id` (TypeString, Optional): 四层代理实例 ID
- `rule_id` (TypeString, Optional): 转发规则 ID
- `line_id` (TypeString, Computed): 线路 ID，由云 API 创建后返回

The resource SHALL support Import via `schema.ImportStatePassthrough`.
The resource SHALL use composite ID format `zone_id#gateway_id#line_id` with `tccommon.FILED_SP` as separator.

#### Scenario: Schema fields are correctly defined
- **WHEN** the resource schema is initialized
- **THEN** `zone_id` and `gateway_id` SHALL be Required with ForceNew
- **AND** `line_type` and `line_address` SHALL be Required without ForceNew
- **AND** `proxy_id` and `rule_id` SHALL be Optional
- **AND** `line_id` SHALL be Computed

### Requirement: Create Multi-Path Gateway Line
The resource SHALL call `CreateMultiPathGatewayLine` API to create a new gateway line.

#### Scenario: Successful creation with custom line type
- **WHEN** a resource is created with `zone_id`, `gateway_id`, `line_type="custom"`, and `line_address`
- **THEN** the system SHALL call `CreateMultiPathGatewayLine` with the provided parameters
- **AND** SHALL set the resource ID to `zone_id#gateway_id#line_id` using the returned `LineId`
- **AND** SHALL call Read to refresh the resource state

#### Scenario: Successful creation with proxy line type
- **WHEN** a resource is created with `line_type="proxy"` along with `proxy_id` and `rule_id`
- **THEN** the system SHALL include `proxy_id` and `rule_id` in the Create request

### Requirement: Read Multi-Path Gateway Line
The resource SHALL call `DescribeMultiPathGatewayLine` API to read the current state of the gateway line. The Read function SHALL get `zone_id` and `gateway_id` from `d.Get()` and `line_id` from the composite ID.

#### Scenario: Resource exists
- **WHEN** the Read function is called for an existing resource
- **THEN** the system SHALL call `DescribeMultiPathGatewayLine` with `zone_id`, `gateway_id`, and `line_id`
- **AND** SHALL populate all schema fields from the `Line` response object including `line_id`, `line_type`, `line_address`, `proxy_id`, and `rule_id`

#### Scenario: Resource not found
- **WHEN** the Read function is called but the resource does not exist
- **THEN** the system SHALL set the resource ID to empty string to indicate removal

### Requirement: Update Multi-Path Gateway Line
The resource SHALL call `ModifyMultiPathGatewayLine` API to update mutable fields when changes are detected. Mutable fields include `line_type`, `line_address`, `proxy_id`, and `rule_id`.

#### Scenario: Update mutable fields
- **WHEN** any of `line_type`, `line_address`, `proxy_id`, or `rule_id` has changed
- **THEN** the system SHALL call `ModifyMultiPathGatewayLine` with `zone_id`, `gateway_id`, `line_id`, and all current field values
- **AND** SHALL call Read to refresh the resource state

#### Scenario: No mutable fields changed
- **WHEN** none of the mutable fields have changed
- **THEN** the system SHALL NOT call ModifyMultiPathGatewayLine
- **AND** SHALL call Read to refresh the resource state

### Requirement: Delete Multi-Path Gateway Line
The resource SHALL call `DeleteMultiPathGatewayLine` API to delete the gateway line.

#### Scenario: Successful deletion
- **WHEN** the Delete function is called
- **THEN** the system SHALL call `DeleteMultiPathGatewayLine` with `zone_id`, `gateway_id`, and `line_id`
- **AND** the resource SHALL be removed from the Terraform state

### Requirement: Resource Registration
The resource SHALL be registered in `provider.go` with the name `tencentcloud_teo_multi_path_gateway_line` and the corresponding function `ResourceTencentCloudTeoMultiPathGatewayLine()`.

#### Scenario: Resource appears in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_multi_path_gateway_line` SHALL be available as a resource type

### Requirement: Retry on API Calls
All CRUD operations SHALL use `resource.Retry` with appropriate timeouts:
- Write operations (Create, Update, Delete) SHALL use `tccommon.WriteRetryTimeout`
- Read operations SHALL use `tccommon.ReadRetryTimeout`
- Errors SHALL be wrapped with `tccommon.RetryError()`

#### Scenario: API call fails transiently
- **WHEN** a cloud API call fails with a retryable error
- **THEN** the system SHALL retry the call within the timeout period
- **AND** SHALL wrap the error with `tccommon.RetryError()`

### Requirement: Unit Tests with Mock
The resource SHALL have unit tests using gomonkey to mock cloud API calls, not using Terraform acceptance test framework.

#### Scenario: Create operation is tested with mock
- **WHEN** unit tests are executed
- **THEN** CreateMultiPathGatewayLine API SHALL be mocked using gomonkey
- **AND** the test SHALL verify the resource ID is correctly set

#### Scenario: Read operation is tested with mock
- **WHEN** unit tests are executed
- **THEN** DescribeMultiPathGatewayLine API SHALL be mocked using gomonkey
- **AND** the test SHALL verify the resource state is correctly populated
