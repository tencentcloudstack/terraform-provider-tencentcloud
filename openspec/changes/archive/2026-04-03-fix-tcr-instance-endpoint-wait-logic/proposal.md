# 变更提案：修复 tencentcloud_tcr_instance Update 中外部端点等待逻辑的位置错误

## 变更类型

**Bug 修复** — 将等待外部端点（External EndPoint）Ready 的逻辑从 `security_policy` 变更块移至 `open_public_operation` 变更块。

## Why

### 问题描述

`resourceTencentCloudTcrInstanceUpdate` 函数中存在逻辑位置错误：

**当前代码流程：**
```
1. open_public_operation 变更 → 调用 ManageTCRExternalEndpoint（Create/Delete）
                                ↓ 无等待，直接结束
2. security_policy 变更      → [等待 ExternalEndpoint 状态为 Opened] ← 位置错误
                                → 执行 security policy add/remove
```

**问题根因：**
- 等待外部端点 Ready（`DescribeExternalEndpointStatus` Retry 轮询）的逻辑放在了 `security_policy` 块的开头（line 504-526）
- 但这个等待实际上是为了确保 `open_public_operation=true` 触发的端口开启操作完成后再进行后续操作
- 如果用户**只变更** `security_policy`（不变更 `open_public_operation`），端点未必处于 Opening 状态，等待会返回 `unexpected external endpoint status` 错误
- 如果用户**只变更** `open_public_operation`（不变更 `security_policy`），操作完成后没有等待，后续依赖端点状态的操作可能失败

**正确流程：**
```
1. open_public_operation 变更 → 调用 ManageTCRExternalEndpoint
                                → [等待 ExternalEndpoint 状态为 Opened] ← 移至此处
2. security_policy 变更      → 直接执行 add/remove（无需等待）
```

### 影响

- 当用户单独修改 `security_policy` 时（不改 `open_public_operation`），若当前端点不是 Opening 状态会报错
- 当用户打开公网访问后立刻修改 `security_policy`，正好省去了等待，可能导致操作在端点未 Ready 时就执行

## What Changes

### 代码变更

**文件**: `tencentcloud/services/tcr/resource_tc_tcr_instance.go`

**修改内容**：将 `security_policy` 块中的等待逻辑（lines 503-526）移至 `open_public_operation` 块末尾（line 499 之后）。

#### 修改前结构（line 483-545）：

```go
if d.HasChange("open_public_operation") {
    // 调用 ManageTCRExternalEndpoint
    // 无等待
}

if d.HasChange("security_policy") {
    // ← 等待 ExternalEndpoint Opened（位置错误）
    err = resource.Retry(5*tccommon.ReadRetryTimeout, func() ...)

    // security policy add/remove
}
```

#### 修改后结构：

```go
if d.HasChange("open_public_operation") {
    operation = d.Get("open_public_operation").(bool)
    // 调用 ManageTCRExternalEndpoint
    
    // 仅当 operation=true（打开公网）时才需要等待
    if operation {
        err = resource.Retry(5*tccommon.ReadRetryTimeout, func() ...)
        // 等待 ExternalEndpoint Opened
    }
}

if d.HasChange("security_policy") {
    // 直接执行 security policy add/remove（无等待）
    o, n := d.GetChange("security_policy")
    // ...
}
```

### 影响范围

- **影响文件**：`tencentcloud/services/tcr/resource_tc_tcr_instance.go`（Update 函数）
- **破坏性变更**：无
- **向后兼容**：完全兼容，仅调整逻辑位置

### 向后兼容性

✅ 完全向后兼容：
- 同时修改两个字段时行为不变（先等待再修改 security_policy）
- 仅修改 `open_public_operation=true` 时新增了等待，更加健壮
- 仅修改 `security_policy` 时去掉了不合理的等待，修复了潜在报错
