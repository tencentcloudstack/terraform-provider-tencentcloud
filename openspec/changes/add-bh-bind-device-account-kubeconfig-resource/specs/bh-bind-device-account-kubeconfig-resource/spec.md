# bh-bind-device-account-kubeconfig-resource Specification

## ADDED Requirements

### Requirement: Resource MUST be registered as `tencentcloud_bh_bind_device_account_kubeconfig`

The provider SHALL register a new config-type resource named `tencentcloud_bh_bind_device_account_kubeconfig` whose Create / Update callbacks invoke the `BindDeviceAccountKubeconfig` API of `tencentcloud-sdk-go/tencentcloud/bh/v20230418`.

#### Scenario: Resource registered in provider map

- **WHEN** the provider is loaded
- **THEN** `provider.go` exposes the resource via key `"tencentcloud_bh_bind_device_account_kubeconfig"` mapped to `bh.ResourceTencentCloudBhBindDeviceAccountKubeconfig()`, alongside existing `tencentcloud_bh_*` resource entries.

#### Scenario: Resource appears in gendoc index

- **WHEN** `tencentcloud/provider.md` is scanned by `make doc`
- **THEN** the Bastion Host(BH) Resource section MUST include `tencentcloud_bh_bind_device_account_kubeconfig` so that `website/docs/r/bh_bind_device_account_kubeconfig.html.markdown` is generated.

### Requirement: Schema MUST mirror the BindDeviceAccountKubeconfig API input (with SDK-reserved-name remap)

The resource schema SHALL declare exactly these top-level argument keys, with semantics matching the SDK request fields of `BindDeviceAccountKubeconfigRequestParams`:

| HCL key | SDK field | Type | Required | ForceNew | Sensitive |
|---|---|---|---|---|---|
| `account_id` | `Id` | TypeInt | Yes | **Yes** | No |
| `kubeconfig` | `Kubeconfig` | TypeString | Yes | No | **Yes** |
| `manage_dimension` | `ManageDimension` | TypeInt | No | No | No |

The HCL key `account_id` maps 1:1 to the SDK request field `Id`. The rename is required because `id` is reserved by Terraform Plugin SDK v2 for the resource's internal identifier and cannot be declared as a top-level schema key. The mapping MUST be documented in the resource markdown via a NOTE.

The schema MUST NOT declare any additional Computed-only fields, since `BindDeviceAccountKubeconfigResponseParams` exposes nothing beyond `RequestId`.

#### Scenario: Required fields enforce on plan

- **WHEN** the user writes a config that omits `account_id` or `kubeconfig`
- **THEN** `terraform plan` SHALL fail validation pointing at the missing required attribute.

#### Scenario: Sensitive flag hides credential output

- **WHEN** `terraform plan` shows changes touching `kubeconfig`
- **THEN** the diff line displays `(sensitive value)` instead of the raw kubeconfig content.

#### Scenario: Changing `account_id` forces replacement

- **GIVEN** state has `account_id = 100`
- **WHEN** the user changes `account_id` to `200` in HCL
- **THEN** Terraform's plan reports a destroy + create cycle (the resource is replaced because `account_id` is ForceNew).

### Requirement: Resource ID MUST be the string form of the HCL `account_id`

After Create, the Terraform resource ID SHALL be `fmt.Sprintf("%d", d.Get("account_id").(int))`. No compound separator is used. Read SHALL NOT touch the resource ID.

#### Scenario: ID is set from the HCL `account_id` field

- **GIVEN** HCL declares `account_id = 12345`
- **WHEN** Create completes
- **THEN** `d.Id()` returns `"12345"`.

### Requirement: Create MUST delegate to Update

To keep the SDK call site DRY (Create and Update share the same `BindDeviceAccountKubeconfig` API), the Create callback SHALL:

1. Read the HCL `account_id`, set `d.SetId(fmt.Sprintf("%d", account_id))`.
2. Return `resourceTencentCloudBhBindDeviceAccountKubeconfigUpdate(d, meta)`.

The Update callback contains the actual SDK call.

#### Scenario: Create produces a single API call

- **GIVEN** HCL with `account_id`, `kubeconfig`, `manage_dimension`
- **WHEN** the user runs `terraform apply`
- **THEN** exactly one `BindDeviceAccountKubeconfig` API request is issued with `Id = account_id`, `Kubeconfig`, and `ManageDimension` populated; `d.Id()` returns the integer string of the supplied `account_id`.

### Requirement: Update MUST call BindDeviceAccountKubeconfig with all current values

The Update callback SHALL build a `BindDeviceAccountKubeconfigRequest` populated from the current `d.Get(...)` values for `account_id` (mapped to `request.Id`), `kubeconfig` (mapped to `request.Kubeconfig`), and (if set) `manage_dimension` (mapped to `request.ManageDimension`), and call `UseBhV20230418Client().BindDeviceAccountKubeconfigWithContext(ctx, request)` inside `resource.Retry(WriteRetryTimeout, ...)`. Any change to `kubeconfig` or `manage_dimension` therefore triggers a full re-bind, which the API supports as overwrite-semantic. Update MUST end by calling Read.

