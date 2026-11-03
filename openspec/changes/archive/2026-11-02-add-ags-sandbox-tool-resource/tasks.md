## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/ags/` directory and `resource_tc_ags_sandbox_tool.go` with schema definition including all fields: `tool_name`, `tool_type`, `network_configuration` (nested: `network_mode`, `vpc_config` → `subnet_ids`, `security_group_ids`), `description`, `default_timeout`, `tags`, `client_token`, `role_arn`, `storage_mounts` (nested), `custom_configuration` (nested), `computer_configuration` (nested: `waa_configuration`, `os_world_configuration`), `log_configuration` (nested: `cls_config`, `log_sources`), `persistent`, and computed fields: `tool_id`, `status`, `default_timeout_seconds`, `create_time`, `update_time`, `status_reason`, `image_digest`
- [x] 1.2 Implement `resourceTencentCloudAgsSandboxToolCreate` function using `CreateSandboxTool` API with retry logic, nil response/id checks, and `d.SetId(tool_id)`
- [x] 1.3 Implement `resourceTencentCloudAgsSandboxToolRead` function using `DescribeSandboxToolList` API, handle empty response with `d.SetId("")`, set all computed fields
- [x] 1.4 Implement `resourceTencentCloudAgsSandboxToolUpdate` function using `UpdateSandboxTool` API, check `d.HasChange()` for mutable fields only
- [x] 1.5 Implement `resourceTencentCloudAgsSandboxToolDelete` function using `DeleteSandboxTool` API with retry logic
- [x] 1.6 Add resource import support with `schema.ImportStatePassthrough`

## 2. Provider Registration

- [x] 2.1 Register `tencentcloud_ags_sandbox_tool` resource in `tencentcloud/provider.go` (import + DataSourcesMap/ResourcesMap)
- [x] 2.2 Add resource entry in `tencentcloud/provider.md`

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/ags/resource_tc_ags_sandbox_tool.md` with usage example covering minimum required fields and VPC network mode

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/ags/resource_tc_ags_sandbox_tool_test.go` with mock-based unit tests for CRUD operations using gomonkey