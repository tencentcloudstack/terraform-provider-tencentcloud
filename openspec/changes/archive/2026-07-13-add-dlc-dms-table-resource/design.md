## Context

Terraform Provider for TencentCloud 已具备 DLC（Data Lake Compute）服务的数据源与部分资源（如 `tencentcloud_dlc_data_engine`、`tencentcloud_dlc_user` 等），但尚未覆盖 DMS 元数据表的管理。DLC DMS 提供了 `CreateDMSTable`、`DescribeDMSTable`、`AlterDMSTable`、`DropDMSTable` 四个云 API，均已 vendored 在 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125` 中，支持对元数据表做完整 CRUD。

该表资源依赖多个复杂嵌套结构：`Asset`（基础对象，含 PermValues/Params/BizParams 等 KVPair 列表）、`DMSSds`（存储描述，含 SerdeParams/Params/Cols/SortColumns 等嵌套）、`DMSColumn`（列定义，含 Params/BizParams）、`DMSPartition`（分区，含 Params 与嵌套 Sds）、`KVPair`（键值对）。

云 API 的特殊性：
- `CreateDMSTable`/`AlterDMSTable`/`DropDMSTable` 的 Response 仅返回 RequestId，无业务 ID，资源需以业务键（db_name + name）作为复合 ID。
- `CreateDMSTable` 不接受 `schema_name`/`catalog`/`keyword`/`pattern`，这些只在 `DescribeDMSTable` 入参中存在；`DescribeDMSTable` 出参会返回 `schema_name`、`retention` 等创建时不可设置的只读字段。
- `AlterDMSTable` 额外需要 `CurrentName`/`CurrentDbName`（旧名称），用于改名场景；其余字段与 Create 对齐。
- `DropDMSTable` 需要 `delete_data`、`env_props`、`datasource_connection_name`，其中 `env_props` 为 `KVPair` 结构。

当前约束：本资源为 RESOURCE_KIND_GENERAL 通用资源，必须支持 import；复合 ID 用 `tccommon.FILED_SP` 分隔；调用云 API 需以 `tccommon.ReadRetryTimeout`/`tccommon.WriteRetryTimeout` 做最终一致性重试。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_dlc_dms_table` 资源，支持对 DLC DMS 元数据表的创建、读取、更新、删除全生命周期管理。
- 支持通过 `terraform import` 导入已存在的表（使用 `db_name` + `name` 复合 ID）。
- 正确映射云 API 的复杂嵌套结构（Asset、DMSSds、DMSColumn、DMSPartition、KVPair）到 Terraform schema。
- 在更新场景下正确处理 `AlterDMSTable` 的 `CurrentName`/`CurrentDbName`（取变更前的旧值）。
- 满足 provider 代码规范：Create 返回值空检查、Read 字段 nil 检查、retry 块仅放接口调用、使用资源蛇形命名做日志。
- 通过 gomonkey mock 云 API 的单元测试覆盖核心 CRUD 业务逻辑。

**Non-Goals:**
- 不实现 DMS 表的权限管理（PermValues 仅作为 Asset 属性透传，不做细粒度权限资源化）。
- 不暴露 `DescribeDMSTable` 专用的查询参数（`keyword`、`pattern`、`catalog`、`schema_name`）作为用户可配置的创建参数；其中 `schema_name` 作为只读 computed 字段从出参回填，`catalog`/`keyword`/`pattern` 不纳入 schema。
- 不实现数据源（datasource）形式的 `tencentcloud_dlc_dms_table`，本变更仅交付 resource。
- 不修改任何已有 DLC 资源的 schema。

## Decisions

### 决策 1：复合 ID 设计 = `db_name` + `name`
- **选择**：使用 `db_name` 与 `name` 通过 `tccommon.FILED_SP` 拼接作为资源 ID。
- **理由**：`CreateDMSTable` 无返回 ID；DMS 表在云端的唯一键即“数据库名 + 表名”。`DescribeDMSTable` 入参也正是 `DbName` + `Name`。二者组合可在 Read/Update/Delete 中拆分还原。
- **备选**：仅用 `name` 作 ID——不可行，因为同名表可存在于不同数据库，无法唯一定位。
- **影响**：Import 需用户提供 `db_name` 与 `name` 的复合 ID；文档需说明。

### 决策 2：Schema 字段分层与 ForceNew 策略
- **选择**：
  - `db_name`、`name`、`datasource_connection_name`：Required，且 `db_name`/`name` 标记 `ForceNew`（改库名/表名走 AlterDMSTable 的改名逻辑而非重建）。但由于 AlterDMSTable 支持改名（通过 CurrentName/CurrentDbName + 新 Name/DbName），故 `db_name`/`name` **不** ForceNew，而在 Update 时通过 d.HasChange 判断并把旧值映射到 `CurrentName`/`CurrentDbName`。
  - `delete_data`、`env_props`：仅 DropDMSTable 使用，作为 Optional 字段（`delete_data` 默认可设为 true 或由用户指定）。
  - `type`、`storage_size`、`record_count`、`life_time`、`data_update_time`、`struct_update_time`、`last_access_time`、`view_original_text`、`view_expanded_text`：Optional，可更新。
  - `asset`、`sds`、`columns`、`partition_keys`、`partitions`：Optional 复杂嵌套结构，可更新。
  - `schema_name`、`retention`：Computed only（从 DescribeDMSTable 出参回填，不可由用户设置）。
