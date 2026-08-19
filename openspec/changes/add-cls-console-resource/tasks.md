# Tasks for `add-cls-console-resource`

## 1. Service layer

- [x] 1.1 Add `DescribeClsConsoleById(ctx context.Context, consoleId string) (ret *cls.Console, errRet error)` to `tencentcloud/services/cls/service_tencentcloud_cls.go`. Build a `DescribeConsolesRequest` with `Limit = 100` and `Filters = [{Key: "ConsoleId", Values: [consoleId]}]`, call `UseClsClient().DescribeConsoles(request)` wrapped in `resource.Retry(tccommon.ReadRetryTimeout, ...)`, return the first element of `Response.Consoles` or `(nil, nil)` when the slice is empty.

## 2. Resource implementation

- [x] 2.1 Create `tencentcloud/services/cls/resource_tc_cls_console.go`. File header, package, imports MUST mirror `tencentcloud/services/igtm/resource_tc_igtm_monitor.go` (same `defer LogElapsed + InconsistentCheck` pattern, same `var (...)` block style).
- [x] 2.2 Declare `ResourceTencentCloudClsConsole() *schema.Resource` with `Create/Read/Update/Delete` callbacks and `Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough}`.
- [x] 2.3 Schema MUST declare every attribute listed in the spec's "Schema MUST mirror the CreateConsole API input 1:1" requirement, with the exact `Required/Optional/ForceNew/Computed` flags shown in the table. Sensitive fields (`password`, `secret_id`, `secret_key`) MUST set `Sensitive: true`.
- [x] 2.4 Implement `resourceTencentCloudClsConsoleCreate`:
  - Build `cls.NewCreateConsoleRequest()`.
  - Map every schema field via `d.GetOk` / `d.GetOkExists`. Use `d.GetOkExists` for `login_mode` and `intranet_type` (TypeInt zero is meaningful: `LoginMode=0` means account-password).
  - For nested list fields, walk each map and populate the SDK nested struct (`Accounts`, `AnonymousLogin`, `AuthRoles`, `AccessControlRules`, `Tags`).
  - For `access_mode`, `hide_params`, `menus` (string slices), iterate the `[]interface{}` and append `helper.String(v.(string))`.
  - Wrap the SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`. Inside the retry, validate `result == nil || result.Response == nil || result.Response.ConsoleId == nil` → `resource.NonRetryableError`.
  - On success, `d.SetId(*response.Response.ConsoleId)`, then return `resourceTencentCloudClsConsoleRead(d, meta)`.
- [x] 2.5 Implement `resourceTencentCloudClsConsoleRead`:
  - Build the typical `logId/ctx/service/consoleId := d.Id()` block.
  - Call `service.DescribeClsConsoleById(ctx, consoleId)`. If err, return err. If `respData == nil`, log warn + `d.SetId("")` + return nil.
  - For every field, guard with `if respData.<Field> != nil` before `d.Set`. For nested structs (e.g. `AnonymousLogin`), check the outer pointer, then build a single-element slice of map.
  - For nested list fields, build `[]map[string]interface{}` with per-element nil checks.
  - Set `console_id`, `domain`, `intranet_domain` (Computed) at the end.
- [x] 2.6 Implement `resourceTencentCloudClsConsoleUpdate`:
  - Use the `mutableArgs` pattern from `resource_tc_igtm_monitor.go`. The list MUST be:
    `["access_mode", "login_mode", "domain_prefix", "accounts", "anonymous_login", "intranet_type", "intranet_region", "vpc_id", "subnet_id", "auth_roles", "hide_params", "access_control_rules", "remarks", "menus"]`.
  - If any of those changed, build `cls.NewModifyConsoleRequest()`, populate ALL mutable fields from current `d.Get(...)` (full overwrite), set `request.ConsoleId = helper.String(d.Id())`, wrap in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.
  - Always end by calling `resourceTencentCloudClsConsoleRead(d, meta)`.
- [x] 2.7 Implement `resourceTencentCloudClsConsoleDelete`:
  - Build `cls.NewDeleteConsoleRequest()`, set `request.ConsoleId = helper.String(d.Id())`, wrap call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.

## 3. Provider registration

- [x] 3.1 In `tencentcloud/provider.go`, locate the existing `tencentcloud_cls_*` resource registrations (around line 2049) and append:
  ```go
  "tencentcloud_cls_console": cls.ResourceTencentCloudClsConsole(),
  ```
  Keep adjacency with the other `cls.Resource...` entries.

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/cls/resource_tc_cls_console.md` containing:
  - One-line summary line.
  - One `Example Usage` HCL block using `login_mode = 0` (account-password).
  - One `Example Usage with anonymous login` HCL block (login_mode = 1).
  - One `Import` section showing `terraform import tencentcloud_cls_console.example ds-xxxxxxxx`.
- [x] 4.2 Run `make doc` to regenerate `website/docs/r/cls_console.html.markdown`. Hand-editing the generated file is forbidden.

## 5. Acceptance test

- [x] 5.1 Create `tencentcloud/services/cls/resource_tc_cls_console_test.go`, package `cls_test`, structure mirroring `resource_tc_config_compliance_pack_test.go`:
  - Test function `TestAccTencentCloudClsConsoleResource_basic` with `t.Parallel()`.
  - `PreCheck: func() { tcacctest.AccPreCheck(t) }`, `Providers: tcacctest.AccProviders`.
  - Steps:
    1. Initial config (basic create): `resource.TestCheckResourceAttrSet("tencentcloud_cls_console.example", "id")`, `resource.TestCheckResourceAttrSet("tencentcloud_cls_console.example", "console_id")`, `resource.TestCheckResourceAttr("tencentcloud_cls_console.example", "login_mode", "0")`.
    2. Updated config: change `remarks`. Assert new value via `resource.TestCheckResourceAttr`.
    3. ImportState step with `ImportStateVerify: true`.
  - Two `const` HCL strings (initial + updated). Use a unique `domain_prefix` referencing test timestamp or random suffix is **not** required at test design time — fixed prefixes are fine.

## 6. Validation

- [x] 6.1 `go build ./tencentcloud/...` clean.
- [x] 6.2 `go vet ./tencentcloud/services/cls/...` clean.
- [x] 6.3 `read_lints` on the new files shows no new errors/warnings.
- [x] 6.4 `openspec validate add-cls-console-resource --strict` passes.
