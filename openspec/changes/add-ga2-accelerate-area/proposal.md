## Why

Tencent Cloud Global Accelerator V2 (`ga2`) already exposes the core object hierarchy as Terraform resources: `tencentcloud_ga2_global_accelerator`, `tencentcloud_ga2_listener`, `tencentcloud_ga2_endpoint_group`, and `tencentcloud_ga2_forwarding_rule`. What is still missing is the **acceleration region** (`AccelerateArea`) tier — the per-region acceleration entry that defines the access region, ISP type, bandwidth and IP version attached to a global accelerator instance. Without it, users must add/remove acceleration regions in the console, breaking the fully-declarative GA2 workflow. Adding `tencentcloud_ga2_accelerate_area` lets operators manage acceleration regions natively in HCL.

## What Changes

- Add a new CRUD resource `tencentcloud_ga2_accelerate_area` backed by the vendored `ga2` v20250115 SDK.
- Map every `CreateAccelerateAreas` request parameter to a schema field (flattening the single `AcceleratorAreas` element the resource manages): `global_accelerator_id`, `accelerate_region`, `bandwidth`, `isp_type`, `ip_version`, `ip_address`.
- Implement async-aware CRUD: `CreateAccelerateAreas`, `ModifyAccelerateAreas`, `DeleteAccelerateAreas` each return only a `TaskId` that MUST be polled via `DescribeTaskResult` until `Status == "SUCCESS"`. Reuse the existing `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` helper.
- Resolve the resource's `AcceleratorAreaId` **after** Create: `CreateAccelerateAreas` does NOT return `AcceleratorAreaId`, so the resource queries `DescribeAccelerateAreas` with `GlobalAcceleratorId` and matches on `AccelerateRegion` to discover the generated `AcceleratorAreaId`.
- Add two service helpers to `service_tencentcloud_ga2.go`:
  - `DescribeGa2AccelerateAreaById(ctx, gaId, areaId) (*ga2v20250115.AcceleratorAreas, error)` — used by Read/Update/Delete, matches on `AcceleratorAreaId`.
  - `DescribeGa2AccelerateAreaByRegion(ctx, gaId, region) (*ga2v20250115.AcceleratorAreas, error)` — used by Create to resolve the new `AcceleratorAreaId` from the natural key `(gaId, region)`.
- Surface read-only computed fields not part of `CreateAccelerateAreas`: `accelerator_area_id`, and the nested `ip_address_info_set` (`ip_address`, `isp_type`).
- Wire the new resource into `tencentcloud/provider.go` under the `ga2` namespace and append it to the `Global Accelerator(GA2)` Resources section of `tencentcloud/provider.md` so `make doc` generates its website page.
- Author resource markdown documentation `resource_tc_ga2_accelerate_area.md` (HCL example + `terraform import` syntax).
- Author acceptance-test scaffolding `resource_tc_ga2_accelerate_area_test.go`.
- Resource ID is the composite `<GlobalAcceleratorId>#<AcceleratorAreaId>` (using `tccommon.FILED_SP`), because `Modify*`/`Delete*` require `AcceleratorAreaId` and `DescribeAccelerateAreas` requires `GlobalAcceleratorId`.
- All SDK calls are wrapped in `resource.Retry` (write paths use `tccommon.WriteRetryTimeout`, read paths use `tccommon.ReadRetryTimeout`); list pages use `Limit=100` (documented maximum).
- `Timeouts` defaults to **5 minutes** for Create/Update/Delete, matching the other GA2 resources.
- `global_accelerator_id` and `accelerate_region` are **ForceNew** — `accelerate_region` is the natural key used to resolve the area ID, and an area cannot be moved between accelerators.

## Capabilities

### New Capabilities
- `ga2-accelerate-area-resource`: Lifecycle management (create / read / update / delete / import) of a Tencent Cloud Global Accelerator V2 acceleration region, including async task polling, post-create `AcceleratorAreaId` resolution by region, full schema parity with `CreateAccelerateAreas`, and exposure of all `AcceleratorAreas` computed fields.

### Modified Capabilities
<!-- None: this change only introduces a new resource; it does not alter requirement-level behavior of any existing capability. -->

## Impact

- **New code**:
  - `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.go` (CRUD + schema + build/flatten helpers, single file, mirroring `tencentcloud_igtm_monitor` style).
  - `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area.md` (resource doc + import syntax, mirroring `resource_tc_config_compliance_pack.md` filename convention).
  - `tencentcloud/services/ga2/resource_tc_ga2_accelerate_area_test.go` (acceptance test skeleton, mirroring `resource_tc_config_compliance_pack_test.go` filename convention).
- **Modified code**:
  - `tencentcloud/services/ga2/service_tencentcloud_ga2.go`: add `DescribeGa2AccelerateAreaById` and `DescribeGa2AccelerateAreaByRegion`.
  - `tencentcloud/provider.go`: register `tencentcloud_ga2_accelerate_area` under the `ga2` namespace block.
  - `tencentcloud/provider.md`: add `tencentcloud_ga2_accelerate_area` under the existing `Global Accelerator(GA2)` Resources section, so `make doc` picks it up.
- **APIs consumed**: `CreateAccelerateAreas`, `DescribeAccelerateAreas`, `ModifyAccelerateAreas`, `DeleteAccelerateAreas`, `DescribeTaskResult` (all already vendored in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/`).
- **No breaking change**: purely additive; no existing schema or state is modified.
- **No SDK upgrade required**: all required APIs are confirmed present in the vendored SDK (`CreateAccelerateAreasRequest`, `DescribeAccelerateAreasRequest`/`Response`, `ModifyAccelerateAreasRequest`, `DeleteAccelerateAreasRequest`, plus the `AcceleratorAreas` and `IpAddressInfoSet` models).
