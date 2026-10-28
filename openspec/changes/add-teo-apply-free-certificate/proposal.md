## Why

腾讯云 EdgeOne (TEO) 提供了 `ApplyFreeCertificate` 接口，用于为 CNAME 接入模式的加速域名申请免费证书。用户在需要使用 DNS 委派验证（`dns_challenge`）或文件验证（`http_challenge`）方式申请免费证书时，需先调用该接口发起申请并获取验证内容，再根据验证内容在域名侧完成配置，最后调用 `CheckFreeCertificateVerification` 接口校验申请结果。当前 Terraform Provider 缺少对该操作的声明式编排能力，用户只能通过控制台或 SDK 手动触发，导致证书申请流程无法纳入 IaC 体系。本次新增基于 Terraform Plugin Framework 的 action 资源 `tencentcloud_teo_apply_free_certificate`，使免费证书申请可通过 Terraform 声明式发起，并在调用后异步轮询 `CheckFreeCertificateVerification` 直到申请成功。

## What Changes

- 新增 Terraform Plugin Framework action 资源 `tencentcloud_teo_apply_free_certificate`，文件 `tencentcloud/services/teo/action_tc_teo_apply_free_certificate.go`，严格参考 `NewTeoConfirmOriginAclUpdate`（`tencentcloud/services/teo/action_tc_teo_confirm_origin_acl_update.go`）与 `NewMongodbRestoreDbInstance`（`tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance.go`）的代码结构：
  - `zone_id`（Required, String）：站点 ID。
  - `domain`（Required, String）：申请免费证书的目标域名。
  - `verification_method`（Required, String）：申请免费证书的验证方式，取值 `http_challenge`（HTTP 访问文件验证）/ `dns_challenge`（DNS 委派验证）。
  - 仅实现 `Invoke`（对应 Create）：调用 `ApplyFreeCertificateWithContext` 发起证书申请，拿到响应中的验证信息（DnsVerification / FileVerification）后记录日志；随后使用 `CheckFreeCertificateVerificationWithContext` 轮询直到证书申请成功（响应中返回证书信息），失败或超时则返回错误。
- 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 新增两个 service 层方法：
  - `ApplyFreeCertificate(ctx, zoneId, domain, verificationMethod string) (dnsVerification *teo.DnsVerification, fileVerification *teo.FileVerification, errRet error)`：封装 `ApplyFreeCertificateWithContext` 调用，内部使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装，失败用 `tccommon.RetryError(e)` 包装；成功后返回 DnsVerification / FileVerification 验证信息。
  - `CheckFreeCertificateVerification(ctx, zoneId, domain string, timeout time.Duration) error`：使用 `resource.Retry` 轮询 `CheckFreeCertificateVerificationWithContext`，当响应返回证书信息（CommonName 非空）时认为申请成功结束轮询；当 API 返回"验证未通过"类错误时继续重试；超时则返回错误。
    - 其中 `dns_challenge` 场景下 API 返回 `FailedOperation.ApplyCertDnsVerificationFailed` + `No DNS delegation record detected.`（DNS 委派记录尚未检测到）属于验证进行中的中间态：service 层通过 `isTeoFreeCertDnsDelegationPending` 显式识别该错误组合，打印 `[DEBUG]` 日志并返回 `resource.RetryableError`（错误信息中提示用户按 DNS 验证内容完成委派记录配置），继续轮询直至检测到记录或轮询超时，而不是立即失败。
- 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 末尾注册 `teo.NewTeoApplyFreeCertificate`（import 块中 `services/teo` 已存在，无需新增）。
- 新增单元测试 `tencentcloud/services/teo/action_tc_teo_apply_free_certificate_test.go`，使用 gomonkey mock 云 API，覆盖 Invoke 成功（含轮询成功）、ApplyFreeCertificate API 错误、CheckFreeCertificateVerification 轮询失败/超时、`FailedOperation.ApplyCertDnsVerificationFailed`（`No DNS delegation record detected.`）后继续轮询成功与持续缺失时超时报错、缺必填输入、client 未配置、Metadata/Schema 校验等路径。
- 新增资源文档 `tencentcloud/services/teo/action_tc_teo_apply_free_certificate.md`（一句话描述带上 TEO 云产品名称；action 语法 Example Usage；NOTE 提示 Terraform 1.14+ 支持；无 Import 段；无手动 Argument/Attribute Reference）。
- 无需修改 `tencentcloud/provider.go`（framework action 资源通过 `framework/registry.go` 统一注册，不在 SDKv2 ResourcesMap 中注册）；需在 `tencentcloud/provider.md` 的 TEO `Action` 索引段追加 `tencentcloud_teo_apply_free_certificate`，供 `gendoc` 生成文档索引。
- connectivity 层无需改动：teo 客户端访问方法 `UseTeoV20220901Client()` 已存在。

## Capabilities

### New Capabilities
- `teo-apply-free-certificate-action`: 提供一次性 action 资源 `tencentcloud_teo_apply_free_certificate`，通过 `ApplyFreeCertificate` API 为 TEO CNAME 接入域名申请免费证书，调用后异步轮询 `CheckFreeCertificateVerification` 直到证书申请成功；操作完成后不持久化任何云端状态。

### Modified Capabilities
无。

## Impact

- 代码：
  - `tencentcloud/services/teo/action_tc_teo_apply_free_certificate.go`（新增 — framework action 资源）
  - `tencentcloud/services/teo/service_tencentcloud_teo.go`（修改 — 新增 `ApplyFreeCertificate` 与 `CheckFreeCertificateVerification` service 方法封装，含 retry 与轮询逻辑）
  - `tencentcloud/services/teo/action_tc_teo_apply_free_certificate_test.go`（新增 — 单元测试）
  - `tencentcloud/services/teo/action_tc_teo_apply_free_certificate.md`（新增 — 文档）
  - `tencentcloud/framework/registry.go`（修改 — 在 `actionFactories` 追加 `teo.NewTeoApplyFreeCertificate`）
  - `tencentcloud/provider.md`（修改 — 在 TEO `Action` 索引段追加 `tencentcloud_teo_apply_free_certificate`）
- 依赖：使用已 vendored 的 `tencentcloud-sdk-go` 中 `teo/v20220901` 包（`ApplyFreeCertificateRequest` / `ApplyFreeCertificateResponse` / `DnsVerification` / `FileVerification` / `CheckFreeCertificateVerificationRequest` / `CheckFreeCertificateVerificationResponse`），无需新增 vendor；复用现有 `TeoService` 结构体。
- 向后兼容：纯新增资源，不影响已有资源配置与 state。
- 文档：需通过收尾阶段 `make doc` 生成 `website/docs/` 下对应文档（禁止手改 website 目录）。