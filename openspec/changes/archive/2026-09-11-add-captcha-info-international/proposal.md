## Why

腾讯云验证码（Captcha）国际站目前没有 Terraform 管理通道，用户只能在网页控制台手工创建/修改/删除验证码配置（验证码信息），无法做基础设施代码化（IaC）。需要新增一个 CRUD 型 Terraform 资源，把 captcha 国际站 4 个新接口（CreateCaptchaInfoInternational / DescribeCaptchaInfoListInternational / ModifyCaptchaInfoInternational / RemoveCaptchaInfoInternational）暴露为 `tencentcloud_captcha_info_international`，让用户能在 HCL 中声明验证码配置并通过 `terraform plan/apply` 管理其全生命周期。

## What Changes

- 新增资源 `tencentcloud_captcha_info_international`（CRUD 型），归属 `tencentcloud/services/captcha/`。
- 新增 Schema 字段，**严格按 `CreateCaptchaInfoInternational` 接口入参 1:1 映射**（13 个字段，snake_case 化）：`app_name`、`channel_info`、`verify_rank`、`user_set_cap_type`、`defend_mode`、`tags`、`disable_invisible_switch`、`verify_domain`、`verify_bundle_id`、`verify_package`、`check_appid_switch`、`check_iv_switch`、`check_box_style`。
- 资源 ID = Create 响应 `Data`（`*int64`，验证码 AppId），通过 `helper.Int64ToStr` 转 string 后 `d.SetId`；并同步映射到 `RemoveCaptchaInfoInternational` / `ModifyCaptchaInfoInternational` / `DescribeCaptchaInfoListInternational` 的 `CaptchaAppId`（`*string`）入参。
- 新增 service 层方法 `DescribeCaptchaInfoInternationalById(ctx, captchaAppId)`，内部调用 `DescribeCaptchaInfoListInternational` + `CaptchaAppId` 精确过滤，分页 `PageSize` 取接口文档最大值。
- Update 路径调用 `ModifyCaptchaInfoInternational`。注意：Modify 接口**不含 `ChannelInfo`** 入参，因此 `channel_info` 字段标 `ForceNew: true`（改通道触发重建）。
- Delete 路径调用 `RemoveCaptchaInfoInternational`，请求体仅 `CaptchaAppId`。
- 所有 SDK 调用必须用 `resource.Retry(...)` 包装，并做空指针保护。
- 新增 `.md` 资源文档（命名与 `resource_tc_config_compliance_pack.md` 同款）+ 验收测试 `_test.go`（命名与 `resource_tc_config_compliance_pack_test.go` 同款）。
- 在 `tencentcloud/provider.go` 注册新资源 `tencentcloud_captcha_info_international`。

## Capabilities

### New Capabilities

- `captcha-info-international-resource`: 新增 `tencentcloud_captcha_info_international` 资源的 schema、CRUD 行为、ID 约定、字段约束、文档与测试规范，作为 captcha 国际站验证码配置的 IaC 入口。

### Modified Capabilities

<!-- 不修改任何已有 capability 的 requirement，仅在 provider.go 增加注册一行 -->

## Impact

- 新文件：
  - `tencentcloud/services/captcha/resource_tc_captcha_info_international.go`
  - `tencentcloud/services/captcha/resource_tc_captcha_info_international.md`
  - `tencentcloud/services/captcha/resource_tc_captcha_info_international_test.go`
  - `website/docs/r/captcha_info_international.html.markdown`（由 `make doc` 生成）
- 既有文件改动：
  - `tencentcloud/services/captcha/service_tencentcloud_captcha.go`：新增 `DescribeCaptchaInfoInternationalById` 方法（不破坏现有方法）
  - `tencentcloud/provider.go`：在资源注册映射中追加 `"tencentcloud_captcha_info_international": captcha.ResourceTencentCloudCaptchaInfoInternational()` 一行
- 依赖：SDK 使用国际站 `tencentcloud-sdk-go-intl-en` 的 `tencentcloud/captcha/v20190722` 包（当前 `go.mod` 已 require `v3.0.1486`），其中已包含 4 个所需接口及配套 model。**注意：本地 vendor 目前尚未引入 captcha 包，实现前需在代码中 import 该包并重新 `go mod vendor`。**
- 不修改任何既有资源的 schema，对现有用户零影响。
