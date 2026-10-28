## Context

`tencentcloud_teo_security_policy_config` 资源封装了 TEO 安全策略的 CRUD，其中 `bot_management.basic_bot_settings` 下的 `source_idc` / `search_engine_bots` / `known_bot_categories` / `ip_reputation.ip_reputation_group` 四个子块，各自包含一个 `action_overrides` 列表（TypeList，元素 schema 由 `botManagementActionOverrideSchema()` 返回）。

当前 `botManagementActionOverrideSchema()` 将 `rule_id` 定义为 `TypeString, Required`（单个 ID），对应的转换函数：
- `buildBotManagementActionOverrideFromMap`：从 schema 读单个 `rule_id` 字符串，构造 `[]*string{helper.String(v)}` 写入 `BotManagementActionOverrides.Ids`
- `flattenBotManagementActionOverride`：从 API 读 `override.Ids[0]` 写入 schema 的 `rule_id`

云 API 现状：`teov20220901.BotManagementActionOverrides.Ids` 类型为 `[]*string`，注释明确说明"Bot 规则组下的具体项"，支持同时改写多条规则项的处置动作。

问题：用户在 HCL 中无法表达多个规则 ID，`apply` 后只有第一个 ID 生效，其余丢失，造成与云端期望状态漂移。

约束：
- 通过保留 `rule_id`（Deprecated）+ 新增 `rule_ids`（Optional）实现，避免破坏性变更，存量配置仍可读
- 四个 `basic_bot_settings` 子块共享同一个 `botManagementActionOverrideSchema()`，改一处即可同时生效
- 已 vendored 的 SDK 已支持 `Ids []*string`，无需变更 vendor

## Goals / Non-Goals

**Goals:**
- 在 `action_overrides` 中新增 `rule_ids`（TypeList, Elem TypeString, Optional）字段，完整支持云 API 的多 ID 能力
- 保留 `rule_id`（TypeString, Optional, Deprecated）作为兼容降级路径，build 函数在 `rule_ids` 为空时回退读取
- 保持四个 `basic_bot_settings` 子块共享同一 schema 的一致性
- 通过单元测试覆盖多 ID 的 build / flatten 双向路径、`rule_id` 回退、以及 `rule_ids` 优先级

**Non-Goals:**
- 不改动 `action_overrides` 中 `action` 子块的 schema 与逻辑（`securityActionSchema()` 不变）
- 不改动 `bot_management` 中其他子块（如 `bot_management_lite`、`custom_rules`、`user_agent_rules` 等）的任何逻辑
- 不改动资源顶层的 schema、CRUD 主流程、ID 构造、retry 逻辑
- 不引入新 vendor 依赖
- 不做 state 自动迁移（`rule_id` 保留为 Deprecated，存量 state 中 `rule_id` 在 Read 时不再回填，改为回填 `rule_ids`，建议用户手动改写 HCL 并刷新 state）

## Decisions

### Decision 1: 保留 `rule_id`（Deprecated）并新增 `rule_ids`，而非直接替换

**选择**：保留 `rule_id`（TypeString, Optional, Deprecated）并新增 `rule_ids`（TypeList, Optional, Elem TypeString）。build 函数优先读 `rule_ids`，为空时回退读 `rule_id`。

**备选**：删除 `rule_id`，新增 `rule_ids`（TypeList, Required, Elem TypeString）。

**理由**：
- 删除 `rule_id` 会导致存量 HCL 配置升级后触发 `unknown attribute` plan 报错，破坏性较强
- 保留 `rule_id` 作为 Deprecated 兼容字段，让存量配置在升级后仍可被 build 回退路径读取，平滑过渡
- `rule_id` 标记 Deprecated 后会在 plan 时提示用户改用 `rule_ids`，引导收敛
- 当 `rule_ids` 与 `rule_id` 同时设置时，`rule_ids` 优先（避免歧义）

### Decision 2: 四个子块共用 `botManagementActionOverrideSchema()`，只改一处

