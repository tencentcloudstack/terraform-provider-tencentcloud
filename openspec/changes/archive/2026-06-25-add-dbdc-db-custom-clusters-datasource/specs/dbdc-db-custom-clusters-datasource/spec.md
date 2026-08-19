## ADDED Requirements

### Requirement: 支持查询 DB Custom 集群列表
数据源 `tencentcloud_dbdc_db_custom_clusters` 必须能够查询腾讯云 DB Custom 集群列表，支持多种过滤条件。

**Rationale**: 用户需要在 Terraform 配置中动态发现和引用已存在的 DB Custom 集群，支持批量查询和资源规划。

#### Scenario: 查询所有集群
- **WHEN** 用户配置了数据源但未指定任何过滤条件
- **THEN** 数据源返回当前账号下所有 DB Custom 集群的列表

**Acceptance Criteria**:
- 调用 `DescribeDBCustomClusters` API 获取集群列表
- 返回的 `cluster_set` 包含所有集群
- 每个集群包含完整的基础信息字段
- 内部自动分页获取所有结果（Limit=100, 递增 Offset）

#### Scenario: 按 Cluster IDs 精确查询
- **WHEN** 用户指定了 `cluster_ids` 参数
- **THEN** 返回匹配指定 ID 列表的集群

**Acceptance Criteria**:
- `cluster_ids` 映射到 API 的 `ClusterIds` 参数
- 支持多个 cluster ID 同时查询（最多100个）
- 返回所有匹配的集群

#### Scenario: 按 Filter 条件查询
- **WHEN** 用户指定了 `filters` 参数
- **THEN** 返回满足过滤条件的集群

**Acceptance Criteria**:
- `filters` 映射到 API 的 `Filters` 参数
- 支持的 filter name: `cluster-name`（精确匹配）、`cluster-status`（Creating, Running, Destroying）
- 每个 filter 包含 `name` 和 `values` 字段

#### Scenario: 按 Tags 查询
- **WHEN** 用户指定了 `tags` 参数
- **THEN** 返回包含所有指定标签的集群

**Acceptance Criteria**:
- `tags` 映射到 API 的 `Tags` 参数
- 支持多标签同时过滤
- 每个 tag 包含 `key` 和 `value` 字段

### Requirement: 完整的集群信息映射
数据源必须返回集群的完整详细信息，涵盖所有 DBCustomCluster 结构体中的关键字段。

**Rationale**: 用户需要完整的集群信息用于资源规划和管理决策。

#### Scenario: 返回基础信息字段
- **WHEN** 查询到集群列表
- **THEN** 每个 `cluster_set` 元素包含以下字段

**Acceptance Criteria**:
- `cluster_id` - 集群唯一标识 (TypeString, Computed)
- `cluster_name` - 集群名称 (TypeString, Computed)
- `region` - 集群支持的地域 (TypeString, Computed)
- `cluster_level` - 集群规模 (TypeString, Computed)
- `cluster_status` - 集群状态 (TypeString, Computed, 枚举: Creating, Running, Destroying)
- `cluster_version` - 集群版本号 (TypeString, Computed)
- `cluster_node_num` - 集群中的节点数量 (TypeInt, Computed)
- `cluster_description` - 集群描述 (TypeString, Computed)
- `created_time` - 创建时间 (TypeString, Computed)
- `tags` - 集群标签信息 (TypeList, Computed, 包含 key 和 value)

### Requirement: 数据类型转换与空值处理
数据源必须正确处理 API 返回的数据类型转换，安全处理空值和 nil 指针。

**Rationale**: 腾讯云 SDK 使用指针类型，DBCustomCluster.Tags 字段可能返回 null，需要安全解引用避免程序崩溃。

#### Scenario: 安全解引用指针字段
- **WHEN** API 返回包含指针类型的字段
- **THEN** 所有指针字段都经过 nil 检查后再设置

**Acceptance Criteria**:
- 在调用 set 方法前，先判断 Response 字段是否为 nil，若为 nil 则不调用 set 方法
- nil 指针字段不会导致 panic
- 遵循项目规则 #8

#### Scenario: Tags 字段空值处理
- **WHEN** API 返回的 Tags 字段为 null
- **THEN** 不调用 Tags 的 set 方法，省略该字段

**Acceptance Criteria**:
- DBCustomCluster.Tags 字段注释标注"可能返回 null"
- nil 检查在 set 之前执行
- 不因 Tags 为 nil 而导致数据源报错

