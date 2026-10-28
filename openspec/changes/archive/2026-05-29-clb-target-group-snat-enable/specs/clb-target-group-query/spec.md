## ADDED Requirements

### Requirement: REQ-CLB-TG-QUERY-005 - 目标组资源支持 SnatEnable 参数

The `tencentcloud_clb_target_group` resource MUST expose a `snat_enable` argument (boolean, optional + computed) that maps to the `SnatEnable` field of the TencentCloud CLB SDK in the `CreateTargetGroup`, `DescribeTargetGroupList`, and `ModifyTargetGroupAttribute` APIs.

`tencentcloud_clb_target_group` 资源必须暴露 `snat_enable` 参数（Bool 类型，Optional + Computed），并将其分别映射到腾讯云 CLB SDK `CreateTargetGroup`、`DescribeTargetGroupList`、`ModifyTargetGroupAttribute` 三个 API 的 `SnatEnable` 字段，以打通源 IP 替换（SNAT）能力的 IaC 管理。

**理由**：
- SDK 三个 API 已原生支持 `SnatEnable`，Provider 仅需透传。
- 用户当前必须借助控制台或外部脚本设置 SNAT，不利于自动化闭环。
- 仅新增 Optional 字段，向后兼容。

**约束**：
- 不修改任何已有字段语义。
- 不引入新的 SDK 依赖。
- `snat_enable` 不设置 `Default`，遵循云端默认行为。
- `snat_enable` 不设 `ForceNew`，支持在线修改。

#### Scenario: 创建目标组时启用 SNAT

- **WHEN** 用户在 HCL 中配置 `tencentcloud_clb_target_group` 资源并设置 `snat_enable = true`，执行 `terraform apply`
- **THEN**
  - Provider 将 `SnatEnable = true` 传入 `CreateTargetGroupRequest`
  - 云端创建的目标组开启 SNAT
  - Terraform state 中 `snat_enable` = `true`
  - `terraform plan` 再次执行无 diff

**验收标准**：
- ✅ API 请求体中包含 `"SnatEnable": true`
- ✅ `terraform state show` 输出 `snat_enable = true`
- ✅ 无重复 plan diff

---

#### Scenario: 创建目标组时未设置 snat_enable

- **WHEN** 用户未在 HCL 中声明 `snat_enable`，执行 `terraform apply`
- **THEN**
  - Provider 不在 `CreateTargetGroupRequest` 中携带 `SnatEnable` 字段
  - 云端按默认行为处理（默认关闭）
  - Read 阶段从 `DescribeTargetGroupList` 回写真实值到 state（通过 Computed 能力）
  - 后续 `terraform plan` 不产生意料之外的 diff

**验收标准**：
- ✅ 创建请求体不包含 `SnatEnable` 字段
- ✅ Read 后 state 中 `snat_enable` 与云端一致
- ✅ 存量未设置 `snat_enable` 的 TF 配置升级后无破坏性变更

---

#### Scenario: Read 流程回写 SnatEnable

- **WHEN** Terraform 触发 `tencentcloud_clb_target_group` 的 Read 操作（refresh / plan）
- **THEN**
  - `ClbService.DescribeTargetGroupList` 返回的 `TargetGroupInfo.SnatEnable` 不为 nil 时
  - Provider 通过 `d.Set("snat_enable", *targetGroup.SnatEnable)` 写入 state
  - 与其他字段（`keepalive_enable` 等）共享相同的 nil 检查模式

**验收标准**：
- ✅ 通过控制台手动开启 SNAT 后，`terraform refresh` 能正确同步到 state
- ✅ 通过控制台手动关闭 SNAT 后，`terraform refresh` 能正确同步到 state
- ✅ API 异常返回 `SnatEnable = nil` 时，state 中字段保持为之前值不被错误清空

---

#### Scenario: Update 切换 SnatEnable 状态

- **WHEN** 用户修改已有 `tencentcloud_clb_target_group` 配置中的 `snat_enable`（如从 `true` 改为 `false`），执行 `terraform apply`
- **THEN**
  - `d.HasChange("snat_enable")` 返回 `true`
  - Provider 调用 `ModifyTargetGroupAttribute`，请求体包含 `SnatEnable = false`
  - 不触发资源重建（Update 而非 Replace）
  - state 同步为新值

**验收标准**：
- ✅ `terraform plan` 显示 `~ snat_enable: true -> false`，且为 in-place update
- ✅ API 调用为 `ModifyTargetGroupAttribute`，请求体包含 `SnatEnable`
- ✅ 资源 ID 不变

---

#### Scenario: Service 层 CreateTargetGroup 接口扩展

- **WHEN** 调用 `ClbService.CreateTargetGroup(...)` 创建目标组
- **THEN**
  - 函数签名末尾新增 `snatEnable *bool` 参数（紧随 `ipVersion string` 之后）
  - 当 `snatEnable != nil` 时，赋值给 `request.SnatEnable`
  - 当 `snatEnable == nil` 时，保持请求体中不携带该字段
  - 仓库内所有调用点（资源 Create 函数）同步更新参数列表

**验收标准**：
- ✅ `go build ./...` 通过
- ✅ `go vet ./...` 通过
- ✅ 没有调用点遗漏（grep `CreateTargetGroup(` 全部覆盖）
