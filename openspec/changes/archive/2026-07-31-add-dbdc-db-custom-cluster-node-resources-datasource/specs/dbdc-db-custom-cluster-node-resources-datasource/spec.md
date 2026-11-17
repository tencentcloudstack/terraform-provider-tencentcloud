## ADDED Requirements

### Requirement: 支持查询 DB Custom 集群节点资源信息
数据源 `tencentcloud_dbdc_db_custom_cluster_node_resources` 必须能够查询指定 DB Custom 集群内节点的资源信息，支持按集群ID和节点ID列表查询。

**Rationale**: 用户需要在 Terraform 配置中动态发现和引用 DB Custom 集群内节点的资源容量、可分配量、申请量、限制量及可调度余量信息，用于容量规���、节点选择和资源调度决策。

#### Scenario: 按 Cluster ID 查询节点资源信息
- **WHEN** 用户指定了 `cluster_id` 参数（必填）
- **THEN** 数据源返回该集群下所有节点的资源信息列表

**Acceptance Criteria**:
- `cluster_id` 为 Required 参数，映射到 API 的 `request.ClusterId`（`*string`）
- 调用 `DescribeDBCustomClusterNodeResources` API 获取节点资源信息
- 该接口为同步接口，无分页参数（Offset/Limit），无需分页循环
- 不暴露 limit/offset 给用户

#### Scenario: 按 Node IDs 查询指定节点资源信息
- **WHEN** 用户指定了 `node_ids` 参数
- **THEN** 仅返回指定节点的资源信息

**Acceptance Criteria**:
- `node_ids` 为 Optional 参数，映射到 API 的 `request.NodeIds`（`[]*string`）
- 每次请求的 `node_ids` 数量上限为 50（由云 API 强制约束）
- schema 描述中说明该上限
- `node_ids` 为 List of String 类型，元素类型为 `schema.TypeString`

### Requirement: 完整的节点资源信息映射
数据源必须返回节点的完整资源详细信息，涵盖 `DBCustomClusterNodeResource` 结构体中的所有字段及每个 `MetaResource` 嵌套字段。

**Rationale**: 用户需要完整的资源指标（容量、可分配、申请、限制、可用余量）用于调度和容量决策。

#### Scenario: 返回节点基础信息字段
- **WHEN** 查询到节点资源列表
- **THEN** 每个 `node_set` 元素包含节点ID字段

**Acceptance Criteria**:
- `node_id` - 节点ID (TypeString, Computed)，映射到 `DBCustomClusterNodeResource.NodeId`

#### Scenario: 返回 Capacity 资源容量字段
- **WHEN** 查询到节点资源列表
- **THEN** 每个 `node_set` 元素包含 `capacity` 嵌套块（节点物理资源总容量）

**Acceptance Criteria**:
- `capacity` 为 TypeList（MaxItems: 1）, Computed
- 嵌套字段:
  - `cpu` - CPU核心数，单位：核 (TypeFloat, Computed)，映射到 `Capacity.Cpu`（`*float64`）
  - `memory` - 内存，单位：GiB (TypeFloat, Computed)，映射到 `Capacity.Memory`（`*float64`）
  - `pods` - POD数量，单位：个 (TypeInt, Computed)，映射到 `Capacity.Pods`（`*uint64`）

#### Scenario: 返回 Allocatable 可分配容量字段
- **WHEN** 查询到节点资源列表
- **THEN** 每个 `node_set` 元素包含 `allocatable` 嵌套块（节点可分配容量 = Capacity - 系统预留）

**Acceptance Criteria**:
- `allocatable` 为 TypeList（MaxItems: 1）, Computed
- 嵌套字段 `cpu`/`memory`/`pods` 同 Capacity，映射到 `Allocatable.Cpu`/`Allocatable.Memory`/`Allocatable.Pods`

#### Scenario: 返回 Requests 申请量字段
- **WHEN** 查询到节点资源列表
- **THEN** 每个 `node_set` 元素包含 `requests` 嵌套块（节点上所有非终态 Pod 的 requests 申请量之和，含系统 Pod）

**Acceptance Criteria**:
- `requests` 为 TypeList（MaxItems: 1）, Computed
- 嵌套字段 `cpu`/`memory`/`pods` 同 Capacity，映射到 `Requests.Cpu`/`Requests.Memory`/`Requests.Pods`

