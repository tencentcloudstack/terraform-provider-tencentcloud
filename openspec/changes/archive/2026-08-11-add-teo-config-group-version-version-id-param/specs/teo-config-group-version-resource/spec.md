## ADDED Requirements

### Requirement: Resource Schema Definition
The system SHALL define a Terraform resource `tencentcloud_teo_config_group_version` with the following schema fields:
- `zone_id` (Required, ForceNew, TypeString): 站点 ID。
- `group_id` (Required, ForceNew, TypeString): 配置组 ID。
- `description` (Optional, ForceNew, TypeString): 版本描述。
- `content` (Required, ForceNew, TypeString): 待导入的配置内容，JSON 格式，UTF-8 编码。
- `version_id` (Computed, TypeString): 配置组版本 ID，由 `DescribeConfigGroupVersionDetail` 接口的 `Response.ConfigGroupVersionInfo.VersionId` 返回。
- `version_number` (Computed, TypeString): 配置组版本号。
- `group_type` (Computed, TypeString): 配置组类型。
- `status` (Computed, TypeString): 版本生效状态。
- `create_time` (Computed, TypeString): 版本创建时间。

#### Scenario: Schema defines version_id as a computed output field
- **WHEN** the resource schema is defined
- **THEN** it SHALL include `version_id` as a `Computed` field of `TypeString`, sourced from `Response.ConfigGroupVersionInfo.VersionId` of the `DescribeConfigGroupVersionDetail` API

#### Scenario: ForceNew fields prevent in-place update
- **WHEN** `zone_id`, `group_id`, `description`, or `content` is changed in the Terraform configuration
- **THEN** the resource SHALL be destroyed and recreated

#### Scenario: version_id is read-only
- **WHEN** a user attempts to set `version_id` in the Terraform configuration
- **THEN** the schema SHALL treat `version_id` as a Computed-only field and SHALL NOT accept user input as a configuration parameter

### Requirement: Resource Read Operation
The system SHALL implement the Read operation by calling the `DescribeConfigGroupVersionDetail` API with the composite ID's `ZoneId`, `GroupId`, and `VersionId` segments, and SHALL populate `version_id` from the `Response.ConfigGroupVersionInfo.VersionId` field when it is not nil.

#### Scenario: Read populates version_id from API response
- **WHEN** the Read operation is invoked on an existing resource
- **AND** the `DescribeConfigGroupVersionDetail` API returns a non-nil `Response.ConfigGroupVersionInfo.VersionId`
- **THEN** the system SHALL set `version_id` in Terraform state to the returned value via `d.Set("version_id", ...)`

#### Scenario: Read handles nil VersionId gracefully
- **WHEN** the Read operation is invoked
- **AND** `Response.ConfigGroupVersionInfo` or `Response.ConfigGroupVersionInfo.VersionId` is nil
- **THEN** the system SHALL NOT call `d.Set("version_id", ...)` and SHALL NOT clear the composite resource id without first logging the id via `log.Printf("[CRUD] ...")`

#### Scenario: API response is empty
- **WHEN** the `DescribeConfigGroupVersionDetail` API returns an empty response (`respData == nil`)
- **THEN** the system SHALL log `[CRUD] teo_config_group_version id=<d.Id()>` to preserve the现场 and then call `d.SetId("")`

### Requirement: Composite Resource ID Format
The system SHALL construct the resource composite ID by joining `ZoneId`, `GroupId`, and `VersionId` with the `tccommon.FILED_SP` separator, and SHALL parse it back in Read/Delete operations.

#### Scenario: Create constructs composite ID
- **WHEN** the Create operation succeeds and the API returns a non-empty `VersionId`
- **THEN** the system SHALL set the resource id to `ZoneId + FILED_SP + GroupId + FILED_SP + VersionId`

#### Scenario: Read parses composite ID
- **WHEN** the Read operation is invoked
- **THEN** the system SHALL split `d.Id()` by `tccommon.FILED_SP` into exactly three segments (`zoneId`, `groupId`, `versionId`) and return an error if the segment count is not three

### Requirement: Unit Tests
The system SHALL provide unit tests in `resource_tc_teo_config_group_version_test.go` using gomonkey to mock the `DescribeConfigGroupVersionDetail` cloud API, testing the Read logic including the `version_id` output field population.

#### Scenario: Read test populates version_id
- **WHEN** a unit test mocks `DescribeConfigGroupVersionDetail` to return a response with `ConfigGroupVersionInfo.VersionId = "ver-2kplomhisdcb"`
- **AND** the Read function is invoked
- **THEN** the test SHALL assert `version_id` is set to `"ver-2kplomhisdcb"` in the resource state

#### Scenario: Read test handles nil ConfigGroupVersionInfo
- **WHEN** a unit test mocks `DescribeConfigGroupVersionDetail` to return a response with nil `ConfigGroupVersionInfo`
- **AND** the Read function is invoked
- **THEN** the test SHALL assert the resource id is cleared (after `[CRUD]` logging) and no panic occurs

### Requirement: Resource Documentation
The system SHALL provide a markdown documentation file `resource_tc_teo_config_group_version.md` with a one-line description mentioning TEO, example usage, and notes that `version_id` is a computed output field.

#### Scenario: Documentation mentions version_id
- **WHEN** the resource documentation is created or updated
- **THEN** the `.md` file SHALL describe `version_id` as a Computed output field sourced from the `DescribeConfigGroupVersionDetail` API
- **AND** the `website/docs/` documentation SHALL be regenerated via `make doc` during the finalize phase (not edited manually)
