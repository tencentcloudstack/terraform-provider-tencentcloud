## 1. Service 层方法实现

- [x] 1.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `ApplyFreeCertificate(ctx, zoneId, domain, verificationMethod string) (dnsVerification *teo.DnsVerification, fileVerification *teo.FileVerification, errRet error)` 方法：构造 `teov20220901.NewApplyFreeCertificateRequest()`，填充 `ZoneId`、`Domain`、`VerificationMethod`；使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `ApplyFreeCertificateWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装；成功后返回 DnsVerification / FileVerification；defer 中记录失败日志
- [x] 1.2 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `CheckFreeCertificateVerification(ctx, zoneId, domain string, timeout time.Duration) (errRet error)` 方法：构造 `teov20220901.NewCheckFreeCertificateVerificationRequest()`，填充 `ZoneId`、`Domain`；使用 `resource.Retry(timeout, ...)` 轮询 `CheckFreeCertificateVerificationWithContext`；当响应 `CommonName` 非空时返回 `nil`（成功）；当 API 返回验证未通过错误时返回 `resource.RetryableError`（继续轮询）；当 API 返回权限/授权错误（如 `UnauthorizedOperation.NoPermission`、`UnauthorizedOperation.CamUnauthorized`）时返回 `resource.NonRetryableError`（立即停止）；defer 中记录失败日志

## 2. Framework Action 资源实现

- [x] 2.1 新增 `tencentcloud/services/teo/action_tc_teo_apply_free_certificate.go`：定义 `TeoApplyFreeCertificate struct { fw.ActionWithConfigure }` 与 `NewTeoApplyFreeCertificate() action.Action` 工厂函数；`var _ action.ActionWithConfigure = &TeoApplyFreeCertificate{}` 编译期断言
- [x] 2.2 实现 `TeoApplyFreeCertificateModel` struct：`ZoneId types.String`、`Domain types.String`、`VerificationMethod types.String`（tfsdk 标签 `zone_id` / `domain` / `verification_method`）
- [x] 2.3 实现 `Metadata`：`resp.TypeName = "tencentcloud_teo_apply_free_certificate"`
- [x] 2.4 实现 `Schema`：定义 `zone_id`（Required, String, 站点 ID）、`domain`（Required, String, 申请免费证书的目标域名）、`verification_method`（Required, String, 验证方式 `http_challenge`/`dns_challenge`）三个属性，带 Description
- [x] 2.5 实现 `Invoke`：解析 config → 校验 `zone_id`/`domain`/`verification_method` 非空 → `a.Client() == nil` 检查 → `service := NewTeoService(a.Client())` → `service.ApplyFreeCertificate(ctx, zoneId, domain, verificationMethod)` 拿到验证信息后记录 `[DEBUG]` 日志 → `service.CheckFreeCertificateVerification(ctx, zoneId, domain, 3*tccommon.ReadRetryTimeout)` 轮询 → 失败用 `resp.Diagnostics.AddError(summary, detail)` 上报

## 3. Action 注册

- [x] 3.1 在 `tencentcloud/framework/registry.go` 的 `actionFactories` 末尾追加 `teo.NewTeoApplyFreeCertificate`（import 块中 `services/teo` 已存在，无需新增）
- [x] 3.2 确认不在 `tencentcloud/provider.go` 的 `ResourcesMap` 中注册该 action
- [x] 3.3 在 `tencentcloud/provider.md` 的 TEO `Action` 索引段追加 `tencentcloud_teo_apply_free_certificate`（供 `gendoc` 生成文档索引）

## 4. 单元测试

- [x] 4.1 新增 `tencentcloud/services/teo/action_tc_teo_apply_free_certificate_test.go`，实现辅助函数：`ptrStringTeoApplyFreeCertificate(s string) *string`、`teoApplyFreeCertificateSchema(t)` 获取 action schema、`newTeoApplyFreeCertificateInvokeRequest(t, zoneId, domain, verificationMethod string)` 构造 InvokeRequest（参考 `action_tc_antiddos_unblock_resources_test.go` 的 tftypes 构造模式）、`setTeoApplyFreeCertificateClient(a *teo.TeoApplyFreeCertificate, client *connectivity.TencentCloudClient)` 注入 mock client
- [x] 4.2 新增 `TestTeoApplyFreeCertificate_Invoke_Success`：mock `ApplyFreeCertificateWithContext` 返回 DnsVerification / FileVerification，mock `CheckFreeCertificateVerificationWithContext` 返回非空 CommonName，断言无 error diagnostic 且 API 被正确参数调用
- [x] 4.3 新增 `TestTeoApplyFreeCertificate_Invoke_ApplyError`：mock `ApplyFreeCertificateWithContext` 返回错误，断言 error diagnostic 包含错误信息
- [x] 4.4 新增 `TestTeoApplyFreeCertificate_Invoke_VerificationTimeout`：mock `ApplyFreeCertificateWithContext` 成功 but `CheckFreeCertificateVerificationWithContext` 始终返回验证未通过，断言 error diagnostic 产生
- [x] 4.5 新增 `TestTeoApplyFreeCertificate_Invoke_MissingInput`：zone_id / domain / verification_method 为空时断言 error diagnostic 且 API 未被调用
- [x] 4.6 新增 `TestTeoApplyFreeCertificate_Invoke_ClientNotConfigured`：provider client 未设置时断言 error diagnostic 提示 "Provider not configured"
- [x] 4.7 新增 `TestTeoApplyFreeCertificate_MetadataAndSchema`：断言 `TypeName == "tencentcloud_teo_apply_free_certificate"` 且所有 required 属性存在且类型正确

## 5. 文档

- [x] 5.1 新增 `tencentcloud/services/teo/action_tc_teo_apply_free_certificate.md`：一句话描述带上 TEO 云产品名称（"Provides an action to apply for TEO free certificate for a domain..."）；Example Usage 使用 `action "tencentcloud_teo_apply_free_certificate" "example" { config { zone_id = ... domain = ... verification_method = ... } }` 语法；NOTE 提示 Terraform 1.14+ 支持；提示用户需在调用后根据日志验证信息配置 DNS/文件验证；提示 `dns_challenge` 下 `FailedOperation.ApplyCertDnsVerificationFailed`（`No DNS delegation record detected.`）会被视为验证进行中的中间态并继续轮询；无 Import 段；无手动 Argument/Attribute Reference
- [x] 5.2 执行 `make doc`，根据 provider 规范重新生成 `website/docs/` 下的 markdown 文档（禁止手改 website 目录，由收尾阶段 tfpacer-finalize skill 执行）

## 6. dns_challenge 中 DNS 委派记录未检测到错误的处理

- [x] 6.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 新增包级常量 `teoFreeCertDnsVerificationFailedCode = "FailedOperation.ApplyCertDnsVerificationFailed"` 与 `teoFreeCertDnsDelegationPendingMessage = "No DNS delegation record detected."`，并新增未导出辅助函数 `isTeoFreeCertDnsDelegationPending(sdkErr *sdkErrors.TencentCloudSDKError) bool`（错误码精确匹配 + 错误信息大小写不敏感包含匹配）
- [x] 6.2 在 `CheckFreeCertificateVerification` 的 retry 闭包错误分支中，权限类错误（`UnauthorizedOperation.NoPermission` / `UnauthorizedOperation.CamUnauthorized`）保持 `resource.NonRetryableError`；命中 `isTeoFreeCertDnsDelegationPending` 时打印 `[DEBUG]` 日志（logId、action、request body、message）并返回 `resource.RetryableError`，错误信息包含"请按 DNS 验证内容完成委派记录配置并等待生效"的指引（`%w` 包装原始 SDK 错误）；其余错误维持原有 `resource.RetryableError(e)` 行为
- [x] 6.3 在 `tencentcloud/services/teo/action_tc_teo_apply_free_certificate_test.go` 新增 `TestTeoApplyFreeCertificate_Invoke_DnsDelegationRecordPending`（首次返回该错误、随后返回非空 CommonName，断言无 error diagnostic 且 `CheckFreeCertificateVerificationWithContext` 至少被调用 2 次）与 `TestTeoApplyFreeCertificate_Invoke_DnsDelegationRecordAlwaysMissing`（持续返回该错误，断言产生 error diagnostic 且 detail 包含 DNS 委派记录指引）
- [x] 6.4 同步更新 `proposal.md`、`design.md`（Decision 4 与 Risks）、`specs/teo-apply-free-certificate-action/spec.md`（轮询方法要求与单元测试/文档场景）以反映该错误处理逻辑
