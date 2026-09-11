# captcha-ip-white-list-international-resource Specification

## ADDED Requirements

### Requirement: Resource MUST be registered as `tencentcloud_captcha_ip_white_list_international`

The provider SHALL register a new CRUD-type resource named `tencentcloud_captcha_ip_white_list_international` whose Create/Read/Update/Delete callbacks invoke the captcha international IP allowlist APIs (`CreateIpWhiteListInternational`, `DescribeIpWhiteListInternational`, `ModifyIpWhiteListInternational`, `DeleteIpWhiteListInternational`) of `tencentcloud-sdk-go-intl-en/tencentcloud/captcha/v20190722`.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_captcha_ip_white_list_international"` mapped to `captcha.ResourceTencentCloudCaptchaIpWhiteListInternational()`.

#### Scenario: Importer is configured

- **WHEN** a user runs `terraform import tencentcloud_captcha_ip_white_list_international.example <CaptchaAppid>#<Id>`
- **THEN** the importer SHALL pass the supplied composite ID straight into the Read callback (using `schema.ImportStatePassthrough`), and Read SHALL hydrate state from `DescribeIpWhiteListInternational` filtered by that CaptchaAppid.

### Requirement: Schema MUST mirror the CreateIpWhiteListInternational API input 1:1

The resource schema SHALL declare exactly these 4 top-level argument keys, with semantics matching the SDK request fields of `CreateIpWhiteListInternationalRequest`:

| HCL key | SDK field | Type | Required | ForceNew |
|---|---|---|---|---|
| `name` | `Name` | TypeString | Yes | No |
| `captcha_appid` | `CaptchaAppid` | TypeInt | Yes | **Yes** |
| `ip` | `Ip` | TypeString | Yes | **Yes** |
| `comment` | `Comment` | TypeString | No | No |

Per the `CreateIpWhiteListInternational` API documentation, `name`, `captcha_appid`, and `ip` are `Required: true`, while `comment` is `Optional: true`.

#### Scenario: Schema field set matches the API input exactly

- **WHEN** the resource schema is inspected
- **THEN** it declares exactly the 4 keys above (no extra field, no missing field), each with the correct SDK type mapping.

#### Scenario: ip and captcha_appid are ForceNew

- **WHEN** the user changes `ip` or `captcha_appid`
- **THEN** Terraform reports a destroy + create cycle (because `ModifyIpWhiteListInternational` has no `Ip` input, and `CaptchaAppid` is part of the resource identity).

### Requirement: Resource ID MUST be the composite `CaptchaAppid` + `IdList[0]`

After Create, the resource's Terraform ID SHALL be set to `<captcha_appid>#<id>` where `<captcha_appid>` is the request `CaptchaAppid` and `<id>` is `*response.Response.IdList[0]`, joined by the project separator `#`. Read/Update/Delete SHALL parse `d.Id()` back into `captchaAppid` and `id`.

#### Scenario: ID is set from CreateIpWhiteListInternational response

- **GIVEN** `captcha_appid = 179000003` and a successful Create returns `IdList = [2100000000]`
- **WHEN** Create completes
- **THEN** `d.Id()` returns `"179000003#2100000000"`.

#### Scenario: Create response missing IdList is fatal

- **WHEN** `CreateIpWhiteListInternational` returns a non-nil response with `len(IdList) == 0` or `IdList[0] == nil`
- **THEN** the Create callback SHALL return a non-retryable error mentioning `IdList is empty`.

### Requirement: Read MUST locate the resource via DescribeIpWhiteListInternational + CaptchaAppid

The service-layer helper `(*CaptchaService).DescribeCaptchaIpWhiteListInternationalById(ctx, captchaAppid, id int64) (*captcha.DescribeCaptchaWhiteListItem, error)` SHALL build a `DescribeIpWhiteListInternationalRequest` with:

- `CaptchaAppid = helper.Int64(captchaAppid)`
- `PageIndex = 1`
- `PageSize = 100` (the API documented maximum)

…and invoke `UseCaptchaClient().DescribeIpWhiteListInternational` inside `resource.Retry(ReadRetryTimeout, ...)`, returning the list item whose `Id == id` or `(nil, nil)` when absent.

If the helper returns nil, the Read callback SHALL call `d.SetId("")` to mark the resource gone.

#### Scenario: Resource exists

- **GIVEN** the API returns a white list item with `CaptchaAppid=179000003` and `Id=2100000000`
- **WHEN** Read runs
- **THEN** all schema attributes are populated from the API response with nil-safe pointer dereferencing.

#### Scenario: Resource removed out-of-band

