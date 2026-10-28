## Context

TencentCloud Terraform Provider 已有多个 MongoDB 数据源（如 `tencentcloud_mongodb_instances`、`tencentcloud_mongodb_instance_connections` 等），但缺少查询 MongoDB 实例节点属性的数据源。云 API `DescribeDBInstanceNodeProperty` 支持查询 MongoDB 实例的 Mongos 节点和副本集节点的详细属性，包括节点名称、地址、角色、状态、优先级、投票权等信息。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_mongodb_db_instance_node_property` 数据源，封装 `DescribeDBInstanceNodeProperty` 云 API
- 支持通过实例 ID、节点 ID、节点角色、Hidden 节点标志、优先级、投票权、节点标签等参数过滤查询
- 返回 Mongos 节点属性列表（`mongos`）和副本集节点信息列表（`replicate_sets`）
- 在 `provider.go` 和 `provider.md` 中注册新数据源
- 提供对应的 `.md` 文档和单元测试

**Non-Goals:**
- 不支持对节点属性的写操作（增删改）
- 不支持分页参数暴露给用户（该 API 无分页参数）

## Decisions

### 决策1：数据源文件组织

遵循现有 MongoDB 数据源的文件组织规范：
- 数据源实现：`tencentcloud/services/mongodb/data_source_tc_mongodb_db_instance_node_property.go`
- 测试文件：`tencentcloud/services/mongodb/data_source_tc_mongodb_db_instance_node_property_test.go`
- 文档文件：`tencentcloud/services/mongodb/data_source_tc_mongodb_db_instance_node_property.md`

### 决策2：Schema 设计

入参 schema（Optional 过滤条件）：
- `instance_id`（Required, String）：实例 ID，必填
- `node_ids`（Optional, List of String）：节点 ID 列表
- `roles`（Optional, List of String）：节点角色列表（PRIMARY/SECONDARY/READONLY/ARBITER）
- `only_hidden`（Optional, Bool）：是否只查询 Hidden 节点
- `priority`（Optional, Int）：节点优先级
- `votes`（Optional, Int）：节点投票权
- `tags`（Optional, List of Object）：节点标签，包含 `tag_key` 和 `tag_value`

出参 schema（Computed）：
- `mongos`（Computed, List of Object）：Mongos 节点属性列表，每个元素包含 NodeProperty 的所有字段
- `replicate_sets`（Computed, List of Object）：副本集节点信息列表，每个元素包含 `nodes` 字段（NodeProperty 列表）

NodeProperty 字段：`zone`、`node_name`、`address`、`wan_service_address`、`role`、`hidden`、`status`、`slave_delay`、`priority`、`votes`、`tags`、`replicate_set_id`

### 决策3：错误处理与重试

遵循 Provider 规范，使用 `tccommon.ReadRetryTimeout` 作为超时时间，在 retry 块中调用云 API，失败时使用 `tccommon.RetryError()` 包装错误。

### 决策4：测试方式

使用 gomonkey mock 云 API 进行单元测试，不依赖真实云环境，使用 `go test -gcflags=all=-l` 运行。

## Risks / Trade-offs

- [风险] 云 API 返回的 `Mongos` 或 `ReplicateSets` 可能为 nil → 在 Read 方法中判断 nil 后再设置字段，避免 panic
- [风险] `NodeProperty.Tags` 字段可能为 nil → 在转换时判断 nil，返回空列表

## Migration Plan

无需迁移，纯新增数据源，不影响现有资源和配置。

## Open Questions

无。
