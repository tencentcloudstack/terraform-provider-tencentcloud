## Context

腾讯云 EdgeOne (TEO) 多通道安全加速网关需要新增 Terraform 资源支持。当前该资源仅能通过控制台或 API 操作，无法通过 Terraform 进行基础设施即代码管理。

当前状态：
- 云 API 已提供完整的 CRUD 接口：`CreateMultiPathGateway`、`DescribeMultiPathGateways`、`ModifyMultiPathGateway`、`DeleteMultiPathGateway`
- SDK 已在 vendor 目录中可用：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
- TEO 服务目录 `tencentcloud/services/teo/` 已存在，可直接添加新资源文件

约束：
- 资源 ID 由 `zone_id` 和 `gateway_id` 联合组成，使用 `tccommon.FILED_SP` 分隔
- 网关类型 `gateway_type` 在创建后不可修改（Create 有此参数，Modify 无此参数）
- `region_id` 仅在创建时需要（Modify 接口无此参数）
- 代码风格参考 `tencentcloud_igtm_strategy` 资源

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_multi_path_gateway` 资源的完整 CRUD 生命周期管理
- 支持资源 Import（通过联合 ID: `zone_id#gateway_id`）
- 在 Read 方法中通过 DescribeMultiPathGateways 查询并匹配到指定网关
- 在 Update 方法中仅支持修改 ModifyMultiPathGateway 接口支持的字段（gateway_name、gateway_ip、gateway_port）
- 正确处理 gateway_type 和 region_id 为 ForceNew 字段
- 添加 gomonkey 单元测试验证业务逻辑
- 在 provider.go 和 provider.md 中注册新资源

**Non-Goals:**
- 不实现多通道安全加速网关线路（Line）管理（有独立的 CreateMultiPathGatewayLine/ModifyMultiPathGatewayLine/DeleteMultiPathGatewayLine 接口）
- 不实现网关密钥管理（有独立的 CreateMultiPathGatewaySecretKey/ModifyMultiPathGatewaySecretKey 接口）
- 不实现网关状态修改（有独立的 ModifyMultiPathGatewayStatus 接口）
- 不添加数据源（仅实现 RESOURCE_KIND_GENERAL 资源）

## Decisions

### 1. 资源 ID 格式
**决策**: 使用 `zone_id`tccommon.FILED_SP`gateway_id` 作为复合 ID
**理由**: Delete 和 Modify 接口均需要 ZoneId 和 GatewayId 两个参数，使用联合 ID 可确保 Read/Update/Delete 方法能正确获取所需参数。参考项目中其他 TE O 资源的模式。

### 2. gateway_type 和 region_id 为 ForceNew
**决策**: `gateway_type` 和 `region_id` 设置为 ForceNew
**理由**: ModifyMultiPathGateway 接口不支持修改这两个字段，修改后只能通过删除重建来实现。

### 3. Read 方法实现
**决策**: 调用 DescribeMultiPathGateways 并通过 GatewayId 过滤匹配到指定网关
**理由**: 没有单独的 DescribeMultiPathGatewayById 接口，需要使用列表查询接口并通过 Filter 过滤。DescribeMultiPathGateways 的 Filters 支持 gateway-id 过滤。

### 4. gateway_port 类型处理
**决策**: gateway_port 在 Terraform schema 中使用 TypeInt，对应云 API 的 int64 类型
**理由**: 云 API 中 GatewayPort 为 int64 类型，Terraform 中 TypeInt 最为匹配。

### 5. Filters 参数
**决策**: DescribeMultiPathGateways 请求中使用 Filters 按 gateway-id 过滤，不将 Filters 暴露给 Terraform 用户
**理由**: Read 方法需要精确查询单个网关，Filters 用于内部查询而非用户配置。

## Risks / Trade-offs

- **[风险] DescribeMultiPathGateways 为列表接口** → 通过 Filters 按 gateway-id 过滤来精确定位目标网关，同时在 Read 中校验返回结果是否唯一
- **[风险] gateway_port 在 Modify 中仅支持 private 类型网关修改** → 在 Update 方法中不限制，由云 API 返回错误；但在 Read 方法中正确回写所有字段
- **[风险] 联合 ID 导入时需要用户提供 zone_id#gateway_id 格式** → 在 .md 文档的 Import 部分明确说明 ID 格式要求
