## Context

腾讯云 EdgeOne (TEO) 产品 SDK 包 `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` 已在 vendor 中提供 `ApplyFreeCertificate` 接口（申请免费证书）和 `CheckFreeCertificateVerification` 接口（检查免费证书申请结果），但 Terraform Provider 当前未接入该操作。

当前状态：
- TEO SDK 包已 vendored，包含以下关键 struct（从 `vendor/.../teo/v20220901/models.go` 按行区间摘录）：
  - `ApplyFreeCertificateRequest` / `ApplyFreeCertificateRequestParams`（models.go ~L697-725）：
    - `ZoneId *string`：站点 ID。
    - `Domain *string`：申请免费证书的目标域名。
    - `VerificationMethod *string`：验证方式，取值 `http_challenge`（HTTP 访问文件验证）/ `dns_challenge`（DNS 委派验证）。
  - `ApplyFreeCertificateResponse` / `ApplyFreeCertificateResponseParams`（models.go ~L749-763）：
    - `DnsVerification *DnsVerification`：当 VerificationMethod 为 dns_challenge 时返回的 DNS 验证信息。
    - `FileVerification *FileVerification`：当 VerificationMethod 为 http_challenge 时返回的文件验证信息。
    - `RequestId *string`。
  - `DnsVerification`（models.go ~L16218-16227）：
    - `Subdomain *string`：主机记录。
    - `RecordType *string`：记录类型。
    - `RecordValue *string`：记录值。
  - `FileVerification`（models.go ~L17224-17230）：
    - `Path *string`：文件验证的 URL Path 部分。
    - `Content *string`：验证文件内容。
  - `CheckFreeCertificateVerificationRequest` / `CheckFreeCertificateVerificationRequestParams`（models.go ~L1766-1802）：
    - `ZoneId *string`：站点 ID。
    - `Domain *string`：加速域名。
  - `CheckFreeCertificateVerificationResponse` / `CheckFreeCertificateVerificationResponseParams`（models.go ~L1805-1823）：
    - `CommonName *string`：证书颁发给的域名（验证通过时返回）。
    - `SignatureAlgorithm *string`：证书签名算法。
    - `ExpireTime *string`：证书过期时间。
    - `RequestId *string`。
- TEO client 方法（从 `vendor/.../teo/v20220901/client.go` 按行区间摘录）：
  - `ApplyFreeCertificateWithContext(ctx, request) (*ApplyFreeCertificateResponse, error)`（client.go ~L115）：发起免费证书申请。
  - `CheckFreeCertificateVerificationWithContext(ctx, request) (*CheckFreeCertificateVerificationResponse, error)`（client.go ~L413）：检查证书申请结果，验证通过返回证书信息，验证失败返回错误。
- `tencentcloud/connectivity/client.go` 中已有 teo 客户端访问方法 `UseTeoV20220901Client()`，无需新增。
- `tencentcloud/services/teo/service_tencentcloud_teo.go` 中已有 `TeoService` 结构体（`NewTeoService(client)` 构造）与 `WaitTeoOriginACLById` 等轮询方法（参考 models.go ~L2215 的轮询模式），本次新增方法沿用同样分层与 `resource.Retry` 轮询风格。
- 该资源为 **RESOURCE_KIND_OPERATION**：一次性操作，操作完成后不需要记录任何状态，只需实现 Create（对应 framework action 的 `Invoke`），RUD 接口为空。
- 用户要求使用 **Terraform Plugin Framework** 的 action 资源类型实现，严格参考 `NewTeoConfirmOriginAclUpdate`（`tencentcloud/services/teo/action_tc_teo_confirm_origin_acl_update.go`）与 `NewMongodbRestoreDbInstance`（`tencentcloud/services/mongodb/action_tc_mongodb_restore_db_instance.go`）的代码结构。
- `ApplyFreeCertificate` 调用后需后端异步处理证书申请；用户需在域名侧配置验证（DNS 记录或验证文件）后，通过轮询 `CheckFreeCertificateVerification` 确认申请结果。用户要求"异步任务根据 CheckFreeCertificateVerification 轮询"。`CheckFreeCertificateVerification` 验证通过时返回 `CommonName` 等证书信息；验证未通过时返回对应验证失败错误码（如 `UnauthorizedOperation.*` 等，见 client.go ~L394-398 的可能错误码）。

