## Context

腾讯云 DLC（Data Lake Compute）的 DMS 元数据服务提供对数据库（Database）的管理能力。当前 terraform-provider-tencentcloud 已存在 DLC 服务目录（`tencentcloud/services/dlc/`）及 `service_tencentcloud_dlc.go` 服务层，已注册多个 DLC 资源/数据源（如 `tencentcloud_dlc_data_engine`、`tencentcloud_dlc_user` 等），但尚无 DMS 元数据库相关资源。

本次新增 `tencentcloud_dlc_dms_database` 资源（RESOURCE_KIND_GENERAL），对应云 SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125` 中的 4 个同步接口：

- `CreateDMSDatabase`（创建库）：入参 `Asset`(*Asset)、`SchemaName`、`Location`、`Name`、`DatasourceConnectionName`。
- `DescribeDMSDatabase`（查询库）：入参 `Name`、`SchemaName`、`Pattern`、`DatasourceConnectionName`；出参 `Name`、`SchemaName`、`Location`、`Asset`(*Asset)。
- `AlterDMSDatabase`（更新库）：入参 `CurrentName`、`SchemaName`、`Location`、`Asset`、`DatasourceConnectionName`。
- `DropDMSDatabase`（删除库）：入参 `Name`、`DeleteData`、`Cascade`、`DatasourceConnectionName`。

这 4 个接口均为同步接口，响应仅返回 `RequestId`（无 `TaskId`），因此 Create/Update/Delete 后直接调用 Read 验证即可，无需异步轮询，也无需在 schema 中声明 `Timeouts` 块。

`Asset` 为嵌套结构体，字段如下（均可能为 nil）：`Id`(*int64)、`Name`、`Guid`、`Catalog`、`Description`、`Owner`、`OwnerAccount`、`PermValues`([]*KVPair)、`Params`([]*KVPair)、`BizParams`([]*KVPair)、`DataVersion`(*int64)、`CreateTime`、`ModifiedTime`、`DatasourceId`(*int64)。其中 `KVPair` 仅含 `Key`、`Value` 两个字符串字段。

## Goals / Non-Goals

**Goals:**

- 实现 `tencentcloud_dlc_dms_database` 资源完整 CRUD，严格参考 `resource_tc_igtm_strategy.go` 代码风格。
- 正确处理 `Asset` 嵌套对象（含 3 个 `[]*KVPair` 列表字段）的序列化/反序列化。
- 通过复合 ID（`name#schema_name#datasource_connection_name`）在 Read/Update/Delete 中精确定位资源。
- 支持资源 Import。
- 在 `provider.go` / `provider.md` 完成注册，并生成资源示例 `.md`。

**Non-Goals:**

- 不实现 DMS 表（Table）、分区（Partition）等其它 DMS 对象的管理。
- 不暴露 `DescribeDMSDatabase` 的 `Pattern` 模糊匹配参数给用户（Read 仅做精确查询，`Pattern` 不纳入 schema）。
- 不处理 DLC 其它子产品（数据引擎、工作组等）的能力。
- 不在本次新增 `service_tencentcloud_dlc.go` 中额外的 service 层封装方法（DMS Database 的 API 调用直接在资源 CRUD 中完成，与 igtm_strategy 一致，避免冗余）。

## Decisions

### 决策 1：复合 ID 设计

**选择**：使用 `name#schema_name#datasource_connection_name`（分隔符为 `tccommon.FILED_SP` = `#`）作为资源 ID。

**理由**：
- `DescribeDMSDatabase` / `DropDMSDatabase` 均需通过 `Name` + `DatasourceConnectionName`（及 SchemaName 上下文）定位资源；`AlterDMSDatabase` 通过 `CurrentName` + `SchemaName` + `DatasourceConnectionName` 定位。
- 云 API 不返回单一主键 ID，因此必须用多字段联合 ID。
- 按项目规范使用 `tccommon.FILED_SP` 分隔，并在 Read/Update/Delete 中 `strings.Split(d.Id(), tccommon.FILED_SP)` 取出各部分。

**Import 说明**：资源支持 Import，import 时需使用联合 ID（`name#schema_name#datasource_connection_name`），在 `.md` 文档中明确说明。

### 决策 2：Schema 字段设计

顶层字段（参考云 API 入参与出参映射）：

| HCL key | SDK 字段 | 类型 | Required | ForceNew | 说明 |
|---|---|---|---|---|---|
| `name` | `Name` | TypeString | Yes | Yes | 数据库名称，创建后作为联合 id 一部分，不可变 |
| `schema_name` | `SchemaName` | TypeString | Yes | Yes | Schema 目录名，联合 id 一部分，不可变 |
| `datasource_connection_name` | `DatasourceConnectionName` | TypeString | Yes | Yes | 数据源连接名，联合 id 一部分，不可变 |
| `location` | `Location` | TypeString | Optional | No | Db 存储路径 |
| `delete_data` | `DeleteData` | TypeBool | Optional | No | 删除时是否删除数据（仅 Delete 使用，默认 false） |
| `cascade` | `Cascade` | TypeBool | Optional | No | 删除时是否级联删除（仅 Delete 使用，默认 false） |
| `asset` | `Asset` | TypeList, MaxItems=1 | Optional | No | 基础元数据对象（嵌套） |