- **GIVEN** the user runs `terraform refresh` after the white list item was deleted in the web console
- **WHEN** Read calls `DescribeIpWhiteListInternational` and no item matches `Id`
- **THEN** Read SHALL `d.SetId("")` and return `nil` (no error), so the next plan proposes a re-create.

### Requirement: Update MUST call ModifyIpWhiteListInternational when any mutable field changes

The Update callback SHALL detect changes against the field set `["name", "comment"]`. `ip` and `captcha_appid` are **not** in this set (both `ForceNew`). When any field in the set has changed, Update SHALL build a `ModifyIpWhiteListInternationalRequest` with `Id` and `CaptchaAppid` parsed from the composite ID, populate `Name`/`Comment` from current `d.Get(...)` (full overwrite), and call `UseCaptchaClient().ModifyIpWhiteListInternational` inside `resource.Retry(WriteRetryTimeout, ...)`. Update MUST end by re-invoking Read.

#### Scenario: Editing only `comment`

- **GIVEN** state has `comment = "old"`
- **WHEN** the user changes `comment` to `"new"`
- **THEN** Update issues exactly one `ModifyIpWhiteListInternational` request with `Id`/`CaptchaAppid` parsed from `d.Id()`, `Name` unchanged, and `Comment = "new"`.

#### Scenario: Editing `ip` triggers replacement

- **GIVEN** state has `ip = "1.2.3.4"`
- **WHEN** the user changes `ip` to `"5.6.7.8"`
- **THEN** Terraform's plan reports a destroy + create cycle (resource is replaced because `ip` is ForceNew). Update is NOT invoked.

### Requirement: Delete MUST call DeleteIpWhiteListInternational

The Delete callback SHALL build a `DeleteIpWhiteListInternationalRequest` with `CaptchaAppid` and `Id` parsed from the composite ID, and call `UseCaptchaClient().DeleteIpWhiteListInternational` inside `resource.Retry(WriteRetryTimeout, ...)`. No additional fields are sent.

#### Scenario: Idempotent delete on a missing resource

- **GIVEN** the resource is already gone in the cloud (e.g. external deletion between Read and Delete)
- **WHEN** Delete is called
- **THEN** the SDK error is propagated to Terraform, which surfaces it to the user; the provider does NOT swallow the error.

### Requirement: Every API call MUST be wrapped in resource.Retry

Every invocation of `UseCaptchaClient().CreateIpWhiteListInternational`, `DescribeIpWhiteListInternational`, `ModifyIpWhiteListInternational`, `DeleteIpWhiteListInternational` SHALL be wrapped in `resource.Retry(...)` with the appropriate timeout (`ReadRetryTimeout` for read paths, `WriteRetryTimeout` for write paths) and SHALL forward errors via `tccommon.RetryError(e)`.

#### Scenario: Transient SDK error is retried

- **GIVEN** the first invocation of `CreateIpWhiteListInternational` returns a retriable error
- **WHEN** the retry callback re-runs
- **THEN** the second attempt's response is observed; the original error is surfaced only after the retry budget is exhausted.

### Requirement: All response field reads MUST be nil-safe

Every `_ = d.Set("<key>", respData.<Field>)` access SHALL be guarded by `if respData.<Field> != nil` before dereference. Slice flatten loops SHALL also nil-check each item before reading.

#### Scenario: Optional field absent in response

- **GIVEN** the API returns `Comment == nil`
- **WHEN** Read runs
- **THEN** `d.Set("comment", ...)` is NOT called for that field; state retains the previous value.

### Requirement: Documentation and acceptance test MUST follow project conventions

- The HCL example markdown SHALL live at `tencentcloud/services/captcha/resource_tc_captcha_ip_white_list_international.md` (mirroring `resource_tc_config_compliance_pack.md`) and contain at least one full HCL example and one Import example.
- An acceptance test SHALL live at `tencentcloud/services/captcha/resource_tc_captcha_ip_white_list_international_test.go` (mirroring `resource_tc_config_compliance_pack_test.go`) covering: basic Create, Update of `comment`, and `ImportState` round-trip.
- Running `make doc` SHALL regenerate `website/docs/r/captcha_ip_white_list_international.html.markdown`.

#### Scenario: Generated website doc lists the resource

- **WHEN** `make doc` runs
- **THEN** `website/docs/r/captcha_ip_white_list_international.html.markdown` exists and lists every schema attribute defined in the Schema requirement above.

#### Scenario: Acceptance test name and structure

- **WHEN** the test file is opened
- **THEN** the package is `captcha_test`, the test function is `TestAccTencentCloudCaptchaIpWhiteListInternationalResource_basic`, and it includes `tcacctest.AccPreCheck`, two HCL configs (initial and updated), and an `ImportState` step with `ImportStateVerify: true`.
