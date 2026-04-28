## Context

TEO (EdgeOne) 是腾讯云的边缘安全加速平台。当前 Terraform Provider 中已有大量 TEO 资源（如 teo_zone、teo_security_api_resource、teo_origin_acl 等），但缺少对多通道安全加速网关（MultiPathGateway）的支持。

多通道安全加速网关允许用户创建云上网关（cloud）和自有网关（private），管理多通道加速线路。云 API 已提供完整的 CRUD 接口：
- `CreateMultiPathGateway`：创建网关，返回 GatewayId
- `DescribeMultiPathGateways`：查询网关列表，支持过滤
- `ModifyMultiPathGateway`：修改网关名称、IP 和端口
- `DeleteMultiPathGateway`：删除网关

现有资源模式参考：`tencentcloud_teo_security_api_resource`（RESOURCE_KIND_GENERAL，同样使用 zone_id + 资源 ID 作为复合 ID）。

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_multi_path_gateway` 资源，支持完整的 CRUD 生命周期管理
- 支持通过 Terraform 创建云上网关（cloud）和自有网关（private）两种类型
- 支持资源导入（Import）
- 遵循现有 TEO 资源的代码模式和约定

**Non-Goals:**
- 不支持网关线路（Lines）的独立管理，线路信息仅作为只读属性展示
- 不支持数据源（datasource），本次仅新增资源
- 不修改任何现有资源的 schema 或行为

## Decisions

### 1. 复合 ID 格式：`zone_id#gateway_id`
**选择**: 使用 `zone_id` + `gateway_id` 以 `#` 分隔符组成复合 ID
**理由**: 所有 TEO 资源都使用 zone_id 作为上下文标识，且所有 CRUD 接口都需要 zone_id 参数。这与 `teo_security_api_resource` 等现有资源保持一致。
**替代方案**: 仅使用 gateway_id 作为 ID — 不可行，因为所有 API 调用都需要 zone_id，且 gateway_id 在不同 zone 下可能重复。

### 2. Schema 字段分类
**Create-only（ForceNew）字段**:
- `zone_id`：站点 ID，所有接口都需要但创建后不可变更
- `gateway_type`：网关类型（cloud/private），Modify 接口不支持修改，因此为 ForceNew
- `region_id`：网关地域，仅 cloud 类型创建时需要，Modify 接口不支持修改，因此为 ForceNew

**Updatable 字段**:
- `gateway_name`：网关名称，Modify 接口支持修改
- `gateway_port`：网关端口，Modify 接口支持修改
- `gateway_ip`：网关地址，Modify 接口支持修改（仅 private 类型）

**Computed（只读）字段**:
- `gateway_id`：网关 ID，创建后由 API 返回
- `status`：网关状态，由 API 返回
- `need_confirm`：回源 IP 列表变化是否需要确认，由 API 返回
- `lines`：线路信息列表，由 API 返回

### 3. Describe 策略：通过 Filters 过滤 GatewayId
**选择**: 在 DescribeMultiPathGateways 请求中使用 Filters 按 gateway-id 过滤指定网关
**理由**: DescribeMultiPathGateways 是列表接口，支持 Filters 过滤。使用 gateway-id 作为过滤条件可以精确查询单个网关，避免遍历整个列表。
**替代方案**: 遍历列表后匹配 GatewayId — 效率较低，Filters 已支持 gateway-id 过滤。

### 4. Lines 字段设计为 TypeList 嵌套结构
**选择**: 将 `lines` 设计为 `TypeList`，每个元素为 `TypeMap` 包含 line_id、line_type、line_address、proxy_id、rule_id 字段
**理由**: MultiPathGatewayLine 是固定结构体，使用嵌套 schema 更符合 Terraform 习惯，且便于用户在 state 中查看线路信息。

### 5. Update 逻辑：仅发送有变更的字段
**选择**: 在 Update 方法中使用 `d.HasChange()` 检查 gateway_name、gateway_ip、gateway_port 是否有变更，仅发送有变更的字段
**理由**: Modify 接口的参数都是可选的，仅发送变更字段减少不必要的 API 调用。这与 igtm_strategy 的模式一致。

## Risks / Trade-offs

- **[API 参数约束]** GatewayType 为 cloud 时，RegionId 必填；为 private 时，GatewayIP 必填。→ 在 schema 中通过 ConflictsWith 或自定义验证处理，或在 Create 方法中添加条件检查逻辑
- **[API 修改限制]** ModifyMultiPathGateway 的 GatewayIP 仅 private 类型可修改，GatewayPort 仅 private 类型可修改。→ 在 Update 方法中添加条件检查，或在 schema description 中说明限制
- **[Describe 列表接口]** DescribeMultiPathGateways 是列表接口，需要使用 Filters 精确查询。→ 使用 gateway-id 过滤条件确保精确匹配
- **[Lines 只读]** Lines 信息由系统管理，用户不能通过 Terraform 修改。→ 将 lines 设为 Computed 只读字段
