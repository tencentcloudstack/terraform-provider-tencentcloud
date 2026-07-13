## ADDED Requirements

### Requirement: 资源注册与命名
系统 SHALL 在 provider 中注册名为 `tencentcloud_dlc_dms_table` 的 Terraform 资源，资源类型为 RESOURCE_KIND_GENERAL，实现文件为 `tencentcloud/services/dlc/resource_tc_tencentcloud_dlc_dms_table.go`，并在 `provider.go` 与 `provider.md` 中完成注册登记。

#### Scenario: provider 可识别新资源
- **WHEN** 执行 `terraform providers schema` 或使用 `tencentcloud_dlc_dms_table` 资源块
- **THEN** provider 能识别并接受该资源类型，无 unknown resource 错误

### Requirement: 资源 schema 定义
系统 SHALL 为 `tencentcloud_dlc_dms_table` 定义 schema，覆盖 CreateDMSTable/AlterDMSTable 入参与 DescribeDMSTable 出参的字段，包括：`asset`、`type`、`db_name`、`storage_size`、`record_count`、`life_time`、`data_update_time`、`struct_update_time`、`last_access_time`、`sds`、`columns`、`partition_keys`、`view_original_text`、`view_expanded_text`、`partitions`、`name`、`datasource_connection_name`、`delete_data`、`env_props`，以及 Computed 只读字段 `schema_name`、`retention`。嵌套结构 `asset`、`sds`、`columns`、`partition_keys`、`partitions`、`env_props` 及其子结构 SHALL 完整映射云 API 的 `Asset`、`DMSSds`、`DMSColumn`、`DMSPartition`、`KVPair` 数据模型。

#### Scenario: 必填字段校验
- **WHEN** 用户在 Terraform 配置中省略 `db_name` 或 `name`
- **THEN** terraform plan 阶段报错，提示必填字段缺失

#### Scenario: 只读字段不可写
- **WHEN** 用户尝试在配置中设置 `schema_name` 或 `retention`
- **THEN** terraform plan 阶段报错，提示该字段为 Computed 不可设置

### Requirement: 创建 DMS 表
系统 SHALL 在资源 Create 时调用 DLC `CreateDMSTable` 接口，将 schema 中 CreateDMSTable 支持的字段映射到请求入参，调用成功后以 `db_name` 与 `name` 通过 `tccommon.FILED_SP` 拼接作为复合 ID 写入 state，并立即调用 Read 回查确认表已存在。

#### Scenario: 创建成功
- **WHEN** 用户提供合法的 `db_name`、`name` 及其它 CreateDMSTable 入参字段
- **THEN** 调用 `CreateDMSTable` 成功，资源 ID 被设置为 `db_name<FILED_SP>name`，state 被回填

#### Scenario: 创建接口返回空
- **WHEN** `CreateDMSTable` 返回 `Response` 为 nil
- **THEN** 返回 NonRetryableError，资源 ID 不被写入，避免空 ID 导致状态混乱

### Requirement: 读取 DMS 表
系统 SHALL 在资源 Read 时从复合 ID 拆分出 `db_name` 与 `name`，调用 DLC `DescribeDMSTable` 接口，并将 `DescribeDMSTableResponse` 中非空的字段回填到 state。当云端返回为空（response/Response 为 nil）时，系统 SHALL 打印包含资源 id 的日志后执行 `d.SetId("")`。

#### Scenario: 读取成功回填字段
- **WHEN** 资源存在且 `DescribeDMSTable` 正常返回
- **THEN** 将 `asset`、`type`、`db_name`、`schema_name`、`storage_size`、`record_count`、`life_time`、`data_update_time`、`struct_update_time`、`last_access_time`、`sds`、`columns`、`partition_keys`、`partitions`、`view_original_text`、`view_expanded_text`、`retention`、`name` 等非空字段回填 state

#### Scenario: 资源不存在
- **WHEN** `DescribeDMSTable` 返回为空或表已被删除
- **THEN** 打印 `[CRUD] dlc_dms_table id=<id>` 日志后执行 `d.SetId("")`，不返回错误

#### Scenario: 读取字段 nil 跳过
- **WHEN** Response 中某字段为 nil
- **THEN** 不对该字段调用 set，避免写入空值覆盖

