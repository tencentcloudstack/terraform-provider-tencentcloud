## Why

Tencent Cloud Global Accelerator V2 (`ga2`) currently exposes `tencentcloud_ga2_endpoint_group` (already shipped) and `tencentcloud_ga2_global_accelerator` (about to ship). The middle tier of the GA2 object hierarchy — the **listener** — is still missing. Without it, users cannot wire `endpoint_group.listener_id` to a Terraform-managed listener, forcing them to create listeners manually in the console. Adding `tencentcloud_ga2_listener` closes this gap and makes the full `Accelerator → Listener → EndpointGroup` chain fully Terraform-native.

## What Changes

- Add a new resource `tencentcloud_ga2_listener` backed by the `ga2` v20250115 SDK.
- Map every `CreateListener` request parameter to a schema field, in particular: `global_accelerator_id`, `name`, `port_ranges`, `description`, `listener_type`, `protocol`, `idle_timeout`, `get_real_ip_type`, `client_affinity`, `request_timeout`, `x_forwarded_for_real_ip`, `certification_type`, `cipher_policy_id`, `server_certificates`, `client_ca_certificates`.
- Implement async-aware CRUD: `CreateListener`, `ModifyListener`, `DeleteListener` all return a `TaskId` that must be polled via `DescribeTaskResult` until `Status == "SUCCESS"`. Reuse the existing `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` helper.
- Add a new service helper `Ga2Service.DescribeGa2ListenerById(ctx, gaId, listenerId) (*ga2v20250115.ListenerSet, error)` that wraps `DescribeListeners` with a `listener-id` filter and pagination capped at the documented maximum (`Limit=100`).
- Surface read-only computed fields from the describe response that are not part of `CreateListener`: `client_affinity_time`, `http_version`, `create_time`, `status`, `endpoint_group_counts`.
- Wire the new resource into `tencentcloud/provider.go` under the `ga2` namespace.
- Author resource markdown documentation `resource_tc_ga2_listener.md` (example HCL snippet + `terraform import` syntax).
- Author acceptance-test scaffolding `resource_tc_ga2_listener_test.go`.
- Resource ID is the composite `<GlobalAcceleratorId>#<ListenerId>` (using `tccommon.FILED_SP`), since `DescribeListeners`/`Modify*`/`Delete*` all require both IDs.
- All SDK calls are wrapped with `resource.Retry` (write paths use `tccommon.WriteRetryTimeout`, read paths use `tccommon.ReadRetryTimeout`).
- `Timeouts` block defaults to **5 minutes** for Create/Update/Delete, matching the existing `tencentcloud_ga2_endpoint_group` and `tencentcloud_ga2_global_accelerator` resources.
- `port_ranges`, `listener_type`, `protocol` are **ForceNew** because `ModifyListener` does not accept them.

## Capabilities

### New Capabilities
- `ga2-listener-resource`: Lifecycle management (create / read / update / delete / import) of a Tencent Cloud Global Accelerator V2 listener, including async task polling, full schema parity with `CreateListener`, and exposure of all `ListenerSet` computed fields.

### Modified Capabilities
<!-- None: this change only introduces a new resource; it does not alter requirement-level behavior of any existing capability. -->

## Impact

- **New code**:
  - `tencentcloud/services/ga2/resource_tc_ga2_listener.go` (CRUD + schema + build/flatten helpers, single file, mirroring `tencentcloud_igtm_monitor` style).
  - `tencentcloud/services/ga2/resource_tc_ga2_listener.md` (resource doc + import syntax, mirroring `resource_tc_config_compliance_pack.md` filename convention).
  - `tencentcloud/services/ga2/resource_tc_ga2_listener_test.go` (acceptance test skeleton, mirroring `resource_tc_config_compliance_pack_test.go` filename convention).
- **Modified code**:
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go`: add `DescribeGa2ListenerById`.
  - `tencentcloud/provider.go`: register `tencentcloud_ga2_listener` under the `ga2` namespace block.
  - `tencentcloud/provider.md`: add `tencentcloud_ga2_listener` under the existing `Global Accelerator(GA2)` Resources section, so `make doc` picks it up.
- **APIs consumed**: `CreateListener`, `DescribeListeners`, `ModifyListener`, `DeleteListener`, `DescribeTaskResult` (all already vendored in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/`).
- **No breaking change**: purely additive; no existing schema or state is modified.
- **No SDK upgrade required**: all required APIs are already present in the vendored SDK.
