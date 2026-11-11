## Context

`tencentcloud_clb_instance` 资源（`tencentcloud/services/clb/resource_tc_clb_instance.go`）在 `resourceTencentCloudClbInstanceUpdate` 中，当 `sla_type` 发生变更时已调用云 API `ModifyLoadBalancerSla`（`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317`）。该接口的入参 `ModifyLoadBalancerSlaRequest.Force *bool` 含义为"是否强制升级，默认否"，vendor 中该字段已存在。但当前 Provider 在调用时未设置 `Force`，导致无法表达"强制升级"意图。

云 API 现状核查（基于 vendor）：
- `ModifyLoadBalancerSla` 接口存在，入参 `Force *bool` 存在（models.go line 12779-12795），无需更新依赖。
- 查询接口 `DescribeLoadBalancers` 返回的 `LoadBalancer` 结构（models.go line 10566+）**不包含** `Force` 字段，即 `Force` 仅是写操作参数，云侧不回显。
- 需求中提到的 `ModifyTags`（`request.tags` → `tags`）经核查在 clb 包中**不存在**，且资源已存在 `tags` 字段，故本次不涉及。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_clb_instance` 资源上新增可选参数 `force`（bool，默认 false），用于在调用 `ModifyLoadBalancerSla` 时透传"是否强制升级"标志。
- 保持向后兼容：现有不设置 `force` 的配置行为不变（默认 false）。
- 补充文档与单元测试。

**Non-Goals:**
- 不新增 `tags` 相关参数（云 API 不支持 `ModifyTags`，且资源已有 `tags`）。
- 不修改 `ModifyLoadBalancerSla` 调用以外的其他逻辑。
- 不改变 `sla_type` 的语义与 ForceNew/Computed 行为。
- 不引入新的 `_extension.go` 文件。

## Decisions

### 决策 1：`force` 参数定义为 Optional、非 Computed
**选择**：`force` 为 `schema.TypeBool`、`Optional: true`、不设 `Computed: true`，默认值 false。

**理由**：云 API `DescribeLoadBalancers` 返回的 `LoadBalancer` 结构不含 `Force` 字段，无法在 Read 中回填。`Force` 是一次性写操作参数（控制升级行为），不是持久化属性。设为非 Computed 可避免 Read 强行回填导致与用户配置产生 diff。

**替代方案**：若设为 `Computed`，需在 Read 中读取云 API 返回，但返回结构无此字段，会导致回填空值覆盖用户配置，产生非预期 diff。因此否决。

### 决策 2：`force` 仅在 `sla_type` 变更调用 `ModifyLoadBalancerSla` 时透传
**选择**：在 `resourceTencentCloudClbInstanceUpdate` 的 `if d.HasChange("sla_type")` 块内，构造 `slaRequest` 后增加 `slaRequest.Force = helper.Bool(d.Get("force").(bool))`，保留现有 retry 与异步任务等待逻辑。

**理由**：`Force` 是 `ModifyLoadBalancerSla` 的入参，只在触发 SLA 升级时有意义。现有 update 已只在 `sla_type` 变化时调用该接口，直接在该调用块透传即可，改动最小、语义最清晰。

**替代方案**：单独为 `force` 建立一个调用分支——不必要，因为 `force` 脱离 `sla_type` 升级场景无意义，单独调用会造成冗余 API 请求。否决。

### 决策 3：不在 Read 中回填 `force`
**选择**：`resourceTencentCloudClbInstanceRead` 不对 `force` 调用 `d.Set("force", ...)`。

**理由**：云 API 不返回 `Force`，无数据可回填。由于 `force` 非 Computed，Terraform 以用户配置为准，符合预期。

### 决策 4：Create 阶段不透传 `force`
**选择**：Create 中不调用 `ModifyLoadBalancerSla`，因此不涉及 `force`。`sla_type` 在 Create 时通过 `CreateLoadBalancer` 的 `SlaType` 入参传入（已有逻辑），与 `Force` 无关。

**理由**：`Force` 仅在"升级"操作（`ModifyLoadBalancerSla`）时有意义，创建时无升级概念。

## Risks / Trade-offs

- **[风险] 用户配置 `force = true` 但未变更 `sla_type`**：`force` 不会触发任何 API 调用，行为为空操作。→ **缓解**：文档明确说明 `force` 仅在 `sla_type` 变更升级时生效；非 Computed 保证不会无意义触发更新。
- **[风险] 强制升级不可回退**：云 API 文档指出"共享型升级为性能容量型实例后，不支持再回退到共享型实例"。用户设 `force=true` 升级后无法通过 Terraform 回退。→ **缓解**：在 schema Description 与文档中提示该不可逆性，引导谨慎使用。
- **[取舍] `force` 非 Computed 导致 import 后 state 无此字段**：import 后用户需在配置中显式声明 `force` 才能在后续升级生效。→ **取舍可接受**，因 `force` 本就是用户意图参数，import 场景默认 false 合理。