### Requirement: 更新 DMS 表
系统 SHALL 在资源 Update 时调用 DLC `AlterDMSTable` 接口。当 `db_name` 或 `name` 发生变更时，系统 SHALL 将变更前的旧值分别映射到 `CurrentDbName` 与 `CurrentName`，将变更后的新值映射到 `DbName` 与 `Name`，实现改名。其余 AlterDMSTable 支持的字段 SHALL 从 schema 读取填入请求。`schema_name`、`retention` 等只在出参出现的字段 SHALL NOT 被加入 AlterDMSTable 请求。

#### Scenario: 更新普通字段
- **WHEN** 用户修改 `type`、`columns`、`sds` 等 AlterDMSTable 支持的字段
- **THEN** 调用 `AlterDMSTable` 成功，state 更新为新值

#### Scenario: 改表名
- **WHEN** 用户修改 `name` 字段
- **THEN** `AlterDMSTable` 请求中 `CurrentName` 为旧表名、`Name` 为新表名，调用成功后复合 ID 更新为新 `db_name<FILED_SP>新name`

#### Scenario: 改库名
- **WHEN** 用户修改 `db_name` 字段
- **THEN** `AlterDMSTable` 请求中 `CurrentDbName` 为旧库名、`DbName` 为新库名，调用成功后复合 ID 更新

### Requirement: 删除 DMS 表
系统 SHALL 在资源 Delete 时从复合 ID 拆分 `db_name` 与 `name`，调用 DLC `DropDMSTable` 接口，将 `db_name`、`name`、`delete_data`、`env_props`、`datasource_connection_name` 映射到请求入参。

#### Scenario: 删除成功
- **WHEN** 用户执行 `terraform destroy` 或删除该资源
- **THEN** 调用 `DropDMSTable` 成功，资源从 state 移除

### Requirement: 导入支持
系统 SHALL 支持 `terraform import`，导入时使用 `db_name` 与 `name` 通过 `tccommon.FILED_SP` 拼接的复合 ID。资源文档 SHALL 在 Import 部分说明需使用复合 ID。

#### Scenario: 导入已存在表
- **WHEN** 用户执行 `terraform import tencentcloud_dlc_dms_table.xxx <db_name>#<name>`
- **THEN** 资源被导入 state，后续 plan 显示无变更

#### Scenario: 导入 ID 格式错误
- **WHEN** 导入 ID 不包含 `FILED_SP` 分隔的两个部分
- **THEN** 返回 "id is broken" 错误

### Requirement: 重试与错误处理
系统 SHALL 在调用 Create/Update/Delete 接口时使用 `tccommon.WriteRetryTimeout` 做 retry，在调用 Read 接口时使用 `tccommon.ReadRetryTimeout` 做 retry。retry 块内 SHALL 仅包含接口调用，接口错误 SHALL 通过 `tccommon.RetryError()` 包装。设置 ID、回填 state 等成功操作 SHALL 在 retry 块外执行。

#### Scenario: 接口瞬时失败重试
- **WHEN** 云 API 返回可重试错误（如限流）
- **THEN** 系统在 retry 超时时间内自动重试，最终成功则正常返回

#### Scenario: 重试耗尽
- **WHEN** 云 API 持续失败直至重试超时
- **THEN** 返回包装后的错误，资源状态不被错误清空

### Requirement: 单元测试覆盖
系统 SHALL 新增 `resource_tc_tencentcloud_dlc_dms_table_test.go`，使用 gomonkey mock DLC 云 API（`CreateDMSTableWithContext`/`DescribeDMSTableWithContext`/`AlterDMSTableWithContext`/`DropDMSTableWithContext`），覆盖 Create/Read/Update/Delete 业务逻辑，并通过 `go test -gcflags=all=-l` 执行通过。测试 SHALL NOT 使用 Terraform 验收测试套件。

#### Scenario: 单元测试通过
- **WHEN** 执行 `go test -gcflags=all=-l ./tencentcloud/services/dlc/ -run TestResourceTencentCloudDlcDmsTable`
- **THEN** 所有 mock 测试用例通过，验证 CRUD 业务逻辑正确
