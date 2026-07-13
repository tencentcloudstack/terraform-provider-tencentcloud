## Why

腾讯云数据湖计算（DLC，Data Lake Compute）的 DMS 元数据服务目前无法通过 Terraform 管理数据库（Database）的完整生命周期。用户只能通过控制台或 SDK 手动创建、查询、修改和删除 DMS 元数据库，无法实现基础设施即代码（IaC）的自动化管理，存在配置漂移、团队协作困难等问题。

为补齐 DLC 产品在 Terraform Provider 中对 DMS 元数据库管理的空白，需要新增 `tencentcloud_dlc_dms_database` 通用资源（RESOURCE_KIND_GENERAL），覆盖 DMS 元数据库的创建、读取、更新、删除全流程。

## What Changes

- 新增 Terraform 资源 `tencentcloud_dlc_dms_database`，类型为 RESOURCE_KIND_GENERAL，实现完整 CRUD：
  - **Create**：调用 DLC SDK `CreateDMSDatabase` 接口创建 DMS 元数据库。
  - **Read**：调用 DLC SDK `DescribeDMSDatabase` 接口查询 DMS 元数据库详情。
  - **Update**：调用 DLC SDK `AlterDMSDatabase` 接口更新 DMS 元数据库（名称、路径、基础对象等）。
  - **Delete**：调用 DLC SDK `DropDMSDatabase` 接口删除 DMS 元数据库。
- 资源 ID 采用复合 ID（`name` + `schema_name` + `datasource_connection_name`，以 `tccommon.FILED_SP` 即 `#` 分隔），以适配云 API 通过多个字段定位资源的特性。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册新资源。
- 生成对应的资源示例 markdown（`resource_tc_dlc_dms_database.md`）并通过 `make doc` 生成 `website/docs/` 文档。
- 新增单元测试文件 `resource_tc_dlc_dms_database_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试。

## Capabilities

### New Capabilities

- `dlc-dms-database-resource`: 通过 Terraform 管理 DLC DMS 元数据库（Database）的完整生命周期资源，包含 schema 定义、CRUD 实现、复合 ID、导入支持、provider 注册与文档示例。

### Modified Capabilities

<!-- 无现有 capability 需要修改 -->

## Impact

- **新增代码文件**：
  - `tencentcloud/services/dlc/resource_tc_dlc_dms_database.go`：资源 CRUD 实现。
  - `tencentcloud/services/dlc/resource_tc_dlc_dms_database_test.go`：单元测试（gomonkey mock）。
  - `tencentcloud/services/dlc/resource_tc_dlc_dms_database.md`：资源示例文档。
- **修改文件**：
  - `tencentcloud/provider.go`：在 DLC 资源组中注册 `tencentcloud_dlc_dms_database`。
  - `tencentcloud/provider.md`：通过 `make doc` 自动生成文档索引条目。
- **依赖的云 API**（SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125`，均已存在于 vendor）：
  - `CreateDMSDatabase`：创建库（入参 Asset/SchemaName/Location/Name/DatasourceConnectionName）。
  - `DescribeDMSDatabase`：查询库（入参 Name/SchemaName/Pattern/DatasourceConnectionName，出参 Name/SchemaName/Location/Asset）。
  - `AlterDMSDatabase`：更新库（入参 CurrentName/SchemaName/Location/Asset/DatasourceConnectionName）。
  - `DropDMSDatabase`：删除库（入参 Name/DeleteData/Cascade/DatasourceConnectionName）。
- **兼容性**：纯新增资源，不影响任何现有资源配置与 state，向后完全兼容。
- **异步说明**：上述 4 个 DMS Database 接口均为同步接口（响应仅返回 RequestId，无 TaskId 异步句柄），Create/Update/Delete 后直接调用 Read 验证即可，无需轮询。
