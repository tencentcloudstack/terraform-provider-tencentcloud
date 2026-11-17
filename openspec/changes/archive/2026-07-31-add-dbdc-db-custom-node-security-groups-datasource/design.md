## Context

DBDC (Database Dedicated Cluster) 产品已支持多个 Terraform 数据源（`tencentcloud_dbdc_db_custom_clusters`、`tencentcloud_dbdc_db_custom_nodes`、`tencentcloud_dbdc_db_custom_cluster_nodes`、`tencentcloud_dbdc_db_custom_images`），但缺少节点安全组查询数据源。云 API `DescribeDBCustomNodeSecurityGroups` 已在 SDK 中可用，可以查询指定节点的安全组绑定信息及其入/出站规则。

当前代码库中 DBDC 服务代码位于 `tencentcloud/services/dbdc/`，使用 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029` SDK 包。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_dbdc_db_custom_node_security_groups` 数据源，通过 `node_id` 查询指定节点的安全组列表
- 将 `SecurityGroup` 和 `PolicyRule` 结构体字段完整映射到 Terraform schema
- 遵循现有 DBDC 数据源代码风格和模式
- 在 `provider.go` 中注册新数据源

**Non-Goals:**
- 不修改现有 DBDC 资源或数据源的 schema
- 不实现安全组的创建/修改/删除操作（仅为只读数据源）
- 不添加分页逻辑（该 API 是单节点查询，不涉及分页）

## Decisions

### 1. 数据源模式：单节点查询（非列表查询）

`DescribeDBCustomNodeSecurityGroups` API 接受单个 `NodeId`，返回 `[]*SecurityGroup`。与 `DescribeDBCustomClusterDetail` 类似，这是按 ID 查询的单条数据源，而非分页列表查询。

**Schema 设计**：
- `node_id`：Required，TypeString，用户指定要查询的节点 ID
- `groups`：Computed，TypeList，包含安全组详情列表，每个元素为嵌套 Resource

### 2. Schema 字段展开策略

按照规范要求，不创建"资源列表型数据"这一层嵌套 schema，而是将 `SecurityGroup` 和 `PolicyRule` 的字段平铺展开：

- `groups` 列表的每个元素包含安全组的所有字段（`security_group_id`、`security_group_name`、`security_group_remark`、`project_id`、`create_time`、`inbound`、`outbound`）
- `inbound` 和 `outbound` 分别是 `PolicyRule` 列表，每个规则包含 `action`、`cidr_ip`、`port_range`、`ip_protocol`、`service_module`、`address_module`、`id`

### 3. 服务层方法

在 `service_tencentcloud_dbdc.go` 中新增 `DescribeDBCustomNodeSecurityGroupsById` 方法：
- 直接调用 `DescribeDBCustomNodeSecurityGroups` API
- 使用 `tccommon.ReadRetryTimeout` 和 retry 机制
- 空响应检查：`result == nil || result.Response == nil` 时返回 `NonRetryableError`
- 在 retry 外检查返回值，失败时打印 `[DATASOURCE]` 日志

### 4. 未使用 `DescribeDBCustomNodeSecurityGroupsByFilter` 模式

该 API 不支持分页和过滤参数，只有一个 `NodeId` 入参。因此不需要实现 `ByFilter` 模式，直接实现 `ById` 模式即可。

## Risks / Trade-offs

- **API 返回空**：当节点不存在或节点未绑定安全组时，API 可能返回空的 `Groups` 列表。按照规范，在 retry 块内检查空响应并返回 `NonRetryableError`，在外层失败路径保留 `log.Printf` 提示。
- **无安全组场景**：如果节点未绑定任何安全组，`Groups` 可能为空列表，此时 `groups` 字段应设置为空列表 `[]`，而非 `nil`。