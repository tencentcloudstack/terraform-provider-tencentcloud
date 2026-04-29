## ADDED Requirements

### Requirement: Resource Schema Definition
The system SHALL define a Terraform resource `tencentcloud_teo_multi_path_gateway` with the following schema fields:
- `zone_id` (Required, ForceNew, TypeString): 站点 ID
- `gateway_type` (Required, ForceNew, TypeString): 网关类型，取值 cloud 或 private
- `gateway_name` (Required, TypeString): 网关名称，16 个字符以内
- `gateway_port` (Optional, Computed, TypeInt): 网关端口，范围 1～65535（除去 8888）
- `region_id` (Optional, Computed, TypeString): 网关地域，GatewayType 为 cloud 时必填
- `gateway_ip` (Optional, Computed, TypeString): 网关地址，GatewayType 为 private 时必填
- `gateway_id` (Computed, TypeString): 网关 ID
- `status` (Computed, TypeString): 网关状态
- `need_confirm` (Computed, TypeString): 是否需要重新确认回源 IP 列表

#### Scenario: Schema defines all CRUD fields
- **WHEN** the resource schema is defined
- **THEN** it SHALL include zone_id, gateway_type, gateway_name, gateway_port, region_id, gateway_ip, gateway_id, status, and need_confirm fields with correct types and constraints

#### Scenario: ForceNew fields prevent in-place update
- **WHEN** zone_id or gateway_type is changed in the Terraform configuration
- **THEN** the resource SHALL be destroyed and recreated

### Requirement: Resource Create Operation
The system SHALL implement the Create operation by calling `CreateMultiPathGateway` API with ZoneId, GatewayType, GatewayName, GatewayPort, RegionId, GatewayIP parameters. Upon success, the resource ID SHALL be set to `ZoneId:GatewayId` using tccommon.FILED_SP as separator.

#### Scenario: Create cloud type gateway
- **WHEN** a resource is created with gateway_type="cloud" and region_id is provided
- **THEN** the system SHALL call CreateMultiPathGateway with ZoneId, GatewayType="cloud", GatewayName, GatewayPort, RegionId
- **AND** set the resource ID to "ZoneId:GatewayId"

#### Scenario: Create private type gateway
- **WHEN** a resource is created with gateway_type="private" and gateway_ip is provided
- **THEN** the system SHALL call CreateMultiPathGateway with ZoneId, GatewayType="private", GatewayName, GatewayPort, GatewayIP
- **AND** set the resource ID to "ZoneId:GatewayId"

#### Scenario: Create returns empty GatewayId
- **WHEN** CreateMultiPathGateway response returns an empty GatewayId
- **THEN** the system SHALL return a NonRetryableError

### Requirement: Resource Read Operation
The system SHALL implement the Read operation by calling `DescribeMultiPathGateways` API with ZoneId and searching the returned Gateways list for a matching GatewayId. If the gateway is not found, the resource SHALL be marked as removed from state.

#### Scenario: Gateway exists
- **WHEN** DescribeMultiPathGateways returns a gateway matching the GatewayId from the resource ID
- **THEN** the system SHALL set gateway_id, gateway_name, gateway_type, gateway_port, status, gateway_ip, region_id, and need_confirm from the response

#### Scenario: Gateway not found
- **WHEN** DescribeMultiPathGateways does not return a gateway matching the GatewayId
- **THEN** the system SHALL call d.SetId("") to remove the resource from state

#### Scenario: Read with retry on API failure
- **WHEN** DescribeMultiPathGateways API call fails
- **THEN** the system SHALL retry with tccommon.ReadRetryTimeout and return tccommon.RetryError

### Requirement: Resource Update Operation
The system SHALL implement the Update operation by calling `ModifyMultiPathGateway` API with ZoneId, GatewayId, GatewayName, GatewayIP, and GatewayPort parameters. Only fields that are present in the ModifyMultiPathGateway API SHALL be updatable.

#### Scenario: Update gateway name and IP
- **WHEN** gateway_name or gateway_ip is changed in the Terraform configuration
- **THEN** the system SHALL call ModifyMultiPathGateway with ZoneId, GatewayId, and the updated fields

#### Scenario: Update gateway port
- **WHEN** gateway_port is changed in the Terraform configuration
- **THEN** the system SHALL call ModifyMultiPathGateway with ZoneId, GatewayId, and the updated GatewayPort

### Requirement: Resource Delete Operation
The system SHALL implement the Delete operation by calling `DeleteMultiPathGateway` API with ZoneId and GatewayId parameters parsed from the composite resource ID.

#### Scenario: Delete existing gateway
- **WHEN** the resource is destroyed
- **THEN** the system SHALL parse ZoneId and GatewayId from the resource ID
- **AND** call DeleteMultiPathGateway with ZoneId and GatewayId

#### Scenario: Delete with retry on API failure
- **WHEN** DeleteMultiPathGateway API call fails
- **THEN** the system SHALL retry with tccommon.ReadRetryTimeout and return tccommon.RetryError

### Requirement: Resource Import Support
The system SHALL support importing existing multi-path gateways. The import ID SHALL be in the format `ZoneId:GatewayId` using tccommon.FILED_SP as separator.

#### Scenario: Import existing gateway
- **WHEN** terraform import is called with ID "zone-xxx:gw-xxx"
- **THEN** the system SHALL parse ZoneId and GatewayId from the ID
- **AND** call Read to populate the resource state

### Requirement: Provider Registration
The system SHALL register the `tencentcloud_teo_multi_path_gateway` resource in `provider.go` and add the resource entry in `provider.md`.

#### Scenario: Resource registered in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_multi_path_gateway` SHALL be available as a resource type

### Requirement: Resource Documentation
The system SHALL provide a markdown documentation file `resource_tc_teo_multi_path_gateway.md` with description, example usage, and import section.

#### Scenario: Documentation exists
- **WHEN** the resource is created
- **THEN** a .md file SHALL exist with a one-line description mentioning TEO, example usage with zone_id and other parameters, and import section showing the composite ID format

### Requirement: Unit Tests
The system SHALL provide unit tests in `resource_tc_teo_multi_path_gateway_test.go` using gomonkey to mock cloud API calls, testing Create, Read, Update, and Delete operations.

#### Scenario: Unit tests pass
- **WHEN** `go test` is run with -gcflags=all=-l on the test file
- **THEN** all test cases for Create, Read, Update, and Delete SHALL pass
