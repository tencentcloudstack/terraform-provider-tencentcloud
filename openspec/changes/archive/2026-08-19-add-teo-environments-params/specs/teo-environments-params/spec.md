## ADDED Requirements

### Requirement: Data source schema field for source_version
The data source `tencentcloud_teo_environments` SHALL provide the following additional output field:
- `source_version` (Computed, string) nested under `env_infos.current_config_group_version_infos`: The source version ID that a config group version was derived from, mapped from `response.EnvInfos[].CurrentConfigGroupVersionInfos[].SourceVersion` of the `DescribeEnvironments` API.

#### Scenario: source_version is populated from API response
- **WHEN** an element of `CurrentConfigGroupVersionInfos` in the `DescribeEnvironments` API response contains a non-nil `SourceVersion`
- **THEN** the data source SHALL set the nested `source_version` field under the corresponding `env_infos.current_config_group_version_infos` element to that string value

#### Scenario: source_version is nil in API response
- **WHEN** an element of `CurrentConfigGroupVersionInfos` in the `DescribeEnvironments` API response contains a nil `SourceVersion`
- **THEN** the data source SHALL leave `source_version` unset for that element and SHALL NOT return an error

### Requirement: Backward compatibility
The addition of `source_version` SHALL be backward compatible. The new field is a `Computed` output-only field and SHALL NOT alter existing schema fields, input parameters, or the data source query behavior.

#### Scenario: Existing configurations remain valid
- **WHEN** a user applies an existing `tencentcloud_teo_environments` data source configuration
- **THEN** the configuration SHALL remain valid and the read SHALL succeed, with the new field simply populated as an additional output

### Requirement: Documentation update
The documentation file `data_source_tc_teo_environments.md` SHALL be updated to reflect the new `source_version` output field, following the project documentation guidelines (no manual `Argument Reference` / `Attribute Reference` sections, which are auto-generated).

#### Scenario: Documentation reflects new field
- **WHEN** the implementation is complete
- **THEN** the documentation file SHALL describe the `source_version` nested field

### Requirement: Unit tests
The test file `data_source_tc_teo_environments_test.go` SHALL verify the data source Read logic populates `source_version` correctly. The tests SHALL be buildable and passable.

#### Scenario: Unit test verifies source_version population
- **WHEN** the `DescribeEnvironments` API returns a response with `SourceVersion` set on a config group version info element
- **THEN** the data source Read SHALL populate the nested `source_version` field with that value
