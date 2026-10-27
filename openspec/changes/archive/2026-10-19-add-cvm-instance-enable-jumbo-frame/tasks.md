## 1. Service layer

- [x] 1.1 Add `ModifyEnableJumboFrame(ctx context.Context, instanceId string, enableJumboFrame bool) error` method to `tencentcloud/services/cvm/service_tencentcloud_cvm.go`, mirroring the existing `ModifyDisableApiTermination` (line ~297): build `cvm.NewModifyInstancesAttributeRequest()`, set `InstanceIds` and `EnableJumboFrame`, run ratelimit check, call `ModifyInstancesAttribute`, log on failure/success.

## 2. Schema definition

- [x] 2.1 In `tencentcloud/services/cvm/resource_tc_instance.go`, add the `enable_jumbo_frame` schema field (TypeBool, Optional, Computed) to the `ResourceTencentCloudInstance()` schema map

## 3. CRUD wiring

- [x] 3.1 In `resourceTencentCloudInstanceCreate`, add `enable_jumbo_frame` handling near the `disable_api_termination` Create block (~line 1066): when `d.GetOkExists("enable_jumbo_frame")` is ok, set `request.EnableJumboFrame = helper.Bool(v.(bool))`.
- [x] 3.2 In `resourceTencentCloudInstanceRead`, add `_ = d.Set("enable_jumbo_frame", instance.EnableJumboFrame)` near the `disable_api_termination` Read set (~line 1383).
- [x] 3.3 In `resourceTencentCloudInstanceUpdate`, add `enable_jumbo_frame` change handling near the `disable_api_termination` Update block (~line 2088): when `d.HasChange("enable_jumbo_frame")`, call `cvmService.ModifyEnableJumboFrame(ctx, instanceId, d.Get("enable_jumbo_frame").(bool))` and return on error.

## 4. Tests

- [x] 4.1 In `tencentcloud/services/cvm/resource_tc_instance_test.go`, add/update test cases covering `enable_jumbo_frame`: creation with the field set, update (change from false to true), and read reflecting the value. Since this modifies an existing resource, use the Terraform acceptance test suite (TF_ACC) consistent with the existing test file style.

## 5. Documentation

- [x] 5.1 Update `tencentcloud/services/cvm/resource_tc_instance.md` to document the `enable_jumbo_frame` parameter and add a usage example.
- [x] 5.2 Run `make doc` to regenerate the `website/docs/` markdown documentation from the resource .md example file (handled by the tfpacer-finalize stage per project rules; do not hand-edit website/docs/files).

## 6. Verification

- [x] 6.1 Verify the code compiles and resource schema/CRUD wiring is correct (build/lint handled by the standard CI flow; do not run `go build`/`go vet` locally per project rules).
- [x] 6.2 Confirm no other files reference a non-existent `EnableJumboFrame` schema key and that `enable_jumbo_frame` is consistently named across schema, Create, Read, Update, and tests.