- **理由**：严格对照云 API Create/Alter 入参字段，确保 Create 入参字段均在 CreateDMSTable 存在、Update 入参字段均在 AlterDMSTable 存在。`schema_name`/`retention` 仅在出参出现，设为 Computed。

### 决策 3：Update 改名通过 CurrentName/CurrentDbName
- **选择**：在 `resourceTencentCloudDlcDmsTableUpdate` 中，先构造 `AlterDMSTableRequest`，对 `db_name`/`name` 使用 `d.GetOk` 取新值填入 `DbName`/`Name`，并使用 `d.GetChange("db_name")`/`d.GetChange("name")` 取旧值填入 `CurrentDbName`/`CurrentName`；其余可变字段从 schema 读取填入。无 immutableArgs 限制（因为 AlterDMSTable 支持所有 Create 字段）。
- **理由**：AlterDMSTable 设计即用于改名与字段更新，符合云 API 语义。
- **备选**：把改名视为 ForceNew 重建——会产生删表再建表的副作用，破坏数据，不可接受。

### 决策 4：嵌套结构 schema 扁平化与 KVPair 复用
- **选择**：对 `KVPair` 定义为可复用的 list-of-map（`key`/`value` 两个字符串字段），在 Asset.PermValues/Params/BizParams、DMSSds.SerdeParams/Params、DMSColumn.Params/BizParams、DMSPartition.Params、DropDMSTable.EnvProps 等位置统一使用该结构。`DMSColumnOrder`（SortCols/SortColumns）定义为含 `col`/`order` 的对象。
- **理由**：保持与 SDK 结构一致，便于双向映射；KVPair 是云 API 通用键值结构。

### 决策 5：Retry 与空值检查遵循 provider 规范
- **选择**：
  - Create/Update/Delete 使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)`，retry 块内仅调用接口，错误用 `tccommon.RetryError(e)` 包装；Create 成功后在 retry 外做 `response == nil || response.Response == nil` 空检查（返回 NonRetryableError），但因 Create Response 无业务字段，仅做 Response 非空校验。
  - Read 使用 `service.DescribeDmsTableById`（内部 `resource.Retry(tccommon.ReadRetryTimeout, ...)`）。在 retry 块内检查 `response == nil || response.Response == nil`，返回 NonRetryableError（不直接 SetId("")），让外层重试耗尽后失败；外层失败路径打印 `log.Printf("[DATASOURCE] read empty, skip SetId")`。注意：本资源为 RESOURCE_KIND_GENERAL 非 DATASOURCE，但 Read 为空时仍遵循“先打印含 id 的日志再 SetId("")”规范。
  - Read 中 set 字段前逐一判断 Response 字段 nil。
- **理由**：符合项目硬约束与代码生成要求第 8/9/14 条。

### 决策 6：单元测试使用 gomonkey mock
- **选择**：新增资源为全新资源，按规范使用 gomonkey mock 云 API（`UseDlcV20210125Client` 及其 `CreateDMSTableWithContext`/`DescribeDMSTableWithContext`/`AlterDMSTableWithContext`/`DropDMSTableWithContext`），仅测业务逻辑，用 `go test -gcflags=all=-l` 跑通。
- **理由**：规范要求新增 terraform 资源测试用 mock，不使用 TF 测试套件。

## Risks / Trade-offs

- [风险] `AlterDMSTable` 改名时若 `CurrentName`/`CurrentDbName` 传错会导致更新失败或操作到错误的表 → 缓解：Update 中严格使用 `d.GetChange` 获取旧值，并在日志中打印 current/new 便于排障。
- [风险] `Asset` 结构层级深、字段多，schema 映射易遗漏或与 SDK 类型不一致 → 缓解：严格按 vendor models.go 逐字段核对类型（int64/string/[]*KVPair/[]*DMSColumn 等），Read 时对每个字段做 nil 判断。
- [风险] `DropDMSTable` 的 `delete_data`（bool）与 `env_props`（KVPair）若未在 schema 暴露，删除时无法控制是否清除底层数据 → 缓解：在 schema 中暴露 `delete_data`（Optional，默认 false 与云 API 一致）与 `env_props`（Optional KVPair 列表）。
- [风险] `DescribeDMSTable` 返回 `schema_name`/`retention` 等只读字段，若误设为可写会导致 update 时云 API 拒绝 → 缓解：这两个字段设为 Computed，不参与 AlterDMSTable 入参构造。
- [权衡] 复杂嵌套结构导致 schema 较大、转换代码冗长 → 接受，以保证功能完整；不抽取通用 helper 以避免影响既有代码。
- [风险] Create Response 无 ID，无法通过返回值判断创建具体实例 → 缓解：Create 成功后用用户输入的 `db_name`+`name` 复合 ID 设置 d.Id()，并立即调用 Read 回查确认表已存在；若 Read 查不到则返回错误。

## Migration Plan

- 纯新增资源，无数据迁移。部署只需合并代码并在 provider 注册新资源。
- 回滚策略：移除 provider.go 中的注册项与资源文件即可，不影响存量 state（该资源类型此前不存在）。
