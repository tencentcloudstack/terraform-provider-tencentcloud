# Tasks for `add-captcha-ip-white-list-international`

> 前置：captcha SDK 包已 vendor（上一个 `tencentcloud_captcha_info_international` 资源已引入 `UseCaptchaClient()`），本 change 复用即可。

## 1. Service layer

- [x] 1.1 Add `DescribeCaptchaIpWhiteListInternationalById(ctx context.Context, captchaAppid, id int64) (ret *captchaintl.DescribeCaptchaWhiteListItem, errRet error)` to `tencentcloud/services/captcha/service_tencentcloud_captcha.go`. Build a `DescribeIpWhiteListInternationalRequest` with `CaptchaAppid = helper.Int64(captchaAppid)`, `PageIndex = 1`, `PageSize = 100` (the API documented maximum — re-verify against the doc), call `UseCaptchaClient().DescribeIpWhiteListInternational(request)` wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)`, iterate `Response.Data.DataList` and return the item whose `Id == id`, or `(nil, nil)` when not found.

## 2. Resource implementation

- [x] 2.1 Create `tencentcloud/services/captcha/resource_tc_captcha_ip_white_list_international.go`. File header, package, imports MUST mirror `tencentcloud/services/igtm/resource_tc_igtm_monitor.go` (same `defer LogElapsed + InconsistentCheck` pattern, same `var (...)` block style).
- [x] 2.2 Declare `ResourceTencentCloudCaptchaIpWhiteListInternational() *schema.Resource` with `Create/Read/Update/Delete` callbacks and `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`.
- [x] 2.3 Schema MUST declare exactly the 4 fields in D1: `name` (Required, TypeString), `captcha_appid` (Required, TypeInt, ForceNew), `ip` (Required, TypeString, ForceNew), `comment` (Optional, TypeString).
- [x] 2.4 Implement `resourceTencentCloudCaptchaIpWhiteListInternationalCreate`:
  - Build `captcha.NewCreateIpWhiteListInternationalRequest()`.
  - Map `name`/`ip`/`comment` via `d.GetOk` (string), `captcha_appid` via `d.GetOkExists` (int, TypeInt).
  - Wrap call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`. Validate `result == nil || result.Response == nil || len(result.Response.IdList) == 0 || result.Response.IdList[0] == nil` → `resource.NonRetryableError`.
  - On success, `d.SetId(fmt.Sprintf("%d#%d", captchaAppid, *response.Response.IdList[0]))`, then return Read.
- [x] 2.5 Implement `resourceTencentCloudCaptchaIpWhiteListInternationalRead`:
  - Parse `d.Id()` into `captchaAppid` and `id` via `strings.Split(d.Id(), tccommon.FILED_SP)` + `helper.StrToInt64` (validate length == 2).
  - Call `service.DescribeCaptchaIpWhiteListInternationalById(ctx, captchaAppid, id)`. If err, return err. If `respData == nil`, log warn + `d.SetId("")` + return nil.
  - Guard every `d.Set` with nil checks: `name`, `captcha_appid`, `ip`, `comment`.
- [x] 2.6 Implement `resourceTencentCloudCaptchaIpWhiteListInternationalUpdate`:
  - Use the `mutableArgs` pattern from D4: `["name", "comment"]`.
  - If changed, build `captcha.NewModifyIpWhiteListInternationalRequest()`, set `Id` and `CaptchaAppid` (parsed from composite ID), populate `Name`/`Comment` from `d.Get(...)`, wrap in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.
  - Always end by calling Read.
- [x] 2.7 Implement `resourceTencentCloudCaptchaIpWhiteListInternationalDelete`:
  - Build `captcha.NewDeleteIpWhiteListInternationalRequest()`, set `CaptchaAppid` and `Id` (parsed from composite ID), wrap call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

## 3. Provider registration

- [x] 3.1 In `tencentcloud/provider.go`, append near the existing captcha entry:
  ```go
  "tencentcloud_captcha_ip_white_list_international": captcha.ResourceTencentCloudCaptchaIpWhiteListInternational(),
  ```
- [x] 3.2 In `tencentcloud/provider.md`, add `tencentcloud_captcha_ip_white_list_international` under the `Captcha` -> `Resource` section.

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/captcha/resource_tc_captcha_ip_white_list_international.md` containing:
  - One-line summary line.
  - One `Example Usage` HCL block.
  - One `Import` section showing `terraform import tencentcloud_captcha_ip_white_list_international.example <CaptchaAppid>#<Id>`.
- [x] 4.2 Run `make doc` to regenerate `website/docs/r/captcha_ip_white_list_international.html.markdown`. Hand-editing the generated file is forbidden.

## 5. Acceptance test

- [x] 5.1 Create `tencentcloud/services/captcha/resource_tc_captcha_ip_white_list_international_test.go`, package `captcha_test`, structure mirroring `resource_tc_config_compliance_pack_test.go`:
  - Test function `TestAccTencentCloudCaptchaIpWhiteListInternationalResource_basic` with `t.Parallel()`.
  - `PreCheck: func() { tcacctest.AccPreCheck(t) }`, `Providers: tcacctest.AccProviders`.
  - Steps: initial config (basic create) with `TestCheckResourceAttrSet(..., "id")`; updated config changing `comment`; `ImportState` step with `ImportStateVerify: true`.

## 6. Validation

- [x] 6.1 `go build ./tencentcloud/...` clean.
- [x] 6.2 `go vet ./tencentcloud/services/captcha/...` clean.
- [x] 6.3 `read_lints` on the new files shows no new errors/warnings.
- [ ] 6.4 `openspec validate add-captcha-ip-white-list-international --strict` passes.
