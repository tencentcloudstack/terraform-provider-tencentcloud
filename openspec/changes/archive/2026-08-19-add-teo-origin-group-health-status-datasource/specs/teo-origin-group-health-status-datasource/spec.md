## ADDED Requirements

### Requirement: Data source registration
数据源 MUST 在 Provider 中注册为 `tencentcloud_teo_origin_group_health_status`，使其可在 Terraform 配置中使用。

#### Scenario: Data source is accessible
- **WHEN** 用户在 Terraform 配置中使用 `data "tencentcloud_teo_origin_group_health_status"`
- **THEN** Terraform 能够成功识别并初始化该数据源

### Requirement: Required parameters
数据源 MUST 要求用户提供站点 ID（`zone_id`）和负载均衡实例 ID（`lb_instance_id`）参数。

#### Scenario: Query with required parameters
- **WHEN** 用户设置 `zone_id = "zone-xxx"` 和 `lb_instance_id = "lb-xxx"`
- **THEN** 数据源调用 `DescribeOriginGroupHealthStatus` API 时传入对应的 ZoneId 和 LBInstanceId 参数

#### Scenario: Missing required parameter
- **WHEN** 用户未设置 `zone_id` 或 `lb_instance_id`
- **THEN** Terraform 在 plan 阶段报错，提示缺少必填参数

### Requirement: Optional origin group filter
数据源 MUST 支持通过 `origin_group_ids` 参数可选地过滤指定源站组的健康状态。

#### Scenario: Filter by origin group IDs
- **WHEN** 用户设置 `origin_group_ids = ["origin-group-1", "origin-group-2"]`
- **THEN** 数据源调用 API 时传入 OriginGroupIds 参数
- **THEN** 只返回指定源站组的健康状态

#### Scenario: Query all origin groups
- **WHEN** 用户未设置 `origin_group_ids` 参数
- **THEN** 数据源调用 API 时不传入 OriginGroupIds 参数
- **THEN** 返回负载均衡下所有源站组的健康状态

### Requirement: Return origin group health status list
数据源 MUST 返回源站组健康状态列表 `origin_group_health_status_list`，包含完整的健康状态信息。

#### Scenario: Origin group basic information
- **WHEN** 查询成功返回健康状态列表
- **THEN** 每个源站组 MUST 包含 `origin_group_id` 字段

#### Scenario: Origin health status
- **WHEN** 查询成功返回健康状态列表
- **THEN** 每个源站组 MUST 包含 `origin_health_status` 列表，包含综合决策的源站健康状态
- **THEN** 每个源站健康状态项 MUST 包含 `origin`（源站）和 `healthy`（健康状态：Healthy/Unhealthy/Undetected）字段

#### Scenario: Check region health status
- **WHEN** 查询成功返回健康状态列表
- **THEN** 每个源站组 MUST 包含 `check_region_health_status` 列表，包含各健康检查区域的源站健康状态
- **THEN** 每个检查区域健康状态项 MUST 包含 `region`（区域）、`healthy`（健康状态）和 `origin_health_status`（源站健康状态列表）字段

### Requirement: API retry for consistency
数据源 MUST 使用重试机制处理 API 最终一致性问题。

#### Scenario: Transient API error
- **WHEN** API 调用返回临时性错误（如限流、服务暂时不可用）
- **THEN** 数据源自动重试 API 调用
- **THEN** 在合理的重试次数内成功后返回结果

#### Scenario: Persistent API error
- **WHEN** API 调用多次重试后仍然失败
- **THEN** 数据源返回错误信息给 Terraform
- **THEN** 不进行无限重试

### Requirement: Empty result handling
数据源 MUST 正确处理 API 返回空结果的情况。

#### Scenario: API returns empty list
- **WHEN** API 返回的 `OriginGroupHealthStatusList` 为空
- **THEN** 数据源不直接 `d.SetId("")`
- **THEN** 返回 `NonRetryableError` 让外层 retry 继续尝试
- **THEN** 外层 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示

### Requirement: Nil field handling
数据源 MUST 正确处理 API 响应中可能为 nil 的字段。

#### Scenario: Nil field in response
- **WHEN** API 响应中某个字段为 nil
- **THEN** 数据源跳过该字段的 setXX() 调用
- **THEN** 不发生空指针错误

### Requirement: Export results to file
数据源 MUST 支持通过 `result_output_file` 参数将查询结果导出到 JSON 文件。

#### Scenario: Export to specified file
- **WHEN** 用户设置 `result_output_file = "/tmp/origin_group_health_status.json"`
- **THEN** 查询结果被写入到指定路径的 JSON 文件
- **THEN** 文件内容包含完整的健康状态列表数据

### Requirement: Documentation
数据源 MUST 提供完整的使用文档。

#### Scenario: Documentation includes examples
- **WHEN** 用户查看数据源文档
- **THEN** 文档包含一句话描述（带上 TEO 产品名称）
- **THEN** 文档包含 Example Usage 示例

#### Scenario: Documentation format
- **WHEN** 生成文档
- **THEN** 文档不包含 Argument Reference 和 Attribute Reference 章节（由工具自动生成）

### Requirement: Unit testing
数据源 MUST 包含单元测试以确保功能正确性。

#### Scenario: Mock-based unit test
- **WHEN** 运行单元测试
- **THEN** 使用 gomonkey 对云 API 进行 mock 处理
- **THEN** 只进行业务代码逻辑的单元测试，不使用 terraform 测试套件
