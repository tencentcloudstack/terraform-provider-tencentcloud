# teo-environments-params Specification

## Purpose
TBD - created by archiving change add-teo-environments-params. Update Purpose after archive.
## Requirements
### Requirement: Data source schema fields for total_count and source_version
The data source `tencentcloud_teo_environments` SHALL provide the following additional output fields:
- `total_count` (Computed, int): Total number of environments for the zone, mapped from `response.TotalCount` of the `DescribeEnvironments` API.
- `source_version` (Computed, string) nested under `env_infos.current_config_group_version_infos`: The source version ID that a config group version was derived from, mapped from `response.EnvInfos[].CurrentConfigGroupVersionInfos[].SourceVersion` of the `DescribeEnvironments` API.

#### Scenario: total_count is populated from API response
- **WHEN** the `DescribeEnvironments` API response contains a non-nil `TotalCount`
- **THEN** the data source SHALL set the top-level `total_count` field to the integer value of `TotalCount`

#### Scenario: total_count is nil in API response
- **WHEN** the `DescribeEnvironments` API response contains a nil `TotalCount`
- **THEN** the data source SHALL leave `total_count` unset (default zero) and SHALL NOT return an error

#### Scenario: source_version is populated from API response
- **WHEN** an element of `CurrentConfigGroupVersionInfos` in the `DescribeEnvironments` API response contains a non-nil `SourceVersion`
- **THEN** the data source SHALL set the nested `source_version` field under the corresponding `env_infos.current_config_group_version_infos` element to that string value

#### Scenario: source_version is nil in API response
- **WHEN** an element of `CurrentConfigGroupVersionInfos` in the `DescribeEnvironments` API response contains a nil `SourceVersion`
- **THEN** the data source SHALL leave `source_version` unset for that element and SHALL NOT return an error

### Requirement: Service layer returns TotalCount
The service layer method `TeoService.DescribeTeoEnvironmentsByFilter` SHALL return the `TotalCount` value from the `DescribeEnvironments` API response in addition to the `EnvInfos` list, so that the data source can surface the total environment count without an additional API call.

#### Scenario: Service layer returns TotalCount when present
- **WHEN** the `DescribeEnvironments` API response contains a non-nil `TotalCount`
- **THEN** the service layer method SHALL return that `TotalCount` value as an additional return value

#### Scenario: Service layer handles nil TotalCount
- **WHEN** the `DescribeEnvironments` API response contains a nil `TotalCount`
- **THEN** the service layer method SHALL return nil for the total count value without error

### Requirement: Backward compatibility
The addition of `total_count` and `source_version` SHALL be backward compatible. Both fields are `Computed` output-only fields and SHALL NOT alter existing schema fields, input parameters, or the data source query behavior.

#### Scenario: Existing configurations remain valid
- **WHEN** a user applies an existing `tencentcloud_teo_environments` data source configuration
- **THEN** the configuration SHALL remain valid and the read SHALL succeed, with the new fields simply populated as additional outputs

### Requirement: Documentation update
The documentation file `data_source_tc_teo_environments.md` SHALL be updated to reflect the new `total_count` and `source_version` output fields, following the project documentation guidelines (no manual `Argument Reference` / `Attribute Reference` sections, which are auto-generated).

#### Scenario: Documentation reflects new fields
- **WHEN** the implementation is complete
- **THEN** the documentation file SHALL describe the `total_count` top-level field and the `source_version` nested field

### Requirement: Unit tests with mock
The test file `data_source_tc_teo_environments_test.go` SHALL use the gomonkey mock approach (not the Terraform test suite) to mock the cloud API and verify the data source Read logic populates `total_count` and `source_version` correctly. The tests SHALL be buildable and passable with `go test -gcflags=all=-l`.

#### Scenario: Unit test verifies total_count population
- **WHEN** the mocked `DescribeEnvironments` returns a response with `TotalCount` set to a specific value
- **THEN** the data source Read SHALL populate `total_count` with that value

#### Scenario: Unit test verifies source_version population
- **WHEN** the mocked `DescribeEnvironments` returns a response with `SourceVersion` set on a config group version info element
- **THEN** the data source Read SHALL populate the nested `source_version` field with that value

#### Scenario: Unit test verifies nil handling
- **WHEN** the mocked `DescribeEnvironments` returns a response with nil `TotalCount` and nil `SourceVersion`
- **THEN** the data source Read SHALL NOT error and SHALL leave the corresponding fields unset