约束：
- framework action 的 `InvokeResponse`（`vendor/.../terraform-plugin-framework/action/invoke.go`）仅暴露 `Diagnostics` 与 `SendProgress`，**不提供**写入输出/Computed 值的机制（action schema attribute 的 `IsComputed()` 始终返回 false，见 `vendor/.../action/schema/string_attribute.go` ~L153）；因此 `ApplyFreeCertificate` 响应中的验证信息（DnsVerification / FileVerification）无法作为 Terraform 输出属性回写。
- framework action 资源通过 `tencentcloud/framework/registry.go` 的 `actionFactories` 统一注册，不在 SDKv2 `provider.go` 的 `ResourcesMap` 中注册，因此无需修改 `provider.go`。

## 代码风格参照文件关键结构说明

### 参照一：`NewTeoConfirmOriginAclUpdate`（`tencentcloud/services/teo/action_tc_teo_confirm_origin_acl_update.go`）

framework action 的骨架结构：
- `var _ action.ActionWithConfigure = &TeoConfirmOriginAclUpdate{}`：编译期接口断言。
- `NewTeoConfirmOriginAclUpdate() action.Action`：工厂函数，返回 `&TeoConfirmOriginAclUpdate{}`。
- `TeoConfirmOriginAclUpdate struct { fw.ActionWithConfigure }`：嵌入 `fw.ActionWithConfigure` 获得 `Client()` / `Configure()`。
- `TeoConfirmOriginAclUpdateModel struct { ZoneId types.String tfsdk:"zone_id" }`：tfsdk 标签映射 schema 属性。
- `Metadata`：`resp.TypeName = "tencentcloud_teo_confirm_origin_acl_update"`。
- `Schema`：`schema.Schema{ Description: ..., Attributes: map[string]schema.Attribute{ "zone_id": schema.StringAttribute{ Required: true, Description: ... } } }`。
- `Invoke`：`req.Config.Get(ctx, &data)` 解析配置 → 校验必填 → `a.Client()==nil` 检查 → `service := NewTeoService(a.Client())` → `service.ConfirmOriginACLUpdate(ctx, zoneId)` → 失败用 `resp.Diagnostics.AddError(summary, detail)` 上报。

关键模式：`fw.ActionWithConfigure` 嵌入（`tencentcloud/framework/fw/action_with_configure.go`）继承 `withMeta.client`，通过 `a.Client()` 获取 `*connectivity.TencentCloudClient`。`sharedmeta.ProviderMeta{Client: client}` 用于测试注入（见 antiddos action test）。

### 参照二：`WaitTeoOriginACLById`（`tencentcloud/services/teo/service_tencentcloud_teo.go` ~L2215）

异步轮询方法的标准模式：
```go
func (me *TeoService) WaitTeoOriginACLById(ctx context.Context, timeout time.Duration, zoneId, status string) (errRet error) {
    logId := tccommon.GetLogId(ctx)
    request := teo.NewDescribeOriginACLRequest()
    request.ZoneId = &zoneId
    defer func() {
        if errRet != nil { log.Printf("[CRITAL]%s api[%s] fail ...", ...) }
    }()
    errRet = resource.Retry(timeout, func() *resource.RetryError {
        result, e := me.client.UseTeoClient().DescribeOriginACLWithContext(ctx, request)
        if e != nil { return tccommon.RetryError(e) }
        log.Printf("[DEBUG]%s api[%s] success ...", ...)
        if result == nil || result.Response == nil || result.Response.OriginACLInfo == nil || result.Response.OriginACLInfo.Status == nil {
            return resource.NonRetryableError(fmt.Errorf("... is nil."))
        }
        if *result.Response.OriginACLInfo.Status == status { return nil }
        return resource.RetryableError(fmt.Errorf("... is still %s.", *result.Response.OriginACLInfo.Status))
    })
    return
}
```

要点：轮询用 `resource.Retry(timeout, ...)`；未达目标状态返回 `resource.RetryableError(...)` 让外层继续重试；达到目标返回 `nil`；API 本身错误用 `tccommon.RetryError(e)` 判断是否可重试。

## Goals / Non-Goals

