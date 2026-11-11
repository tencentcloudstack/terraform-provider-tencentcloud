## 1. Resource Schema Changes

- [x] 1.1 Add `disk_encrypt` schema field (TypeInt, Optional, immutable after creation, description "Disk encryption.") to `ResourceTencentCloudCdwdorisInstance()` in `tencentcloud/services/cdwdoris/resource_tc_cdwdoris_instance.go`
- [x] 1.2 Add `encrypt` schema field (TypeInt, Computed, description "Whether disk encryption is enabled. Read from InstanceInfo.MasterSummary.Encrypt / CoreSummary.Encrypt.") to `ResourceTencentCloudCdwdorisInstance()`

## 2. Create Function Changes

- [x] 2.1 In `resourceTencentCloudCdwdorisInstanceCreate`, read `disk_encrypt` from schema data (use `d.GetOkExists`) and set `request.DiskEncrypt = helper.IntInt64(v.(int))` when present, placed before the `resource.Retry` block (alongside the other request field assignments)
- [x] 2.2 Verify `request.DiskEncrypt` assignment is placed outside the retry block; the retry block only contains the API call and response nil-check

## 3. Read Function Changes

- [x] 3.1 In `resourceTencentCloudCdwdorisInstanceRead`, after retrieving `respData`, refresh the Computed `encrypt` field: if `respData.MasterSummary != nil && respData.MasterSummary.Encrypt != nil`, set `encrypt` to `respData.MasterSummary.Encrypt`; else if `respData.CoreSummary != nil && respData.CoreSummary.Encrypt != nil`, set `encrypt` to `respData.CoreSummary.Encrypt`. Only call `d.Set("encrypt", ...)` when a value is available.

## 4. Update Function Changes

- [x] 4.1 Add `disk_encrypt` to the `immutableArgs` array in `resourceTencentCloudCdwdorisInstanceUpdate` so changes are rejected with `argument 'disk_encrypt' cannot be changed`

## 5. Unit Test Changes

- [x] 5.1 Add unit test cases in `tencentcloud/services/cdwdoris/resource_tc_cdwdoris_instance_is_ssc_test.go` using gomonkey to mock the CDWDoris cloud API calls, covering: create with `disk_encrypt` set (assert `request.DiskEncrypt` is passed), create without `disk_encrypt` (assert nil), read populating `encrypt` from `MasterSummary.Encrypt`, read populating `encrypt` from `CoreSummary.Encrypt` fallback, and update returning immutable error when `disk_encrypt` changes

## 6. Documentation

- [x] 6.1 Documentation (`website/` and service `.md`) is generated via `make doc` during the finalize phase, not edited directly.

## 7. Validation

- [x] 7.1 Verify the code compiles successfully
- [x] 7.2 Verify no lint errors