#### Scenario: Editing `kubeconfig` triggers exactly one Bind call

- **GIVEN** state has `kubeconfig = "old-content"`
- **WHEN** the user changes it to `"new-content"`
- **THEN** Update issues exactly one `BindDeviceAccountKubeconfig` request with `Id = <state account_id>`, `Kubeconfig = "new-content"`, `ManageDimension = <current value or omitted>`.

#### Scenario: Transient SDK error is retried

- **GIVEN** the first invocation of `BindDeviceAccountKubeconfig` returns a retriable error (e.g. internal error)
- **WHEN** the retry callback re-runs
- **THEN** the second attempt's response is observed; the original error is surfaced only after the retry budget is exhausted.

### Requirement: Read MUST be a no-op

The Read callback SHALL `return nil` immediately. It MUST NOT call any SDK API. It MUST NOT modify state. It MUST NOT clear `d.Id()`.

This is by design because the BH API does not expose a query endpoint to inspect an existing kubeconfig binding; refreshing state would be impossible.

#### Scenario: Read leaves state untouched

- **GIVEN** state has `account_id = 12345`, `kubeconfig = "abc"`, `manage_dimension = 1`
- **WHEN** `terraform refresh` runs
- **THEN** state remains `account_id = 12345`, `kubeconfig = "abc"`, `manage_dimension = 1`, and no API call is made.

### Requirement: Delete MUST be a no-op

The Delete callback SHALL `return nil` immediately. It MUST NOT call any SDK API.

The BH API does not currently expose an `UnbindDeviceAccountKubeconfig` (or equivalent) endpoint. Deleting from Terraform only removes the resource from state; the binding on the backend is preserved. The resource markdown MUST document this clearly so users understand `terraform destroy` does NOT remove the underlying binding.

#### Scenario: terraform destroy removes only state

- **GIVEN** state has the resource bound
- **WHEN** the user runs `terraform destroy`
- **THEN** Terraform removes the resource from state without making any SDK call; the kubeconfig binding on the backend remains intact.

### Requirement: Every API call MUST be wrapped in resource.Retry

The single SDK invocation `UseBhV20230418Client().BindDeviceAccountKubeconfigWithContext(ctx, request)` SHALL be wrapped in `resource.Retry(tccommon.WriteRetryTimeout, ...)` and forward errors via `tccommon.RetryError(e)`.

#### Scenario: Retry wrapper present

- **WHEN** a code reviewer inspects the Update callback
- **THEN** the SDK call is inside a `resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError { ... })` block, not a bare invocation.

### Requirement: Response field reads MUST be nil-safe

Even though the SDK response only exposes `RequestId`, the implementation SHALL guard `result == nil || result.Response == nil` and return a non-retryable error when either is nil, before attempting any further use of `result`.

#### Scenario: Nil response is detected

- **GIVEN** the SDK returns `(result == nil, err == nil)` (defensive case)
- **WHEN** the retry callback runs
- **THEN** the callback returns `resource.NonRetryableError(...)` with a clear message and does NOT panic.

### Requirement: Documentation and acceptance test MUST follow project conventions

- The HCL example markdown SHALL live at `tencentcloud/services/bh/resource_tc_bh_bind_device_account_kubeconfig.md` (mirroring `resource_tc_config_compliance_pack.md`) and contain at least one full HCL example, plus an explicit NOTE explaining (a) the absence of a query API (state is authoritative; drift is invisible) and (b) the absence of an unbind API (`terraform destroy` does NOT remove the backend binding).
- An acceptance test SHALL live at `tencentcloud/services/bh/resource_tc_bh_bind_device_account_kubeconfig_test.go` (mirroring `resource_tc_config_compliance_pack_test.go`) covering: basic Create, Update of `kubeconfig`. **No** `ImportState` step is included (no query API ⇒ import would produce empty state).
- Running `make doc` SHALL regenerate `website/docs/r/bh_bind_device_account_kubeconfig.html.markdown`.

#### Scenario: Generated website doc lists the resource

- **WHEN** `make doc` runs
- **THEN** `website/docs/r/bh_bind_device_account_kubeconfig.html.markdown` exists and lists every schema attribute defined in the Schema requirement above.

#### Scenario: Acceptance test name and structure

- **WHEN** the test file is opened
- **THEN** the package is `bh_test`, the test function is `TestAccTencentCloudBhBindDeviceAccountKubeconfigResource_basic`, and it includes `tcacctest.AccPreCheck`, two HCL configs (initial and updated `kubeconfig`), but does NOT include an `ImportState` step.