**Goals:**
- 新增 framework action 资源 `tencentcloud_teo_apply_free_certificate`，通过 `ApplyFreeCertificate` API 声明式发起 TEO 免费证书申请。
- 调用成功后通过 `CheckFreeCertificateVerification` 轮询直到证书申请成功（返回证书信息），保证证书实际生效。
- 通过单元测试覆盖 Invoke 成功（含轮询成功）、ApplyFreeCertificate API 错误、CheckFreeCertificateVerification 轮询失败/超时、缺必填输入、client 未配置、Metadata/Schema 校验等路径（gomonkey mock 云 API）。
- 保持与参考实现 `NewTeoConfirmOriginAclUpdate` / `NewMongodbRestoreDbInstance` 一致的代码风格与结构。

**Non-Goals:**
- 不持久化任何云端状态（不设置 id、不实现 Read/Update/Delete、不回写验证信息为输出属性）。
- 不对 `verification_method` 做取值校验，由云 API 返回错误。
- 不新增 `provider.go` 注册条目（framework action 走 registry.go）。
- 不新增 `_extension.go` 文件。
- 不修改 connectivity 层（teo 客户端已存在）。
- 不在 action 中暴露验证信息（DnsVerification / FileVerification）作为 Terraform 输出属性（framework action 技术上无法实现）；验证信息仅记录到日志供排障使用。

## Decisions

### Decision 1: 使用 Terraform Plugin Framework action 类型，而非 SDKv2 operation 资源

**选择**：按用户要求，使用 framework action（`action.Action` 接口）实现，文件命名为 `action_tc_teo_apply_free_certificate.go`，工厂函数 `NewTeoApplyFreeCertificate()` 返回 `action.Action`，在 `framework/registry.go` 的 `actionFactories` 注册（`services/teo` import 已存在）。

**备选**：使用 SDKv2 的 operation 资源模式（`resource_tc_teo_apply_free_certificate_operation.go`，在 `provider.go` ResourcesMap 注册）。

**理由**：
- 用户明确要求使用 framework plugin 框架实现 action 资源类型，参考 `NewTeoConfirmOriginAclUpdate` 与 `NewMongodbRestoreDbInstance`，忽略其他 SDKv2 设计。
- framework action 是表达"一次性操作"语义的标准方式，与参考实现一致。

### Decision 2: schema 只含输入属性，不暴露验证信息为输出

**选择**：schema 定义 `zone_id`（Required, String）、`domain`（Required, String）、`verification_method`（Required, String）三个输入属性。`Invoke` 中调用 `ApplyFreeCertificateWithContext`，拿到 DnsVerification / FileVerification 后记录日志；随后轮询 `CheckFreeCertificateVerification`。成功后直接返回（仅通过 Diagnostics 上报错误），不设置任何输出属性。

**备选**：使用 SDKv2 operation 资源（`resource_tc_teo_apply_free_certificate_operation.go`），将 DnsVerification / FileVerification 作为 Computed 嵌套字段回写，使用户可引用验证信息。

**理由**：
- framework action 的 `InvokeResponse` 不提供写入输出/Computed 值的机制（仅 `Diagnostics` 与 `SendProgress`），技术上无法回写。
- 需求明确"一次性操作，操作完不需要记录任何状态"，与参考实现 `NewTeoConfirmOriginAclUpdate`、`NewDlcInitializeTCLake`（只读输入、不写输出，验证信息仅记录日志）一致。
- 验证信息（DnsVerification / FileVerification）记录到 `log.Printf("[DEBUG]...")` 供排障使用，用户可在 Apply 输出日志中查到。

### Decision 3: 通过 service 层方法封装 ApplyFreeCertificate 的 retry，返回验证信息

