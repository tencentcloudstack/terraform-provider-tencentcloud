## ADDED Requirements

### Requirement: Resource Schema Definition
The resource `tencentcloud_teo_multi_path_gateway` SHALL define the following schema fields:

- `zone_id` (TypeString, Required, ForceNew): 站点 ID
- `gateway_type` (TypeString, Required, ForceNew): 网关类型，取值 cloud 或 private
- `gateway_name` (TypeString, Required): 网关名称
- `gateway_port` (TypeInt, Required): 网关端口，范围 1~65535（除去 8888）
- `region_id` (TypeString, Optional, ForceNew): 网关地域，cloud 类型网关必填
- `gateway_ip` (TypeString, Optional): 网关地址，private 类型网关必填
- `gateway_id` (TypeString, Computed): 网关 ID，由创建接口返回
- `status` (TypeString, Computed): 网关状态
- `need_confirm` (TypeString, Computed): 回源 IP 列表是否需要确认

The resource SHALL support Import with composite ID format `zone_id#gateway_id`.

#### Scenario: Schema fields match cloud API parameters
- **WHEN** the resource schema is defined
- **THEN** all Create API input parameters SHALL have corresponding schema fields
- **THEN** Create API output parameters SHALL be mapped to Computed fields
- **THEN** gateway_type and region_id SHALL be ForceNew since Modify API does not support updating them

#### Scenario: Composite ID for import
- **WHEN** a user imports the resource with ID `zone-123#gateway-456`
- **THEN** zone_id SHALL be parsed as `zone-123`
- **THEN** gateway_id SHALL be parsed as `gateway-456`

### Requirement: Resource Create
The resource Create method SHALL call `CreateMultiPathGateway` API to create a multi-path gateway.

#### Scenario: Successful creation with cloud gateway
- **WHEN** creating a resource with gateway_type=cloud, zone_id, gateway_name, gateway_port, and region_id
- **THEN** the system SHALL call CreateMultiPathGateway with the provided parameters
- **THEN** the system SHALL set the resource ID to `zone_id` + FILED_SP + `gateway_id`
- **THEN** the system SHALL check if the response GatewayId is empty and return NonRetryableError if so

#### Scenario: Successful creation with private gateway
- **WHEN** creating a resource with gateway_type=private, zone_id, gateway_name, gateway_port, and gateway_ip
- **THEN** the system SHALL call CreateMultiPathGateway with the provided parameters
- **THEN** the system SHALL set the resource ID to `zone_id` + FILED_SP + `gateway_id`

#### Scenario: Create API call failure
- **WHEN** the CreateMultiPathGateway API call fails
- **THEN** the system SHALL wrap the error using tccommon.RetryError and return

### Requirement: Resource Read
The resource Read method SHALL call `DescribeMultiPathGateways` API to query the gateway details.

#### Scenario: Gateway exists
- **WHEN** reading a resource and the gateway is found
- **THEN** the system SHALL call DescribeMultiPathGateways with ZoneId and Filter by gateway-id
- **THEN** the system SHALL set all schema fields from the response including gateway_id, gateway_name, gateway_type, gateway_port, status, gateway_ip, region_id, need_confirm
- **THEN** the system SHALL NOT call set methods for nil response fields

#### Scenario: Gateway not found
- **WHEN** reading a resource and the gateway is not found (empty list)
- **THEN** the system SHALL set the resource ID to empty string to signal resource removal

### Requirement: Resource Update
The resource Update method SHALL call `ModifyMultiPathGateway` API to update the gateway information.

#### Scenario: Successful update
- **WHEN** updating gateway_name, gateway_ip, or gateway_port
- **THEN** the system SHALL call ModifyMultiPathGateway with zone_id, gateway_id, and the changed fields
- **THEN** the system SHALL call Read to refresh the state after update

#### Scenario: Update API call failure
- **WHEN** the ModifyMultiPathGateway API call fails
- **THEN** the system SHALL wrap the error using tccommon.RetryError and return

### Requirement: Resource Delete
The resource Delete method SHALL call `DeleteMultiPathGateway` API to delete the gateway.

#### Scenario: Successful deletion
- **WHEN** deleting a resource
- **THEN** the system SHALL call DeleteMultiPathGateway with zone_id and gateway_id
- **THEN** the system SHALL set the resource ID to empty string after successful deletion

#### Scenario: Delete API call failure
- **WHEN** the DeleteMultiPathGateway API call fails
- **THEN** the system SHALL wrap the error using tccommon.RetryError and return

### Requirement: Provider Registration
The resource SHALL be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`.

#### Scenario: Provider registration
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_multi_path_gateway` SHALL be available as a resource
- **THEN** the resource factory function SHALL be `ResourceTencentCloudTeoMultiPathGateway`

### Requirement: Resource Documentation
The resource SHALL have a corresponding `.md` documentation file following the project's documentation format.

#### Scenario: Documentation content
- **WHEN** the .md file is generated
- **THEN** it SHALL contain a one-line description mentioning the TEO product
- **THEN** it SHALL contain Example Usage section with HCL examples
- **THEN** it SHALL contain Import section explaining the composite ID format
- **THEN** it SHALL NOT contain Argument Reference or Attribute Reference sections

### Requirement: Unit Tests
The resource SHALL have unit tests using gomonkey mock approach.

#### Scenario: Unit test coverage
- **WHEN** unit tests are implemented
- **THEN** they SHALL use gomonkey to mock cloud API calls
- **THEN** they SHALL test Create, Read, Update, Delete business logic
- **THEN** they SHALL be runnable with `go test -gcflags=all=-l`
