## 1. Service-layer plumbing

- [x] 1.1 Add `DescribeGa2GlobalAcceleratorById(ctx, gaId string) (*ga2v20250115.GlobalAcceleratorSet, error)` to `tencentcloud/services/ga2/service_tencentcloud_ga2.go`. Reuse the existing pattern from `DescribeGa2EndpointGroupById`: build the `DescribeGlobalAcceleratorsRequest` and `Filters` slice **outside** the for-loop; iterate with `Limit=100` (literal, no new constant); strict-equal on `*item.GlobalAcceleratorId == gaId`; wrap each SDK page in `resource.Retry(tccommon.ReadRetryTimeout, ...)`; return `(nil, nil)` when not found; emit `[CRITAL]` log on retry failure.
- [x] 1.2 Confirm `connectivity.UseGa2V20250115Client` is already wired (no change required) and that the SDK package import path remains `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`.

## 2. Resource implementation

- [x] 2.1 Create `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.go` (single file). Top-level layout: package + imports → `ResourceTencentCloudGa2GlobalAccelerator()` schema → `resourceTencentCloudGa2GlobalAcceleratorCreate/Read/Update/Delete` → `buildGa2GlobalAcceleratorTags` helper → no other helpers needed.
- [x] 2.2 Schema fields per the spec, in this declaration order: `name`, `instance_charge_type` (ForceNew), `description`, `cross_border_type`, `cross_border_promise_flag`, `tags`, `global_accelerator_id` (computed), `state`, `status`, `cname`, `ddos_id`, `create_time`, `listener_counts`, `accelerator_area_counts`. Add a `Timeouts` block with `Create/Update/Delete: 5 * time.Minute` defaults.
- [x] 2.3 Implement `Create`: build `CreateGlobalAcceleratorRequest` from schema (translate `tags` map into `[]*ga2v20250115.Tag` inline); wrap SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; defend against nil `Response` / `GlobalAcceleratorId` / `TaskId`; assert `TaskId != ""` and call `Ga2Service.WaitForGa2TaskFinish(ctx, *resp.Response.TaskId, d.Timeout(schema.TimeoutCreate))`; finally `d.SetId(*resp.Response.GlobalAcceleratorId)` and return `Read`.
- [x] 2.4 Implement `Read`: call `DescribeGa2GlobalAcceleratorById`; on `(nil, nil)` log `[WARN]` and `d.SetId("")` then return; otherwise `_ = d.Set(...)` for every input field and every computed field; populate `tags` from `set.TagSet` via the standard provider tag-helper pattern (build `map[string]string`, then `d.Set("tags", tagsMap)`).
- [x] 2.5 Implement `Update`: parse old/new `tags`; if any of `name` / `description` / `cross_border_type` / `cross_border_promise_flag` changed, call `ModifyGlobalAcceleratorWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`, then `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate))`; if `tags` changed, reconcile via `svctag.NewTagService(...).ModifyTags(ctx, tccommon.BuildTagResourceName("ga2", "globalAccelerator", region, gaId), replaceTags, deleteTags)`. Skip the Modify call when only tags changed.
- [x] 2.6 Implement `Delete`: call `DeleteGlobalAcceleratorWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`; defend against nil `Response` / `TaskId`; call `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete))`.
- [x] 2.7 At the top of every CRUD function, add `defer tccommon.LogElapsed("resource.tencentcloud_ga2_global_accelerator.<op>")()` and `defer tccommon.InconsistentCheck(d, meta)()`.
- [x] 2.8 Importer: use `schema.ImportStatePassthrough` (matches the existing `tencentcloud_ga2_endpoint_group` resource for consistency).

## 3. Provider registration

- [x] 3.1 In `tencentcloud/provider.go`, register `"tencentcloud_ga2_global_accelerator": ga2.ResourceTencentCloudGa2GlobalAccelerator()` in `ResourcesMap`. Place it adjacent to the existing `tencentcloud_ga2_endpoint_group` entry to keep the `ga2` namespace contiguous.

## 4. Documentation example

- [x] 4.1 Create `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator.md` (filename pattern matches `resource_tc_config_compliance_pack.md`). Content sections: short description; HCL example showing all primary input fields plus a `tags` map; an `Import` section showing `terraform import tencentcloud_ga2_global_accelerator.example ga-xxxxxxxx`. **Do not** hand-edit `website/docs/r/ga2_global_accelerator.html.markdown` — that file will be (re)generated from the resource Schema/Description by `make doc`.

## 5. Acceptance test scaffolding

- [x] 5.1 Create `tencentcloud/services/ga2/resource_tc_ga2_global_accelerator_test.go` (filename pattern matches `resource_tc_config_compliance_pack_test.go`). Include `TestAccTencentCloudGa2GlobalAcceleratorResource_basic` with `resource.Test` containing at minimum: a create step, an update step (e.g. change `description` and `tags`), and an `ImportState` verification step. Use `tcacctest.AccPreCheck` and the standard provider factories.

## 6. Build & verification

- [x] 6.1 Run `go build ./tencentcloud/...` and confirm clean compilation.
- [x] 6.2 Run `go vet ./tencentcloud/services/ga2/...` and confirm no new findings.
- [x] 6.3 Run `read_lints` (or `make lint` if available) on the two changed files and address any newly introduced (non-pre-existing) issues.
- [x] 6.4 Run `make doc` to regenerate `website/docs/r/ga2_global_accelerator.html.markdown` from the resource Schema/Description and the `.md` example file. Verify the generated doc renders the example block, all fields, computed fields, and the import syntax. **Do not** hand-edit the generated file.
- [ ] 6.5 (Optional, requires `TENCENTCLOUD_SECRET_ID` / `TENCENTCLOUD_SECRET_KEY`) Run the acceptance test: `TF_ACC=1 go test -timeout 30m -run TestAccTencentCloudGa2GlobalAcceleratorResource_basic ./tencentcloud/services/ga2/...`.
