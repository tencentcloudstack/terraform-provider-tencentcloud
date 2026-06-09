## Context

`tencentcloud_cdn_domain` 是 CDN 域名管理资源，支持数十个配置字段。当前 Read 函数中存在两类守卫逻辑，导致 import 时字段无法写入 state：

1. **`checkCdnInfoWritable` 守卫**（影响 19 个字段）：该函数通过 `helper.InterfacesHeadMap(d, key)` 检查 state 中是否已存在该 key 的 map 数据。import 时 state 为空，`ok` 始终为 false，所有受保护字段均被跳过。
2. **`d.GetOk` 守卫**（影响 5 个 switch 字段）：`response_header_cache_switch`、`seo_switch`、`video_seek_switch`、`offline_cache_switch`、`quic_switch` 在 Read 中使用 `if _, ok := d.GetOk("xxx"); ok && dc.Xxx != nil` 判断，import 时 state 为空，`GetOk` 返回 false，字段被跳过。

约束：不能破坏存量用户的正常 apply 流程（向后兼容）。

## Goals / Non-Goals

**Goals:**
- import 后所有 API 可返回的字段均能正确写入 state
- 存量用户正常 apply 流程不受影响
- 不修改 schema 定义

**Non-Goals:**
- 修复 `https_config`（含私钥，API 不返回）
- 修复 `full_url_cache`（Deprecated bool 字段，需要特殊转换逻辑）
- 修复 `authentication`（含鉴权密钥，API 不返回）
- 新增字段或修改现有字段类型

## Decisions

### 决策 1：修改 `checkCdnInfoWritable` 而非在每个调用处单独修改

**选择**：将 `checkCdnInfoWritable` 改为 `return val != nil`，一次性修复所有 19 个受影响字段。

**理由**：该函数的原始设计意图是"只有用户配置了该字段才写入"，但这与 Terraform 的 import 语义冲突——import 应该读取资源的完整状态。修改函数本身比在 19 处调用点逐一修改更安全、更一致。

**向后兼容性**：存量用户 apply 时，state 中已有这些字段的值，API 返回非 nil，`val != nil` 为 true，行为与之前完全一致（之前 `ok` 也为 true）。

**备选方案**：在每个调用处添加 `|| d.Id() != ""` 判断（import 时 ID 已设置）。但这会使代码更复杂，且语义不清晰。

### 决策 2：直接移除 5 个 switch 字段的 `d.GetOk` 守卫

**选择**：将 `if _, ok := d.GetOk("xxx"); ok && dc.Xxx != nil` 改为 `if dc.Xxx != nil`。

**理由**：这 5 个字段的 `d.GetOk` 守卫没有实际业务意义——只要 API 返回了该字段，就应该写入 state。守卫的存在只会在 import 时造成数据丢失。

## Risks / Trade-offs

- **[风险] 某些字段 API 返回默认值与用户配置不同** → 缓解：这是 Terraform provider 的标准行为，用户可通过 `ignore_changes` 处理；且该问题在修改前同样存在（只是 import 时被掩盖）
- **[风险] `response_header` 和 `status_code_cache` 字段格式不完全匹配** → 缓解：这两个字段不受本次修改影响（它们已经使用 `if dc.Xxx != nil` 直接写入），保留在 `ImportStateVerifyIgnore` 中
