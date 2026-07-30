## Context

`BindDeviceAccountKubeconfig`（bh v20230418）是腾讯云堡垒机的"配置型"接口：它对一个**已经存在**的容器账号（由 `Id` 标识）绑定 kubeconfig 凭据。本资源不创建容器账号本身（容器账号由其他流程/资源创建），只对账号附加 kubeconfig。

接口语义关键点：

- 入参 3 个：`Id` (uint64, required), `Kubeconfig` (string, required), `ManageDimension` (uint64, optional, 当前枚举仅 1=集群)。
- 出参仅有 `RequestId`，无业务返回字段。
- API 是"覆盖式绑定"：再次调用同一个 `Id` 会用新 `Kubeconfig` 覆盖旧的；这意味着 Create 与 Update 共用同一个 SDK 调用。
- 当前业务侧**没有**对应查询接口（如 `DescribeDeviceAccountKubeconfig`），也**没有**解绑接口（如 `UnbindDeviceAccountKubeconfig`）。

参考：本 change 的代码风格严格对齐 `tencentcloud_waf_owasp_rule_status_config` —— 一个同样的"无独立查询/无独立删除"配置型资源，在 Provider 内已有标准实现样板（Create→SetId+Update、Update 真正调 SDK、Delete return nil）。

## Goals / Non-Goals

**Goals:**

- Schema 字段名与 `BindDeviceAccountKubeconfig` 接口入参 1:1 映射（snake_case 化）。
- 资源 ID = `fmt.Sprintf("%d", id)`，由 HCL `id` 字段（容器账号 Id）转换而来，与接口入参对齐。
- Create 与 Update 共用同一个 SDK 调用 `BindDeviceAccountKubeconfigWithContext`；按参考资源习惯，Create 仅 `d.SetId(...)` 后转调 Update。
- Read 直接返回 nil（用户明确要求；业务暂无查询接口）。
- Delete 直接返回 nil（业务暂无解绑接口）。
- 全部 SDK 调用包裹 `resource.Retry(WriteRetryTimeout, ...)`。
- 全部接口返回值取值（包括 `result.Response`）做空指针保护。
- `kubeconfig` 字段标 `Sensitive: true`（凭据敏感）。
- 资源 doc 命名 `resource_tc_bh_bind_device_account_kubeconfig.md`，测试 `resource_tc_bh_bind_device_account_kubeconfig_test.go`，与参考样板一致。

**Non-Goals:**

- 不实现配套数据源（`tencentcloud_bh_bind_device_account_kubeconfigs`）—— 没有查询接口，无法实现。
- 不实现 import：没有查询接口，import 后 Read 会返回空 state，与 HCL 不一致。后续若 SDK 增加查询接口，再单独 PR 补 import。
- 不在 schema 加 `manage_dimension` 的强枚举校验：当前 API 文档只声明 `1=集群` 一个值，但后续可能扩展，不做收紧。
- 不修改任何既有资源/数据源/service 方法。

## Decisions

### D1 — Schema 字段映射

| HCL 字段 | SDK 字段 | 类型 | 必填 | ForceNew | Computed | Sensitive |
|---|---|---|---|---|---|---|
| `account_id` | `Id` | `TypeInt` | Required | **Yes** | No | No |
| `kubeconfig` | `Kubeconfig` | `TypeString` | Required | No | No | **Yes** |
| `manage_dimension` | `ManageDimension` | `TypeInt` | Optional | No | No | No |

**理由**：
- `account_id` 是绑定的目标账号身份键；改 `account_id` 等于"换一个账号绑定 kubeconfig"，必须 ForceNew。
- HCL 字段名为 `account_id` 而非 `id`，因为 **`id` 是 Terraform Plugin SDK v2 顶层 schema 的保留字段**（每个 `*schema.ResourceData` 内置 Resource ID）。`account_id` 在 Create 中通过 `d.SetId(fmt.Sprintf("%d", v.(int)))` 投影到 Resource ID。`account_id ↔ Id` 的映射关系在文档 NOTE 中明确说明。
- `kubeconfig` 与 `manage_dimension` 在同一个账号上可以多次覆盖式更新，所以 No-ForceNew、走 Update 重新 Bind。
- `kubeconfig` 是访问 K8s 集群的完整凭据（含 token / cert），按 Provider 内其他敏感字段（如 cls 资源里的 password / secret_key）的惯例标 Sensitive。

**字段类型选择**：尽管 SDK 是 `*uint64`，HCL 用 `TypeInt`（按 igtm / waf 等同类资源惯例），转换时用 `helper.IntUint64(v.(int))`。

### D2 — Create 走 "SetId + Update" 双跳

