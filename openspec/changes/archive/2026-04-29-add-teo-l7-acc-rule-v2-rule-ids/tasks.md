## 1. 移除 rule_ids

- [x] 1.1 从 `resource_tc_teo_l7_acc_rule_v2.go` 的 Schema 中移除 `rule_ids` 参数定义
- [x] 1.2 从 Create 函数中移除 `ruleIds` 变量和收集/设置逻辑
- [x] 1.3 从 Read 函数中移除 `ruleIds` 的收集和设置逻辑
- [x] 1.4 从 `resource_tc_teo_l7_acc_rule_v2_test.go` 中移除 `rule_ids` 相关单元测试及辅助函数

## 2. 添加 Vary action

- [x] 2.1 在 `resource_tc_teo_l7_acc_rule_extension.go` 的 `TencentTeoL7RuleBranchBasicInfo` schema 中添加 `vary_parameters` 字段（TypeList, Optional, MaxItems:1, 包含 `switch` 字段）
- [x] 2.2 在 `resourceTencentCloudTeoL7AccRuleGetBranchs` 函数中添加 `vary_parameters` 到 `VaryParameters` 的映射
- [x] 2.3 在 `resourceTencentCloudTeoL7AccRuleSetBranchs` 函数中添加 `VaryParameters` 到 `vary_parameters` 的扁平化映射
- [x] 2.4 修正 `name` 字段的 Description：将 `SetContentIdentifierParameters` 改为 `SetContentIdentifier`，并补充 `Vary`、`ContentCompression`、`OriginAuthentication` action（移除 `Shield`）

## 3. 添加 OriginAuthentication action

- [x] 3.1 在 schema 中添加 `origin_authentication_parameters` 字段（TypeList, Optional, MaxItems:1, 包含 `request_properties` 列表，每项有 `type`/`name`/`value`）
- [x] 3.2 在 Get 函数中添加 `origin_authentication_parameters` 到 `OriginAuthenticationParameters` 的映射
- [x] 3.3 在 Set 函数中添加 `OriginAuthenticationParameters` 到 `origin_authentication_parameters` 的扁平化映射

## 4. 文档

- [x] 4.1 更新 `resource_tc_teo_l7_acc_rule_v2.md` 示例，添加 `Vary` 和 `OriginAuthentication` action 配置块
