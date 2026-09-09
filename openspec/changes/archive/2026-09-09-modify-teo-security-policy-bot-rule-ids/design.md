## Context

`tencentcloud_teo_security_policy_config` 资源封装了 TEO 安全策略的 CRUD，其中 `bot_management.basic_bot_settings` 下的 `source_idc` / `search_engine_bots` / `known_bot_categories` / `ip_reputation.ip_reputation_group` 四个子块，各自包含一个 `action_overrides` 列表（TypeList，元素 schema 由 `botManagementActionOverrideSchema()` 返回）。

当前 `botManagementActionOverrideSchema()` 将 `rule_id` 定义为 `TypeString, Required`（单个 ID），对应的转换函数：
- `buildBotManagementActionOverrideFromMap`：从 schema 读单个 `rule_id` 字符串，构造 `[]*string{helper.String(v)}` 写入 `BotManagementActionOverrides.Ids`
- `flattenBotManagementActionOverride`：从 API 读 `override.Ids[0]` 写入 schema 的 `rule_id`

云 API 现状：`teov20220901.BotManagementActionOverrides.Ids` 类型为 `[]*string`，注释明确说明"Bot 规则组下的具体项"，支持同时改写多条规则项的处置动作。

问题：用户在 HCL 中无法表达多个规则 ID，`apply` 后只有第一个 ID 生效，其余丢失，造成与云端期望状态漂移。

约束：
- 该字段位于深层嵌套块，是 schema 字段名 + 类型的双重变更，无法纯加法兼容
- 四个 `basic_bot_settings` 子块共享同一个 `botManagementActionOverrideSchema()`，改一处即可同时生效
- 已 vendored 的 SDK 已支持 `Ids []*string`，无需变更 vendor

## Goals / Non-Goals

**Goals:**
- 将 `action_overrides` 的规则 ID 字段从单值 `rule_id`（TypeString）改为列表 `rule_ids`（TypeList, Elem TypeString），完整支持云 API 的多 ID 能力
- 保持四个 `basic_bot_settings` 子块共享同一 schema 的一致性
- 通过单元测试覆盖多 ID 的 build / flatten 双向路径

**Non-Goals:**
- 不改动 `action_overrides` 中 `action` 子块的 schema 与逻辑（`securityActionSchema()` 不变）
- 不改动 `bot_management` 中其他子块（如 `bot_management_lite`、`custom_rules`、`user_agent_rules` 等）的任何逻辑
- 不改动资源顶层的 schema、CRUD 主流程、ID 构造、retry 逻辑
- 不引入新 vendor 依赖
- 不做 state 自动迁移（字段名变更无法自动迁移，需用户手动改写 HCL）

## Decisions

### Decision 1: 直接将 `rule_id` 替换为 `rule_ids`，而非保留旧字段并新增

**选择**：删除 `rule_id`，新增 `rule_ids`（TypeList, Required, Elem TypeString）。

**备选**：保留 `rule_id`（标记 Deprecated）并新增 `rule_ids`，二选一。

**理由**：
- `rule_id` 与 `rule_ids` 语义完全重叠且同时存在会让用户困惑（到底哪个生效？build 函数需处理两者并存）
- 当前 `rule_id` 本身就是 bug（只取第一个 ID 导致数据丢失），保留一个已知错误字段并标 Deprecated 反而误导
- 该字段在深层嵌套块内，使用面相对窄，破坏性影响可控
- 字段名从单数到复数、类型从 string 到 list，是表达"多 ID"语义最清晰的方式，符合 provider 中已有的 `block_rule_ids` / `observe_rule_ids`（同为 TypeList 多 ID）的命名惯例

### Decision 2: 四个子块共用 `botManagementActionOverrideSchema()`，只改一处

**选择**：仅修改 `botManagementActionOverrideSchema()` 这一个函数的返回值，四个调用点（`source_idc` / `search_engine_bots` / `known_bot_categories` / `ip_reputation` 的 `action_overrides` Elem）自动生效。

**理由**：
- 现有代码已通过共享函数复用 schema，改一处即可保持四个子块一致
- 避免四处重复修改引入不一致风险

### Decision 3: flatten 函数完整展开 `Ids` 切片

**选择**：`flattenBotManagementActionOverride` 遍历 `override.Ids`，将每个非 nil 元素加入 `[]interface{}`，写入 `m["rule_ids"]`。

**备选**：保留对 `Ids[0]` 的处理作为兼容。

**理由**：
- 保留单元素兼容会与 `rule_ids` 列表语义冲突，且无法表达多 ID
- 完整展开才是对齐 API 的正确实现

### Decision 4: build 函数遍历 `rule_ids` 列表构造 `[]*string`

**选择**：`buildBotManagementActionOverrideFromMap` 从 `m["rule_ids"].([]interface{})` 遍历，对每个非空字符串元素 `helper.String(v)`，append 到 `[]*string`，赋值给 `override.Ids`。

**理由**：
- 与 provider 中 `block_rule_ids` / `observe_rule_ids` 的 build 模式一致
- 空列表时不设置 `Ids`（保持 nil），与 API 语义一致

### Decision 5: 文档示例改写，不做 state 迁移

**选择**：在 `resource_tc_teo_security_policy_config.md` 中将 `action_overrides` 示例由 `rule_id = "..."` 改为 `rule_ids = [...]`，并在变更说明中提示用户需手动改写存量 HCL。

**理由**：
- 字段名变更无法通过 schema 层自动迁移（Terraform Plugin SDK 不支持字段重命名迁移）
- 影响范围限于使用 `bot_management.action_overrides` 的用户，属于 bug 修复范畴

## Risks / Trade-offs

- **Risk**：存量用户配置中使用 `rule_id` 的部分升级后会触发 `unknown attribute "rule_id"` plan 报错 → **Mitigation**：这是修复 schema/API 不一致的必要破坏性变更；在 changelog 与文档中明确提示改写为 `rule_ids = ["原值"]`
- **Risk**：存量 state 中 `action_overrides` 下的 `rule_id` 键在升级后 Read 时不再被回填（改为 `rule_ids`），导致 plan 显示 diff → **Mitigation**：用户执行一次 `terraform apply`（或 `terraform refresh`）即可让 state 收敛到 `rule_ids`；因 Read 改为完整展开 `Ids`，state 会正确反映云端真实的多 ID 配置
- **Trade-off**：破坏向后兼容，但换取了与云 API 能力对齐的正确行为，避免数据丢失

## Migration Plan

- 升级步骤：
  1. 升级 provider 版本
  2. 将 HCL 中 `action_overrides` 块内的 `rule_id = "xxx"` 改写为 `rule_ids = ["xxx"]`
  3. 执行 `terraform plan`，确认无 `unknown attribute` 报错
  4. 执行 `terraform apply`（或先 `terraform refresh`）使 state 收敛
- 回滚：若需回退，将 schema 改回 `rule_id`（TypeString）并恢复 build/flatten 的单 ID 逻辑；但回退后多 ID 配置仍会丢失，不建议回退
- 文档：在 `resource_tc_teo_security_policy_config.md` 示例中体现 `rule_ids` 用法

## Open Questions

- 无