参考 `tencentcloud_waf_owasp_rule_status_config` 的写法（其 Create 仅 SetId 后 `return resourceTencentCloudWafOwaspRuleStatusConfigUpdate(d, meta)`），本资源同样：

```go
func resourceTencentCloudBhBindDeviceAccountKubeconfigCreate(d, meta) error {
    if v, ok := d.GetOkExists("account_id"); ok {
        d.SetId(fmt.Sprintf("%d", v.(int)))
    }
    return resourceTencentCloudBhBindDeviceAccountKubeconfigUpdate(d, meta)
}
```

**理由**：
- 因 Create 与 Update 都调用同一个 `BindDeviceAccountKubeconfig`，避免 SDK 调用代码两份重复。
- 与参考资源风格保持完全一致，便于后续维护 / code review。

### D3 — Read 直接 return nil

本资源没有对应查询接口。Read 不调任何 API，**也不**清空 state；返回 nil。

**理由**：
- 用户明确要求："当前暂无查询接口，所以 read 模块可以直接 return nil"。
- 这是配置型资源在无查询接口时的标准模式（参考 cls_alarm 早期实现、tag 类资源等）。
- 行为后果：state 中的字段值仅由 Create / Update 写入；plan 永远不会因为外部 drift 触发 diff（drift 不可见）。

### D4 — Delete 直接 return nil

API 未提供 `UnbindDeviceAccountKubeconfig` 等独立解绑接口。Delete 仅在 Terraform state 中移除资源，不调任何 API。

**理由**：
- 业务上"删除绑定"的语义需要后端配合，目前无接口。强行调用 BindDeviceAccountKubeconfig 传空 kubeconfig 不是 API 设计意图。
- 与 `tencentcloud_waf_owasp_rule_status_config` 的 Delete 实现一致（其 Delete 也是 `return nil`）。
- 用户在 `terraform destroy` 后，堡垒机后端的绑定关系仍然保留；如需真正解绑，需要在控制台或调用其他业务接口。该行为应在 `.md` 文档中明确说明。

### D5 — Update 全字段重传

Update 路径中无论 `d.HasChange` 命中哪个字段，都把当前 HCL 的全部字段重新填入请求并调用一次 `BindDeviceAccountKubeconfigWithContext`：

```go
request.Id = helper.IntUint64(<account_id>)
request.Kubeconfig = helper.String(<kubeconfig>)
if v, ok := d.GetOkExists("manage_dimension"); ok {
    request.ManageDimension = helper.IntUint64(v.(int))
}
```

**理由**：API 是覆盖式语义，一次调用即可。无需做精细化 diff。

### D6 — 文档与测试命名

- `resource_tc_bh_bind_device_account_kubeconfig.md`：含 1 个 HCL example + 1 段 NOTE 说明"无查询/无解绑"语义。
- `resource_tc_bh_bind_device_account_kubeconfig_test.go`：基础 Create + Update 两步验证（无 ImportState 步骤，因为没有查询接口，import 后 plan 必然漂移），包名 `bh_test`，测试函数 `TestAccTencentCloudBhBindDeviceAccountKubeconfigResource_basic`。

均与 `resource_tc_config_compliance_pack.md` / `_test.go` 同款结构（区别仅在去掉 import 步骤）。

### D7 — Provider 注册位置

在 `tencentcloud/provider.go` 既有 bh 资源注册段（紧邻其他 `bh.Resource...`）追加：

```go
"tencentcloud_bh_bind_device_account_kubeconfig": bh.ResourceTencentCloudBhBindDeviceAccountKubeconfig(),
```

同时在 `tencentcloud/provider.md` 的 Bastion Host(BH) Resource 列表追加资源名一行，确保 gendoc 能扫描并生成 website doc。

## Risks / Trade-offs

- **Risk**: 用户改了 kubeconfig 后没及时 apply，控制台又被人改回旧 kubeconfig → Terraform 看不到 drift → 下次 apply 仍以为是旧 state，把 HCL 写的"未变化"值不重传 → state 与后端不一致 → Mitigation: 这是无查询接口的固有限制，文档明确说明"drift 不可见"。建议用户在变更 kubeconfig 后立即 apply。
- **Risk**: Delete 后重建（同 `id`）会重新绑定，与之前的状态关联清晰；但跨 workspace 同时管理同一个 `id` 会发生覆盖竞争 → Mitigation: 文档说明"同一个容器账号同时只能由一份 HCL 管理"。
- **Trade-off**: 不实现 import。理由：无查询接口 = import 后 Read 拿不到字段，state 必然为空；强行 import 会导致下次 plan 全字段重建，体验更差。后续若 SDK 增加查询接口，再做 import 支持。
- **Trade-off**: Read 不调 API = 完全信任 state。这个权衡符合用户明确要求与配置型资源的惯例，但牺牲了 drift detection。
