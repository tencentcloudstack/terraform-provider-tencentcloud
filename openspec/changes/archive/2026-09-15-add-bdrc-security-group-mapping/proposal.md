## Why

业务灾难恢复中心（BDRC，Business Disaster Recovery Center）产品目前尚未接入 Terraform Provider，用户无法通过 Terraform 管理站点对（SitePair）下生产端与容灾端之间的安全组映射。安全组映射用于在容灾切换时保证安全组规则的同步，是 BDRC 实例类容灾的核心配套能力。本次需要新增 `tencentcloud_bdrc_security_group_mapping` 资源（RESOURCE_KIND_GENERAL），让用户可以以基础设施即代码的方式创建、查询、删除安全组映射。

## What Changes

- 新增资源 `tencentcloud_bdrc_security_group_mapping`（通用资源，管理完整生命周期）：
  - **Create**：调用 `CreateSecurityGroupMapping` 创建安全组映射（入参：`src_security_group_id`、`target_security_group_id`、`site_pair_id`）。该接口为异步接口且不返回映射 ID，Create 后需轮询 `DescribeSecurityGroupMappings` 直到新建映射生效出现，并获取其 `SecurityGroupMappingId` 作为资源 ID 的一部分。
  - **Read**：调用 `DescribeSecurityGroupMappings`（按 `SitePairId` + `SecurityGroupMappingId` 过滤）查询单条映射详情并回填 state。
  - **Update**：云 API 无安全组映射的更新接口（只有 Create/Describe/Delete），属于"只有 CRD 接口"的资源。除 `id` 外所有顶层参数加入 `immutableArgs`，若变更则返回 error，强制用户删除重建。
  - **Delete**：调用 `DeleteSecurityGroupMapping` 删除安全组映射（入参：`site_pair_id`、`security_group_mapping_ids`）。
- 新增服务目录 `tencentcloud/services/bdrc/`，包含：
  - `service_tencentcloud_bdrc.go`：服务层，封装 `DescribeSecurityGroupMappings` 按 ID 查询/轮询逻辑。
  - `resource_tc_bdrc_security_group_mapping.go`：资源实现（含单元测试 `resource_tc_bdrc_security_group_mapping_test.go`，使用 gomonkey mock 云 API）。
  - `resource_tc_bdrc_security_group_mapping.md`：资源使用文档（由 `make doc` 生成 website 文档的输入）。
- 在 `tencentcloud/connectivity/client.go` 中新增 `UseBdrcV20260330Client()` 客户端访问器及对应连接字段（当前 provider 中尚无 bdrc 客户端）。
- 在 `tencentcloud/provider.go` 中注册 `tencentcloud_bdrc_security_group_mapping` 资源，并在 `tencentcloud/provider.md` 中补充资源条目。
- 资源 ID 采用联合 ID：`site_pair_id # security_group_mapping_id`（分隔符为 `tccommon.FILED_SP`，即 `#`），Read/Update/Delete 均从 `d.Id()` 解析这两部分作为请求参数。

## Capabilities

### New Capabilities
- `bdrc-security-group-mapping-resource`: BDRC 安全组映射资源（`tencentcloud_bdrc_security_group_mapping`）的生命周期管理能力，覆盖基于 `CreateSecurityGroupMapping` / `DescribeSecurityGroupMappings` / `DeleteSecurityGroupMapping` 三个云 API 的增删查实现、异步创建轮询、不可变字段约束以及联合 ID 设计。

### Modified Capabilities
<!-- 无：本次为全新资源，不修改任何既有能力的需求。 -->

## Impact

- **新增代码**：
  - `tencentcloud/services/bdrc/resource_tc_bdrc_security_group_mapping.go`（资源 CRUD 实现）
  - `tencentcloud/services/bdrc/service_tencentcloud_bdrc.go`（服务层查询/轮询封装）
  - `tencentcloud/services/bdrc/resource_tc_bdrc_security_group_mapping_test.go`（gomonkey 单元测试）
  - `tencentcloud/services/bdrc/resource_tc_bdrc_security_group_mapping.md`（资源文档）
- **修改代码**：
  - `tencentcloud/connectivity/client.go`：新增 `bdrcv20260330Conn` 字段与 `UseBdrcV20260330Client()` 方法、bdrc SDK import。
  - `tencentcloud/provider.go`：注册资源 `tencentcloud_bdrc_security_group_mapping`。
  - `tencentcloud/provider.md`：补充资源条目。
- **依赖**：`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330`（已在 vendor 目录中就绪，无需变更 go.mod）。
- **云 API**：
  - `CreateSecurityGroupMapping`（bdrc v20260330，创建安全组映射，异步，无返回 ID）
  - `DescribeSecurityGroupMappings`（bdrc v20260330，查询安全组映射列表，支持 `Filters`（`src-security-group-id` / `target-security-group-id`）、`Offset`/`Limit` 分页、`Order`/`OrderField` 排序）
  - `DeleteSecurityGroupMapping`（bdrc v20260330，批量删除安全组映射）
- **兼容性**：纯新增资源，不影响任何现有资源的 schema 与 state，向后兼容。
