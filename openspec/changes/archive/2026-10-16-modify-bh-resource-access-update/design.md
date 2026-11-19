## Context

`tencentcloud_bh_resource`（堡垒机服务实例）资源当前在 `tencentcloud/services/bh/resource_tc_bh_resource.go` 中实现。该资源的大多数"访问开关"类字段（`intranet_access`、`external_access`）已经支持在 update 阶段通过 Enable/Disable 接口修改并轮询 `DescribeResources` 等待生效；唯独 `web_access` 与 `client_access` 两个字段仍带有 `ForceNew: true`，意味着用户一旦想调整这两个开关就必须重建实例。而堡垒机资源（`Delete`）明确返回"不支持删除，需工单处理"，所以这两个字段在实际使用中形同"写死"。

腾讯云 bh v20230418 SDK 已经提供了独立的开关接口（vendor 中已存在并经核对）：

- `EnableWebAccess` / `DisableWebAccess`：入参仅 `ResourceId *string`。
- `EnableClientTcpAccess` / `DisableClientTcpAccess`：入参仅 `ResourceId *string`。

同时 `DescribeResources` 返回的 `Resource` 结构体中 `WebAccess`、`ClientAccess` 均为 `*uint64`，其取值含义为：`0` 关闭 / `1` 开启 / `2` 开通中 / `3` 关闭中。该中间态语义与已有的 `IntranetAccess`(`0/1/2/3`)、`ExternalAccess`(`0/1/2/3`) 完全一致，因此可以复用同资源中 `intranet_access`/`external_access` 的 update + wait 写法，不引入新的模式。

参考样板：本资源 `resourceTencentCloudBhResourceUpdate` 中第 556–779 行的 `intranet_access` / `external_access` 处理块即为本变更的直接对照实现。

## Goals / Non-Goals

**Goals:**

- 移除 `web_access`、`client_access` 两个字段 schema 的 `ForceNew: true`，使其可在 update 阶段修改。
- 在 `resourceTencentCloudBhResourceUpdate` 中新增 `web_access` 变更分支：值变为 `1` 调 `EnableWebAccess`，值变为 `0` 调 `DisableWebAccess`；调用成功后用 `DescribeResources` 轮询 `WebAccess` 字段直至稳定到目标值（`1` 或 `0`）。
- 在 `resourceTencentCloudBhResourceUpdate` 中新增 `client_access` 变更分支：值变为 `1` 调 `EnableClientTcpAccess`，值变为 `0` 调 `DisableClientTcpAccess`；调用成功后用 `DescribeResources` 轮询 `ClientAccess` 字段直至稳定到目标值。
- 复用既有 wait 模式：`resource.Retry(tccommon.ReadRetryTimeout*7, ...)`，与 `intranet_access`/`external_access` 一致。
- 保持向后兼容：字段类型、含义、Create/Deploy 行为不变，只放开 update 能力。
- 在 `resource_tc_bh_resource_test.go` 中补充 update 变更单元测试用例（使用 gomonkey mock 云 API）。

**Non-Goals:**

- 不修改 `web_access` / `client_access` 在 Create 阶段的 `DeployResource` 调用逻辑（仍由 deploy 参数下发）。
- 不调整 `intranet_access` / `external_access` 等其他字段的既有关逻辑。
- 不新增数据源、不新增资源、不修改 service 层。
- 不为该资源引入 `Timeouts` 块（既有实现未使用 Timeouts，保持一致；异步等待用 `ReadRetryTimeout*7` 的 retry 实现）。
- 不在 schema 中对 `web_access`/`client_access` 增加 `ValidateFunc` 枚举校验（与现有同类字段一致，不下沉）。

## Decisions

### D1 — Schema 变更：仅移除 ForceNew

`web_access` 与 `client_access` 的 schema 仅删除 `ForceNew: true` 一行，其余属性（`Type: schema.TypeInt`、`Optional: true`、`Computed: true`、`Description`）保持不变。

- 字段依然是 `Optional + Computed`：Create 阶段由用户显式提供或由云侧返回默认值；Read 阶段从 `DescribeResources` 回填。
- 不设 `Required`：与现有 `intranet_access`/`external_access` 字段保持一致，避免破坏既有配置。

**为什么不在 Create 阶段也调用 Enable/Disable？** Create 阶段这两个字段已经通过 `DeployResource` 请求的 `WebAccess`/`ClientAccess` 入参一次性下发（见现有代码第 292–298 行），无需另调开关接口；只在 update 阶段补齐开关能力即可。

### D2 — update 分支结构对齐 intranet_access/external_access

在 `resourceTencentCloudBhResourceUpdate` 中，于现有 `external_access` 处理块之后、`return resourceTencentCloudBhResourceRead(d, meta)` 之前，新增两个并列的 `if d.HasChange(...)` 分支：

