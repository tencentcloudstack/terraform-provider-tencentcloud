## Why

`tencentcloud_teo_security_policy_config` 资源当前已支持例外规则（`exception_rules`）配置，包括安全防护模块（`web_security_modules_for_exception`）、托管规则与规则组等跳过维度。但 TEO 云 API `ExceptionRule` 结构体新增了 `WebSecuritySubmodulesForException` 字段，用于指定安全防护**子模块**（如托管规则的规则集/高频扫描防护、速率限制规则、自定义规则、HTTP DDoS 防护的子功能、高级/基础 Bot 管理子功能等）作为例外目标，提供比模块级更细粒度的例外控制。当前 Terraform 资源未暴露该字段，用户无法通过 Terraform 配置子模块级例外规则，只能借助控制台或 SDK 手动操作，导致声明式运维流程断裂。

## What Changes

- 在 `tencentcloud_teo_security_policy_config` 资源的 `security_policy.exception_rules.rules` 块下新增 `web_security_submodules_for_exception` 参数（TypeSet, Optional, Elem 为 TypeString），与现有的 `web_security_modules_for_exception` 同级。
- 在资源 Read 逻辑中，新增对 `DescribeSecurityPolicy` 返回的 `SecurityPolicy.ExceptionRules.Rules[].WebSecuritySubmodulesForException` 字段的读取与 nil 检查，回填到 Terraform state。
- 在资源 Create/Update 逻辑中，新增对 `ModifySecurityPolicy` 请求的 `SecurityPolicy.ExceptionRules.Rules[].WebSecuritySubmodulesForException` 字段的填充，将 TypeSet 转换为 `[]*string`。
- 在 `_test.go` 中补充覆盖该新参数的单元测试用例（使用 gomonkey mock 云 API）。
- 同步更新资源文档 `resource_tc_teo_security_policy_config.md` 的示例（在收尾阶段通过 `make doc` 生成）。

非破坏性：仅新增 Optional 字段，不影响已有 state 与 TF 配置，向后完全兼容。

## Capabilities

### New Capabilities
<!-- 无新增 capability，本次仅在现有资源上新增一个参数 -->

### Modified Capabilities
- `teo-security-policy-exception-submodules`: 为 `tencentcloud_teo_security_policy_config` 资源的例外规则新增 `web_security_submodules_for_exception` 子模块级例外参数，覆盖 schema 定义、Read/Create/Update 逻辑、单元测试与文档。

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go`（schema 新增字段、Read 读取逻辑、Create/Update 填充逻辑）
  - `tencentcloud/services/teo/resource_tc_teo_security_policy_config_test.go`（新增覆盖 `web_security_submodules_for_exception` 的单测用例）
  - `tencentcloud/services/teo/resource_tc_teo_security_policy_config.md`（示例文档，收尾阶段通过 `make doc` 生成）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901.ExceptionRule.WebSecuritySubmodulesForException`（`[]*string`），无需变更 vendor。
- 向后兼容：新增 Optional 字段，不影响已有 state 与 TF 配置。
- 文档：需要同步更新 website docs（由收尾阶段 `make doc` 自动生成读取 `.md` 文件）。
