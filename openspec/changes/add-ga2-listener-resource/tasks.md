## 1. Service-layer plumbing

- [x] 1.1 Add `DescribeGa2ListenerById(ctx, gaId, listenerId string) (*ga2v20250115.ListenerSet, error)` to `tencentcloud/services/ga2/service_tencentcloud_ga2.go`. Reuse the existing pattern from `DescribeGa2GlobalAcceleratorById`: build the `DescribeListenersRequest`, set `request.GlobalAcceleratorId`, and the `Filters` slice (`[{Name:"listener-id", Values:[listenerId]}]`) **outside** the for-loop; iterate with `Limit=100` (literal, no new constant); strict-equal on both `*item.ListenerId == listenerId` and `*item.GlobalAcceleratorId == gaId`; wrap each SDK page in `resource.Retry(tccommon.ReadRetryTimeout, ...)`; return `(nil, nil)` when not found; emit `[CRITAL]` log on retry failure.
- [x] 1.2 Confirm `connectivity.UseGa2V20250115Client` is already wired (no change required) and that the SDK package import path remains `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`.

## 2. Resource implementation

- [x] 2.1 Create `tencentcloud/services/ga2/resource_tc_ga2_listener.go` (single file). Top-level layout: package + imports → `ResourceTencentCloudGa2Listener()` schema → `resourceTencentCloudGa2ListenerCreate/Read/Update/Delete` → `parseGa2ListenerId` → `buildPortRanges` / `buildStringPointers` build helpers → `flattenPortRanges` / flatten helpers.
- [x] 2.2 Schema fields per the spec, in this declaration order: `global_accelerator_id` (Required, ForceNew), `name`, `port_ranges` (Required, ForceNew, MaxItems=1, nested `from_port`/`to_port` Required Ints), `description`, `listener_type` (ForceNew), `protocol` (ForceNew), `idle_timeout`, `get_real_ip_type`, `client_affinity`, `client_affinity_time`, `request_timeout`, `x_forwarded_for_real_ip`, `certification_type`, `cipher_policy_id`, `server_certificates` (TypeSet of String), `client_ca_certificates` (TypeSet of String), and computed-only: `listener_id`, `http_version`, `create_time`, `status`, `endpoint_group_counts`. Add a `Timeouts` block with `Create/Update/Delete: 5 * time.Minute`.
- [x] 2.3 Implement `Create`: build `CreateListenerRequest` from schema (translate `port_ranges[0]` into `*ga2v20250115.PortRanges`; convert `server_certificates`/`client_ca_certificates` Sets via `(*schema.Set).List()` then to `[]*string`); wrap SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; defend against nil `Response` / `ListenerId` / `TaskId`; assert `TaskId != ""` and call `Ga2Service.WaitForGa2TaskFinish(ctx, *resp.Response.TaskId, d.Timeout(schema.TimeoutCreate))`; finally `d.SetId(strings.Join([]string{gaId, listenerId}, tccommon.FILED_SP))` and return `Read`.
- [x] 2.4 Implement `Read`: parse composite ID into `(gaId, listenerId)` via `parseGa2ListenerId`; call `DescribeGa2ListenerById`; on `(nil, nil)` log `[WARN]` and `d.SetId("")` then return; otherwise `_ = d.Set(...)` for every input field (including the nested `port_ranges` block flattened) and every computed field. Use the "set only when non-nil" idiom (no `else` branches).
- [x] 2.5 Implement `Update`: short-circuit if no Modify-supported field changed (skip Modify call); otherwise build `ModifyListenerRequest` populated with `GlobalAcceleratorId` + `ListenerId` + every Modify-supported field whose schema getter returns non-zero. Wrap SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; defend against nil response/TaskId; call `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate))`.
- [x] 2.6 Implement `Delete`: call `DeleteListenerWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`; populate both `GlobalAcceleratorId` and `ListenerId` from the parsed composite ID; defend against nil `Response` / `TaskId`; call `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete))`.
- [x] 2.7 At the top of every CRUD function, add `defer tccommon.LogElapsed("resource.tencentcloud_ga2_listener.<op>")()` and `defer tccommon.InconsistentCheck(d, meta)()`.
- [x] 2.8 Importer: use `schema.ImportStatePassthrough` (matches `tencentcloud_ga2_endpoint_group` and `tencentcloud_ga2_global_accelerator`).
- [x] 2.9 ID parser: `parseGa2ListenerId(id string) (gaId, listenerId string, err error)` — split on `tccommon.FILED_SP`, expect exactly 2 non-empty parts, otherwise return a descriptive error.

## 3. Provider registration

- [x] 3.1 In `tencentcloud/provider.go`, register `"tencentcloud_ga2_listener": ga2.ResourceTencentCloudGa2Listener()` in `ResourcesMap`. Place it adjacent to `tencentcloud_ga2_endpoint_group` and `tencentcloud_ga2_global_accelerator` to keep the `ga2` namespace contiguous.
- [x] 3.2 In `tencentcloud/provider.md`, append `tencentcloud_ga2_listener` to the existing `Global Accelerator(GA2)` Resources section so `gendoc` picks it up on the next `make doc`.

## 4. Documentation example

- [x] 4.1 Create `tencentcloud/services/ga2/resource_tc_ga2_listener.md` (filename pattern matches `resource_tc_config_compliance_pack.md`). Content sections: short description; HCL example showing `global_accelerator_id`, `name`, `protocol`, a `port_ranges` block, plus a representative subset of optional fields; an `Import` section showing `terraform import tencentcloud_ga2_listener.example ga-xxxxxxxx#lsr-yyyyyyy`. **Do not** hand-edit `website/docs/r/ga2_listener.html.markdown` — that file is regenerated from the resource Schema/Description by `make doc`.

## 5. Acceptance test scaffolding

- [x] 5.1 Create `tencentcloud/services/ga2/resource_tc_ga2_listener_test.go` (filename pattern matches `resource_tc_config_compliance_pack_test.go`). Include `TestAccTencentCloudGa2ListenerResource_basic` with `resource.Test` containing at minimum: a create step, an update step (e.g. change `description`), and an `ImportState` verification step. Use `tcacctest.AccPreCheck` and the standard provider factories.

## 6. Build & verification

- [x] 6.1 Run `go build ./tencentcloud/...` and confirm clean compilation.
- [x] 6.2 Run `go vet ./tencentcloud/services/ga2/...` and confirm no new findings.
- [x] 6.3 Run `read_lints` on the changed files and confirm no new ERROR/WARN beyond the project-wide deprecated HINTs (`resource.Retry`, `schema.ImportStatePassthrough`, `d.GetOkExists`).
- [x] 6.4 Run `make doc` to regenerate `website/docs/r/ga2_listener.html.markdown` from the resource Schema/Description and the `.md` example file. Verify the generated doc renders the example block, all input fields, computed fields, the `port_ranges` nested block, and the import syntax. **Do not** hand-edit the generated file.
- [ ] 6.5 (Optional, requires `TENCENTCLOUD_SECRET_ID` / `TENCENTCLOUD_SECRET_KEY`) Run the acceptance test: `TF_ACC=1 go test -timeout 30m -run TestAccTencentCloudGa2ListenerResource_basic ./tencentcloud/services/ga2/...`.
