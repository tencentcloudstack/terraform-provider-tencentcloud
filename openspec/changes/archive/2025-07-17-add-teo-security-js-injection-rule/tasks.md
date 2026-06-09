## 1. Schema 定义与资源框架

- [x] 1.1 创建资源文件 `tencentcloud/services/teo/resource_tc_teo_security_js_injection_rule.go`，定义 `ResourceTencentCloudTeoSecurityJsInjectionRule()` 函数，包含完整的 schema 定义（zone_id、js_injection_rules 嵌套列表、js_injection_rule_ids）和 CRUD 函数入口
- [x] 1.2 实现 schema 中 `js_injection_rules` 嵌套结构，包含 rule_id（Computed）、name（Required）、priority（Optional/Computed）、condition（Required）、inject_js（Optional/Computed）字段
- [x] 1.3 在 schema 中配置 Importer 支持（schema.ImportStatePassthrough），并声明 Timeouts 块

## 2. CRUD 函数实现

- [x] 2.1 实现 Create 函数 `resourceTencentCloudTeoSecurityJsInjectionRuleCreate`：构造 CreateSecurityJSInjectionRule 请求，调用 API，设置 d.SetId(zone_id)，存储 js_injection_rule_ids，调用 Read 刷新状态
- [x] 2.2 实现 Read 函数 `resourceTencentCloudTeoSecurityJsInjectionRuleRead`：调用 DescribeSecurityJSInjectionRule（Limit=100 分页），解析响应填充 js_injection_rules 和 js_injection_rule_ids 到 state；资源不存在时设置 d.SetId("")
- [x] 2.3 实现 Update 函数 `resourceTencentCloudTeoSecurityJsInjectionRuleUpdate`：检测 js_injection_rules 变更，调用 ModifySecurityJSInjectionRule，调用 Read 刷新状态
- [x] 2.4 实现 Delete 函数 `resourceTencentCloudTeoSecurityJsInjectionRuleDelete`：从 state 读取 js_injection_rule_ids，调用 DeleteSecurityJSInjectionRule 删除所有规则

## 3. Provider 注册

- [x] 3.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册 `"tencentcloud_teo_security_js_injection_rule": teo.ResourceTencentCloudTeoSecurityJsInjectionRule()`
- [x] 3.2 在 `tencentcloud/provider.md` 的 TEO 分类下添加 `tencentcloud_teo_security_js_injection_rule` 条目

## 4. 资源文档

- [x] 4.1 创建 `tencentcloud/services/teo/resource_tc_teo_security_js_injection_rule.md` 文件，包含资源描述、Example Usage（含 jsonencode 用法示例）和 Import 说明

## 5. 单元测试

- [x] 5.1 创建测试文件 `tencentcloud/services/teo/resource_tc_teo_security_js_injection_rule_test.go`，使用 gomonkey mock 方式编写 Create 函数单元测试
- [x] 5.2 编写 Read 函数单元测试，mock DescribeSecurityJSInjectionRule API
- [x] 5.3 编写 Update 函数单元测试，mock ModifySecurityJSInjectionRule 和 DescribeSecurityJSInjectionRule API
- [x] 5.4 编写 Delete 函数单元测试，mock DeleteSecurityJSInjectionRule API

## 6. 验证

- [x] 6.1 使用 `go test -gcflags=all=-l` 运行单元测试，确保所有测试通过
