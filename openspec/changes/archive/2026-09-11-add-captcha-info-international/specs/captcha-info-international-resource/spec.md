# captcha-info-international-resource Specification

## ADDED Requirements

### Requirement: Resource MUST be registered as `tencentcloud_captcha_info_international`

The provider SHALL register a new CRUD-type resource named `tencentcloud_captcha_info_international` whose Create/Read/Update/Delete callbacks invoke the captcha international APIs (`CreateCaptchaInfoInternational`, `DescribeCaptchaInfoListInternational`, `ModifyCaptchaInfoInternational`, `RemoveCaptchaInfoInternational`) of `tencentcloud-sdk-go-intl-en/tencentcloud/captcha/v20190722`.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_captcha_info_international"` mapped to `captcha.ResourceTencentCloudCaptchaInfoInternational()`.

#### Scenario: Importer is configured

- **WHEN** a user runs `terraform import tencentcloud_captcha_info_international.example <CaptchaAppId>`
- **THEN** the importer SHALL pass the supplied `CaptchaAppId` straight into the Read callback (using `schema.ImportStatePassthrough`), and Read SHALL hydrate state from `DescribeCaptchaInfoListInternational` filtered by that CaptchaAppId.

### Requirement: Schema MUST mirror the CreateCaptchaInfoInternational API input 1:1

The resource schema SHALL declare exactly these 13 top-level argument keys, with semantics matching the SDK request fields of `CreateCaptchaInfoInternationalRequest`:

| HCL key | SDK field | Type | ForceNew |
|---|---|---|---|
| `app_name` | `AppName` | TypeString | No |
| `channel_info` | `ChannelInfo` | TypeString | **Yes** |
| `verify_rank` | `VerifyRank` | TypeString | No |
| `user_set_cap_type` | `UserSetCapType` | TypeString | No |
| `defend_mode` | `DefendMode` | TypeString | No |
| `tags` | `Tags` | TypeList of TypeString | No |
| `disable_invisible_switch` | `DisableInvisibleSwitch` | TypeString | No |
| `verify_domain` | `VerifyDomain` | TypeString | No |
| `verify_bundle_id` | `VerifyBundleId` | TypeString | No |
| `verify_package` | `VerifyPackage` | TypeString | No |
| `check_appid_switch` | `CheckAppidSwitch` | TypeInt | No |
| `check_iv_switch` | `CheckIvSwitch` | TypeInt | No |
| `check_box_style` | `CheckBoxStyle` | TypeString | No |

Per the `CreateCaptchaInfoInternational` API documentation, `app_name` and `channel_info` are `Required: true`, while the remaining 11 fields are `Optional: true`.

#### Scenario: Schema field set matches the API input exactly

- **WHEN** the resource schema is inspected
- **THEN** it declares exactly the 13 keys above (no extra field, no missing field), each with the correct SDK type mapping.

#### Scenario: channel_info is ForceNew

- **WHEN** the user changes `channel_info`
- **THEN** Terraform reports a destroy + create cycle (because `channel_info` is ForceNew, since `ModifyCaptchaInfoInternational` has no `ChannelInfo` input).

### Requirement: Resource ID MUST be the Create response `Data` mapped to `CaptchaAppId`

After Create, the resource's Terraform ID SHALL be set to the decimal string of `*response.Response.Data` (an `*int64`), via `helper.Int64ToStr`. Read/Update/Delete SHALL use `d.Id()` directly as the `CaptchaAppId` (`*string`) input.

#### Scenario: ID is set from CreateCaptchaInfoInternational response

- **GIVEN** a successful `CreateCaptchaInfoInternational` returns `Data = 20210610001`
- **WHEN** Create completes
- **THEN** `d.Id()` returns `"20210610001"`.

#### Scenario: Create response missing Data is fatal

- **WHEN** `CreateCaptchaInfoInternational` returns a non-nil response with `Data == nil`
- **THEN** the Create callback SHALL return a non-retryable error mentioning `Data is nil`.

### Requirement: Read MUST locate the resource via DescribeCaptchaInfoListInternational + CaptchaAppId

The service-layer helper `(*CaptchaService).DescribeCaptchaInfoInternationalById(ctx, captchaAppId) (*captcha.DescribeCaptchaConsoleSubDataInternational, error)` SHALL build a `DescribeCaptchaInfoListInternationalRequest` with:

- `CaptchaAppId = helper.String(captchaAppId)`
- `PageIndex = 1`
- `PageSize = 100` (the API documented maximum)

…and invoke `UseCaptchaClient().DescribeCaptchaInfoListInternational` inside `resource.Retry(ReadRetryTimeout, ...)`.

If `Response.Data.DataList` is empty, the helper SHALL return `(nil, nil)`. The Read callback SHALL detect that case and call `d.SetId("")` to mark the resource gone.

#### Scenario: Resource exists

- **GIVEN** the API holds a captcha info with `CaptchaAppId="20210610001"`
- **WHEN** Read runs
- **THEN** all schema attributes are populated from the API response with nil-safe pointer dereferencing.