#### Scenario: 返回 Limits 限制量字段
- **WHEN** 查询到节点资源列表
- **THEN** 每个 `node_set` 元素包含 `limits` 嵌套块（节点上所有非终态 Pod 的 limits 上限之和，含系统 Pod）

**Acceptance Criteria**:
- `limits` 为 TypeList（MaxItems: 1）, Computed
- 嵌套字段 `cpu`/`memory`/`pods` 同 Capacity，映射到 `Limits.Cpu`/`Limits.Memory`/`Limits.Pods`

#### Scenario: 返回 Available 可调度余量字段
- **WHEN** 查询到节点资源列表
- **THEN** 每个 `node_set` 元素包含 `available` 嵌套块（节点可再调度余量 = max(0, Allocatable - Requests)）

**Acceptance Criteria**:
- `available` 为 TypeList（MaxItems: 1）, Computed
- 嵌套字段 `cpu`/`memory`/`pods` 同 Capacity，映射到 `Available.Cpu`/`Available.Memory`/`Available.Pods`

### Requirement: 数据类型转换与空值处理
数据源必须正确处理 API 返回的数据类型转换，安全处理空值和 nil 指针。

**Rationale**: 腾讯云 SDK 使用指针类型，且 `MetaResource` 各字段可能返回 null（表示取不到有效值），需要安全解引用避免程序崩溃。

#### Scenario: 安全解引用指针字段
- **WHEN** API 返回包含指针类型的字段
- **THEN** 所有指针字段都经过 nil 检查后再设置

**Acceptance Criteria**:
- 在调用 set 方法前，先判断 Response 字段是否为 nil，若为 nil 则不调用 set 方法
- 每个 `MetaResource` 块（Capacity/Allocatable/Requests/Limits/Available）先判断是否为 nil，再判断其内部 `Cpu`/`Memory`/`Pods` 是否为 nil
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
- **WHEN** API 返回空结果（response 为 nil、Response 为 nil、NodeSet 为 nil）
- **THEN** 不直接 `d.SetId("")`，而是返回 `NonRetryableError`

**Acceptance Criteria**:
- 在 retry 块内检查 API 返回是否为空（`result == nil || result.Response == nil || result.Response.NodeSet == nil`）
- 若为空，返回 `resource.NonRetryableError`
- 在 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")`
- 遵循项目规则 #14 for RESOURCE_KIND_DATASOURCE

#### Scenario: API 返回错误时记录日志
- **WHEN** API 调用失败
- **THEN** 记录详细的错误日志

**Acceptance Criteria**:
- 使用 `log.Printf` 记录错误信息
- 日志中使用资源名称 `dbdc_db_custom_cluster_node_resources` 而非模糊措辞

### Requirement: 代码质量与规范
数据源代码必须符合项目规范，遵循命名和结构约定。

**Rationale**: 保持代码库一致性和可维护性。

#### Scenario: 遵循命名规范
- **WHEN** 审查数据源代码
- **THEN** 命名符合项目规范

**Acceptance Criteria**:
- 文件名: `data_source_tc_dbdc_db_custom_cluster_node_resources.go`
- 数据源名: `tencentcloud_dbdc_db_custom_cluster_node_resources`
- 函数名: `DataSourceTencentCloudDbdcDbCustomClusterNodeResources`
- 服务方法名: `DescribeDBCustomClusterNodeResourcesByFilter`
- 在日志/打印/错误信息中使用资源名称 `dbdc_db_custom_cluster_node_resources`
- 不在资源 go 文件开头添加注释

#### Scenario: 遵循文件组织规范
- **WHEN** 审查文件结构
- **THEN** 文件放置在正确的位置

**Acceptance Criteria**:
- 数据源文件: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_node_resources.go`
- 服务层方法: `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` (新增方法)
- 测试文件: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_node_resources_test.go`
- 文档文件: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_cluster_node_resources.md`
- 在 `tencentcloud/provider.go` 中注册数据源
- 在 `tencentcloud/provider.md` 中添加数据源条目

### Requirement: Service Layer 实现
必须在 `service_tencentcloud_dbdc.go` 中新增 `DescribeDBCustomClusterNodeResourcesByFilter` 方法。

**Rationale**: Service 层封装 API 调用，提供复用性。

#### Scenario: Service 方法处理参数转换
- **WHEN** 传入参数 map 包含 `ClusterId`、`NodeIds`
- **THEN** 正确转换为 SDK 请求参数

**Acceptance Criteria**:
- `ClusterId` 转换为 `*string` 并设置到 `request.ClusterId`
- `NodeIds` 转换为 `[]*string` 并设置到 `request.NodeIds`
- 该接口无分页参数，不进行分页循环

#### Scenario: Service 方法处理空响应
- **WHEN** API 返回空结果（response 为 nil 或 Response 为 nil 或 NodeSet 为 nil）
- **THEN** 返回 NonRetryableError

**Acceptance Criteria**:
- 在 retry 块内检查 `result == nil || result.Response == nil || result.Response.NodeSet == nil`
- 若为空，返回 `resource.NonRetryableError(fmt.Errorf("Describe dbdc_db_custom_cluster_node_resources failed, Response is nil."))`
- 在 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")`
- 成功时返回 `response.Response.NodeSet`

