## Context

TencentCloud EdgeOne (TEO) 提供共享 CNAME 功能，允许用户创建共享 CNAME 并将多个加速域名绑定到同一个 CNAME 记录。当前 provider 已有多个 TEO 资源（如 `tencentcloud_teo_zone`、`tencentcloud_teo_function` 等），本次新增 `tencentcloud_teo_shared_cname` 资源遵循相同的代码组织模式。

云 API SDK 已在 vendor 中可用：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`，包含完整的 CRUD 接口：
- `CreateSharedCNAME`：创建共享 CNAME
- `DescribeSharedCNAME`：查询共享 CNAME 列表
- `ModifySharedCNAME`：修改共享 CNAME
- `DeleteSharedCNAME`：删除共享 CNAME

所有接口均为同步接口，无需轮询等待。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_shared_cname` 资源的完整 CRUD 生命周期
- 支持通过 `zone_id` + `shared_cname` 联合 ID 进行 import
- 支持修改 description 和 ipssl_setting
- 遵循现有 provider 代码模式（retry、错误处理、日志等）
- 提供单元测试（使用 gomonkey mock 云 API）

**Non-Goals:**
- 不实现 `BindSharedCNAME` 接口（绑定/解绑加速域名到共享 CNAME 的操作属于独立资源）
- 不实现数据源 `tencentcloud_teo_shared_cname_list`（本次仅实现 RESOURCE_KIND_GENERAL）
- 不暴露 DescribeSharedCNAME 的分页参数给用户

## Decisions

### 1. 资源 ID 设计
**决策**：使用 `zone_id` + `shared_cname` 的联合 ID，以 `tccommon.FILED_SP` 作为分隔符。

**理由**：`shared_cname` 是由 API 创建后返回的完整 CNAME 值，是资源的唯一标识。但 Delete/Modify/Describe 接口都需要 `zone_id`，因此需要联合 ID 来存储两个值。

**替代方案**：仅使用 `shared_cname` 作为 ID —— 但这样 Read/Update/Delete 时无法获取 `zone_id`。

### 2. ForceNew 字段设计
**决策**：`zone_id` 和 `shared_cname_prefix` 设为 ForceNew。

**理由**：创建后 `zone_id` 不可变，`shared_cname_prefix` 仅在创建时使用（创建后 API 返回完整的 `shared_cname`，无法修改前缀）。

### 3. Read 实现方式
**决策**：使用 `DescribeSharedCNAME` 接口，通过 Filters 中的 `shared-cname` 字段精确过滤目标资源。

**理由**：API 没有提供单个资源的 Get 接口，只有列表查询接口。通过 filter 精确匹配 `shared_cname` 值可以获取单个资源的详情。

### 4. Update 实现方式
**决策**：Update 方法支持修改 `description` 和 `ipssl_setting` 两个字段，使用 `ModifySharedCNAME` 接口。

**理由**：这是 ModifySharedCNAME 接口支持的全部可修改参数。

### 5. Schema 中 ipssl_setting 的设计
**决策**：`ipssl_setting` 作为 Optional 的嵌套 block，包含 `operation`（bind/unbind）和 `associated_domain` 两个字段。

**理由**：与 SDK 中 `IPSSLSetting` 结构体一致，用于设置 IP SSL 类型共享 CNAME 的绑定关系。

### 6. 测试方式
**决策**：使用 gomonkey mock 云 API 进行单元测试，不使用 Terraform 验收测试套件。

**理由**：按照项目要求，新增资源使用 mock 方式进行单元测试。

## Risks / Trade-offs

- [Risk] DescribeSharedCNAME 返回空列表时资源可能已被外部删除 → 在 Read 中检测到空结果时打印日志并 SetId("") 标记资源已删除
- [Risk] CreateSharedCNAME 返回的 SharedCNAME 为空 → 返回 NonRetryableError，避免写入空 ID
- [Risk] IPSSLSetting 的 bind/unbind 操作可能存在异步延迟 → 由于 API 文档未标注为异步接口，按同步处理
