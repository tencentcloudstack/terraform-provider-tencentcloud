## ADDED Requirements

### Requirement: Data source schema definition
The system SHALL provide a data source `tencentcloud_teo_multi_path_gateway` with the following schema:
- `zone_id` (TypeString, Required): 站点 ID
- `filters` (TypeList, Optional): 过滤条件，包含 `name`（TypeString, Required）和 `values`（TypeSet of TypeString, Required）子字段
- `gateways` (TypeList, Computed): 网关详情列表，包含以下子字段：
  - `gateway_id` (TypeString, Computed): 网关 ID
  - `gateway_name` (TypeString, Computed): 网关名
  - `gateway_type` (TypeString, Computed): 网关类型
  - `gateway_port` (TypeInt, Computed): 网关端口
  - `status` (TypeString, Computed): 网关状态
  - `gateway_ip` (TypeString, Computed): 网关 IP
  - `region_id` (TypeString, Computed): 网关地域 ID
  - `need_confirm` (TypeString, Computed): 是否需要确认回源 IP 变化
- `result_output_file` (TypeString, Optional): 结果输出文件路径

#### Scenario: Query all gateways by zone_id
- **WHEN** user provides only `zone_id` without `filters`
- **THEN** system SHALL call `DescribeMultiPathGateways` with the given ZoneId and return all gateways in that zone

#### Scenario: Query gateways with filters
- **WHEN** user provides `zone_id` and `filters` with name="gateway-type" and values=["cloud"]
- **THEN** system SHALL call `DescribeMultiPathGateways` with the given ZoneId and Filters, returning only matching gateways

### Requirement: Service layer implementation
The system SHALL implement a service method `DescribeTeoMultiPathGatewaysByFilter` in `service_tencentcloud_teo.go` that:
- Accepts `context.Context` and `map[string]interface{}` parameters
- Maps `ZoneId` and `Filters` from the parameter map to the SDK request
- Uses pagination with Limit=1000 and automatic offset increment
- Wraps API calls in `resource.Retry` with `tccommon.ReadRetryTimeout`
- Returns `[]*teov20220901.MultiPathGateway` slice

#### Scenario: Paginated query returns all gateways
- **WHEN** the total number of gateways exceeds one page (1000 items)
- **THEN** system SHALL automatically paginate through all pages and return the complete list

#### Scenario: API call fails with retryable error
- **WHEN** the `DescribeMultiPathGateways` API call returns a retryable error
- **THEN** system SHALL retry the call using `tccommon.RetryError` within `tccommon.ReadRetryTimeout`

### Requirement: Provider registration
The system SHALL register the data source `tencentcloud_teo_multi_path_gateway` in both `provider.go` DataSourcesMap and `provider.md`.

#### Scenario: Data source is available in Terraform
- **WHEN** user writes a Terraform configuration using `data.tencentcloud_teo_multi_path_gateway`
- **THEN** Terraform SHALL recognize the data source and be able to plan and apply it

### Requirement: Documentation
The system SHALL provide a `.md` documentation file at `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway.md` with usage examples showing query by zone_id and optional filters.

#### Scenario: Documentation exists
- **WHEN** the data source is implemented
- **THEN** a corresponding `.md` file SHALL exist with at least one example usage

### Requirement: Lines field excluded
The system SHALL NOT include the `Lines` field in the gateways output, as it is not returned by the `DescribeMultiPathGateways` list API.

#### Scenario: Lines field not present in output
- **WHEN** the data source reads gateway information
- **THEN** the `gateways` output SHALL NOT contain a `lines` field