`asset` 嵌套块字段（全部 Optional）：
- `id`（TypeInt，对应 `Asset.Id`，Computed，创建后由云返回）
- `name`（TypeString，对应 `Asset.Name`）
- `guid`（TypeString，对应 `Asset.Guid`，Computed）
- `catalog`（TypeString，对应 `Asset.Catalog`）
- `description`（TypeString，对应 `Asset.Description`）
- `owner`（TypeString，对应 `Asset.Owner`）
- `owner_account`（TypeString，对应 `Asset.OwnerAccount`）
- `perm_values`（TypeList，对应 `Asset.PermValues`，KVPair 列表）
- `params`（TypeList，对应 `Asset.Params`，KVPair 列表）
- `biz_params`（TypeList，对应 `Asset.BizParams`，KVPair 列表）
- `data_version`（TypeInt，对应 `Asset.DataVersion`，Computed）
- `create_time`（TypeString，对应 `Asset.CreateTime`，Computed）
- `modified_time`（TypeString，对应 `Asset.ModifiedTime`，Computed）
- `datasource_id`（TypeInt，对应 `Asset.DatasourceId`，Computed）

`perm_values` / `params` / `biz_params` 为 KVPair 列表，每项含 `key`（TypeString, Required）、`value`（TypeString, Optional）。

**ForceNew 策略**：`name`/`schema_name`/`datasource_connection_name` 为 ForceNew（联合 id 组成部分，改变需重建）；其余字段可更新。

### 决策 3：Update 实现策略

`AlterDMSDatabase` 入参为 `CurrentName`（变更前库名）、`SchemaName`、`Location`、`Asset`、`DatasourceConnectionName`。

- `name` 标记为 ForceNew，因此 Update 时库名不会改变，`CurrentName` 取值为当前 `d.Get("name")`。
- 当 `location` 或 `asset` 发生变更时，调用 `AlterDMSDatabase`。
- 按 RESOURCE_KIND_GENERAL 规范，资源支持 update，将 `name`/`schema_name`/`datasource_connection_name` 列入 `immutableArgs`，若发生变更返回 error（虽然它们已是 ForceNew，此处作为双保险）。
- 由于 `name` 为 ForceNew，正常情况下不会进入 update 的不可变检查分支。

### 决策 4：Read 实现

调用 `DescribeDMSDatabase`，使用 `name` + `schema_name` + `datasource_connection_name` 精确查询（不传 `Pattern`，避免模糊匹配返回多条）。

- 在 retry 块内检查 `response == nil || response.Response == nil`，返回 `NonRetryableError`。
- 响应字段 nil 安全：每个 `setXX` 前判断对应字段是否为 nil。
- 若查询返回为空（`response.Response.Name == nil`），先打印 `log.Printf("[CRUD] dlc_dms_database id=%s", d.Id())`，再 `d.SetId("")`。

### 决策 5：Delete 实现

调用 `DropDMSDatabase`，传入 `name`、`datasource_connection_name`，以及用户配置的 `delete_data`、`cascade`。

### 决策 6：Create 后处理

`CreateDMSDatabase` 响应仅含 `RequestId`，无返回 ID。Create 成功后：
1. 打印 `logId` 与 `d.Get("name")`。
2. 检查响应 `response == nil || response.Response == nil`，为空则返回 `NonRetryableError`。
3. 设置复合 ID：`d.SetId(name + tccommon.FILED_SP + schemaName + tccommon.FILED_SP + datasourceConnectionName)`。
4. 调用 Read 验证。

### 决策 7：测试策略

新增资源的单元测试使用 gomonkey mock 云 API（不使用 terraform 测试套件），仅测试业务逻辑。使用 `go test -gcflags=all=-l` 跑通涉及的测试文件。

## Risks / Trade-offs

- **[风险] Asset 嵌套结构与 KVPair 列表的序列化复杂** → 在 Read/Set 时对每个嵌套字段做 nil 判断，KVPair 列表用 `schema.TypeList` + Resource 元素；在 Create/Update 构建请求时仅设置非空字段。需仔细处理 `Params`/`BizParams`/`PermValues` 与 `Asset` 的对应关系。
- **[风险] `DescribeDMSDatabase` 可能因短暂波动返回空导致 state 中 id 被清空** → 严格按规范：retry 内不直接 `d.SetId("")`，而是返回 `NonRetryableError` 让外层继续重试；仅在确认资源确实不存在（重试耗尽后）时清空 id，并先打印 `[CRUD]` 日志保留现场。
- **[风险] 复合 ID 包含用户可控字符串，若值含 `#` 会导致 split 异常** → DLC 数据库名/schema 名/数据源连接名按云 API 约定不含 `#`，风险可控；如未来出现冲突再行评估转义方案。
- **[权衡] 不在 service 层新增封装方法** → DMS Database 调用逻辑直接放在资源 CRUD 中，与 igtm_strategy 风格一致，减少冗余代码；若后续有 datasource 复用再抽取。
- **[权衡] `delete_data`/`cascade` 暴露为顶层 Optional 字段** → 这两个字段仅 Delete 接口使用，但为让用户能在 HCL 中声明删除行为，置于顶层 schema；Read 接口不返回这些字段，故不设 Computed。
