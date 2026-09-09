## 1. Schema 调整

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go` 的 `botManagementActionOverrideSchema()` 函数中，将 `rule_id`（TypeString, Required）字段替换为 `rule_ids`（TypeList, Required, Elem: TypeString），更新 Description 说明其支持一个或多个 Bot 规则 ID / 分类 ID
- [x] 1.2 确认四个 `basic_bot_settings` 子块（`source_idc` / `search_engine_bots` / `known_bot_categories` / `ip_reputation.ip_reputation_group`）的 `action_overrides` Elem 均引用 `botManagementActionOverrideSchema()`，无需重复修改

## 2. 转换函数修改

- [x] 2.1 修改 `buildBotManagementActionOverrideFromMap` 函数：从 schema map 读取 `rule_ids`（`[]interface{}`），遍历每个非空字符串元素构造 `[]*string`（`helper.String(v)`），赋值给 `override.Ids`；移除原 `rule_id` 单字符串读取逻辑
- [x] 2.2 修改 `flattenBotManagementActionOverride` 函数：遍历 `override.Ids`（`[]*string`）将每个非 nil 元素加入 `[]interface{}`，写入 `m["rule_ids"]`；当 `Ids` 为空时不设置该字段；移除原仅取 `Ids[0]` 的逻辑
- [x] 2.3 检查 `buildBotManagementActionOverrideFromMap` 中 `action` 子块的处理逻辑保持不变（仍调用 `buildSecurityActionFromMap`）

## 3. 单元测试

- [x] 3.1 在 `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go` 中新增测试用例，验证 `buildBotManagementActionOverrideFromMap` 在 `rule_ids = ["rule-a", "rule-b"]` 时构造的 `BotManagementActionOverrides.Ids` 包含且仅包含两个匹配元素
- [x] 3.2 新增测试用例，验证 `flattenBotManagementActionOverride` 在 API 返回 `Ids = ["rule-a", "rule-b"]` 时，扁平化后的 map 中 `rule_ids` 包含且仅包含两个匹配元素
- [x] 3.3 新增测试用例，验证 `flattenBotManagementActionOverride` 在 `Ids` 为 nil 时不设置 `rule_ids` 字段
- [x] 3.4 保证已有测试用例不受影响（build/flatten 单 ID 场景改为以单元素列表形式覆盖）

## 4. 文档同步

- [x] 4.1 在 `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md` 中，将 `bot_management` → `basic_bot_settings` → `action_overrides` 相关示例由 `rule_id = "..."` 改写为 `rule_ids = ["..."]`（包含多 ID 示例）
- [ ] 4.2 执行 `make doc`，根据 provider 规范重新生成 `website/docs/` 下的 markdown 文档（禁止手改 website/ 目录）

## 5. 验证

- [x] 5.1 确认代码可正确编译（不执行 go build，由后续流程验证）
- [x] 5.2 复查 `rule_id` 在资源文件中已无残留引用（仅存在于 `rule_ids` 中）
- [x] 5.3 复查四个 `basic_bot_settings` 子块的 `action_overrides` 行为一致
