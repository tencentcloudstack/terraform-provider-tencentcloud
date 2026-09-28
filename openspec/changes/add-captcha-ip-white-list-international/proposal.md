## Why

腾讯云验证码（Captcha）国际站的 IP 白名单目前没有 Terraform 管理通道，用户只能在网页控制台手工创建/修改/删除 IP 白名单，无法做基础设施代码化（IaC）。需要新增一个 CRUD 型 Terraform 资源，把 captcha 国际站 4 个 IP 白名单接口（CreateIpWhiteListInternational / DescribeIpWhiteListInternational / ModifyIpWhiteListInternational / DeleteIpWhiteListInternational）暴露为 `tencentcloud_captcha_ip_white_list_international`，让用户能在 HCL 中声明 IP 白名单并通过 `terraform plan/apply` 管理其全生命周期。

## What Changes

- 新增资源 `tencentcloud_captcha_ip_white_list_international`（CRUD 型），归属 `tencentcloud/services/captcha/`。
- 新增 Schema 字段，**严格按 `CreateIpWhiteListInternational` 接口入参 1:1 映射**（4 个字段，snake_case 化）：`name`、`captcha_appid`、`ip`、`comment`。
- 资源 ID = Create 响应 `CaptchaAppid + IdList[0]`（复合 ID，`#` 分隔），同步映射 `DeleteIpWhiteListInternational` / `ModifyIpWhiteListInternational` 的 `CaptchaAppid` + `Id` 入参，以及 `DescribeIpWhiteListInternational` 的 `CaptchaAppid` 过滤入参。
- 新增 service 层方法 `DescribeCaptchaIpWhiteListInternationalById(ctx, captchaAppid int64, id int64)`，内部调用 `DescribeIpWhiteListInternational` + `CaptchaAppid` 过滤，分页 `PageSize` 取接口文档最大值，并从返回 `DataList` 中按 `Id` 精确匹配。
- Update 路径调用 `ModifyIpWhiteListInternational`（Name / Id / CaptchaAppid / Status / Comment）。注意：Modify 接口**不含 `Ip`**，因此 `ip` 字段标 `ForceNew: true`；`captcha_appid` 作为资源 ID 组成部分同样标 `ForceNew: true`。
- Delete 路径调用 `DeleteIpWhiteListInternational`，请求体 `CaptchaAppid` + `Id`。
- 所有 SDK 调用必须用 `resource.Retry(...)` 包装，并做空指针保护。
- 新增 `.md` 资源文档（命名与 `resource_tc_config_compliance_pack.md` 同款）+ 验收测试 `_test.go`（命名与 `resource_tc_config_compliance_pack_test.go` 同款）。
- 在 `tencentcloud/provider.go` 注册新资源 `tencentcloud_captcha_ip_white_list_international`。

## Capabilities

### New Capabilities

- `captcha-ip-white-list-international-resource`: 新增 `tencentcloud_captcha_ip_white_list_international` 资源的 schema、CRUD 行为、ID 约定（复合 ID）、字段约束、文档与测试规范，作为 captcha 国际站 IP 白名单的 IaC 入口。

### Modified Capabilities

<!-- 不修改任何已有 capability 的 requirement，仅在 provider.go 增加注册一行 -->

## Impact

- 新文件：
  - `tencentcloud/services/captcha/resource_tc_captcha_ip_white_list_international.go`
  - `tencentcloud/services/captcha/resource_tc_captcha_ip_white_list_international.md`
  - `tencentcloud/services/captcha/resource_tc_captcha_ip_white_list_international_test.go`
  - `website/docs/r/captcha_ip_white_list_international.html.markdown`（由 `make doc` 生成）
- 既有文件改动：
  - `tencentcloud/services/captcha/service_tencentcloud_captcha.go`：新增 `DescribeCaptchaIpWhiteListInternationalById` 方法（不破坏现有方法）
  - `tencentcloud/provider.go`：在资源注册映射中追加 `"tencentcloud_captcha_ip_white_list_international": captcha.ResourceTencentCloudCaptchaIpWhiteListInternational()` 一行
  - `tencentcloud/provider.md`：在 `Captcha` 段追加资源名
- 依赖：SDK 使用国际站 `tencentcloud-sdk-go-intl-en` 的 `tencentcloud/captcha/v20190722` 包（已 vendor），其中已包含 4 个所需接口及配套 model（`CreateIpWhiteListInternationalRequest/Response`、`DescribeIpWhiteListInternationalRequest/Response`、`ModifyIpWhiteListInternationalRequest`、`DeleteIpWhiteListInternationalRequest`、`DescribeCaptchaWhiteListItem`），无需升级 SDK。
- 不修改任何既有资源的 schema，对现有用户零影响。
