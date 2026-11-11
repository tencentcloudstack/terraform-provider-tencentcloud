# Tasks for `add-captcha-info-international`

## 0. SDK vendor (prerequisite)

- [x] 0.1 In `tencentcloud/services/captcha/resource_tc_captcha_info_international.go`, import `github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/captcha/v20190722` (aliased as `captcha`).
- [x] 0.2 Run `go mod vendor` to vendor the captcha package into `vendor/`.

## 1. Service layer

- [x] 1.1 Add `DescribeCaptchaInfoInternationalById(ctx context.Context, captchaAppId string) (ret *captcha.DescribeCaptchaConsoleSubDataInternational, errRet error)` to `tencentcloud/services/captcha/service_tencentcloud_captcha.go`. Build a `DescribeCaptchaInfoListInternationalRequest` with `CaptchaAppId = helper.String(captchaAppId)`, `PageIndex = 1`, `PageSize = 100` (the API documented maximum — re-verify against the doc during implementation), call `UseCaptchaClient().DescribeCaptchaInfoListInternational(request)` wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)`, return the first element of `Response.Data.DataList` or `(nil, nil)` when the slice is empty.

## 2. Resource implementation

- [x] 2.1 Create `tencentcloud/services/captcha/resource_tc_captcha_info_international.go`. File header, package, imports MUST mirror `tencentcloud/services/igtm/resource_tc_igtm_monitor.go` (same `defer LogElapsed + InconsistentCheck` pattern, same `var (...)` block style).
- [x] 2.2 Declare `ResourceTencentCloudCaptchaInfoInternational() *schema.Resource` with `Create/Read/Update/Delete` callbacks and `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`.
- [x] 2.3 Schema MUST declare exactly the 13 fields listed in D1 (no more, no fewer), with `Required/Optional` strictly validated against the `CreateCaptchaInfoInternational` API documentation, and types per D1. `channel_info` MUST set `ForceNew: true`.
- [x] 2.4 Implement `resourceTencentCloudCaptchaInfoInternationalCreate`:
  - Build `captcha.NewCreateCaptchaInfoInternationalRequest()`.
  - Map every schema field via `d.GetOk`. For `check_appid_switch` / `check_iv_switch` (TypeInt), use `d.GetOkExists` (zero is meaningful).
  - For `tags` (string slice), iterate the `[]interface{}` and append `helper.String(v.(string))`.
  - Wrap the SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`. Inside the retry, validate `result == nil || result.Response == nil || result.Response.Data == nil` → `resource.NonRetryableError`.
  - On success, `d.SetId(helper.Int64ToStr(*response.Response.Data))`, then return `resourceTencentCloudCaptchaInfoInternationalRead(d, meta)`.
- [x] 2.5 Implement `resourceTencentCloudCaptchaInfoInternationalRead`:
  - Build the typical `logId/ctx/service/captchaAppId := d.Id()` block.
  - Call `service.DescribeCaptchaInfoInternationalById(ctx, captchaAppId)`. If err, return err. If `respData == nil`, log warn + `d.SetId("")` + return nil.
  - For every field, guard with `if respData.<Field> != nil` before `d.Set`.
- [x] 2.6 Implement `resourceTencentCloudCaptchaInfoInternationalUpdate`:
  - Use the `mutableArgs` pattern from D4: `["app_name", "verify_rank", "user_set_cap_type", "defend_mode", "disable_invisible_switch", "verify_domain", "verify_bundle_id", "verify_package", "check_appid_switch", "check_iv_switch", "check_box_style", "tags"]`.
  - If any of those changed, build `captcha.NewModifyCaptchaInfoInternationalRequest()`, set `request.CaptchaAppId = helper.String(d.Id())`, populate ALL mutable fields from current `d.Get(...)` (full overwrite), wrap in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.
  - Always end by calling `resourceTencentCloudCaptchaInfoInternationalRead(d, meta)`.
- [x] 2.7 Implement `resourceTencentCloudCaptchaInfoInternationalDelete`:
  - Build `captcha.NewRemoveCaptchaInfoInternationalRequest()`, set `request.CaptchaAppId = helper.String(d.Id())`, wrap call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

## 3. Provider registration

- [x] 3.1 In `tencentcloud/provider.go`, locate a suitable spot and append:
  ```go
  "tencentcloud_captcha_info_international": captcha.ResourceTencentCloudCaptchaInfoInternational(),
  ```

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/captcha/resource_tc_captcha_info_international.md` containing:
  - One-line summary line.
  - One `Example Usage` HCL block.
  - One `Import` section showing `terraform import tencentcloud_captcha_info_international.example <CaptchaAppId>`.
- [x] 4.2 Run `make doc` to regenerate `website/docs/r/captcha_info_international.html.markdown`. Hand-editing the generated file is forbidden.

## 5. Acceptance test

- [x] 5.1 Create `tencentcloud/services/captcha/resource_tc_captcha_info_international_test.go`, package `captcha_test`, structure mirroring `resource_tc_config_compliance_pack_test.go`:
  - Test function `TestAccTencentCloudCaptchaInfoInternationalResource_basic` with `t.Parallel()`.
  - `PreCheck: func() { tcacctest.AccPreCheck(t) }`, `Providers: tcacctest.AccProviders`.
  - Steps: initial config (basic create) with `TestCheckResourceAttrSet(..., "id")`; updated config changing a mutable field (e.g. `app_name`); `ImportState` step with `ImportStateVerify: true`.

## 6. Validation

- [x] 6.1 `go build ./tencentcloud/...` clean.
- [x] 6.2 `go vet ./tencentcloud/services/captcha/...` clean.
- [x] 6.3 `read_lints` on the new files shows no new errors/warnings.
- [ ] 6.4 `openspec validate add-captcha-info-international --strict` passes.
