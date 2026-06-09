## Context

Terraform Provider for TencentCloud 需要支持 EdgeOne (TEO) 的预热回源限速配置管理。该功能允许用户对预热任务回源带宽进行限速控制，避免预热任务占用过多源站带宽。

当前 TEO 模块已有多种 CONFIG 类型资源（如 `tencentcloud_teo_security_policy_config`、`tencentcloud_teo_ddos_protection_config`、`tencentcloud_teo_certificate_config`），本资源遵循相同模式。

云 API 接口分析：
- `ModifyPrefetchOriginLimit`：创建/更新限速配置，入参包含 ZoneId、DomainName、Area、Bandwidth、Enabled
- `DescribePrefetchOriginLimit`：查询限速配置，返回 PrefetchOriginLimit 列表，每项包含 ZoneId、DomainName、Area、Bandwidth、CreateTime、UpdateTime
- 注意：DescribePrefetchOriginLimit 不返回 Enabled 字段，Enabled 仅在 ModifyPrefetchOriginLimit 时使用（on=启用，off=删除）

关键发现：
- Enabled 字段在云 API 中为 on/off 开关，on 表示启用限速，off 表示删除限速
- DescribePrefetchOriginLimit 返回的 PrefetchOriginLimit 结构体不包含 Enabled 字段
- 同一 ZoneId 下可以有多个 DomainName + Area 组合的限速配置
- 联合 ID 采用 zone_id + domain_name + area（FILED_SP 分隔）

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_prefetch_origin_limit` CONFIG 类型资源，支持预热回源限速配置的完整生命周期管理
- 支持 Create（调用 ModifyPrefetchOriginLimit 设置 Enabled=on）、Read、Update、Delete（调用 ModifyPrefetchOriginLimit 设置 Enabled=off）
- 支持 Terraform Import
- 遵循现有 TEO CONFIG 资源的代码风格和架构模式

**Non-Goals:**
- 不实现数据源（datasource）
- 不暴露 DescribePrefetchOriginLimit 的分页参数（Offset/Limit）给用户
- 不修改其他现有资源的 schema 或行为

## Decisions

### 1. 资源类型：CONFIG（无 Create/Delete 云API，使用 Modify 替代）
**决定**：Create 方法调用 ModifyPrefetchOriginLimit（Enabled=on），Delete 方法调用 ModifyPrefetchOriginLimit（Enabled=off）
**原因**：该资源为配置类资源，云 API 只提供 Modify 和 Describe 接口，符合 RESOURCE_KIND_CONFIG 的定义。Enabled=off 在云 API 中表示删除该限速配置。

### 2. 联合 ID 设计
**决定**：使用 `zone_id + domain_name + area` 作为联合 ID，使用 `tccommon.FILED_SP` 作为分隔符
**原因**：同一站点下可以按域名和加速区域组合配置不同的限速策略，需要三个字段唯一标识一条限速配置。

### 3. Enabled 字段处理
**决定**：Enabled 字段作为 Required 参数暴露给用户，取值为 "on" 或 "off"
**原因**：Enabled 在 ModifyPrefetchOriginLimit 中是核心参数，控制限速配置的启停。用户需要显式控制启用状态。

### 4. Read 方法中匹配配置
**决定**：Read 方法调用 DescribePrefetchOriginLimit 并通过 Filters 过滤 domain-name 和 area，从返回列表中匹配对应配置
**原因**：DescribePrefetchOriginLimit 返回的是列表，需要通过 Filters 精确匹配到对应的配置项。

### 5. Bandwidth 类型
**决定**：Bandwidth 在 Terraform schema 中使用 TypeInt，对应云 API 的 int64 类型
**原因**：云 API 中 Bandwidth 取值范围 100-100000 Mbps，为整数值。

### 6. Delete 时清理资源状态
**决定**：Delete 方法调用 ModifyPrefetchOriginLimit 设置 Enabled=off 后，调用 Read 验证配置已被删除
**原因**：确保删除操作生效，符合 Terraform 资源删除的最终一致性要求。

## Risks / Trade-offs

- **[DescribePrefetchOriginLimit 不返回 Enabled]** → 在 Read 方法中，不设置 enabled 字段（因为 API 不返回），仅在其他参数不变时保持 enabled 为用户设置的值。Read 操作后 enabled 保持 Create/Update 时用户设置的值。
- **[配置删除后 Read 仍然可能返回数据]** → Delete 设置 Enabled=off 后，配置可能需要短暂时间才能完全删除。在 Read 中检查匹配结果，若未找到则设置 d.SetId("") 标记资源已删除。
- **[分页限制]** → DescribePrefetchOriginLimit 有分页限制，但配置项数量通常很少（按 domain+area 组合），设置 Limit=100 足够覆盖。
