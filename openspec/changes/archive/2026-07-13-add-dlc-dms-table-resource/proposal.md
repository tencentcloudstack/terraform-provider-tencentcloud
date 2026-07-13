## Why

Terraform TencentCloud provider 目前缺少对 DLC（Data Lake Compute）DMS 元数据表的管理能力。用户无法通过 Terraform 声明式地创建、查询、更新和删除 DMS 元数据表（`tencentcloud_dlc_dms_table`），导致涉及 DLC 湖仓元数据的基础设施编排必须依赖手工控制台操作或脚本，无法纳入 IaC 统一管理。新增该资源可以让用户以 Terraform 管理表结构、分区、列、存储描述等完整生命周期，提升湖仓元数据治理的自动化与一致性。

## What Changes

- 新增 Terraform 资源 `tencentcloud_dlc_dms_table`（RESOURCE_KIND_GENERAL），覆盖 DMS 元数据表的完整 CRUD 生命周期：
  - 创建：调用 DLC `CreateDMSTable` 接口。
  - 查询：调用 DLC `DescribeDMSTable` 接口。
  - 更新：调用 DLC `AlterDMSTable` 接口。
  - 删除：调用 DLC `DropDMSTable` 接口。
- 资源采用复合 ID（数据库名 + 表名），通过 `tccommon` 的 `FILED_SP` 分隔符组合，以支持导入与唯一标识。
- 在 `tencentcloud/provider.go` 与 `tencentcloud/provider.md` 中注册新资源。
- 新增资源文档 `tencentcloud/services/dlc/resource_tc_tencentcloud_dlc_dms_table.md`。
- 新增单元测试 `tencentcloud/services/dlc/resource_tc_tencentcloud_dlc_dms_table_test.go`，使用 gomonkey mock 云 API 进行业务逻辑测试。

## Capabilities

### New Capabilities
- `dlc-dms-table-resource`: 通过 Terraform 管理 DLC DMS 元数据表资源的完整生命周期（创建、读取、更新、删除），包含表基础信息、列、分区键、分区、存储描述（Sds）、视图文本及数据源连接等参数。

### Modified Capabilities
<!-- 无现有 spec 需要修改 -->

## Impact

- **新增代码文件**：
  - `tencentcloud/services/dlc/resource_tc_tencentcloud_dlc_dms_table.go`（资源 CRUD 实现）
  - `tencentcloud/services/dlc/resource_tc_tencentcloud_dlc_dms_table_test.go`（单元测试）
  - `tencentcloud/services/dlc/resource_tc_tencentcloud_dlc_dms_table.md`（资源文档）
- **修改代码文件**：
  - `tencentcloud/provider.go`：注册 `tencentcloud_dlc_dms_table` 资源。
  - `tencentcloud/provider.md`：新增资源文档条目。
- **依赖**：依赖 `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125` 中已有的 `CreateDMSTable`、`DescribeDMSTable`、`AlterDMSTable`、`DropDMSTable` 接口及 `Asset`、`DMSSds`、`DMSColumn`、`DMSPartition`、`KVPair` 等数据结构（已 vendored）。
- **向后兼容**：纯新增资源，不修改任何现有资源 schema，完全向后兼容。
