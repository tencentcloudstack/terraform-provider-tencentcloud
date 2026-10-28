## 1. Schema 定义

- [x] 1.1 在 `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go` 的 `security_policy.exception_rules.rules` 块下，紧邻 `web_security_modules_for_exception` 之后，新增 `web_security_submodules_for_exception` 字段（TypeSet, Optional, Elem 为 TypeString），并在 Description 中说明仅当 `skip_scope` 为 `WebSecuritySubmodules` 时有效，给出子模块合法取值示例（托管规则子模块、速率限制子模块、自定义规则子模块、HTTP DDoS 防护子模块、高级/基础 Bot 管理子模块等）
- [x] 1.2 确认新增字段为纯 Optional（无 Computed、无 ForceNew），保证向后兼容不破坏已有 state 与 TF 配置

## 2. Read 逻辑

- [x] 2.1 在 Read 函数遍历 `respData.ExceptionRules.Rules` 的循环中，紧接 `WebSecurityModulesForException` 的 nil 检查与 set 之后，新增对 `rules.WebSecuritySubmodulesForException` 的 nil 检查：若非 nil 则 `rulesMap["web_security_submodules_for_exception"] = rules.WebSecuritySubmodulesForException`
- [x] 2.2 确认云端返回 nil 时不调用 set，避免覆盖已有配置；保持与现有 `web_security_modules_for_exception` 读取模式一致

## 3. Create / Update 逻辑

- [x] 3.1 在 Create/Update 遍历 `exception_rules.rules` 配置的循环中，紧接 `web_security_modules_for_exception` 的扩展逻辑之后，新增对 `web_security_submodules_for_exception` 的处理：从 `rulesMap` 取 `*schema.Set`，遍历其 List，对每个非 nil 元素断言为 string 后取地址追加到 `exceptionRule.WebSecuritySubmodulesForException`
- [x] 3.2 确认空集合或未配置时 `exceptionRule.WebSecuritySubmodulesForException` 保持 nil，与同级字段处理一致
- [x] 3.3 确认 Create 与 Update 共用同一构造逻辑，无需重复实现；Update 末尾仍调用 read 回写最新状态

## 4. 单元测试

- [x] 4.1 在 `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go` 中，使用 gomonkey mock 云 API（`ModifySecurityPolicyWithContext` 与 `DescribeSecurityPolicyWithContext` 或对应 service 层封装），新增测试用例覆盖 Create 时 `web_security_submodules_for_exception` 被正确填充到 `ExceptionRule.WebSecuritySubmodulesForException`
- [x] 4.2 新增测试用例覆盖 Read 回读：mock `DescribeSecurityPolicy` 返回包含 `WebSecuritySubmodulesForException` 的响应，断言 state 中 `web_security_submodules_for_exception` 与预期一致
- [x] 4.3 新增测试用例覆盖 Update 变更：变更 `web_security_submodules_for_exception` 集合内容，断言 mock `ModifySecurityPolicyWithContext` 被调用且请求字段更新
- [x] 4.4 保证已有 Create/Read/Update/Delete 测试用例继续通过（不修改其行为）

## 5. 文档同步

- [x] 5.1 在 `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md` 的 Example Usage 的 `exception_rules.rules` 块中，补充 `web_security_submodules_for_exception` 字段示例与可选取值说明
- [ ] 5.2 在收尾阶段通过 `make doc` 自动生成 `website/docs/` 下的 markdown 文档（禁止手动编辑 `website/` 目录）

## 6. 代码正确性检查

- [x] 6.1 核对新增参数 `web_security_submodules_for_exception` 在 `DescribeSecurityPolicy` 出参路径 `response.SecurityPolicy.ExceptionRules.Rules.WebSecuritySubmodulesForException` 与 SDK `teov20220901.ExceptionRule.WebSecuritySubmodulesForException` 一致
- [x] 6.2 核对新增参数在 `ModifySecurityPolicy` 入参路径 `request.SecurityPolicy.ExceptionRules.Rules.WebSecuritySubmodulesForException` 与 SDK 一致，确认 CRUD 接口均支持该字段
- [x] 6.3 确认所有新增代码中函数返回的 error 均被检查处理（无需处理的必成功调用用 `_ =` 接收）
- [x] 6.4 确认未在资源 go 文件开头添加注释，未生成多余的 `_extension.go` 文件
