## ADDED Requirements

### Requirement: Resource schema definition
The system SHALL define a Terraform resource `tencentcloud_teo_multi_path_gateway` with the following schema fields:
- `zone_id` (string, required, ForceNew): 站点 ID
- `gateway_type` (string, required, ForceNew): 网关类型，取值 cloud 或 private
- `gateway_name` (string, required): 网关名称，16个字符以内
- `gateway_port` (int, required): 网关端口，范围 1~65535（除去 8888）
- `region_id` (string, optional, ForceNew): 网关地域，GatewayType 为 cloud 时必填
- `gateway_ip` (string, optional): 网关地址，GatewayType 为 private 时必填
- `gateway_id` (string, computed): 网关 ID
- `status` (string, computed): 网关状态
- `need_confirm` (string, computed): 回源 IP 列表变化是否需要确认
- `lines` (list, computed): 线路信息列表，每个元素包含 line_id、line_type、line_address、proxy_id、rule_id

#### Scenario: Schema defines all required and computed fields
- **WHEN** the resource schema is initialized
- **THEN** all fields SHALL be present with correct types, and gateway_type, zone_id, region_id SHALL have ForceNew set to true

### Requirement: Resource Create
The system SHALL implement a Create function that calls `CreateMultiPathGateway` API with zone_id, gateway_type, gateway_name, gateway_port, region_id (if provided), and gateway_ip (if provided).

#### Scenario: Create cloud type gateway
- **WHEN** user creates a resource with gateway_type="cloud", zone_id, gateway_name, gateway_port, and region_id
- **THEN** the system SHALL call CreateMultiPathGateway with all provided fields, set the resource ID to `zone_id#gateway_id`, and refresh state via Read

#### Scenario: Create private type gateway
- **WHEN** user creates a resource with gateway_type="private", zone_id, gateway_name, gateway_port, and gateway_ip
- **THEN** the system SHALL call CreateMultiPathGateway with all provided fields, set the resource ID to `zone_id#gateway_id`, and refresh state via Read

#### Scenario: Create returns empty GatewayId
- **WHEN** CreateMultiPathGateway API returns an empty GatewayId
- **THEN** the system SHALL return a NonRetryableError

### Requirement: Resource Read
The system SHALL implement a Read function that parses the composite ID (`zone_id#gateway_id`), calls `DescribeMultiPathGateways` with Filters to filter by gateway-id, and sets all schema fields from the response.

#### Scenario: Read existing gateway
- **WHEN** the resource ID is set and the gateway exists
- **THEN** the system SHALL call DescribeMultiPathGateways with ZoneId and gateway-id Filter, find the matching gateway, and set all schema fields from the response

#### Scenario: Read deleted gateway
- **WHEN** the gateway no longer exists (not found in the response list)
- **THEN** the system SHALL call `d.SetId("")` to signal resource deletion to Terraform

#### Scenario: Read with nil response fields
- **WHEN** the API response contains nil fields
- **THEN** the system SHALL skip calling set methods for those nil fields

### Requirement: Resource Update
The system SHALL implement an Update function that calls `ModifyMultiPathGateway` API with zone_id, gateway_id, and the changed fields (gateway_name, gateway_ip, gateway_port).

#### Scenario: Update gateway name
- **WHEN** user changes gateway_name
- **THEN** the system SHALL call ModifyMultiPathGateway with the new gateway_name, zone_id, and gateway_id

#### Scenario: Update gateway port
- **WHEN** user changes gateway_port
- **THEN** the system SHALL call ModifyMultiPathGateway with the new gateway_port, zone_id, and gateway_id

#### Scenario: Update gateway IP
- **WHEN** user changes gateway_ip
- **THEN** the system SHALL call ModifyMultiPathGateway with the new gateway_ip, zone_id, and gateway_id

#### Scenario: No changes to updateable fields
- **WHEN** no updateable fields have changed
- **THEN** the system SHALL skip the Modify API call and only refresh state via Read

### Requirement: Resource Delete
The system SHALL implement a Delete function that calls `DeleteMultiPathGateway` API with zone_id and gateway_id.

#### Scenario: Delete existing gateway
- **WHEN** user destroys the resource
- **THEN** the system SHALL call DeleteMultiPathGateway with zone_id and gateway_id parsed from the composite ID, and return nil on success

### Requirement: Resource Import
The system SHALL support resource import with the composite ID format `zone_id#gateway_id`.

#### Scenario: Import existing gateway
- **WHEN** user runs `terraform import tencentcloud_teo_multi_path_gateway.example zone-xxx#gw-xxx`
- **THEN** the system SHALL parse the composite ID, call Read to populate the state, and the resource SHALL be managed by Terraform

### Requirement: Service layer Describe method
The system SHALL add a `DescribeTeoMultiPathGatewayById` method to the teo service layer that calls `DescribeMultiPathGateways` with Filters to query a specific gateway by its ID.

#### Scenario: Query gateway by ID
- **WHEN** the service method is called with zoneId and gatewayId
- **THEN** it SHALL call DescribeMultiPathGateways with ZoneId and a gateway-id Filter, and return the matching MultiPathGateway object

#### Scenario: Gateway not found
- **WHEN** the gateway is not found in the response
- **THEN** the method SHALL return an error indicating the resource was not found

### Requirement: Provider registration
The system SHALL register the new resource `tencentcloud_teo_multi_path_gateway` in `provider.go` ResourcesMap and add an entry in `provider.md`.

#### Scenario: Resource is registered in provider
- **WHEN** the provider is initialized
- **THEN** `tencentcloud_teo_multi_path_gateway` SHALL be available as a resource with the constructor `teo.ResourceTencentCloudTeoMultiPathGateway()`

### Requirement: Documentation
The system SHALL generate a `.md` documentation file for the resource following the gendoc format.

#### Scenario: Documentation file exists
- **WHEN** the resource is created
- **THEN** a `resource_tc_teo_multi_path_gateway.md` file SHALL exist in `tencentcloud/services/teo/` with a one-line description, example usage, and import section

### Requirement: Unit tests
The system SHALL create unit test file `resource_tc_teo_multi_path_gateway_test.go` using gomonkey mock approach for testing business logic.

#### Scenario: Unit tests cover CRUD operations
- **WHEN** the test file is executed with `go test -gcflags=all=-l`
- **THEN** tests SHALL cover Create, Read, Update, and Delete operations using gomonkey to mock cloud API calls
