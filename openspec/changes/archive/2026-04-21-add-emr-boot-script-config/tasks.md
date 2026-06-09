## 1. Resource Implementation

- [x] 1.1 Create `resource_tc_emr_boot_script.go`:
  - Schema: `instance_id` (Required/ForceNew), `boot_type` (Required/ForceNew), `pre_executed_file_settings` (Optional, List of `PreExecuteFileSetting` fields)
  - Create handler: `d.SetId(instanceId + "#" + bootType)`, call Update
  - Read handler: `DescribeBootScript` via service layer, read the correct slice by `boot_type`, set `pre_executed_file_settings`
  - Update handler: `ModifyBootScript` wrapped in `resource.Retry(WriteRetryTimeout)`, then call Read
  - Delete handler: `ModifyBootScript` with empty `PreExecutedFileSettings` list

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_emr_boot_script` in `provider.go`

## 3. Documentation & Tests

- [x] 3.1 Create `resource_tc_emr_boot_script.md`
- [x] 3.2 Create `resource_tc_emr_boot_script_test.go`

## 4. Refactor

- [x] 4.1 Rename from `tencentcloud_emr_boot_script_config` to `tencentcloud_emr_boot_script` (general resource, not config resource)
- [x] 4.2 Implement real Delete: call `ModifyBootScript` with empty `PreExecutedFileSettings`
