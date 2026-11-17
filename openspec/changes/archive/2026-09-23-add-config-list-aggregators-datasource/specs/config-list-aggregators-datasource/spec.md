## ADDED Requirements

### Requirement: 查询 Config 账号组列表数据源
系统 SHALL 提供数据源 `tencentcloud_config_list_aggregators`，调用云 API `ListAggregators`（config v20220802）查询账号组列表，并返回账号组总数与账号组详情。

#### Scenario: 查询全部账号组
- **WHEN** 用户使用 `tencentcloud_config_list_aggregators` 数据源且不传入任何过滤参数
- **THEN** 系统内部自动分页拉取全部账号组，并在 state 中填充 `total` 与 `items` 列表

#### Scenario: 账号组列表字段完整返回
- **WHEN** 云 API 返回账号组列表
- **THEN** 每个 `items` 元素 SHALL 包含字段：`name`、`description`、`owner_uin`、`create_time`、`account_count`、`type`、`account_group_id`、`aggregator_status`、`member_name`，仅当对应云 API 字段非 nil 时 set

#### Scenario: 云 API 返回空时不清空数据源
- **WHEN** 数据源 Read 调用 `ListAggregators` 返回空（`items` 为 nil 且 `total` 为 0）
- **THEN** 系统 SHALL 在 retry 块内返回 `NonRetryableError`，不得直接 `d.SetId("")`，并保留 `[DATASOURCE] read empty, skip SetId` 日志

### Requirement: 数据源注册与文档
系统 SHALL 在 `provider.go` 中注册数据源 `tencentcloud_config_list_aggregators`，并通过 `make doc` 生成对应文档。

#### Scenario: provider 注册
- **WHEN** 实现完成
- **THEN** `provider.go` 数据源 map 中 SHALL 包含键 `tencentcloud_config_list_aggregators` 指向 `config.DataSourceTencentCloudConfigListAggregators()`

### Requirement: 数据源单元测试
系统 SHALL 提供使用 gomonkey mock 云 API 的单元测试，覆盖数据源 Read 的 flatten 逻辑。

#### Scenario: mock 返回账号组列表
- **WHEN** 单元测试 mock `ListAggregatorsWithContext` 返回构造的 `Aggregator` 列表
- **THEN** 测试 SHALL 验证 `items` 与 `total` 被正确 set 到 ResourceData