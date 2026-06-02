# Tasks for `add-bh-bind-device-account-kubeconfig-resource`

## 1. Resource implementation

- [x] 1.1 Create `tencentcloud/services/bh/resource_tc_bh_bind_device_account_kubeconfig.go`. File header, package, imports MUST mirror `tencentcloud/services/waf/resource_tc_waf_owasp_rule_status_config.go` style: package `bh`, imports of `context`, `fmt`, `log`, terraform plugin sdk v2 helpers, the bh v20230418 SDK package, `tccommon`, `helper`. Avoid importing `strings` (no compound ID).
- [x] 1.2 Declare `ResourceTencentCloudBhBindDeviceAccountKubeconfig() *schema.Resource` with `Create/Read/Update/Delete` callbacks. Do NOT declare an `Importer` (no query API ⇒ import would yield empty state).
- [x] 1.3 Schema MUST declare exactly:
  - `account_id` (TypeInt, Required, ForceNew, Description: "Container account Id. Maps to the SDK request field `Id`. Renamed in HCL because `id` is reserved by the Terraform Plugin SDK as the resource's internal identifier.")
  - `kubeconfig` (TypeString, Required, Sensitive, Description: "Container account kubeconfig credential.")
  - `manage_dimension` (TypeInt, Optional, Description: "Manage dimension. 1 means cluster.")
- [x] 1.4 Implement `resourceTencentCloudBhBindDeviceAccountKubeconfigCreate`:
  - `defer tccommon.LogElapsed(...)` and `defer tccommon.InconsistentCheck(d, meta)()` at the top.
  - Read HCL `account_id` via `d.GetOkExists("account_id")`; if absent, return an error.
  - `d.SetId(fmt.Sprintf("%d", account_id))`.
  - Return `resourceTencentCloudBhBindDeviceAccountKubeconfigUpdate(d, meta)` (delegate to Update for the actual SDK call).
- [x] 1.5 Implement `resourceTencentCloudBhBindDeviceAccountKubeconfigRead`:
  - `defer tccommon.LogElapsed(...)` and `defer tccommon.InconsistentCheck(d, meta)()`.
  - `return nil` immediately.
- [x] 1.6 Implement `resourceTencentCloudBhBindDeviceAccountKubeconfigUpdate`:
  - `defer tccommon.LogElapsed(...)` and `defer tccommon.InconsistentCheck(d, meta)()`.
  - Build `var (logId, ctx)` block per project conventions (`tccommon.GetLogId`, `tccommon.NewResourceLifeCycleHandleFuncContext`).
  - Build `request := bhv20230418.NewBindDeviceAccountKubeconfigRequest()`.
  - Populate `request.Id` from `d.GetOkExists("account_id")` via `helper.IntUint64(v.(int))`.
  - Populate `request.Kubeconfig` from `d.GetOk("kubeconfig")` via `helper.String(v.(string))`.
  - Populate `request.ManageDimension` from `d.GetOkExists("manage_dimension")` via `helper.IntUint64(v.(int))` (only if present).
  - Wrap `UseBhV20230418Client().BindDeviceAccountKubeconfigWithContext(ctx, request)` in `resource.Retry(tccommon.WriteRetryTimeout, ...)`. Inside the retry, validate `result == nil || result.Response == nil` → `resource.NonRetryableError`. Log `[DEBUG]` with `request.GetAction()`, `request.ToJsonString()`, `result.ToJsonString()` on success (matching the project pattern).
  - On retry failure, log `[CRITAL]` (matching project pattern of misspelling `CRITICAL` as `CRITAL`) and return the error.
  - End with `return resourceTencentCloudBhBindDeviceAccountKubeconfigRead(d, meta)`.
- [x] 1.7 Implement `resourceTencentCloudBhBindDeviceAccountKubeconfigDelete`:
  - `defer tccommon.LogElapsed(...)` and `defer tccommon.InconsistentCheck(d, meta)()`.
  - `return nil` immediately.

## 2. Provider registration

- [x] 2.1 In `tencentcloud/provider.go`, locate the existing `tencentcloud_bh_*` resource registrations and append:
  ```go
  "tencentcloud_bh_bind_device_account_kubeconfig": bh.ResourceTencentCloudBhBindDeviceAccountKubeconfig(),
  ```
  Keep adjacency with the other `bh.Resource...` entries.
- [x] 2.2 In `tencentcloud/provider.md`, locate the `Bastion Host(BH)` Resource section and append a new line `tencentcloud_bh_bind_device_account_kubeconfig` so gendoc picks it up.

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/bh/resource_tc_bh_bind_device_account_kubeconfig.md` containing:
  - One-line summary line.
  - One `~> **NOTE:**` paragraph explicitly stating: (a) the API does not provide a query endpoint, so Read is a no-op and external drift is invisible; (b) the API does not provide an unbind endpoint, so `terraform destroy` only removes state and does NOT unbind on the backend.
  - One `Example Usage` HCL block with all three fields populated.
  - Do NOT include an `Import` section (the resource intentionally does not support import).
- [x] 3.2 Run `make doc` to regenerate `website/docs/r/bh_bind_device_account_kubeconfig.html.markdown`. Hand-editing the generated file is forbidden.

## 4. Acceptance test

- [x] 4.1 Create `tencentcloud/services/bh/resource_tc_bh_bind_device_account_kubeconfig_test.go`, package `bh_test`, structure mirroring `resource_tc_config_compliance_pack_test.go`:
  - Test function `TestAccTencentCloudBhBindDeviceAccountKubeconfigResource_basic` with `t.Parallel()`.
  - `PreCheck: func() { tcacctest.AccPreCheck(t) }`, `Providers: tcacctest.AccProviders`.
  - Two test steps:
    1. Initial config with `account_id`, `kubeconfig = "test-kubeconfig-v1"`, `manage_dimension = 1`. Assert `resource.TestCheckResourceAttrSet("tencentcloud_bh_bind_device_account_kubeconfig.example", "id")` and `resource.TestCheckResourceAttr(..., "manage_dimension", "1")`.
    2. Updated config: change `kubeconfig = "test-kubeconfig-v2"`. Assert with `resource.TestCheckResourceAttr(..., "kubeconfig", "test-kubeconfig-v2")`.
  - Do NOT add an `ImportState` step.

## 5. Validation

- [x] 5.1 `go build ./tencentcloud/...` clean.
- [x] 5.2 `go vet ./tencentcloud/services/bh/...` clean.
- [x] 5.3 `read_lints` on the new files shows no new errors/warnings beyond the project-wide deprecated hints (resource.Retry, etc.).
- [x] 5.4 `openspec validate add-bh-bind-device-account-kubeconfig-resource --strict` passes.
