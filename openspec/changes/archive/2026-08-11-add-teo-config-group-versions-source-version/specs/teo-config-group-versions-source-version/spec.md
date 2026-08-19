## ADDED Requirements

### Requirement: Source version output field in config_group_version_infos
The `tencentcloud_teo_config_group_versions` data source SHALL expose a `source_version` computed field within each element of the `config_group_version_infos` nested schema block, mapped from the `SourceVersion` field of the `ConfigGroupVersionInfo` structure returned by the `DescribeConfigGroupVersions` API.

The `source_version` field SHALL be of type string, representing the source version ID that a configuration group version was derived from (format: `ver-xxxxxxxx`, e.g., `ver-2kplomhisdcb`).

#### Scenario: SourceVersion is present in API response
- **WHEN** the `DescribeConfigGroupVersions` API returns a `ConfigGroupVersionInfo` element with a non-nil `SourceVersion` value
- **THEN** the data source Read function SHALL populate the `source_version` field in the corresponding `config_group_version_infos` block element with that value

#### Scenario: SourceVersion is nil in API response
- **WHEN** the `DescribeConfigGroupVersions` API returns a `ConfigGroupVersionInfo` element with a nil `SourceVersion` value
- **THEN** the data source Read function SHALL omit the `source_version` field from the element map (i.e., not call set for that field)
- **AND** no error SHALL be raised

#### Scenario: Field schema definition
- **WHEN** the `tencentcloud_teo_config_group_versions` data source schema is defined
- **THEN** the `config_group_version_infos` nested schema SHALL include a `source_version` field of type `schema.TypeString`
- **AND** the field SHALL be declared as Optional (consistent with sibling output fields in the same nested block)

### Requirement: Backward compatibility of the data source
The addition of the `source_version` field SHALL NOT modify, remove, or alter the behavior of any existing schema field, input parameter, query logic, retry behavior, or state format of the `tencentcloud_teo_config_group_versions` data source.

#### Scenario: Existing configurations remain functional
- **WHEN** a user has an existing Terraform configuration using the `tencentcloud_teo_config_group_versions` data source that references only previously supported fields
- **THEN** the configuration SHALL continue to work without modification
- **AND** `terraform plan` SHALL not show unexpected changes for existing fields
- **AND** the data source output SHALL contain the same values as before for all existing fields
