## Why

`tencentcloud_teo_security_client_attester` 资源当前仅暴露了 `ClientAttester` 的基础字段（name、attester_source、attester_duration 及各认证选项），未覆盖云 API `ClientAttester` 结构体中已支持的使用限制（`UsageLimit`、`MaxUsageCount`）以及 `TCEOCaptchaOption` 中的移动端人机校验模式（`MobileCaptchaMode`）。用户无法通过 Terraform 声明式地配置认证结果使用次数限制和移动端 native app 的验证交互模式，必须借助控制台或 SDK 手动操作，导致声明式运维流程断裂。本次变更将这些字段补齐到资源 schema 中。

## What Changes

- 在 `tencentcloud_teo_security_client_attester` 资源的 `client_attesters` 嵌套 schema 中新增 `usage_limit`（Optional, TypeString）字段，映射云 API `ClientAttesters[].UsageLimit`（取值 `on`/`off`）。
- 在 `client_attesters` 嵌套 schema 中新增 `max_usage_count`（Optional, TypeInt）字段，映射云 API `ClientAttesters[].MaxUsageCount`（取值范围 1-100000，`UsageLimit=on` 时必填）。
- 在 `client_attesters.tc_eo_captcha_option` 嵌套 schema 中新增 `mobile_captcha_mode`（Optional, TypeString）字段，映射云 API `ClientAttesters[].TCEOCaptchaOption.MobileCaptchaMode`（取值 `Invisible`/`Visible`）。
- 更新 `buildClientAttesterFromMap` 构造函数，使其在 Create/Update 时填充新增的 `UsageLimit`、`MaxUsageCount`、`TCEOCaptchaOption.MobileCaptchaMode` 字段。
- 更新 Read 操作中的 `attesterMap` 回填逻辑，将 `UsageLimit`、`MaxUsageCount`、`TCEOCaptchaOption.MobileCaptchaMode` 回写到 state。
- 在 `resource_tc_teo_security_client_attester_test.go` 中补充覆盖新增字段的单元测试用例（使用 gomonkey mock 云 API）。
- 更新资源文档 `resource_tc_teo_security_client_attester.md` 的示例，展示新增字段的使用方式（最终 website docs 由 `make doc` 生成）。

非破坏性：所有新增字段均为 Optional，不影响已有 state 与 TF 配置，向后兼容。

## Capabilities

### New Capabilities
<!-- 无新增 capability，本次仅修改现有 capability -->

### Modified Capabilities
- `teo-security-client-attester`: 在 `client_attesters` 嵌套 schema 中新增 `usage_limit`、`max_usage_count` 字段，并在 `tc_eo_captcha_option` 子块中新增 `mobile_captcha_mode` 字段；Read/Create/Update 逻辑同步支持回填与传递这些字段。

## Impact

- 代码：
  - `tencentcloud/services/teo/resource_tc_teo_security_client_attester.go`（schema 新增字段、Create/Read/Update 逻辑、`buildClientAttesterFromMap` 构造函数）
  - `tencentcloud/services/teo/resource_tc_teo_security_client_attester_test.go`（新增覆盖字段用例）
  - `tencentcloud/services/teo/resource_tc_teo_security_client_attester.md`（示例补充）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teov20220901` 包，`ClientAttester` 与 `TCEOCaptchaOption` 结构体已包含 `UsageLimit`/`MaxUsageCount`/`MobileCaptchaMode` 字段，无需变更 vendor。
- 向后兼容：新增 Optional 字段，不影响已有 state 与 TF 配置。
- 文档：需要同步更新 website docs（由 `make doc` 生成）。