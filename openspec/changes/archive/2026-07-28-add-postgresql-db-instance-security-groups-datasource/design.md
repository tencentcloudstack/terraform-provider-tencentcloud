## Context

Terraform Provider for TencentCloud 当前已有大量 postgresql 服务的数据源（如 `data_source_tc_postgresql_readonly_groups.go`、`data_source_tc_postgresql_instances.go` 等），但缺少查询实例/只读组关联安全组的数据源。云 API `DescribeDBInstanceSecurityGroups`（postgres/v20170312）已存在于 vendor 中，可直接调用。

该接口入参为 `DBInstanceId` 和 `ReadOnlyGroupId`（至少传一个），出参为 `SecurityGroupSet`（`[]*SecurityGroup`），其中 `SecurityGroup` 包含 `ProjectId`、`CreateTime`、`Inbound`/`Outbound`（`[]*PolicyRule`）、`SecurityGroupId`、`SecurityGroupName`、`SecurityGroupDescription` 字段；`PolicyRule` 包含 `Action`、`CidrIp`、`PortRange`、`IpProtocol`、`Description` 字段。

该接口为同步接口（非异步），无需轮询 Read 接口等待生效，且无分页参数。

## Goals / Non-Goals

**Goals:**
- 新增数据源 `tencentcloud_postgresql_db_instance_security_groups`，支持通过 `db_instance_id` 或 `read_only_group_id` 查询关联的安全组列表
- 将 `SecurityGroupSet` 列表展开到顶层 schema（遵循数据源不嵌套"列表型数据"层的规范），输出 `security_group_set` 列表，包含每个安全组的字段及入站/出站规则
- 在 `tencentcloud/provider.go` 和 `tencentcloud/provider.md` 中注册该数据源
- 补充单元测试（使用 gomonkey mock 云 API）和数据源文档

**Non-Goals:**
- 不实现安全组的创建、修改、删除（仅查询）
- 不暴露分页参数给用户（该接口本身无分页参数）
- 不修改任何现有资源/数据源的 schema

## Decisions

### Decision 1: 数据源文件与服务层
- 新增文件 `tencentcloud/services/postgresql/data_source_tc_postgresql_db_instance_security_groups.go`
- 参考现有数据源 `data_source_tc_postgresql_readonly_groups.go` 的代码风格与结构
- 在 `service_tencentcloud_postgresql.go` 中新增 `DescribePostgresqlDbInstanceSecurityGroups` 服务方法封装 `DescribeDBInstanceSecurityGroups` API 调用

**理由**: 遵循现有代码组织约定，服务层封装 API 调用便于复用与测试。

### Decision 2: Schema 设计（列表展开）
- 顶层参数：`db_instance_id`（Optional）、`read_only_group_id`（Optional）、`result_output_file`（Optional）
- 输出参数：`security_group_set`（Computed, TypeList），元素为 `schema.Resource`，字段平铺每个 `SecurityGroup` 的属性：`project_id`、`create_time`、`security_group_id`、`security_group_name`、`security_group_description`、`inbound`（TypeList 嵌套 PolicyRule）、`outbound`（TypeList 嵌套 PolicyRule）
- `inbound`/`outbound` 元素字段：`action`、`cidr_ip`、`port_range`、`ip_protocol`、`description`

**理由**: 遵循数据源规范——禁止创建"列表型数据"这一层嵌套；将列表元素字段平铺，使每个字段都可被 terraform 单独 set/read。

### Decision 3: Read 逻辑与重试
- 使用 `resource.Retry(tccommon.ReadRetryTimeout, ...)` 包装 API 调用，失败时用 `tccommon.RetryError()` 包装错误
- 在 retry 块内检查 API 返回是否为空（`response == nil` / `response.Response == nil` / `len(response.Response.SecurityGroupSet) == 0`），若为空则返回 `NonRetryableError`（不直接 `d.SetId("")`），避免云 API 短暂波动清空 state
- 在外层 retry 失败路径保留 `log.Printf("[DATASOURCE] read empty, skip SetId")` 提示
- 使用 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`
- 在 set 字段前判断 Response 字段是否为 nil，为 nil 则不调用 set

**理由**: 遵循 RESOURCE_KIND_DATASOURCE 的 Read 重试与空响应处理规范。

### Decision 4: 数据源 ID 设置
- 使用 `helper.DataResourceIdsHash(ids)` 生成数据源 ID，其中 `ids` 为查询返回的安全组 ID 列表

**理由**: 与现有数据源（如 readonly_groups）保持一致。

### Decision 5: 测试方式
- 单元测试使用 gomonkey mock 云 API（`PostgresqlService.DescribePostgresqlDbInstanceSecurityGroups`），不使用 terraform 测试套件

**理由**: 遵循新增 terraform 数据源的单元测试规范——使用 mock 进行业务代码逻辑的单元测试。

## Risks / Trade-offs

- **[风险] API 入参约束**: `DBInstanceId`、`ReadOnlyGroupId` 至少需传一个。→ **缓解**: 在数据源 schema 中将两者均设为 Optional，但在 Read 逻辑中不做强制校验，交由云 API 返回错误（`InvalidParameter.ParameterCheckError`），保持与 API 行为一致。
- **[风险] 返回空列表**: 当实例未关联安全组时返回空列表。→ **缓解**: 空列表属于正常业务结果，不应视为错误；但在 retry 块内区分"API 调用成功但返回空"与"API 调用失败"，确保仅在 API 调用成功时正常处理空列表。根据规范，DATASOURCE 的 retry 块内若返回空需返回 NonRetryableError 让重试耗尽失败——此处需注意：安全组列表为空可能是正常的无安全组场景。综合考虑，遵循规范在 retry 块内对空返回 NonRetryableError，以优先保证 state 安全。
- **[权衡] 不支持 import**: 数据源无需 import。