```
if d.HasChange("web_access") {
    if v, ok := d.GetOkExists("web_access"); ok {
        if v.(int) == 1 {
            // 调 EnableWebAccess，入参 ResourceId
            // wait: DescribeResources 轮询 WebAccess == 1
        } else {
            // 调 DisableWebAccess，入参 ResourceId
            // wait: DescribeResources 轮询 WebAccess == 0
        }
    }
}

if d.HasChange("client_access") {
    if v, ok := d.GetOkExists("client_access"); ok {
        if v.(int) == 1 {
            // 调 EnableClientTcpAccess，入参 ResourceId
            // wait: DescribeResources 轮询 ClientAccess == 1
        } else {
            // 调 DisableClientTcpAccess，入参 ResourceId
            // wait: DescribeResources 轮询 ClientAccess == 0
        }
    }
}
```

**入参确认**：vendor SDK 中四个接口的 Request 结构体均只有 `ResourceId *string` 一个字段（已核对 `EnableWebAccessRequest`、`DisableWebAccessRequest`、`EnableClientTcpAccessRequest`、`DisableClientTcpAccessRequest`），因此无需像 `intranet_access` 那样额外传 `VpcId`/`SubnetId`/`VpcCidrBlock`，只需 `request.ResourceId = helper.String(resourceId)`。

### D3 — wait 轮询判定条件

参照 `intranet_access`：Enable 后等待字段变为 `1`（`*intranetAccess == 1` 返回 nil，否则 `RetryableError`），Disable 后等待字段变为 `0`。本变更对 `web_access`/`client_access` 采用相同判定：

- Enable：轮询 `result.Response.ResourceSet[0].WebAccess`（或 `ClientAccess`），`== 1` 即完成；中间态 `2`(开通中)/`3`(关闭中) 继续重试。
- Disable：轮询至 `== 0` 即完成。

空指针保护沿用既有写法：`result == nil || result.Response == nil || result.Response.ResourceSet == nil || len(...) == 0` → `NonRetryableError`；字段为 nil → `NonRetryableError`。

retry 超时使用 `tccommon.ReadRetryTimeout*7`，与同资源其他 wait 块一致。

### D4 — Create 阶段的 web_access/client_access 保持原样

Create 流程（`CreateResource` → `DeployResource` → wait Status==1 → Read）不变。`web_access`/`client_access` 仍作为 `DeployResource` 的入参下发。这样 update 新增的分支只在"实例已创建后用户修改这两个字段"时触发，符合 Terraform `d.HasChange` 语义。

### D5 — mutableArgs 不动

现有 `mutableArgs := []string{"resource_edition", "resource_node", "auto_renew_flag"}` 仅控制 `ModifyResource` 调用。`web_access`/`client_access` 走独立的 Enable/Disable 接口而非 `ModifyResource`，因此**不**将它们加入 `mutableArgs`，避免误调 `ModifyResource`（该接口也没有这两个字段）。各访问开关字段（`intranet_access`/`external_access`/`web_access`/`client_access`）都各自独立 `d.HasChange` 判定，互不影响。

### D6 — 测试方式：gomonkey mock 云 API

按项目规范，对**新增字段 update 能力**在 `resource_tc_bh_resource_test.go` 中补充测试用例，使用 gomonkey 对云 API client 方法（`EnableWebAccessWithContext`、`DisableWebAccessWithContext`、`EnableClientTcpAccessWithContext`、`DisableClientTcpAccessWithContext`、`DescribeResourcesWithContext`）进行 mock，仅做业务逻辑单测，不使用 terraform 测试套件。参照同文件已有测试风格。

## Risks / Trade-offs

- **[中间态轮询超时]** `WebAccess`/`ClientAccess` 在开通/关闭期间可能停留在 `2`/`3` 较久，若超过 `ReadRetryTimeout*7` 会导致 update 失败。→ 与 `intranet_access`/`external_access` 完全相同的既有行为，沿用既有超时倍数，不单独放大；若用户环境确实较慢，可后续统一调整。
- **[并发修改同一实例]** 用户在同一次 `terraform apply` 中同时改 `web_access` 和 `client_access`，两个分支会串行执行（先后调各自接口 + 各自 wait），不存在并发冲突；等待均基于 `DescribeResources` 的同名字段，互不影响。→ 维持串行，不做合并。
- **[向后兼容]** 移除 ForceNew 后，旧 state 中的 `web_access`/`client_access` 值不变，首次 `terraform plan` 不会产生 diff（值未变）。→ 仅当用户主动修改时才触发 update 分支，无静默变更风险。
- **[Enable/Disable 接口入参极简]** 四个接口仅接收 `ResourceId`，无法在请求中携带目标值，依赖云侧根据接口名执行对应动作。→ 这是云 API 设计，本变更严格遵循；通过选对接口名（Enable→开启，Disable→关闭）保证语义正确。