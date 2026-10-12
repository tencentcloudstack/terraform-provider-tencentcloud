## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/config/resource_tc_config_update_config_deliver.go` with `ResourceTencentCloudConfigUpdateConfigDeliver()` schema definition (status required int; deliver_name, target_arn, deliver_prefix, deliver_type optional strings; deliver_content_type optional int; create_time computed string) + Importer passthrough
- [x] 1.2 Implement Create handler: `d.SetId(helper.BuildToken())` then delegate to Update handler
- [x] 1.3 Implement Read handler: call `ConfigService.DescribeConfigDeliver(ctx)`, nil-check each response field before `d.Set`, on nil response log `[WARN]...tencentcloud_config_update_config_deliver [<id>] not found` and `d.SetId("")`
- [x] 1.4 Implement Update handler: build `UpdateConfigDeliverRequest` (GetOkExists for int fields status/deliver_content_type, GetOk for strings), wrap call in `resource.Retry(tccommon.WriteRetryTimeout, ...)` with `tccommon.RetryError`, then call Read
- [x] 1.5 Implement Delete handler: no-op (return nil)

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_config_update_config_deliver` → `config.ResourceTencentCloudConfigUpdateConfigDeliver()` in `tencentcloud/provider.go` ResourcesMap
- [x] 2.2 Add `tencentcloud_config_update_config_deliver` to the Config Resource comment index in `tencentcloud/provider.go`

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/config/resource_tc_config_update_config_deliver.md` with one-line description (mentioning Tencent Cloud Config) and Example Usage HCL block (do NOT add Argument/Attribute Reference sections)

## 4. Tests

- [x] 4.1 Create `tencentcloud/services/config/resource_tc_config_update_config_deliver_test.go` using gomonkey/gomock to mock the cloud API (DescribeConfigDeliver, UpdateConfigDeliver) and test Read/Update/Read-back business logic (NOT the TF acceptance test suite, per new-resource rules)

## 5. Verification (separate from code changes)

- [x] 5.1 Run `gofmt` on new/modified Go files (handled in finalize phase via tfpacer-finalize skill)
- [x] 5.2 Run `make doc` to generate website/docs documentation (handled in finalize phase via tfpacer-finalize skill)
- [x] 5.3 Create `.changelog` entry (handled in finalize phase via tfpacer-finalize skill)