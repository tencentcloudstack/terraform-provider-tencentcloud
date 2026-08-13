## ADDED Requirements

### Requirement: 支持查询 DB Custom 集群节点列表
数据源 `tencentcloud_dbdc_db_custom_cluster_nodes` 必须能够查询指定 DB Custom 集群中的节点列表，支持按集群ID和过滤条件查询。

**Rationale**: 用户需要在 Terraform 配置中动态发现和引用特定 DB Custom 集群内的节点信息，用于基础设施规划和节点管理。

#### Scenario: 按 Cluster ID 查询节点列表
- **WHEN** 用户指定了 `cluster_id` 参数（必填）
- **THEN** 数据源返回该集群下所有节点列表

**Acceptance Criteria**:
- `cluster_id` 为 Required 参数，映射到 API 的 `ClusterId`
- 调用 `DescribeDBCustomClusterNodes` API 获取节点列表
- 内部自动分页获取所有结果（Limit=100, 递增 Offset）
- 不暴露 limit/offset 给用户

#### Scenario: 按 Filter 条件查询节点
- **WHEN** 用户指定了 `filters` 参数
- **THEN** 返回满足过滤条件的节点

**Acceptance Criteria**:
- `filters` 为 Optional 参数，映射到 API 的 `Filters`
- 支持的 filter name: `node-name`（DB Custom 节点名称）
- 每个 filter 包含 `name` 和 `values` 字段
- `name` 为 Required String，`values` 为 Required List of String
- 使用 `[]*dbdcv20201029.Filter` 类型

### Requirement: 完整的节点信息映射
数据源必须返回节点的完整详细信息，涵盖所有 `DBCustomClusterNode` 结构体中的字段。

**Rationale**: 用户需要完整的节点信息用于资源规划和管理决策，包括节点标识、网络配置、状态和规格。

#### Scenario: 返回节点基础信息字段
- **WHEN** 查询到节点列表
- **THEN** 每个 `node_set` 元素包含以下字段

**Acceptance Criteria**:
- `node_id` - 节点ID (TypeString, Computed)
- `node_name` - 节点名称 (TypeString, Computed)
- `lan_ip` - 节点内网IP地址 (TypeString, Computed)
- `ssh_endpoint` - 节点SSH访问Endpoint，格式为IP:Port (TypeString, Computed)
- `status` - 节点在集群中的实例状态 (TypeString, Computed)
- `zone` - 节点所属地域 (TypeString, Computed)
- `node_type` - 节点类型 (TypeString, Computed，枚举值包括 DB.AT5.32XLARGE512, DB.AT5.64XLARGE1152, DB.AT5.128XLARGE2304, DB.AT5.16XLARGE256, DB.AT5.8XLARGE128)

#### Scenario: 返回总数量
- **WHEN** 查询到节点列表
- **THEN** 返回集群下总的节点数量

**Acceptance Criteria**:
- `total_count` - 集群下总的节点数量 (TypeInt, Computed)
- 映射到 API 的 `Response.TotalCount`

### Requirement: 数据类型转换与空值处理
数据源必须正确处理 API 返回的数据类型转换，安全处理空值和 nil 指针。

**Rationale**: 腾讯云 SDK 使用指针类型，需要安全解引用避免程序崩溃。

#### Scenario: 安全解引用指针字段
- **WHEN** API 返回包含指针类型的字段
- **THEN** 所有指针字段都经过 nil 检查后再设置

**Acceptance Criteria**:
- 在调用 set 方法前，先判断 Response 字段是否为 nil，若为 nil 则不调用 set 方法
- nil 指针字段不会导致 panic
- 遵循项目规则 #8

### Requirement: 输出文件支持
数据源支持将查询结果输出到 JSON 文件，方便用户审查和分析。

**Rationale**: 用户可能需要将查询结果导出用于离线分析或审计。

#### Scenario: 输出结果到指定文件
- **WHEN** 用户指定了 `result_output_file` 参数
- **THEN** 将结果序列化为 JSON 并写入文件

**Acceptance Criteria**:
- `result_output_file` 参数为 Optional String 类型
- 使用 `tccommon.WriteToFile` 写入

### Requirement: 错误处理与重试
数据源必须正确处理 API 错误，实现重试逻辑以应对临时性故障。

**Rationale**: 云 API 调用可能因为网络、限流等原因失败，需要重试机制。

#### Scenario: API 调用失败时重试
- **WHEN** API 调用返回可重试错误（如限流）
- **THEN** 自动重试直到成功或超时

**Acceptance Criteria**:
- 使用 `resource.Retry` 包装 API 调用
- 设置合理的重试超时（使用 `tccommon.ReadRetryTimeout`）
- 对可重试错误使用 `tccommon.RetryError`
- 对不可重试错误使用 `resource.NonRetryableError`

#### Scenario: API 返回空结果时返回 NonRetryableError
- **WHEN** API 返回空结果（response 为 nil、Response 为 nil、NodeSet 长度为 0）
- **THEN** 不直接 `d.SetId("")`，而是返回 `NonRetryableError`

**Acceptance Criteria**:
- 在 retry 块内检查 API 返回是否为空
- 若为空，返回 `resource.NonRetryableError`
- 在 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")`
- 遵循项目规则 #14 for RESOURCE_KIND_DATASOURCE

#### Scenario: API 返回错误时记录日志
- **WHEN** API 调用失败
- **THEN** 记录详细的错误日志

**Acceptance Criteria**:
- 使用 `log.Printf` 记录错误信息
- 日志中使用资源名称 `dbdc_db_custom_cluster_nodes` 而非模糊措辞