#### Scenario: Service 方法使用 ratelimit 与 retry
- **WHEN** 调用 `DescribeDBCustomClusterNodeResourcesByFilter`
- **THEN** 内部使用 `resource.Retry` + `ratelimit.Check` 包装 SDK 调用

**Acceptance Criteria**:
- 使用 `tccommon.ReadRetryTimeout` 作为超时
- 调用 `ratelimit.Check(request.GetAction())`
- 调用 `me.client.UseDbdcV20201029Client().DescribeDBCustomClusterNodeResources(request)`

### Requirement: Provider 注册
必须在 `tencentcloud/provider.go` 中注册数据源，并在 `tencentcloud/provider.md` 中添加数据源条目。

**Rationale**: 数据源必须在 Provider 中注册才能被 Terraform 使用。

#### Scenario: 数据源在 Terraform 中可用
- **WHEN** 用户在 Terraform 配置中使用 `data "tencentcloud_dbdc_db_custom_cluster_node_resources"`
- **THEN** Terraform 正常识别和使用该数据源

**Acceptance Criteria**:
- `provider.go` 的 DataSourcesMap 中添加 `"tencentcloud_dbdc_db_custom_cluster_node_resources": dbdc.DataSourceTencentCloudDbdcDbCustomClusterNodeResources()`
- `provider.md` 中 DBDC Data Source 部分添加 `tencentcloud_dbdc_db_custom_cluster_node_resources`

### Requirement: 文档
必须提供完整的 .md 文档文件。

**Rationale**: 每个数据源必须有文档说明使用方法。

#### Scenario: 用户查看文档并使用数据源
- **WHEN** 用户查看 `data_source_tc_dbdc_db_custom_cluster_node_resources.md`
- **THEN** 文档包含完整的使用说明

**Acceptance Criteria**:
- 一句话描述包含云产品名称（DBDC），格式: "Use this data source to query detailed information of DB Custom cluster node resources"
- Example Usage 部分，包含按 cluster_id 查询和按 node_ids 查询的示例
- 不包含 Argument Reference 和 Attribute Reference（这些由工具自动生成）

### Requirement: 单元测试
必须提供 gomonkey mock 方式的单元测试，验证数据源的 Read 功能。

**Rationale**: 项目规则要求新资源使用 mock 方式进行单元测试。

#### Scenario: 基本 Read 测试
- **WHEN** 运行 `TestDbdcDbCustomClusterNodeResourcesDS_ReadBasic`
- **THEN** 数据源正确读取并映射 API 返回的节点资源列表

**Acceptance Criteria**:
- 使用 gomonkey mock `DescribeDBCustomClusterNodeResources` API 返回
- 验证 `node_set` 包含正确的字段值，包括嵌套的 capacity/allocatable/requests/limits/available 块
- 验证 `d.Id()` 不为空
- 使用 `go test -gcflags="all=-l"` 参数运行

#### Scenario: Schema 结构验证测试
- **WHEN** 运行 `TestDbdcDbCustomClusterNodeResourcesDS_Schema`
- **THEN** 数据源 schema 结构符合规范

**Acceptance Criteria**:
- 验证所有必需的 schema 字段存在（cluster_id, node_ids, node_set, result_output_file）
- 验证字段类型和属性正确（cluster_id Required, node_ids Optional, node_set Computed, etc.）
- 验证 node_set 嵌套元素包含 node_id, capacity, allocatable, requests, limits, available 字段

#### Scenario: 空响应处理测试
- **WHEN** API 返回空结果
- **THEN** 数据源返回错误（NonRetryableError）

**Acceptance Criteria**:
- mock API 返回空/nil NodeSet
- 验证 Read 方法返回错误
