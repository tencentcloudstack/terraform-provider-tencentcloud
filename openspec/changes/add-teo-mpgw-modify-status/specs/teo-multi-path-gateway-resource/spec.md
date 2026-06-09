## MODIFIED Requirements

### Requirement: Resource Schema Definition
The system SHALL define a Terraform resource `tencentcloud_teo_multi_path_gateway` with the following schema fields:
- `zone_id` (Required, ForceNew, TypeString): 站点 ID
- `gateway_type` (Required, ForceNew, TypeString): 网关类型，取值 cloud 或 private
- `gateway_name` (Required, TypeString): 网关名称，16 个字符以内
- `gateway_port` (Optional, Computed, TypeInt): 网关端口，范围 1～65535（除去 8888）
- `region_id` (Optional, Computed, ForceNew, TypeString): 网关地域，GatewayType 为 cloud 时必填
- `gateway_ip` (Optional, Computed, TypeString): 网关地址，GatewayType 为 private 时必填
- `gateway_id` (Computed, TypeString): 网关 ID
- `status` (Optional, Computed, TypeString): 网关状态。合法取值：`online`（启用）、`offline`（停用）。未显式配置时由云端回填。
- `need_confirm` (Computed, TypeString): 是否需要重新确认回源 IP 列表

#### Scenario: Schema defines all CRUD fields
- **WHEN** the resource schema is defined
- **THEN** it SHALL include zone_id, gateway_type, gateway_name, gateway_port, region_id, gateway_ip, gateway_id, status, and need_confirm fields with correct types and constraints

#### Scenario: ForceNew fields prevent in-place update
- **WHEN** zone_id, gateway_type, or region_id is changed in the Terraform configuration
- **THEN** the resource SHALL be destroyed and recreated

#### Scenario: Status field is optional and computed
- **WHEN** a user does not set `status` in the Terraform configuration
- **THEN** the system SHALL populate `status` from the API response (Computed behavior) without triggering a plan diff

#### Scenario: Status field accepts user input
- **WHEN** a user sets `status = "online"` or `status = "offline"` in the Terraform configuration
- **THEN** the schema SHALL accept the value and mark it as a candidate for update

### Requirement: Resource Update Operation
The system SHALL implement the Update operation by calling `ModifyMultiPathGateway` API with ZoneId, GatewayId, GatewayName, GatewayIP, and GatewayPort parameters when `gateway_name`, `gateway_ip`, or `gateway_port` changes. Additionally, when `status` changes and is explicitly set in the configuration, the system SHALL call `ModifyMultiPathGatewayStatus` API with ZoneId, GatewayId, and GatewayStatus parameters, and then poll `DescribeMultiPathGateways` until the gateway status stabilizes.

#### Scenario: Update gateway name and IP
- **WHEN** gateway_name or gateway_ip is changed in the Terraform configuration
- **THEN** the system SHALL call ModifyMultiPathGateway with ZoneId, GatewayId, and the updated fields

#### Scenario: Update gateway port
- **WHEN** gateway_port is changed in the Terraform configuration
- **THEN** the system SHALL call ModifyMultiPathGateway with ZoneId, GatewayId, and the updated GatewayPort

#### Scenario: Enable gateway via status change
- **WHEN** `status` changes to `"online"` in the Terraform configuration
- **THEN** the system SHALL call ModifyMultiPathGatewayStatus with ZoneId, GatewayId, GatewayStatus="online"
- **AND** poll DescribeMultiPathGateways until the gateway Status reaches a stable (non-transient) value or the update timeout is reached

#### Scenario: Disable gateway via status change
- **WHEN** `status` changes to `"offline"` in the Terraform configuration
- **THEN** the system SHALL call ModifyMultiPathGatewayStatus with ZoneId, GatewayId, GatewayStatus="offline"
- **AND** poll DescribeMultiPathGateways until the gateway Status reaches a stable value or the update timeout is reached

#### Scenario: Status not set in configuration
- **WHEN** `status` is not present (or removed) in the Terraform configuration and no other mutable fields change
- **THEN** the system SHALL NOT call ModifyMultiPathGatewayStatus

#### Scenario: ModifyMultiPathGatewayStatus API failure
- **WHEN** ModifyMultiPathGatewayStatus returns a retryable error
- **THEN** the system SHALL retry using tccommon.WriteRetryTimeout via resource.Retry
- **AND** return tccommon.RetryError on non-retryable errors

### Requirement: Unit Tests
The system SHALL provide unit tests in `resource_tc_teo_multi_path_gateway_test.go` using gomonkey to mock cloud API calls, testing Create, Read, Update (including status change), and Delete operations.

#### Scenario: Unit tests pass
- **WHEN** `go test` is run with -gcflags=all=-l on the test file
- **THEN** all test cases for Create, Read, Update, status change (online/offline), and Delete SHALL pass

#### Scenario: Status change test case
- **WHEN** a test simulates changing `status` from `"online"` to `"offline"`
- **THEN** the mocked ModifyMultiPathGatewayStatus SHALL be invoked with GatewayStatus="offline"
- **AND** the test SHALL assert the API was called exactly once

### Requirement: Resource Documentation
The system SHALL provide a markdown documentation file `resource_tc_teo_multi_path_gateway.md` with description, example usage (including an example of the `status` field usage), and import section.

#### Scenario: Documentation exists
- **WHEN** the resource is created
- **THEN** a .md file SHALL exist with a one-line description mentioning TEO, example usage with zone_id and other parameters (including `status`), and import section showing the composite ID format