**选择**：仅修改 `botManagementActionOverrideSchema()` 这一个函数的返回值，四个调用点（`source_idc` / `search_engine_bots` / `known_bot_categories` / `ip_reputation` 的 `action_overrides` Elem）自动生效。

**理由**：
- 现有代码已通过共享函数复用 schema，改一处即可保持四个子块一致
- 避免四处重复修改引入不一致风险

### Decision 3: flatten 函数完整展开 `Ids` 切片写入 `rule_ids`，不再回填 `rule_id`

**选择**：`flattenBotManagementActionOverride` 遍历 `override.Ids`，将每个非 nil 元素加入 `[]interface{}`，写入 `m["rule_ids"]`；不再回填已废弃的 `rule_id` 字段。

**备选**：同时回填 `rule_id = Ids[0]`。

**理由**：
- `rule_id` 已废弃，不应作为权威字段被 Read 回填
- 同时回填两者会让 state 中 `rule_ids` 与 `rule_id` 并存，产生持续 diff
- 仅回填 `rule_ids` 使 state 收敛到新字段，配合文档提示用户迁移

### Decision 4: build 函数优先 `rule_ids`，回退 `rule_id`

**选择**：`buildBotManagementActionOverrideFromMap` 先从 `m["rule_ids"].([]interface{})` 遍历构造 `[]*string`；若结果为空，再回退读 `m["rule_id"].(string)` 构造单元素 `[]*string`，赋值给 `override.Ids`。

**理由**：
- 与 provider 中 `block_rule_ids` / `observe_rule_ids` 的 build 模式一致
- 回退路径保证存量配置仍可工作
- 空列表且无 `rule_id` 时不设置 `Ids`（保持 nil），与 API 语义一致

### Decision 5: 文档示例改写，不做 state 迁移

**选择**：在 `resource_tc_teo_security_policy_config.md` 中将 `action_overrides` 示例由 `rule_id = "..."` 改为 `rule_ids = ["..."]`，并在变更说明中提示用户存量配置可保留 `rule_id`（已废弃）或迁移到 `rule_ids`。

**理由**：
- `rule_id` 保留为 Deprecated，字段名未变，存量配置无需立即改写即可工作
- state 中 `rule_id` 在 Read 时不再回填（改为回填 `rule_ids`），用户执行一次 `terraform apply`（或 `terraform refresh`）即可让 state 收敛到 `rule_ids`
- 影响范围限于使用 `bot_management.action_overrides` 的用户

## Risks / Trade-offs

- **Risk**：存量 state 中 `action_overrides` 下的 `rule_id` 键在升级后 Read 时不再回填（改为回填 `rule_ids`），导致 plan 显示 diff → **Mitigation**：用户执行一次 `terraform apply`（或 `terraform refresh`）即可让 state 收敛到 `rule_ids`；因 Read 改为完整展开 `Ids`，state 会正确反映云端真实的多 ID 配置
- **Risk**：同时设置 `rule_id` 与 `rule_ids` 时，`rule_ids` 优先，`rule_id` 被忽略 → **Mitigation**：`rule_id` 已标记 Deprecated，plan 会提示用户改用 `rule_ids`，引导收敛到单一字段
- **Trade-off**：保留 Deprecated 字段增加少量维护成本，但换取了向后兼容与平滑过渡

## Migration Plan

- 升级步骤（推荐迁移路径）：
  1. 升级 provider 版本（存量配置中 `rule_id = "xxx"` 仍可正常工作，build 回退路径读取）
  2.（可选迁移）将 HCL 中 `action_overrides` 块内的 `rule_id = "xxx"` 改写为 `rule_ids = ["xxx"]`
  3. 执行 `terraform plan`，确认无报错
  4. 执行 `terraform apply`（或先 `terraform refresh`）使 state 收敛到 `rule_ids`
- 回滚：若需回退，移除 `rule_ids` 字段并恢复 `rule_id` 为 Required；但回退后多 ID 配置仍会丢失，不建议回退
- 文档：在 `resource_tc_teo_security_policy_config.md` 示例中体现 `rule_ids` 用法

## Open Questions

- 无
