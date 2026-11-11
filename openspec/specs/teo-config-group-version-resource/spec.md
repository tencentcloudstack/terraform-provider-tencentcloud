# teo-config-group-version-resource Specification

## Purpose
TBD - created by archiving change add-teo-config-group-version-source-version. Update Purpose after archive.
## Requirements
### Requirement: Resource Schema Definition

The resource `tencentcloud_teo_config_group_version` SHALL define the following schema fields:
- `zone_id` (Required, ForceNew, String): 站点 ID
- `group_id` (Required, ForceNew, String): 待新建版本的配置组 ID
- `description` (Optional, ForceNew, String): 版本描述
- `content` (Required, ForceNew, String): 待导入的配置内容，JSON 格式
- `source_version` (Optional, ForceNew, String): 新版本所基于的来源版本 ID，新版本将在该来源版本的配置基础上派生创建；未传入时默认采用当前生产环境正在生效的版本作为来源版本
- `version_id` (Computed, String): 版本 ID
- `version_number` (Computed, String): 版本号
- `group_type` (Computed, String): 配置组类型
- `status` (Computed, String): 版本状态
- `create_time` (Computed, String): 版本创建时间

The resource ID SHALL be composed of `zone_id`, `group_id`, and `version_id` joined by `tccommon.FILED_SP` separator. The resource SHALL support import using the composite ID format `zone_id#group_id#version_id`.

#### Scenario: Schema defines source_version as optional and ForceNew
- **WHEN** the resource schema is defined
- **THEN** it SHALL include `source_version` field with type `schema.TypeString`
- **AND** the field SHALL be `Optional: true` and `ForceNew: true`
- **AND** the Description SHALL explain that it specifies the source version ID from which the new version is derived

#### Scenario: source_version not set uses cloud default behavior
- **WHEN** a user creates a `tencentcloud_teo_config_group_version` resource without setting `source_version`
- **THEN** the system SHALL NOT set `request.SourceVersion` (leave it nil)
- **AND** the cloud API SHALL use the currently active production version as the source version by default

### Requirement: Create Operation

The resource Create method SHALL call the `CreateConfigGroupVersion` API with the following parameter mapping:
- `zone_id` → `request.ZoneId`
- `group_id` → `request.GroupId`
- `content` → `request.Content`
- `description` → `request.Description` (if set)
- `source_version` → `request.SourceVersion` (if set)

The Create method SHALL use `resource.Retry(tccommon.WriteRetryTimeout, ...)` for retry logic. After the API call succeeds, the Create method SHALL check that `response.Response.VersionId` is not empty; if empty, return `NonRetryableError`. The resource ID SHALL be set to `zone_id#group_id#version_id`.

#### Scenario: Successful creation with source_version
- **WHEN** user applies a configuration with a valid `source_version`
- **THEN** the resource SHALL call `CreateConfigGroupVersion` API with `request.SourceVersion` set to the configured value
- **AND** after success, the resource SHALL set the composite ID from `zone_id`, `group_id`, and the returned `version_id`

#### Scenario: Successful creation without source_version
- **WHEN** user applies a configuration without setting `source_version`
- **THEN** the resource SHALL call `CreateConfigGroupVersion` API without setting `request.SourceVersion`
- **AND** the cloud API SHALL use the default source version

#### Scenario: Create API returns error
- **WHEN** `CreateConfigGroupVersion` API returns a retryable error
- **THEN** the resource SHALL wrap the error with `tccommon.RetryError()` and retry within `tccommon.WriteRetryTimeout`
- **AND** on non-retryable errors, the resource SHALL return the error directly

#### Scenario: Create returns empty VersionId
- **WHEN** the Create API succeeds but `response.Response.VersionId` is nil or empty string
- **THEN** the resource SHALL return a `NonRetryableError` to avoid writing an empty ID

### Requirement: Read Operation

The resource Read method SHALL parse the composite ID to extract `zone_id`, `group_id`, and `version_id`, then call `DescribeConfigGroupVersionDetail` API via `TeoService.DescribeTeoConfigGroupVersionById`. The Read method SHALL read fields from `ConfigGroupVersionInfo` including `SourceVersion`, performing nil checks before each `d.Set` call.

#### Scenario: Successful read with source_version
- **WHEN** the Read method queries a config group version that has a `SourceVersion` value
- **THEN** the resource SHALL check `respData.ConfigGroupVersionInfo.SourceVersion` is not nil
- **AND** set `source_version` in state via `d.Set("source_version", ...)`

#### Scenario: Read with nil SourceVersion
- **WHEN** the Read method queries a config group version where `ConfigGroupVersionInfo.SourceVersion` is nil
- **THEN** the resource SHALL NOT call `d.Set("source_version", ...)`
- **AND** no error SHALL be returned for the nil field

#### Scenario: Resource not found
- **WHEN** `DescribeConfigGroupVersionDetail` returns empty response or `ConfigGroupVersionInfo` is nil
- **THEN** the resource SHALL log `log.Printf("[CRUD] teo_config_group_version id=%s", d.Id())` before clearing the ID
- **AND** call `d.SetId("")` to mark the resource as deleted

### Requirement: Delete Operation

The resource Delete method SHALL be a no-op (the config group version resource does not support cloud-side deletion via Terraform).

#### Scenario: Delete is a no-op
- **WHEN** user destroys the `tencentcloud_teo_config_group_version` resource
- **THEN** the Delete method SHALL return nil without calling any cloud API

### Requirement: Provider Registration

The resource `tencentcloud_teo_config_group_version` SHALL already be registered in `tencentcloud/provider.go` and documented in `tencentcloud/provider.md`. No new registration is required for this parameter addition.

#### Scenario: Resource available in provider
- **WHEN** user references `tencentcloud_teo_config_group_version` in their configuration
- **THEN** Terraform SHALL recognize it as a valid resource type

### Requirement: Unit Tests with gomonkey mock

The resource SHALL have unit tests in `resource_tc_teo_config_group_version_test.go` that use gomonkey to mock the cloud API calls (`CreateConfigGroupVersionWithContext` and `DescribeConfigGroupVersionDetail`), covering the `source_version` parameter scenarios.

#### Scenario: Unit test for Create with source_version
- **WHEN** the unit test for Create with `source_version` set is executed
- **THEN** it SHALL mock `CreateConfigGroupVersionWithContext` to return a valid `VersionId`
- **AND** mock `DescribeConfigGroupVersionDetail` to return the `SourceVersion` value
- **AND** verify the resource is created correctly with `source_version` populated in state

#### Scenario: Unit test for Create without source_version
- **WHEN** the unit test for Create without `source_version` is executed
- **THEN** it SHALL mock the APIs and verify the resource is created successfully
- **AND** verify `request.SourceVersion` is not set (nil)

#### Scenario: Unit test for Read with nil SourceVersion
- **WHEN** the unit test for Read where the API response has nil `SourceVersion` is executed
- **THEN** it SHALL mock `DescribeConfigGroupVersionDetail` with nil `SourceVersion`
- **AND** verify no error occurs and the resource state is populated with other fields

### Requirement: Resource Documentation

The system SHALL update the markdown documentation file `resource_tc_teo_config_group_version.md` to include `source_version` in the example usage. The documentation SHALL NOT include `Argument Reference` or `Attribute Reference` sections (auto-generated by tooling).

#### Scenario: Documentation includes source_version example
- **WHEN** the resource documentation is updated
- **THEN** the example usage SHALL demonstrate the `source_version` field
- **AND** the one-line description SHALL mention TEO (EdgeOne)
- **AND** an import section SHALL be present showing the composite ID format `zone_id#group_id#version_id`

