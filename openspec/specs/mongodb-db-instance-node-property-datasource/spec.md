## Requirements

### Requirement: 数据源注册
数据源 `tencentcloud_mongodb_db_instance_node_property` SHALL 在 `tencentcloud/provider.go` 中注册，并在 `tencentcloud/provider.md` 中添加引用。

#### Scenario: 数据源可被 Terraform 识别
- **WHEN** 用户在 Terraform 配置中使用 `data "tencentcloud_mongodb_db_instance_node_property"` 块
- **THEN** Terraform 能够识别并调用该数据源

### Requirement: 查询 MongoDB 实例节点属性
数据源 SHALL 通过调用云 API `DescribeDBInstanceNodeProperty` 查询指定 MongoDB 实例的节点属性信息。

#### Scenario: 必填参数 instance_id 查询
- **WHEN** 用户提供有效的 `instance_id`
- **THEN** 数据源调用 `DescribeDBInstanceNodeProperty` 并返回该实例的 Mongos 节点列表和副本集节点列表

#### Scenario: 可选过滤参数 node_ids
- **WHEN** 用户提供 `node_ids` 参数
- **THEN** 数据源将 `node_ids` 传入云 API 请求，过滤返回指定节点的属性

#### Scenario: 可选过滤参数 roles
- **WHEN** 用户提供 `roles` 参数（如 `["PRIMARY", "SECONDARY"]`）
- **THEN** 数据源将 `roles` 传入云 API 请求，过滤返回指定角色的节点属性

#### Scenario: 可选过滤参数 only_hidden
- **WHEN** 用户提供 `only_hidden = true`
- **THEN** 数据源将 `only_hidden` 传入云 API 请求，只返回 Hidden 节点属性

#### Scenario: 可选过滤参数 priority
- **WHEN** 用户提供 `priority` 参数
- **THEN** 数据源将 `priority` 传入云 API 请求，过滤返回指定优先级的节点属性

#### Scenario: 可选过滤参数 votes
- **WHEN** 用户提供 `votes` 参数
- **THEN** 数据源将 `votes` 传入云 API 请求，过滤返回指定投票权的节点属性

#### Scenario: 可选过滤参数 tags
- **WHEN** 用户提供 `tags` 参数（包含 `tag_key` 和 `tag_value`）
- **THEN** 数据源将 `tags` 传入云 API 请求，过滤返回带有指定标签的节点属性

### Requirement: 返回 Mongos 节点属性
数据源 SHALL 将云 API 返回的 `Mongos` 字段映射到 `mongos` 属性，每个元素包含 NodeProperty 的完整字段。

#### Scenario: 返回 Mongos 节点列表
- **WHEN** 云 API 返回非空的 `Mongos` 列表
- **THEN** 数据源的 `mongos` 属性包含所有 Mongos 节点的属性，包括 `zone`、`node_name`、`address`、`wan_service_address`、`role`、`hidden`、`status`、`slave_delay`、`priority`、`votes`、`tags`、`replicate_set_id`

#### Scenario: Mongos 节点列表为空
- **WHEN** 云 API 返回空的 `Mongos` 列表或 nil
- **THEN** 数据源的 `mongos` 属性为空列表，不报错

### Requirement: 返回副本集节点信息
数据源 SHALL 将云 API 返回的 `ReplicateSets` 字段映射到 `replicate_sets` 属性，每个元素包含 `nodes` 字段（NodeProperty 列表）。

#### Scenario: 返回副本集节点列表
- **WHEN** 云 API 返回非空的 `ReplicateSets` 列表
- **THEN** 数据源的 `replicate_sets` 属性包含所有副本集信息，每个副本集的 `nodes` 字段包含该副本集中所有节点的属性

#### Scenario: 副本集节点列表为空
- **WHEN** 云 API 返回空的 `ReplicateSets` 列表或 nil
- **THEN** 数据源的 `replicate_sets` 属性为空列表，不报错

### Requirement: 错误处理与重试
数据源 SHALL 使用 `tccommon.ReadRetryTimeout` 作为超时时间，在 retry 块中调用云 API，失败时使用 `tccommon.RetryError()` 包装错误。

#### Scenario: 云 API 调用失败时重试
- **WHEN** 云 API 调用返回可重试错误
- **THEN** 数据源在超时时间内自动重试，超时后返回错误

#### Scenario: 云 API 返回空响应
- **WHEN** 云 API 返回 nil response 或 nil Response 字段
- **THEN** 数据源返回 NonRetryableError，不清空 id

### Requirement: 单元测试
数据源 SHALL 有对应的单元测试文件，使用 gomonkey mock 云 API，通过 `go test -gcflags=all=-l` 运行。

#### Scenario: 单元测试覆盖正常查询流程
- **WHEN** 运行单元测试
- **THEN** 测试覆盖正常查询、空结果、错误处理等场景，且测试通过

### Requirement: 文档
数据源 SHALL 有对应的 `.md` 文档文件，包含一句话描述、Example Usage 和必要说明。

#### Scenario: 文档包含示例配置
- **WHEN** 用户查看数据源文档
- **THEN** 文档包含完整的 Terraform 配置示例，展示如何使用该数据源查询 MongoDB 实例节点属性
