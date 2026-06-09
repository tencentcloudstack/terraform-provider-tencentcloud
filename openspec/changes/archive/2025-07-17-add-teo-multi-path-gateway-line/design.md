## Context

腾讯云 EdgeOne (TEO) 多通道安全加速网关线路资源需要通过 Terraform 进行管理。目前 provider 中没有此资源的实现。该资源需要支持完整的 CRUD 操作，通过 TEO SDK 的 `CreateMultiPathGatewayLine`、`DescribeMultiPathGatewayLine`、`ModifyMultiPathGatewayLine`、`DeleteMultiPathGatewayLine` 四个接口实现。

资源使用复合 ID（zone_id#gateway_id#line_id）标识唯一实例，参考 `tencentcloud_igtm_strategy` 资源的实现模式。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_multi_path_gateway_line` 资源的完整 CRUD 功能
- 支持资源导入（Import）
- 在 provider.go 和 provider.md 中正确注册新资源
- 编写单元测试，使用 gomonkey mock 云 API
- 生成资源文档（.md 文件）

**Non-Goals:**
- 不实现数据源（data source），本资源仅为 RESOURCE_KIND_GENERAL
- 不修改已有的 TEO 资源
- 不支持异步操作轮询（云 API 接口均为同步接口）

## Decisions

### 1. 复合 ID 设计
**决策**: 使用 `zone_id#gateway_id#line_id` 作为复合 ID，以 `tccommon.FILED_SP`（`#`）作为分隔符。

**理由**: 资源的唯一标识需要 zone_id、gateway_id 和 line_id 三个字段共同确定，与 DescribeMultiPathGatewayLine 接口的入参一致。参考 igtm_strategy 使用 `instanceId#strategyId` 的模式。

### 2. Schema 字段设计
**决策**:
- `zone_id` 和 `gateway_id`：Required + ForceNew（创建后不可修改）
- `line_id`：Computed（由云 API 创建后返回，不由用户指定）
- `line_type`：Required（线路类型，创建时必填，支持更新）
- `line_address`：Required（线路地址，格式为 ip:port，支持更新）
- `proxy_id`：Optional（四层代理实例 ID，LineType 为 proxy 时必传，支持更新）
- `rule_id`：Optional（转发规则 ID，LineType 为 proxy 时必传，支持更新）

**理由**: 根据 CreateMultiPathGatewayLine 入参，zone_id、gateway_id 为必填；line_type 和 line_address 为必填；proxy_id 和 rule_id 在 LineType 为 proxy 时才需要。Modify 接口支持修改 line_type、line_address、proxy_id、rule_id。

### 3. CRUD 函数设计
**决策**:
- Create: 调用 CreateMultiPathGatewayLine，返回 line_id，设置复合 ID
- Read: 解析复合 ID，调用 DescribeMultiPathGatewayLine 获取线路详情，通过 ReadRetryTimeout 重试
- Update: 检测可变字段变化（line_type、line_address、proxy_id、rule_id），调用 ModifyMultiPathGatewayLine
- Delete: 解析复合 ID，调用 DeleteMultiPathGatewayLine

**理由**: 遵循标准 RESOURCE_KIND_GENERAL 模式，参考 igtm_strategy 资源实现。所有 API 调用均添加 retry 处理。

### 4. Read 函数中从 d.Get() 获取 ID 字段
**决策**: 在 Read、Update、Delete 函数中，不直接从 d.Id() 解析 ID，而是从 d.Get() 获取 zone_id、gateway_id，line_id 则从 d.Id() 解析。

**理由**: 根据代码生成要求，使用多个字段作为联合 id 时，在对应的 read、update、delete 方法中，不要直接从 d.Id() 方法中获取 id，而是从 d.Get() 方法中获取 id 的各个字段，作为 request 的一部分。但 line_id 在创建前不存在于 schema 中，因此需要从 d.Id() 解析。

## Risks / Trade-offs

- **[API 参数不一致风险]**: Create 和 Modify 接口的参数略有不同（Create 没有 LineId，Modify 有 LineId），需要注意区分。→ 通过在 Create 中从 response 获取 line_id 并设置到复合 ID 中来缓解
- **[线路类型约束]**: 不同线路类型（direct/proxy/custom）对修改和删除操作有不同的限制。→ 在 Terraform 层面不做限制，由云 API 返回错误信息
- **[proxy_id 和 rule_id 的条件必填]**: 当 line_type 为 proxy 时，proxy_id 和 rule_id 为必填，但 Terraform schema 中标记为 Optional。→ 在 Create/Update 函数中进行条件校验
