## Why

TEO (EdgeOne) 平台目前不支持通过 Terraform 管理 JavaScript 注入规则（Security JS Injection Rule）。JavaScript 注入规则是 TEO 安全防护的重要组成部分，用于在匹配条件下注入 JS SDK 进行客户端验证（如 TC-RCE 和 TC-CAPTCHA 认证），帮助防护 Bot 攻击和恶意流量。用户需要在基础设施即代码的流程中管理这些安全规则的创建、查询、修改和删除。

## What Changes

- 新增 Terraform 通用资源 `tencentcloud_teo_security_js_injection_rule`，支持完整的 CRUD 操作：
  - **Create**: 调用 `CreateSecurityJSInjectionRule` 接口创建 JS 注入规则
  - **Read**: 调用 `DescribeSecurityJSInjectionRule` 接口查询 JS 注入规则
  - **Update**: 调用 `ModifySecurityJSInjectionRule` 接口修改 JS 注入规则
  - **Delete**: 调用 `DeleteSecurityJSInjectionRule` 接口删除 JS 注入规则
- 新增资源注册代码到 `tencentcloud/provider.go`
- 新增资源文档到 `tencentcloud/provider.md`
- 新增资源说明 `.md` 文件
- 新增单元测试文件，使用 gomonkey mock 方式测试业务逻辑

## Capabilities

### New Capabilities
- `teo-security-js-injection-rule-resource`: TEO JavaScript 注入规则资源，支持通过 Terraform 管理站点级别的 JS 注入安全规则，包括规则名称、优先级、匹配条件和注入选项的完整生命周期管理

### Modified Capabilities
<!-- 无需修改现有能力的规格要求 -->

## Impact

- **新增文件**:
  - `tencentcloud/services/teo/resource_tc_teo_security_js_injection_rule.go` - 资源实现
  - `tencentcloud/services/teo/resource_tc_teo_security_js_injection_rule_test.go` - 单元测试
  - `tencentcloud/services/teo/resource_tc_teo_security_js_injection_rule.md` - 资源文档
- **修改文件**:
  - `tencentcloud/provider.go` - 注册新资源
  - `tencentcloud/provider.md` - 添加资源索引条目
- **依赖**: 使用现有 `teo/v20220901` SDK 包，无需新增 vendor 依赖
- **API 调用**: 4 个 TEO 云 API 接口（Create/Describe/Modify/DeleteSecurityJSInjectionRule）
