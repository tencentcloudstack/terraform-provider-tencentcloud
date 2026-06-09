# 变更提案：tencentcloud_sqlserver_account 新增自定义 Timeouts 功能

## 变更类型

**功能增强** — 为 `tencentcloud_sqlserver_account` 资源的 Create 和 Delete 操作增加自定义超时配置。

## Why

当前代码中 Create 和 Delete 的等待逻辑使用硬编码的全局常量：

| 操作 | 当前使用常量 | 实际时长 |
|------|------------|--------|
| Create | `tccommon.WriteRetryTimeout` | 5 分钟 |
| Delete（确认存在）| `tccommon.ReadRetryTimeout` | 3 分钟 |
| Delete（执行删除）| `tccommon.WriteRetryTimeout` | 5 分钟 |
| Delete（等待消失）| `tccommon.ReadRetryTimeout` | 3 分钟 |

用户无法根据实际环境（如大实例、网络慢等）自定义等待时间，遇到超时只能修改全局环境变量，影响所有资源。

## What Changes

### Timeouts 块

```hcl
timeouts {
  create = "30m"  # 默认 30 分钟
  delete = "10m"  # 默认 10 分钟
}
```

### 修改位置

| 文件 | 修改内容 |
|------|---------|
| `tencentcloud/services/sqlserver/resource_tc_sqlserver_account.go` | 新增 `Timeouts` 块（Create 默认 30m，Delete 默认 10m）；Create/Delete 中的 retry 超时改用 `d.Timeout(schema.TimeoutCreate/Delete)` |

### 各操作等待逻辑调整

**Create**：
- `CreateSqlserverAccount` retry → `d.Timeout(schema.TimeoutCreate)`

**Delete**：
- 确认存在的 retry → `d.Timeout(schema.TimeoutDelete)`
- `DeleteSqlserverAccount` retry → `d.Timeout(schema.TimeoutDelete)`
- 等待删除完成的 retry → `d.Timeout(schema.TimeoutDelete)`

### 向后兼容性

✅ 完全向后兼容：不设置 `timeouts` 块时使用默认值（30m/10m），行为符合预期。
