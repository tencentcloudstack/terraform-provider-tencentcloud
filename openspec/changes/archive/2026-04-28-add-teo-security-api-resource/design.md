## Context

TEO（EdgeOne 边缘安全加速平台）是腾讯云的边缘安全产品，提供 DDoS 防护、WAF、Bot 管理等安全能力。API 资源（Security API Resource）是 TEO 安全防护的核心配置单元，用于定义需要保护的 API 端点及其关联的 API 服务。

当前 Terraform Provider 中已存在大量 TEO 资源（如 `tencentcloud_teo_acceleration_domain` 等），但缺少对 API 资源的管理能力。用户需要通过 Terraform 管理 API 资源的完整生命周期。

云 API 特点：
- CreateSecurityAPIResource: 创建 API 资源，传入 APIResources 列表（每次仅1个），返回 APIResourceIds 列表
- DescribeSecurityAPIResource: 按 ZoneId 查询，支持分页（Limit 最大 100），返回 TotalCount 和 APIResources 列表，不支持 Filter
- ModifySecurityAPIResource: 修改 API 资源，传入 APIResources 列表（需包含 Id 字段标识，每次仅1个）
- DeleteSecurityAPIResource: 删除 API 资源，传入 APIResourceIds 列表

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_security_api_resource` 资源，支持 TEO API 资源的 CRUD 生命周期管理
- 使用复合 ID `zoneId#apiResourceId` 标识每个 API 资源实例
- 支持分页查询 API 资源列表并按 apiResourceId 匹配特定资源
- 遵循现有 TEO 资源的代码风格和模式
- 在 provider.go/provider.md 中正确注册资源
- 生成对应的 .md 文档

**Non-Goals:**
- 不修改已有 TEO 资源的 schema 或行为
- 不实现数据源（data source）功能
- 不处理异步接口轮询（所有 CRUD 接口均为同步接口）

## Decisions

### 1. 资源 ID 设计：使用复合 ID `zoneId#apiResourceId`

**决策**: 使用 `zoneId#apiResourceId`（FILED_SP 分隔）作为 Terraform 资源 ID

**理由**:
- 每个 Terraform 资源实例对应一个 API 资源，而非整个站点下的所有 API 资源
- Create 返回的 APIResourceIds[0] 与 ZoneId 组合成复合 ID
- Read/Update/Delete 操作都从复合 ID 中解析出 zoneId 和 apiResourceId
- Delete 使用 APIResourceIds 参数传入单个 apiResourceId
- 支持标准的 Terraform Import，导入时需提供 `zoneId#apiResourceId`

### 2. api_resources 字段设计：TypeList 嵌套块，MaxItems: 1

**决策**: `api_resources` 使用 `schema.TypeList` + `schema.Resource` 嵌套块，MaxItems: 1

**理由**:
- 云 API 每次请求仅操作一个 API 资源
- 每个 APIResource 包含 name(Required)/path(Required)/api_service_ids(Optional)/methods(Optional)/request_constraint(Optional)/id(Computed)
- path 设为 Required，因为 API 资源必须定义路径
- Id 字段在创建时不需要（由服务端生成），在查询时返回，设为 Computed
- MaxItems: 1 限制每个资源实例仅管理一个 API 资源

### 3. Read 操作使用分页查询匹配

**决策**: DescribeTeoSecurityAPIResourceById 按 ZoneId 分页查询所有 API 资源（Limit=100），遍历结果匹配 apiResourceId

**理由**:
- 云 API 不支持 Filter 参数，只能按 ZoneId 查询所有 API 资源
- 需要分页获取所有资源后，遍历匹配目标 apiResourceId
- 遵循项目规则中"分页字段给定值为云 API 注释中标注的最大值"的要求
- 返回单个 `*APIResource` 而非列表，与其他 Describe*ById 方法保持一致

### 4. Update 操作：单资源修改模式

**决策**: ModifySecurityAPIResource 传入单个 APIResource（包含 Id 字段）

**理由**:
- 修改接口要求传入 APIResources 列表，每个 item 需带 Id
- MaxItems: 1 限制每次仅修改一个 API 资源
- 使用 buildSecurityAPIResourceFromMap 辅助函数构建请求参数，Update 时传入 apiResourceId

### 5. 服务层方法

**决策**: 在 `service_tencentcloud_teo.go` 中新增 `DescribeTeoSecurityAPIResourceById` 方法

**理由**:
- 遵循现有模式，每个资源都有对应的 Describe*ById 服务方法
- 封装分页查询逻辑和按 ID 匹配逻辑，返回指定的单个 APIResource
- 与其他 TEO 资源的 Describe*ById 方法签名保持一致

### 6. 辅助函数

**决策**: 新增 `buildSecurityAPIResourceFromMap` 辅助函数

**理由**:
- Create 和 Update 操作都需要构建 APIResource 对象
- Create 时 id 参数传空字符串（API 不接受 Id），Update 时传实际 apiResourceId
- 统一参数构建逻辑，减少重复代码

## Risks / Trade-offs

- **分页查询开销**: Read 操作需要分页获取所有 API 资源再匹配，当站点下 API 资源很多时可能有性能开销 → 使用 Limit=100 减少分页次数
- **并发修改**: 多个 Terraform 资源实例可能同时操作同一 ZoneId 下的 API 资源 → 这是 Terraform 的通用限制，通过 state lock 机制缓解
- **Import 复合 ID**: 用户需要知道 zoneId 和 apiResourceId 才能导入 → 在文档中明确说明导入格式
