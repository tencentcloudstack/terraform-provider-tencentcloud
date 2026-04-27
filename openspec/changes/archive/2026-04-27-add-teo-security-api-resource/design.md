## Context

TEO（EdgeOne 边缘安全加速平台）是腾讯云的边缘安全产品，提供 DDoS 防护、WAF、Bot 管理等安全能力。API 资源（Security API Resource）是 TEO 安全防护的核心配置单元，用于定义需要保护的 API 端点及其关联的 API 服务。

当前 Terraform Provider 中已存在大量 TEO 资源（如 `tencentcloud_teo_acceleration_domain` 等），但缺少对 API 资源的管理能力。用户需要通过 Terraform 管理 API 资源的完整生命周期。

云 API 特点：
- CreateSecurityAPIResource: 批量创建，传入 APIResources 列表，返回 APIResourceIds 列表
- DescribeSecurityAPIResource: 按 ZoneId 查询，支持分页（Limit 最大 100），返回 TotalCount 和 APIResources 列表
- ModifySecurityAPIResource: 批量修改，传入 APIResources 列表（需包含 Id 字段标识）
- DeleteSecurityAPIResource: 批量删除，传入 APIResourceIds 列表

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_security_api_resource` 资源，支持 TEO API 资源的 CRUD 生命周期管理
- 以 ZoneId 为维度管理站点下的所有 API 资源，资源 ID 使用 zone_id
- 支持分页查询 API 资源列表
- 遵循现有 TEO 资源的代码风格和模式
- 在 provider.go/provider.md 中正确注册资源
- 生成对应的 .md 文档

**Non-Goals:**
- 不修改已有 TEO 资源的 schema 或行为
- 不实现数据源（data source）功能
- 不处理异步接口轮询（所有 CRUD 接口均为同步接口）

## Decisions

### 1. 资源 ID 设计：使用 zone_id 作为资源 ID

**决策**: 使用 `zone_id` 作为 Terraform 资源 ID（`d.SetId(zoneId)`）

**理由**:
- 该资源的所有 CRUD 操作都以 ZoneId 为入参
- 一个 Terraform 资源实例管理一个站点下的所有 API 资源
- Create 返回的 APIResourceIds 作为 computed 字段 `api_resource_ids` 存储
- Delete 需要 ZoneId + APIResourceIds，其中 APIResourceIds 从 `api_resource_ids` 字段获取

**备选方案**: 使用 `zone_id#api_resource_id_1#api_resource_id_2` 复合 ID
- 不采用原因：API 资源以 ZoneId 为维度批量管理，不适合拆分为多个独立资源

### 2. api_resources 字段设计：TypeList 嵌套块

**决策**: `api_resources` 使用 `schema.TypeList` + `schema.TypeMap` 嵌套块，每个 item 包含 id/name/api_service_ids/path/methods/request_constraint

**理由**:
- 云 API 的 APIResources 字段为 `[]*APIResource` 列表结构
- 每个 APIResource 包含 6 个字段：Id, Name, APIServiceIds, Path, Methods, RequestConstraint
- Id 字段在创建时不需要（由服务端生成），在查询时返回，设为 computed
- Name 为必填字段，其余为可选

### 3. Read 操作使用分页查询

**决策**: DescribeSecurityAPIResource 使用 Limit=100（最大值）进行分页查询，循环获取所有 API 资源

**理由**:
- 云 API 支持分页，Limit 最大值为 100
- 需要确保获取站点下所有 API 资源
- 遵循项目规则中"分页字段给定值为云 API 注释中标注的最大值"的要求

### 4. Update 操作：全量替换模式

**决策**: ModifySecurityAPIResource 传入完整的 APIResources 列表（包含 Id 字段），实现全量更新

**理由**:
- 修改接口要求传入完整的 APIResources 列表，每个 item 需带 Id
- Read 时已经获取了所有 API 资源（包含 Id），Update 时将用户配置的 api_resources 与已有的 Id 对应后传入

### 5. 服务层方法

**决策**: 在 `service_tencentcloud_teo.go` 中新增 `DescribeSecurityAPIResourceById` 方法

**理由**:
- 遵循现有模式，每个资源都有对应的 Describe*ById 服务方法
- 封装分页查询逻辑，返回指定 ZoneId 下的所有 API 资源

## Risks / Trade-offs

- **批量操作风险**: Create/Delete 都是批量操作，如果部分失败可能导致状态不一致 → 通过 Terraform 的 Read 操作重新同步状态来缓解
- **ID 对应关系**: Create 返回的 APIResourceIds 与传入的 APIResources 顺序对应，但修改时用户可能改变列表顺序 → Update 时需要通过 Read 获取当前状态中的 Id 信息
- **并发修改**: 多个 Terraform 资源实例可能同时操作同一 ZoneId 下的 API 资源 → 这是 Terraform 的通用限制，通过 state lock 机制缓解