**选择**：在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中新增 `TeoService.ApplyFreeCertificate(ctx, zoneId, domain, verificationMethod string) (dnsVerification *teo.DnsVerification, fileVerification *teo.FileVerification, errRet error)` 方法，内部使用 `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包装 `ApplyFreeCertificateWithContext` 调用，失败用 `tccommon.RetryError(e)` 包装；成功后提取并返回 DnsVerification / FileVerification（供 Invoke 记录日志），defer 中记录失败日志。

**备选**：在 `Invoke` 中直接调用 SDK client（不经 service 层），并在 Invoke 内同时完成 retry 与轮询。

**理由**：
- 参考实现 `TeoConfirmOriginAclUpdate` 通过 `service.ConfirmOriginACLUpdate(ctx, zoneId)` 封装 retry 逻辑，保持 action 文件简洁；本次沿用同样分层。
- ApplyFreeCertificate 的 retry 仅用于 API 调用本身（触发证书申请任务），轮询任务完成是独立步骤，不应放在同一个 retry 块中（避免 retry 语义混淆），由 Invoke 在拿到验证信息后单独调用轮询方法。

### Decision 4: 通过 service 层方法封装 CheckFreeCertificateVerification 的轮询

**选择**：在 `TeoService` 中新增 `TeoService.CheckFreeCertificateVerification(ctx, zoneId, domain string, timeout time.Duration) (errRet error)` 方法，使用 `resource.Retry(timeout, ...)` 轮询 `CheckFreeCertificateVerificationWithContext`，当响应返回证书信息（`CommonName` 非空）时认为申请成功返回 `nil`；当 API 返回"验证未通过"类错误时返回 `resource.RetryableError` 让外层继续重试；超时则返回错误。

**备选**：在 Invoke 中直接内联轮询循环。

**理由**：
- 参考实现 `WaitTeoOriginACLById` 使用 `resource.Retry` + `resource.RetryableError` 模式轮询状态，本次沿用。

### Decision 5: 异步轮询超时取 3 * ReadRetryTimeout

**选择**：在 `Invoke` 中调用 `TeoService.CheckFreeCertificateVerification(ctx, zoneId, domain, 3*tccommon.ReadRetryTimeout)` 轮询证书申请结果，与 `NewMongodbRestoreDbInstance`（3 * ReadRetryTimeout）等异步操作一致。

**备选**：轮询超时取 20 * ReadRetryTimeout（如 `tencentcloud_teo_zone` 创建用）。

**理由**：
- 证书申请验证需要用户手动配置 DNS/文件后再轮询，3 * ReadRetryTimeout（默认 9 分钟）与现有异步操作一致；如遇超时，用户可重试。

### Decision 6: 不注册到 SDKv2 provider.go，仅在 framework/registry.go 注册

**选择**：仅在 `tencentcloud/framework/registry.go` 的 `actionFactories` 追加 `teo.NewTeoApplyFreeCertificate`（import 块 `services/teo` 已存在）；不修改 `tencentcloud/provider.go`。

**理由**：
- framework action 资源通过 `framework/registry.go` 统一注册，与 teo、bdrc action 一致；在 SDKv2 `provider.go` ResourcesMap 注册会导致重复注册与 mux 冲突。

### Decision 7: 文档命名与格式

**选择**：文档文件 `action_tc_teo_apply_free_certificate.md`，格式为一句话描述（带上云产品名称 TEO）+ Example Usage（action 语法 `action "tencentcloud_teo_apply_free_certificate" "example" { config { ... } }`）+ NOTE 提示 Terraform 1.14+ 支持。无 Import 段、无手动 Argument/Attribute Reference。

**理由**：与参考 framework action 文档 `action_tc_teo_confirm_origin_acl_update.md` 格式一致。

## Risks / Trade-offs

- **Risk**：framework action 功能依赖 Terraform 1.14+ → **Mitigation**：仓库已存在 framework mux 机制，teo action 已注册并通过，本次仅追加一条 factory 条目，风险可控。
- **Risk**：`CheckFreeCertificateVerification` 轮询需要用户先在域名侧配置好验证（DNS 记录或验证文件），若用户未配置则轮询将超时失败 → **Mitigation**：在文档 NOTE 中提示用户需在调用 action 前或调用后立即配置验证内容（根据 ApplyFreeCertificate 返回的验证信息），验证配置完成后重试 apply 以使轮询成功。
- **Risk**：`CheckFreeCertificateVerification` API 错误中哪些属于"验证未通过"可重试、哪些属不可重试，SDK 文档未列出全部错误码 → **Mitigation**：采用保守策略，将非 `UnauthorizedOperation.NoPermission`/`UnauthorizedOperation.CamUnauthorized` 等权限类错误视为可重试（转 `resource.RetryableError`），让外层 retry 在超时范围内持续尝试；权限类错误直接 NonRetryableError 上报。

## Migration Plan

- 纯新增资源，无 state 迁移需求。
- 存量配置不受影响（不修改任何已有资源 schema）。
- 回滚：删除 `action_tc_teo_apply_free_certificate.go`、测试、文档，并移除 `framework/registry.go` 中的 teo action 条目即可。

## Open Questions

- 无