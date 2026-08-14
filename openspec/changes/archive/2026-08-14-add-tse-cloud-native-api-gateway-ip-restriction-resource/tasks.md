# Tasks for `add-tse-cloud-native-api-gateway-ip-restriction-resource`

## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/tse/resource_tc_tse_cloud_native_api_gateway_ip_restriction.go`. File header, package (`tse`), imports MUST mirror `tencentcloud/services/tse/resource_tc_tse_cngw_certificate.go` style: imports of `context`, `fmt`, `log`, `strings`, terraform plugin sdk v2 helpers (`resource`, `schema`), the tse v20201207 SDK package (`tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"`), `tccommon`, `helper`.
- [x] 1.2 Declare `ResourceTencentCloudTseCloudNativeAPIGatewayIPRestriction() *schema.Resource` with `Create/Read/Update/Delete` callbacks and an `Importer` using `schema.ImportStatePassthrough`.
- [x] 1.3 Schema MUST declare exactly:
  - `gateway_id` (TypeString, Required, ForceNew, Description: "Gateway ID.")
  - `source_type` (TypeString, Required, ForceNew, Description: "Resource type bound to the IP restriction plugin: route|service.")
  - `source_id` (TypeString, Required, ForceNew, Description: "Route or service ID.")
  - `enabled` (TypeBool, Required, Description: "Whether to enable the plugin.")
  - `restriction_type` (TypeString, Required, Description: "IP restriction type: whiteList|blackList.")
  - `address_list` (TypeList of TypeString, Required, Description: "IP/CIDR address list.")
- [x] 1.4 Implement `resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionCreate`:
  - `defer tccommon.LogElapsed("resource.tencentcloud_tse_cloud_native_api_gateway_ip_restriction.create")()` and `defer tccommon.InconsistentCheck(d, meta)()`.
  - Build `var (logId, ctx)` block: `logId = tccommon.GetLogId(tccommon.ContextNil)`, `ctx = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)`.
  - Read `gateway_id`, `source_type`, `source_id` via `d.GetOk(...)`.
  - `d.SetId(strings.Join([]string{gatewayId, sourceType, sourceId}, tccommon.FILED_SP))`.
  - Return `resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionUpdate(d, meta)` (delegate to Update for the actual SDK call).
