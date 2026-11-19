## Context

`tencentcloud_teo_load_balancer` 资源封装了 TEO Load Balancer 的 CRUD 操作。当前 `health_checker` 嵌套块暴露了健康检查策略的各项参数（type/port/interval/timeout 等），但尚未包含云 API 新增的 `ProbeCluster`（探测集群区域）字段。同时顶层 `type` 字段新增了 `HTTP_V2`（HTTP 专用型 V2）取值，该取值支持选择发起探测的区域。

当前状态：
- `resource_tc_teo_load_balancer.go` 中 schema 已声明 `health_checker` 为 TypeList（MaxItems: 1）嵌套块，含 14 个子字段
- `buildTeoHealthChecker` 函数负责 schema → `*teo.HealthChecker` 转换（Create/Update 共用）
- `flattenTeoHealthChecker` 函数负责 `*teo.HealthChecker` → schema 转换（Read 使用）
- 顶层 `type` 字段当前 Description 仅包含 `HTTP` 和 `GENERAL` 两个取值，缺少 `HTTP_V2`

vendor SDK 已新增字段（`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`）：

`HealthChecker` struct（models.go:17651-17693）末尾新增：
```go
// 健康检查探测集群。指定本负载均衡实例发起健康探测的集群区域范围...仅 HTTP 专用型（V2）版本负载均衡实例支持设置。取值有：
// global：由全球所有区域的探测集群发起探测...
// mainland_china：仅由中国大陆区域的探测集群发起探测...
// overseas：仅由非中国大陆区域的探测集群发起探测...
// 默认为 global。
ProbeCluster *string `json:"ProbeCluster,omitnil,omitempty" name:"ProbeCluster"`
```

`CreateLoadBalancerRequest`（models.go:4224-4247）的 `HealthChecker` 字段为 `*HealthChecker` 类型，因此 `ProbeCluster` 通过该嵌套结构传入。

`ModifyLoadBalancerRequest`（models.go:21553-21579）的 `HealthChecker` 字段同为 `*HealthChecker` 类型，`ProbeCluster` 同样可用。

`LoadBalancer` struct（models.go:18996-19029，DescribeLoadBalancerList 响应元素）的 `HealthChecker` 字段为 `*HealthChecker` 类型，Read 时可通过 `respData.HealthChecker.ProbeCluster` 读取。

`LoadBalancer.Type` 字段（models.go:19003）及 `CreateLoadBalancerRequest.Type` 字段（models.go:4233）的注释均已更新为包含 `HTTP_V2`、`HTTP`、`GENERAL` 三个取值的完整描述。

约束：
- 向后兼容：`probe_cluster` 为新增 Optional 字段，未配置时不应触发 plan diff
- `ProbeCluster` 仅在 `type=HTTP_V2` 时生效；其他类型实例设置该字段由云 API 服务端校验，provider 层面不做强制限制（与现有字段风格一致）
- `type` 字段为 ForceNew，描述更新不改变 schema 行为

## Goals / Non-Goals

**Goals:**
- 在 `health_checker` 嵌套块中新增 `probe_cluster` 可选字段，支持声明式配置探测集群区域
- 在 `buildTeoHealthChecker` / `flattenTeoHealthChecker` 中补充该字段的读写逻辑
- 更新 `type` 字段 Description，包含 `HTTP_V2` 取值说明
- 通过单元测试覆盖新增字段

**Non-Goals:**
- 不改变 Create/Read/Update/Delete 的主流程结构（仅在现有 health_checker 转换函数中新增一个字段分支）
- 不新增独立的探测集群资源
- 不在 schema 层面用 ValidateFunc 限制 `probe_cluster` 取值（由云 API 服务端校验）
- 不修改 `type` 字段的 ForceNew 属性或其他顶层字段

## Decisions

### Decision 1: `probe_cluster` 作为 `health_checker` 嵌套块的可选子字段

**选择**：在 `health_checker` 嵌套块 schema 中新增 `probe_cluster`（TypeString, Optional），与云 API `HealthChecker.ProbeCluster` 一一对应。

**备选**：作为顶层字段。

**理由**：
- `ProbeCluster` 在云 API 中属于 `HealthChecker` 结构的子字段，语义上是健康检查策略的一部分
- 与现有 `health_checker` 嵌套块的其他字段（type/port/interval 等）结构一致，无需引入额外层级

### Decision 2: `buildTeoHealthChecker` 中仅在值非空时设置 `ProbeCluster`

**选择**：在 `buildTeoHealthChecker` 中使用 `if val, ok := m["probe_cluster"].(string); ok && val != "" { hc.ProbeCluster = helper.String(val) }` 模式，与现有所有字段保持一致。

**理由**：
- 未配置时传 nil，云 API 使用默认值 `global`，保持向后兼容
- 与资源内现有字段的写入模式完全一致

### Decision 3: `flattenTeoHealthChecker` 中仅在值非 nil 时回写

**选择**：在 `flattenTeoHealthChecker` 中使用 `if hc.ProbeCluster != nil { m["probe_cluster"] = hc.ProbeCluster }` 模式。

**理由**：
- 遵循资源 Read 的 nil 检查规范，避免写入 nil 指针
- 与现有所有字段的读取模式完全一致

### Decision 4: `type` 字段描述更新为包含 `HTTP_V2` 的完整说明

**选择**：更新 `type` 字段的 Description，采用用户提供的新版释义：`HTTP_V2`（HTTP 专用型 V2，支持选择探测区域，灵敏度更低）、`HTTP`（HTTP 专用型 V1，不支持选择探测区域，灵敏度更高）、`GENERAL`（通用型，不支持选择探测区域，灵敏度更高）。

**理由**：
- 与云 API 最新文档保持一致，使 `probe_cluster` 字段的有效前置条件（仅 V2 支持）在描述中可见
- 描述更新不影响 schema 行为（仍为 Required + ForceNew + TypeString）

### Decision 5: 不使用 ValidateFunc 限制 `probe_cluster` 取值

**选择**：schema 层面不添加 `ValidateFunc`，合法取值（`global`/`mainland_china`/`overseas`）仅在 Description 中说明。

**理由**：
- 与资源内现有字段的风格一致（如 `steering_policy`、`failover_policy` 均无 ValidateFunc）
- 减少与未来 API 取值扩展的耦合，错误由云 API 返回即可

## Risks / Trade-offs

- **Risk**：`probe_cluster` 仅在 `type=HTTP_V2` 时生效，用户在 `GENERAL`/`HTTP` 类型实例上设置该字段可能被云 API 忽略或报错 → **Mitigation**：Description 中明确说明仅 V2 支持，错误由服务端返回并经 `tccommon.RetryError` 透传给用户
- **Risk**：存量 state 中无 `probe_cluster` 字段，升级后 Read 回填云端默认值 `global` 可能产生 plan diff → **Mitigation**：字段为 Optional 且仅在非 nil 时回写，若用户未在配置中声明则不产生 diff（Optional 字段标准行为）
- **Trade-off**：`type` 描述更新使文档变长，但提高了 `HTTP_V2` 与 `probe_cluster` 关联关系的可见性，可接受

## Migration Plan

- 新增字段为纯加法（Optional），无 state 迁移需求
- 存量资源：Terraform state 中无 `probe_cluster`，升级后 `terraform plan` 对未在 HCL 配置该字段的资源不产生 diff
- 文档更新：由收尾阶段 `make doc` 自动生成
- 回滚：若需回退，移除 schema 中 `probe_cluster` 字段及 build/flatten 中的对应分支即可，state 值不会丢失（字段移除后 Terraform 自动忽略）

## Open Questions

- 无
