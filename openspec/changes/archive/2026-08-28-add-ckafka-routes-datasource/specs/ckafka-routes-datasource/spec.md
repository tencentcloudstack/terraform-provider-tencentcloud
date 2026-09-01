## ADDED Requirements

### Requirement: Data Source Schema Definition

`tencentcloud_ckafka_routes` 数据源 SHALL 支持以下输入参数与输出属性。

**Input Parameters:**
- `instance_id` (String, Required): CKafka 实例 ID
- `route_id` (Int, Optional): 路由 ID，用于精确查询单条路由
- `main_route_flag` (Bool, Optional): 是否显示主路由
- `result_output_file` (String, Optional): 输出结果到文件

**Output Attributes:**
- `routers` (List): 路由信息列表，每个元素包含：
  - `access_type` (Int): 实例接入方式（0: PLAINTEXT, 1: SASL_PLAINTEXT, 2: SSL, 3: SASL_SSL）
  - `route_id` (Int): 路由 ID
  - `vip_type` (Int): 路由网络类型（3: vpc 路由, 7: 内部支撑路由, 1: 公网路由）
  - `vip_list` (List): 虚拟 IP 列表，每个元素含 `vip` (String)、`vport` (String)
  - `domain` (String): 域名
  - `domain_port` (Int): 域名端口
  - `delete_timestamp` (String): 删除时间戳
  - `subnet` (String): 子网 ID
  - `broker_vip_list` (List): broker 虚拟 IP 列表，每个元素含 `vip` (String)、`vport` (String)
  - `vpc_id` (String): 私有网络 ID
  - `note` (String): 备注信息
  - `status` (Int): 路由状态（1: 创建中, 2: 创建成功, 3: 创建失败, 4: 删除中, 6: 删除失败）

#### Scenario: Query all routes by instance id

```hcl
data "tencentcloud_ckafka_routes" "example" {
  instance_id = "ckafka-bqwlyrg8"
}

output "routes" {
  value = data.tencentcloud_ckafka_routes.example.routers
}
```
- **WHEN** user provides only `instance_id`
- **THEN** the data source calls `DescribeRoute` with that instance ID
- **AND** returns all routes of the instance in `routers`

#### Scenario: Query a specific route by route id

- **WHEN** user provides `instance_id` and `route_id`
- **THEN** the data source calls `DescribeRoute` with both parameters
- **AND** returns the matching route in `routers`

#### Scenario: Query routes including main route

- **WHEN** user provides `instance_id` and sets `main_route_flag = true`
- **THEN** the data source calls `DescribeRoute` with `MainRouteFlag` set to true
- **AND** returns the route list including the main route created at instance creation

#### Scenario: Output to file

- **WHEN** user specifies `result_output_file` parameter
- **THEN** the route information is written to the specified file in JSON format

### Requirement: API Integration

数据源 SHALL 集成 CKafka API v20190819 的 `DescribeRoute` 接口。

#### Scenario: API call execution

- **WHEN** the data source is read
- **THEN** it calls the service layer method `DescribeCkafkaRouteByFilter` which invokes `DescribeRoute`
- **AND** wraps the API call in `resource.Retry(tccommon.ReadRetryTimeout, ...)` for retry handling
- **AND** uses `tccommon.RetryError()` to wrap API errors for retry classification
- **AND** logs API calls using `tccommon.LogElapsed` and debug logging

### Requirement: Error Handling

数据源 SHALL 正确处理空响应与各类错误情况。

#### Scenario: API returns empty response

- **WHEN** the `DescribeRoute` response is nil, or `response.Response` is nil, or `response.Response.Result` is nil, or `len(response.Response.Result.Routers) == 0`
- **THEN** the data source returns a `NonRetryableError` within the retry block
- **AND** does NOT call `d.SetId("")` to avoid clearing state on transient API fluctuation
- **AND** logs `log.Printf("[DATASOURCE] read empty, skip SetId")` on the retry failure path

#### Scenario: Instance not found

- **WHEN** the specified CKafka instance does not exist
- **THEN** the data source returns the API error message indicating the instance was not found

#### Scenario: Network or service error

- **WHEN** the API call fails due to network issues or service unavailability
- **THEN** the data source retries according to the configured retry policy
- **AND** returns an appropriate error if all retries fail
