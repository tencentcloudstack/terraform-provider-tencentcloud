## 1. Resource Update Logic

- [x] 1.1 In `tencentcloud/services/emr/resource_tc_serverless_hbase_instance.go`, remove `disk_size` from the `immutableArgs` slice in `resourceTencentCloudServerlessHbaseInstanceUpdate`
- [x] 1.2 Add a new `if d.HasChange("disk_size")` branch (placed before the `zone_settings`/`node_num` block) that: constructs `emr.NewModifySLInstanceRequest()`, sets `InstanceId` from `d.Id()` and `DiskSize` from `d.Get("disk_size")`
- [x] 1.3 Inside the `disk_size` branch, wrap the `ModifySLInstanceWithContext` call in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with errors wrapped via `tccommon.RetryError`; on failure log `[CRITAL]%s update serverless hbase instance disk_size failed, reason:%+v` with `logId` and return the error
- [x] 1.4 After the retry block returns nil, build a state change conf via `tccommon.BuildStateChangeConf([]string{"22"}, []string{"2"}, 10*tccommon.ReadRetryTimeout, time.Second, emrService.SLInstanceStateRefreshFunc(instanceId, []string{}))` and call `WaitForState`; return its error on failure (mirror existing `node_num` wait at lines 488–494)
- [x] 1.5 Ensure the `disk_size` branch does not set `Zone`/`NodeNum`, and the existing `node_num` loop is unchanged, so the two change paths remain independent

## 2. Unit Tests

- [x] 2.1 In `tencentcloud/services/emr/resource_tc_serverless_hbase_instance_test.go`, add a test case that mocks the EMR client's `ModifySLInstanceWithContext` (via gomonkey) and `DescribeSLInstanceList` (status refresh) to verify the `disk_size` update calls `ModifySLInstance` with `InstanceId`+`DiskSize` and waits for status `2`
- [x] 2.2 Add a test case asserting that changing `disk_size` no longer triggers the `argument 'disk_size' cannot be changed` error (i.e. `disk_size` is no longer in `immutableArgs`)

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/emr/resource_tc_serverless_hbase_instance.md` to note that `disk_size` supports in-place expansion (storage capacity can be expanded; the call is async and waits until the instance returns to running state) if the example/description should reflect the new behavior
- [x] 3.2 Run `make doc` (in the finalize phase via tfpacer-finalize skill) to regenerate `website/docs/` markdown — do NOT manually edit website/ files

## 4. Verification (build / lint / test)

- [x] 4.1 Run `gofmt` on modified `.go` files (finalize phase via tfpacer-finalize skill)
- [x] 4.2 Run unit tests for the emr serverless hbase instance package to confirm the new disk_size update path works (note: do not run via `go test` during generation; verify by code review that tests are buildable)