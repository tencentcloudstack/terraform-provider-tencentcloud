## Context

TEO (TencentCloud EdgeOne) 是腾讯云的边缘安全加速平台，当前 Terraform Provider 已支持 TEO 的站点管理、源站组、四层代理、加速域名等资源，但尚未支持负载均衡实例(Load Balancer)资源的管理。

负载均衡实例是 TEO 站点加速场景中的核心组件，用于实现源站组间的流量调度、健康检查和故障转移。云 API 已提供完整的 CRUD 接口：
- `CreateLoadBalancer`：创建负载均衡实例
- `DescribeLoadBalancerList`：查询负载均衡实例列表（支持按 InstanceId/InstanceName 过滤）
- `ModifyLoadBalancer`：修改负载均衡实例配置
- `DeleteLoadBalancer`：删除负载均衡实例

现有代码模式参考：`tencentcloud/services/teo/resource_tc_teo_origin_group.go`（同为 teo 下的资源，使用 ZoneId 作为联合 ID 的一部分）

## Goals / Non-Goals

**Goals:**
- 新增 `tencentcloud_teo_load_balancer` 资源，支持负载均衡实例的完整 CRUD 生命周期
- 支持所有 Create/Modify 接口参数的配置（zone_id, name, type, origin_groups, health_checker, steering_policy, failover_policy）
- 支持资源导入（Import）
- 使用 ZoneId + InstanceId 作为联合 ID（FILED_SP 分隔）
- 生成完整的资源文档和单元测试
- 在 provider.go 和 provider.md 中正确注册资源

**Non-Goals:**
- 不新增数据源(data source)，本变更仅处理通用资源
- 不支持负载均衡实例的状态查询（OriginGroupHealthStatus、L4UsedList、L7UsedList 等只读字段作为 computed 属性）
- 不修改已有资源的 schema 或行为

## Decisions

### 1. 资源 ID 使用联合 ID（ZoneId:InstanceId）
**决策**：使用 `tccommon.FILED_SP` 分隔的 `ZoneId:InstanceId` 作为资源 ID
**理由**：TEO 的 LoadBalancer API 均要求传入 ZoneId 和 InstanceId，这是 TEO 资源的标准模式（参考 teo_origin_group）。Import 时用户需要提供联合 ID。
**替代方案**：仅使用 InstanceId 作为 ID —— 不可行，因为 DescribeLoadBalancerList 等接口均需要 ZoneId 参数。

### 2. Read 方法使用 DescribeLoadBalancerList + Filter
**决策**：Read 操作通过 DescribeLoadBalancerList 接口并使用 InstanceId Filter 来查询单个实例
**理由**：云 API 未提供 DescribeLoadBalancer 单实例查询接口，只能通过列表接口 + Filter 实现。Limit 设为 1，Filter 使用 InstanceId 精确匹配。
**替代方案**：无单实例查询接口可用。

### 3. Type 字段设置为 ForceNew
**决策**：Type（实例类型）字段设置为 ForceNew: true
**理由**：ModifyLoadBalancer 接口不支持修改 Type 字段（该字段不在 Update 接口的参数列表中），因此类型变更需要重建资源。

### 4. HealthChecker 使用嵌套 TypeList 结构
**决策**：health_checker 字段使用 TypeList（MaxItems: 1）的嵌套结构
**理由**：云 API 的 HealthChecker 是一个对象而非数组，使用 TypeList + MaxItems:1 可以更好地表达一对一关系。参考 teo_origin_group 中类似嵌套结构的设计。

### 5. OriginGroups 使用 TypeList 结构
**决策**：origin_groups 字段使用 TypeList，包含 priority 和 origin_group_id 两个必需子字段
**理由**：OriginGroupInLoadBalancer 结构体包含 Priority 和 OriginGroupId，需要保持顺序（优先级），TypeList 比 TypeSet 更合适。

### 6. Computed 只读属性
**决策**：LoadBalancer 结构体中的 Status、OriginGroupHealthStatus、L4UsedList、L7UsedList、References 设为 computed 只读属性
**理由**：这些字段仅在 DescribeLoadBalancerList 响应中返回，不支持通过 Create/Modify 接口设置。

### 7. Name 字段在 Create 中必填，在 Modify 中可选
**决策**：name 在 Terraform schema 中设为 Required，Modify 时传入当前值
**理由**：CreateLoadBalancer 接口要求 Name 必填，ModifyLoadBalancer 中 Name 可选（不填则维持原有配置）。Terraform 侧设为 Required 可以确保状态一致性。

## Risks / Trade-offs

- **[DescribeLoadBalancerList 分页查询]** → 使用 InstanceId Filter + Limit=1 精确查询，避免分页问题
- **[Type 字段不可变更]** → 设置 ForceNew: true，类型变更时重建资源，用户需注意销毁重建的影响
- **[CreateLoadBalancer 响应异步]** → Create 后需要 Read 轮询确认资源生效，使用 tccommon.ReadRetryTimeout 进行重试
- **[ModifyLoadBalancer 响应异步]** → Update 后需要 Read 轮询确认配置生效
