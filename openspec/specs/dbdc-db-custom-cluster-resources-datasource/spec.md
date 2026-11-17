## ADDED Requirements

### Requirement: 查询 DB Custom 集群资源信息
数据源 `tencentcloud_dbdc_db_custom_cluster_resources` 必须能够通过 `cluster_id` 查询指定 DB Custom 集群的资源汇总信息，调用 `DescribeDBCustomClusterResources` API 并返回集群的资源容量、可分配量、申请量、上限和可用余量。

**Rationale**: 用户需要在 Terraform 配置中动态获取集群资源使用情况，用于容量规划和调度决策。

#### Scenario: 按集群 ID 查询资源信息
- **WHEN** 用户指定了 `cluster_id` 参数
- **THEN** 数据源调用 `DescribeDBCustomClusterResources` API（入参 `ClusterId`）获取集群资源汇总
- **AND** 返回 `node_count`、`capacity`、`allocatable`、`requests`、`limits`、`available` 等字段

#### Scenario: 集群 ID 为必选参数
- **WHEN** 用户未指定 `cluster_id` 参数
- **THEN** Terraform 在 plan 阶段报错，提示 `cluster_id` 为必选字段

### Requirement: 返回集群资源汇总字段
数据源必须返回集群资源汇总的完整详细信息，包括节点总数和各类资源（CPU、内存、Pods）的容量、可分配量、申请量、上限和可用余量。

**Rationale**: 用户需要完整的资源汇总信息用于容量规划和管理决策。

#### Scenario: 返回节点总数
- **WHEN** 查询到集群资源信息
- **THEN** 返回 `node_count` 字段（TypeInt），表示参与汇总的工作节点总数（不含控制面节点）

#### Scenario: 返回资源物理总容量
- **WHEN** 查询到集群资源信息
- **THEN** 返回 `capacity` 嵌套结构，包含 `cpu`（TypeFloat，单位核）、`memory`（TypeFloat，单位 GiB）、`pods`（TypeInt，单位个）

#### Scenario: 返回可分配容量
- **WHEN** 查询到集群资源信息
- **THEN** 返回 `allocatable` 嵌套结构（= Capacity - 系统预留），包含 `cpu`、`memory`、`pods`

#### Scenario: 返回资源申请量
- **WHEN** 查询到集群资源信息
- **THEN** 返回 `requests` 嵌套结构，表示集群所有非终态 Pod 的 requests 申请量之和（含系统 Pod），包含 `cpu`、`memory`、`pods`

#### Scenario: 返回资源上限
- **WHEN** 查询到集群资源信息
- **THEN** 返回 `limits` 嵌套结构，表示集群所有非终态 Pod 的 limits 上限之和（含系统 Pod），包含 `cpu`、`memory`、`pods`

#### Scenario: 返回可用余量
- **WHEN** 查询到集群资源信息
- **THEN** 返回 `available` 嵌套结构，表示集群可再调度余量（所有节点 max(0, Allocatable - Requests) 累加求和），包含 `cpu`、`memory`、`pods`

### Requirement: 空值安全处理
数据源必须正确处理 API 返回的指针类型和可能的 nil 值，避免程序崩溃。

**Rationale**: 腾讯云 SDK 使用指针类型，`MetaResource` 嵌套对象及其字段可能为 nil，需要安全解引用。

#### Scenario: 安全处理 nil 嵌套对象
- **WHEN** API 返回的 `Capacity`/`Allocatable`/`Requests`/`Limits`/`Available` 为 nil
- **THEN** 跳过对应嵌套字段的 set，不报错

#### Scenario: 安全处理 nil 指针字段
- **WHEN** `MetaResource` 中的 `Cpu`/`Memory`/`Pods` 为 nil
- **THEN** 跳过对应字段的 set，不报错

### Requirement: 错误处理与重试
数据源必须正确处理 API 错误，实现重试逻辑以应对临时性故障，并在响应为空时返回不可重试错误。

**Rationale**: 云 API 调用可能因网络、限流等原因失败或返回空响应，需要重试机制和正确的空响应处理。

#### Scenario: API 调用失败时重试
- **WHEN** API 调用返回可重试错误
- **THEN** 使用 `resource.Retry`（`tccommon.ReadRetryTimeout`）包装调用，自动重试

#### Scenario: API 返回空响应时返回不可重试错误
- **WHEN** API 返回 `response == nil` 或 `response.Response == nil`
- **THEN** 不执行 `d.SetId("")`，而是返回 `resource.NonRetryableError`
- **AND** 外层 retry 失败路径打印 `[DATASOURCE] read empty, skip SetId` 日志

### Requirement: 输出文件支持
数据源支持将查询结果输出到 JSON 文件，方便用户审查和分析。

**Rationale**: 用户可能需要将查询结果导出用于离线分析或审计。

#### Scenario: 输出结果到指定文件
- **WHEN** 用户指定了 `result_output_file` 参数
- **THEN** 将结果序列化为 JSON 并写入文件

#### Scenario: 未指定输出文件
- **WHEN** 用户未指定 `result_output_file` 参数
- **THEN** 跳过文件写入，仅设置 state

### Requirement: Provider 注册
数据源必须在 `provider.go` 中注册，并在 `provider.md` 中声明。

**Rationale**: 数据源需要通过 provider 暴露给 Terraform 用户使用。

#### Scenario: 数据源在 provider 中注册
- **WHEN** 用户在 Terraform 配置中声明 `data "tencentcloud_dbdc_db_custom_cluster_resources"`
- **THEN** provider 能正确识别并处理该数据源

#### Scenario: 数据源在 provider.md 中声明
- **WHEN** 查看 `provider.md` 文档
- **THEN** 在 DB Dedicated Cloud (DBDC) 部分的 Data Source 节中包含 `tencentcloud_dbdc_db_custom_cluster_resources`
