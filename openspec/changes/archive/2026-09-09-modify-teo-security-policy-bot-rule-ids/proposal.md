## Why

`tencentcloud_teo_security_policy_config` 资源中 `bot_management` → `basic_bot_settings`（`source_idc` / `search_engine_bots` / `known_bot_categories` / `ip_reputation`）下的 `action_overrides` 块，当前 schema 仅暴露 `rule_id`（TypeString，单个 ID）。但对应的腾讯云 TEO API `BotManagementActionOverrides.Ids` 字段类型为 `[]*string`，原生支持多个规则 ID 同时改写处置动作。由于 schema 与 API 能力不一致，用户在 HCL 中配置多条规则 ID 时，实际 `apply` 后只会写入第一个 ID，其余 ID 被丢弃，导致策略配置不完整、与云端期望状态漂移。需要新增 `rule_ids`（TypeList，支持多个 ID），并保留 `rule_id` 标记为 Deprecated，使 Terraform 与云 API 能力对齐，同时不破坏存量配置的可读性。

## What Changes

- 新增 `botManagementActionOverrideSchema()` 中的字段 `rule_ids`（TypeList, Optional, Elem: TypeString），以匹配云 API `BotManagementActionOverrides.Ids` 的多 ID 语义；保留原 `rule_id`（TypeString, Optional）并标记为 Deprecated，提示用户改用 `rule_ids`。
- 修改 `flattenBotManagementActionOverride` 函数：将云 API 返回的 `override.Ids`（`[]*string`）完整展开写入 `rule_ids` 列表，而非仅取 `Ids[0]` 写入 `rule_id`。
- 修改 `buildBotManagementActionOverrideFromMap` 函数：优先从 schema 的 `rule_ids` 列表读取多个 ID 构造 `[]*string` 赋值给 `override.Ids`；当 `rule_ids` 为空时，回退读取已废弃的 `rule_id` 单字符串作为兼容降级路径。
- 新增单元测试覆盖多 ID 场景的 flatten / build 逻辑、`rule_id` 回退路径、以及 `rule_ids` 优先级场景（使用 gomonkey mock 云 API）。
- 更新资源文档 `resource_tc_teo_security_policy_config.md`，将 `action_overrides` 示例由 `rule_id = "..."` 改为 `rule_ids = ["..."]`。

兼容性说明：保留 `rule_id` 并标记 Deprecated，存量用户配置中使用了 `rule_id` 的部分在升级后仍可正常被 build 回退路径读取，不会立即触发 plan 报错（unknown attribute）；同时通过 `rule_ids` 提供多 ID 能力。推荐用户将 `rule_id = "xxx"` 迁移改写为 `rule_ids = ["xxx"]` 以获取完整能力。

## Capabilities

### New Capabilities
- `teo-security-policy-bot-action-override-rule-ids`: 修正 `tencentcloud_teo_security_policy_config` 资源中 `bot_management.basic_bot_settings.*.action_overrides` 块的规则 ID 字段，新增列表 `rule_ids` 对齐云 API `BotManagementActionOverrides.Ids` 的多 ID 语义，同时保留 `rule_id` 标记为 Deprecated。

### Modified Capabilities
<!-- 无现有 spec 描述 action_overrides 的 rule_id 字段，本次为新增 capability -->

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go`（`botManagementActionOverrideSchema` schema、`flattenBotManagementActionOverride`、`buildBotManagementActionOverrideFromMap` 三个函数）
  - `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go`（新增多 ID、rule_id 回退、优先级单测）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.BotManagementActionOverrides.Ids`（`[]*string`），无需变更 vendor。
- 向后兼容：**兼容**——保留 `rule_id`（Deprecated），存量 HCL 配置仍可被 build 回退路径读取，不会立即报错；`rule_ids` 为新增能力字段。
- 文档：需要同步更新 `resource_tc_teo_security_policy_config.md` 示例。