### Requirement: 输出文件支持
数据源支持将查询结果输出到 JSON 文件，方便用户审查和分析。

**Rationale**: 用户可能需要将查询结果导出用于离线分析或审计。

#### Scenario: 输出结果到指定文件
- **WHEN** 用户指定了 `result_output_file` 参数
- **THEN** 将结果序列化为 JSON 并写入文件

**Acceptance Criteria**:
- `result_output_file` 参数为 Optional String 类型
- 结果以可读的 JSON 格式输出
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
- **WHEN** API 返回空结果（response 为 nil、Response 为 nil、ClusterSet 长度为 0）
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
- 日志包含 logId 和错误原因
- 日志中使用资源名称 `dbdc_db_custom_clusters` 而非模糊措辞

### Requirement: 代码质量与规范
数据源代码必须符合项目规范，遵循命名和结构约定。

**Rationale**: 保持代码库一致性和可维护性。

#### Scenario: 遵循命名规范
- **WHEN** 审查数据源代码
- **THEN** 命名符合项目规范

**Acceptance Criteria**:
- 文件名: `data_source_tc_dbdc_db_custom_clusters.go`
- 数据源名: `tencentcloud_dbdc_db_custom_clusters`
- 函数名: `DataSourceTencentCloudDbdcDbCustomClusters`
- 服务方法名: `DescribeDBCustomClustersByFilter`
- 在日志/打印/错误信息中使用资源名称 `dbdc_db_custom_clusters`
- 不在资源 go 文件开头添加注释

#### Scenario: 遵循文件组织规范
- **WHEN** 审查文件结构
- **THEN** 文件放置在正确的位置

**Acceptance Criteria**:
- 数据源文件: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_clusters.go`
- 服务层文件: `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go`
- 测试文件: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_clusters_test.go`
- 文档文件: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_clusters.md`
- 在 `tencentcloud/provider.go` 中注册数据源
- 在 `tencentcloud/provider.md` 中添加数据源条目

### Requirement: Service Layer 实现
必须在 `service_tencentcloud_dbdc.go` 中实现 `DescribeDBCustomClustersByFilter` 方法。

**Rationale**: Service 层封装 API 调用和分页逻辑，提供复用性。

#### Scenario: Service 方法处理分页
- **WHEN** 调用 `DescribeDBCustomClustersByFilter`
- **THEN** 内部自动处理分页获取所有数据

**Acceptance Criteria**:
- 内部 Limit 设置为 API 最大值 100
- Offset 递增直到返回结果数小于 Limit
- 合并所有分页结果返回完整列表
- 不暴露 limit/offset 给用户

#### Scenario: Service 方法处理参数转换
- **WHEN** 传入参数 map 包含 `ClusterIds`、`Filters`、`Tags`
- **THEN** 正确转换为 SDK 请求参数

**Acceptance Criteria**:
- `ClusterIds` 转换为 `[]*string`
- `Filters` 转换为 `[]*dbdc.Filter`，包含 Name 和 Values
- `Tags` 转换为 `[]*dbdc.Tag`，包含 Key 和 Value

### Requirement: Provider 注册
必须在 `tencentcloud/provider.go` 中注册数据源，并在 `tencentcloud/provider.md` 中添加数据源条目。

**Rationale**: 数据源必须在 Provider 中注册才能被 Terraform 使用。

#### Scenario: 数据源在 Terraform 中可用
- **WHEN** 用户在 Terraform 配置中使用 `data "tencentcloud_dbdc_db_custom_clusters"`
- **THEN** Terraform 正常识别和使用该数据源

**Acceptance Criteria**:
- `provider.go` 的 DataSourcesMap 中添加 `"tencentcloud_dbdc_db_custom_clusters": dbdc.DataSourceTencentCloudDbdcDbCustomClusters()`
- `provider.md` 中添加数据源条目

### Requirement: 文档
必须提供完整的 .md 文档文件。

**Rationale**: 每个数据源必须有文档说明使用方法。

#### Scenario: 用户查看文档并使用数据源
- **WHEN** 用户查看 `data_source_tc_dbdc_db_custom_clusters.md`
- **THEN** 文档包含完整的使用说明

**Acceptance Criteria**:
- 一句话描述，包含云产品名称 (DB Custom)，格式: "Use this data source to query ..."
- Example Usage 部分
- 不包含 Argument Reference 和 Attribute Reference（这些由工具自动生成）