### Requirement: 代码质量与规范
数据源代码必须符合项目规范，遵循命名和结构约定。

**Rationale**: 保持代码库一致性和可维护性。

#### Scenario: 遵循命名规范
- **WHEN** 审查数据源代码
- **THEN** 命名符合项目规范

**Acceptance Criteria**:
- 文件名: `data_source_tc_dbdc_db_custom_cluster_nodes.go`
- 数据源名: `tencentcloud_dbdc_db_custom_cluster_nodes`
- 函数名: `DataSourceTencentCloudDbdcDbCustomClusterNodes`
- 服务方法名: `DescribeDBCustomClusterNodesByFilter`
- 在日志/打印/错误信息中使用资源名称 `dbdc_db_custom_cluster_nodes`
- 不在资源 go 文件开头添加注释

#### Scenario: 遵循文件组织规范
- **WHEN** 审查文件结构
- **THEN** 文件放置在正确的位置

**Acceptance Criteria**:
- 数据源文件: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_nodes.go`
- 服务层方法: `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` (新增方法)
- 测试文件: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_nodes_test.go`
- 文档文件: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_nodes.md`
- 在 `tencentcloud/provider.go` 中注册数据源
- 在 `tencentcloud/provider.md` 中添加数据源条目

### Requirement: Service Layer 实现
必须在 `service_tencentcloud_dbdc.go` 中新增 `DescribeDBCustomClusterNodesByFilter` 方法。

**Rationale**: Service 层封装 API 调用和分页逻辑，提供复用性。

#### Scenario: Service 方法处理分页
- **WHEN** 调用 `DescribeDBCustomClusterNodesByFilter`
- **THEN** 内部自动处理分页获取所有数据

**Acceptance Criteria**:
- 内部 Limit 设置为 API 最大值 100
- Offset 递增直到返回结果数小于 Limit
- 合并所有分页结果返回完整列表
- 不暴露 limit/offset 给用户

#### Scenario: Service 方法处理参数转换
- **WHEN** 传入参数 map 包含 `ClusterId`、`Filters`
- **THEN** 正确转换为 SDK 请求参数

**Acceptance Criteria**:
- `ClusterId` 转换为 `*string` 并设置到 request
- `Filters` 转换为 `[]*dbdcv20201029.Filter`，包含 Name 和 Values

#### Scenario: Service 方法处理空响应
- **WHEN** API 返回空结果（response 为 nil 或 NodeSet 为 nil/空）
- **THEN** 返回 NonRetryableError

**Acceptance Criteria**:
- 在 retry 块内检查 response != nil && response.Response != nil && response.Response.NodeSet != nil
- 若为空，返回 `resource.NonRetryableError(fmt.Errorf("Describe dbdc_db_custom_cluster_nodes failed, Response is nil."))`

### Requirement: Provider 注册
必须在 `tencentcloud/provider.go` 中注册数据源，并在 `tencentcloud/provider.md` 中添加数据源条目。

**Rationale**: 数据源必须在 Provider 中注册才能被 Terraform 使用。

#### Scenario: 数据源在 Terraform 中可用
- **WHEN** 用户在 Terraform 配置中使用 `data "tencentcloud_dbdc_db_custom_cluster_nodes"`
- **THEN** Terraform 正常识别和使用该数据源

**Acceptance Criteria**:
- `provider.go` 的 DataSourcesMap 中添加 `"tencentcloud_dbdc_db_custom_cluster_nodes": dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodes()`
- `provider.md` 中 DBDC Data Source 部分添加 `tencentcloud_dbdc_db_custom_cluster_nodes`

### Requirement: 文档
必须提供完整的 .md 文档文件。

**Rationale**: 每个数据源必须有文档说明使用方法。

#### Scenario: 用户查看文档并使用数据源
- **WHEN** 用户查看 `data_source_tc_dbdc_db_custom_cluster_nodes.md`
- **THEN** 文档包含完整的使用说明

**Acceptance Criteria**:
- 一句话描述包含云产品名称（DBDC），格式: "Use this data source to query detailed information of DB Custom cluster nodes"
- Example Usage 部分，包含按 cluster_id 查询和按 filters 查询的示例
- 不包含 Argument Reference 和 Attribute Reference（这些由工具自动生成）

### Requirement: 单元测试
必须提供 gomonkey mock 方式的单元测试，验证数据源的 Read 功能。

**Rationale**: 项目规则要求新资源使用 mock 方式进行单元测试。

#### Scenario: 基本 Read 测试
- **WHEN** 运行 `TestDbdcDbCustomClusterNodesDS_ReadBasic`
- **THEN** 数据源正确读取并映射 API 返回的节点列表

**Acceptance Criteria**:
- 使用 gomonkey mock `DescribeDBCustomClusterNodes` API 返回
- 验证 `node_set` 包含正确的字段值
- 验证 `d.Id()` 不为空
- 使用 `go test -gcflags="all=-l"` 参数运行

#### Scenario: Schema 结构验证测试
- **WHEN** 运行 `TestDbdcDbCustomClusterNodesDS_Schema`
- **THEN** 数据源 schema 结构符合规范

**Acceptance Criteria**:
- 验证所有必需的 schema 字段存在（cluster_id, filters, node_set, total_count, result_output_file）
- 验证字段类型和属性正确（cluster_id Required, node_set Computed, etc.）

#### Scenario: 空响应处理测试
- **WHEN** API 返回空结果
- **THEN** 数据源返回错误（NonRetryableError）

**Acceptance Criteria**:
- mock API 返回空 NodeSet
- 验证 Read 方法返回错误