- [x] 1.5 Implement `resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionRead`:
  - `defer tccommon.LogElapsed(...)` and `defer tccommon.InconsistentCheck(d, meta)()`.
  - Split `d.Id()` by `tccommon.FILED_SP`; if `len != 3` return `fmt.Errorf("id is broken,%s", d.Id())`.
  - Assign `gatewayId`, `sourceType`, `sourceId`.
  - Build `request := tse.NewDescribeCloudNativeAPIGatewayIPRestrictionRequest()`, set `GatewayId`, `SourceType`, `SourceId`.
  - Wrap `UseTseClient().DescribeCloudNativeAPIGatewayIPRestrictionWithContext(ctx, request)` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`. Inside retry: on error return `tccommon.RetryError(e)`; on success log `[DEBUG]`. Guard `result == nil || result.Response == nil` → `resource.NonRetryableError`.
  - After retry: if `reqErr != nil`, log `[CRITAL]` and return. If `result.Response.Result == nil`, log `[CRUD] tse_cloud_native_api_gateway_ip_restriction id=%s` with `d.Id()` then `d.SetId("")` and `return nil`.
  - Set `gateway_id`, `source_type`, `source_id` from the split ID. Set `enabled`, `restriction_type`, `address_list` from `result.Response.Result` when each field is non-nil.
- [x] 1.6 Implement `resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionUpdate`:
  - `defer tccommon.LogElapsed(...)` and `defer tccommon.InconsistentCheck(d, meta)()`.
  - Build `var (logId, ctx)`.
  - Build `request := tse.NewCreateOrModifyCloudNativeAPIGatewayIPRestrictionRequest()`.
  - Populate `request.GatewayId` from `d.GetOk("gateway_id")`, `request.SourceType` from `d.GetOk("source_type")`, `request.SourceId` from `d.GetOk("source_id")`, `request.Enabled` from `d.GetOk("enabled")` (as `*bool`), `request.RestrictionType` from `d.GetOk("restriction_type")`, `request.AddressList` from `d.Get("address_list").([]interface{})` converting each element to `*string`.
  - Wrap `UseTseClient().CreateOrModifyCloudNativeAPIGatewayIPRestrictionWithContext(ctx, request)` in `resource.Retry(tccommon.WriteRetryTimeout, ...)`. Inside retry: on error `tccommon.RetryError(e)`; on success log `[DEBUG]`. Guard `result == nil || result.Response == nil` → `resource.NonRetryableError`.
  - On retry failure, log `[CRITAL]` and return the error.
  - End with `return resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionRead(d, meta)`.
- [x] 1.7 Implement `resourceTencentCloudTseCloudNativeAPIGatewayIPRestrictionDelete`:
  - `defer tccommon.LogElapsed(...)` and `defer tccommon.InconsistentCheck(d, meta)()`.
  - Build `var (logId, ctx)`.
  - Split `d.Id()` by `tccommon.FILED_SP`; if `len != 3` return broken-id error.
  - Build `request := tse.NewDeleteCloudNativeAPIGatewayIPRestrictionRequest()`, set `GatewayId`, `SourceType`, `SourceId`.
  - Wrap `UseTseClient().DeleteCloudNativeAPIGatewayIPRestrictionWithContext(ctx, request)` in `resource.Retry(tccommon.WriteRetryTimeout, ...)`. Inside retry: on error `tccommon.RetryError(e)`; on success log `[DEBUG]`. Guard `result == nil || result.Response == nil` → `resource.NonRetryableError`.
  - On retry failure, log `[CRITAL]` and return the error.

## 2. Provider Registration

- [x] 2.1 In `tencentcloud/provider.go`, locate the existing `tencentcloud_tse_*` resource registrations and append:
  ```go
  "tencentcloud_tse_cloud_native_api_gateway_ip_restriction": tse.ResourceTencentCloudTseCloudNativeAPIGatewayIPRestriction(),
  ```
  Keep adjacency with the other `tse.Resource...` entries.
- [x] 2.2 In `tencentcloud/provider.md`, locate the TSE Resource section and append a new line `tencentcloud_tse_cloud_native_api_gateway_ip_restriction` so gendoc picks it up.

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/tse/resource_tc_tse_cloud_native_api_gateway_ip_restriction.md` containing:
  - One-line summary starting with "Provides a resource to ..." and mentioning TSE (e.g. "Provides a resource to manage a TSE cloud native API gateway IP restriction.").
  - `Example Usage` HCL block with all six fields populated.
  - `Import` section stating the resource can be imported using the composite id `gateway_id#source_type#source_id`.
  - Do NOT include `Argument Reference` or `Attribute Reference` sections (auto-generated).

## 4. Unit Test

- [x] 4.1 Create `tencentcloud/services/tse/resource_tc_tse_cloud_native_api_gateway_ip_restriction_test.go`, package `tse_test`, using gomonkey to mock the TSE client (mirror `tencentcloud/services/bh/resource_tc_bh_bind_device_resource_test.go` pattern):
  - Declare a `mockMeta` implementing `tccommon.ProviderMeta` with a `*connectivity.TencentCloudClient`.
  - Mock `UseTseClient` to return `&tse.Client{}`.
  - `Test...Create`: mock `CreateOrModifyCloudNativeAPIGatewayIPRestrictionWithContext` returning a valid response; mock `DescribeCloudNativeAPIGatewayIPRestrictionWithContext` returning `Result{Enabled, RestrictionType, AddressList}`; assert `d.Id()` equals the composite `gw-1#route#r-1` and state fields set.
  - `Test...Read`: pre-set `d.SetId("gw-1#route#r-1")`, mock Describe returning a populated `Result`; assert state fields and id preserved.
  - `Test...ReadNotFound`: mock Describe returning nil Result; assert `d.Id()` becomes `""`.
  - `Test...Update`: mock `CreateOrModify...WithContext` and `Describe...WithContext`; assert no error.
  - `Test...Delete`: mock `DeleteCloudNativeAPIGatewayIPRestrictionWithContext`; assert no error.

## 5. Validation

- [x] 5.1 Ensure all function error returns are checked; functions that cannot fail use `_ = func()` to avoid unused-variable errors.
- [x] 5.2 Confirm the new code compiles cleanly (no `go build`/`go vet` execution required; review by inspection).
- [x] 5.3 Confirm the `.md` example file and test file exist and follow the conventions in `gendoc/README.md` and sibling resources.
