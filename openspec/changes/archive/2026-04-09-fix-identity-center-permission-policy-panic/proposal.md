# 变更提案：修复 identity_center_role_configuration_permission_policy_attachment Read 函数 nil pointer panic

## 变更类型

**紧急 Bug 修复（P0）** — 修复 `resourceTencentCloudIdentityCenterRoleConfigurationPermissionPolicyAttachmentRead` 函数中的 nil pointer dereference，该 panic 会导致 provider 进程崩溃。

## Why

### Panic 根因

堆栈指向 `resource_tc_identity_center_role_configuration_permission_policy_attachment.go:206`：

```go
// 第 197~221 行
if respData.RolePolicies != nil {
    var rolePolicie *organization.RolePolicie   // 初始值为 nil
    for _, r := range respData.RolePolicies {
        if *r.RolePolicyId == rolePolicyId {
            rolePolicie = r
            break
        }
    }
    // 问题：若循环结束后未找到匹配项，rolePolicie 仍为 nil
    if rolePolicie.RolePolicyName != nil {       // ← 第 206 行：nil pointer dereference!
```

**触发条件**：当 `respData.RolePolicies` 不为空，但其中没有 `RolePolicyId` 等于当前资源的 `rolePolicyId` 的条目时（如策略已被手动解绑但资源尚未删除），循环结束后 `rolePolicie` 仍为 `nil`，对其成员字段直接访问即触发 panic，导致 provider 进程崩溃。

### 同类隐患

同一个 for 循环中还有另一处潜在问题：
- `*r.RolePolicyId` 在 `r.RolePolicyId == nil` 时也会 panic

## What Changes

### 修复方案

**方案 A（推荐）**：在循环结束后判断 `rolePolicie == nil`，将其处理为资源已不存在（`d.SetId("")`，与其他资源保持一致）。

**方案 B**：仅加 nil 守卫继续执行。不推荐，因为找不到对应策略说明该 attachment 已失效，应该清除 state。

采用方案 A：找不到匹配条目 → 视为资源已删除 → `d.SetId("")` + return nil。

同时修复 `*r.RolePolicyId` 的 nil 解引用风险，访问前先判断 `r.RolePolicyId != nil`。

### 修改位置

| 文件 | 行数 | 修改内容 |
|------|------|---------|
| `tencentcloud/services/tco/resource_tc_identity_center_role_configuration_permission_policy_attachment.go` | 197-222 | 循环中加 `r.RolePolicyId != nil` 守卫；循环后加 `rolePolicie == nil` 判断并 SetId("") |
