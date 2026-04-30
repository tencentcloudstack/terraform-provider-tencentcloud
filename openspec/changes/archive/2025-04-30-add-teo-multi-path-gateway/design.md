## Context

Terraform Provider for TencentCloud 需要新增 `tencentcloud_teo_multi_path_gateway` 资源，支持 TEO (EdgeOne) 多通道安全加速网关的完整生命周期管理。当前 TEO 多通道安全加速网关的线路管理(`teo-multi-path-gateway-line-resource`)、区域查询(`teo-multi-path-gateway-region-datasource`)和密钥配置(`teo-multi-path-gateway-secret-key-config`)已有 spec，但网关本身的 CRUD 资源尚无支持。

云 API 提供四个核心接口：
- `CreateMultiPathGateway`: 创建网关，需要 ZoneId，根据 GatewayType 区分云上网关(cloud)和自有网关(private)
- `DescribeMultiPathGateways`: 查询网关列表，返回 MultiPathGateway 结构体
- `ModifyMultiPathGateway`: 修改网关名称、IP、端口
- `DeleteMultiPathGateway`: 删除网关

关键约束：资源使用 ZoneId + GatewayId 作为复合 ID，支持 import。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_multi_path_gateway` 资源的完整 CRUD 操作
- 使用 ZoneId + GatewayId 作为复合 ID（FILED_SP 分隔），支持 import
- 正确处理云上网关(cloud)和自有网关(private)的差异化参数（cloud 需要 RegionId，private 需要 GatewayIP）
- 在 Read 中通过 DescribeMultiPathGateways 按 GatewayId 过滤获取网关详情
- 在 Update 中调用 ModifyMultiPathGateway 更新可修改字段
- 遵循 igtm_strategy 资源代码风格

**Non-Goals:**
- 不管理网关线路(MultiPathGatewayLine)，该能力由 teo-multi-path-gateway-line-resource 覆盖
- 不管理网关密钥，该能力由 teo-multi-path-gateway-secret-key-config 覆盖
- 不管理网关状态切换(ModifyMultiPathGatewayStatus)，仅管理网关资源配置
- 不实现 datasource（仅实现 RESOURCE_KIND_GENERAL 资源）

## Decisions

### 1. 复合 ID 设计
**决策**: 使用 `ZoneId:GatewayId` 作为 Terraform 资源 ID（使用 tccommon.FILED_SP 分隔符）

**理由**: DeleteMultiPathGateway 和 ModifyMultiPathGateway 都需要 ZoneId + GatewayId 两个参数，DescribeMultiPathGateways 也需要 ZoneId。使用复合 ID 可以在 Read/Update/Delete 中解析出两个必需参数。

**替代方案**: 仅使用 GatewayId 作为 ID — 不可行，因为所有 API 都需要 ZoneId。

### 2. Schema 设计
**决策**:
- `zone_id`: Required, ForceNew — 创建后不可变更
- `gateway_type`: Required, ForceNew — 网关类型创建后不可变更
- `gateway_name`: Required — 可通过 ModifyMultiPathGateway 修改
- `gateway_port`: Optional, Computed — 可通过 ModifyMultiPathGateway 修改（仅 private 类型）
- `region_id`: Optional, Computed, ForceNew — 仅 cloud 类型需要，创建后不可变更（ModifyMultiPathGateway 不接受 RegionId 参数）
- `gateway_ip`: Optional, Computed — 仅 private 类型需要，可通过 ModifyMultiPathGateway 修改
- `gateway_id`: Computed — 由云 API 返回
- `status`: Computed — 由云 API 返回
- `need_confirm`: Computed — 由云 API 返回

**理由**: 根据 CreateMultiPathGateway 入参与 ModifyMultiPathGateway 入参的对比：
- Create 有但 Modify 没有的参数：GatewayType, RegionId → 设为 ForceNew
- Create 和 Modify 都有的参数：GatewayName, GatewayIP, GatewayPort → 允许更新
- 仅 Modify 有的参数：GatewayId → 已在复合 ID 中

### 3. Read 实现策略
**决策**: 调用 DescribeMultiPathGateways 并通过 Filters 按 GatewayId 过滤

**理由**: DescribeMultiPathGateways 支持 Filters 过滤（gateway-type, keyword），但没有直接的 GatewayId 过滤。需要在返回的 Gateways 列表中匹配 GatewayId 来找到目标网关。

### 4. Update 实现策略
**决策**: 直接调用 ModifyMultiPathGateway 传入修改后的字段

**理由**: ModifyMultiPathGateway 接口接受 ZoneId、GatewayId、GatewayName、GatewayIP、GatewayPort 参数，覆盖了所有可更新字段。

## Risks / Trade-offs

- **[风险] DescribeMultiPathGateways 无 GatewayId 直接过滤** → 通过遍历返回的 Gateways 列表匹配 GatewayId 来缓解，使用 Limit=1000 获取最大数量减少遗漏
- **[风险] 修改 GatewayPort 仅支持 private 类型网关** → ModifyMultiPathGateway API 注释说明仅 private 类型可修改端口，Terraform 侧不做额外限制，由 API 侧返回错误
- **[风险] GatewayType 取值为 cloud 时 RegionId 必填、private 时 GatewayIP 必填** → 使用 ConflictingWith 或在 Create 中根据 GatewayType 校验必填字段