#### Scenario: Resource removed out-of-band

- **GIVEN** the user runs `terraform refresh` after the captcha info was deleted in the web console
- **WHEN** Read calls `DescribeCaptchaInfoListInternational` and `len(Response.Data.DataList) == 0`
- **THEN** Read SHALL `d.SetId("")` and return `nil` (no error), so the next plan proposes a re-create.

### Requirement: Update MUST call ModifyCaptchaInfoInternational when any mutable field changes

The Update callback SHALL detect changes against the field set:

```
app_name, verify_rank, user_set_cap_type, defend_mode,
disable_invisible_switch, verify_domain, verify_bundle_id,
verify_package, check_appid_switch, check_iv_switch,
check_box_style, tags
```

`channel_info` is **not** in this set (it is `ForceNew`). When any field in the set has changed, Update SHALL build a single `ModifyCaptchaInfoInternationalRequest` populated from current `d.Get(...)` values for ALL mutable fields (full overwrite) and call `UseCaptchaClient().ModifyCaptchaInfoInternational` inside `resource.Retry(WriteRetryTimeout, ...)`. Each call MUST include `CaptchaAppId = d.Id()`. Update MUST end by re-invoking Read.

#### Scenario: Editing only `app_name`

- **GIVEN** state has `app_name = "old"`
- **WHEN** the user changes `app_name` to `"new"`
- **THEN** Update issues exactly one `ModifyCaptchaInfoInternational` request with `CaptchaAppId = d.Id()`, `AppName = "new"`, and the unchanged values of all other mutable fields.

#### Scenario: Editing `channel_info` triggers replacement

- **GIVEN** state has `channel_info = "a"`
- **WHEN** the user changes `channel_info` to `"b"`
- **THEN** Terraform's plan reports a destroy + create cycle (resource is replaced because `channel_info` is ForceNew). Update is NOT invoked.

### Requirement: Delete MUST call RemoveCaptchaInfoInternational

The Delete callback SHALL build a `RemoveCaptchaInfoInternationalRequest` with `CaptchaAppId = d.Id()` and call `UseCaptchaClient().RemoveCaptchaInfoInternational` inside `resource.Retry(WriteRetryTimeout, ...)`. No additional fields are sent.

#### Scenario: Idempotent delete on a missing resource

- **GIVEN** the resource is already gone in the cloud (e.g. external deletion between Read and Delete)
- **WHEN** Delete is called
- **THEN** the SDK error is propagated to Terraform, which surfaces it to the user; the provider does NOT swallow the error.

### Requirement: Every API call MUST be wrapped in resource.Retry

Every invocation of `UseCaptchaClient().CreateCaptchaInfoInternational`, `DescribeCaptchaInfoListInternational`, `ModifyCaptchaInfoInternational`, `RemoveCaptchaInfoInternational` SHALL be wrapped in `resource.Retry(...)` with the appropriate timeout (`ReadRetryTimeout` for read paths, `WriteRetryTimeout` for write paths) and SHALL forward errors via `tccommon.RetryError(e)`.

#### Scenario: Transient SDK error is retried

- **GIVEN** the first invocation of `CreateCaptchaInfoInternational` returns a retriable error
- **WHEN** the retry callback re-runs
- **THEN** the second attempt's response is observed; the original error is surfaced only after the retry budget is exhausted.

### Requirement: All response field reads MUST be nil-safe

Every `_ = d.Set("<key>", respData.<Field>)` access SHALL be guarded by `if respData.<Field> != nil` before dereference. Slice flatten loops SHALL also nil-check each item before reading.

#### Scenario: Optional field absent in response

- **GIVEN** the API returns `VerifyPackage == nil`
- **WHEN** Read runs
- **THEN** `d.Set("verify_package", ...)` is NOT called for that field; state retains the previous value.

### Requirement: Documentation and acceptance test MUST follow project conventions

- The HCL example markdown SHALL live at `tencentcloud/services/captcha/resource_tc_captcha_info_international.md` (mirroring `resource_tc_config_compliance_pack.md`) and contain at least one full HCL example and one Import example.
- An acceptance test SHALL live at `tencentcloud/services/captcha/resource_tc_captcha_info_international_test.go` (mirroring `resource_tc_config_compliance_pack_test.go`) covering: basic Create, Update of a mutable field, and `ImportState` round-trip.
- Running `make doc` SHALL regenerate `website/docs/r/captcha_info_international.html.markdown`.

#### Scenario: Generated website doc lists the resource

- **WHEN** `make doc` runs
- **THEN** `website/docs/r/captcha_info_international.html.markdown` exists and lists every schema attribute defined in the Schema requirement above.

#### Scenario: Acceptance test name and structure

- **WHEN** the test file is opened
- **THEN** the package is `captcha_test`, the test function is `TestAccTencentCloudCaptchaInfoInternationalResource_basic`, and it includes `tcacctest.AccPreCheck`, two HCL configs (initial and updated), and an `ImportState` step with `ImportStateVerify: true`.
