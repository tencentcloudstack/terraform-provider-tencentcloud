# teo-config-group-version-source-version Specification

## Purpose
Defines that the `tencentcloud_teo_config_group_version_detail` data source SHALL expose a `source_version` computed field within the `config_group_version_info` block, mapping to the cloud API response field `ConfigGroupVersionInfo.SourceVersion` of the `DescribeConfigGroupVersionDetail` API, to surface the source version ID that the config group version was derived from.

## Requirements
### Requirement: Source version field in config group version detail data source
The `tencentcloud_teo_config_group_version_detail` data source SHALL expose a `source_version` computed field within the `config_group_version_info` block, mapping to the cloud API response field `response.ConfigGroupVersionInfo.SourceVersion` of the `DescribeConfigGroupVersionDetail` API. The field represents the source version ID that the config group version was derived from.

#### Scenario: Read config group version detail with source version populated
- **WHEN** the provider reads a `tencentcloud_teo_config_group_version_detail` data source and the `DescribeConfigGroupVersionDetail` API returns a non-nil `ConfigGroupVersionInfo.SourceVersion`
- **THEN** the provider SHALL populate `config_group_version_info.0.source_version` in state with the value of `ConfigGroupVersionInfo.SourceVersion`

#### Scenario: Read config group version detail with source version nil
- **WHEN** the provider reads a `tencentcloud_teo_config_group_version_detail` data source and the `DescribeConfigGroupVersionDetail` API returns a nil `ConfigGroupVersionInfo.SourceVersion` (e.g., the initial version with no source)
- **THEN** the provider SHALL NOT set `source_version` into state (omit the field from the map), consistent with the nil-check pattern used for sibling fields

#### Scenario: Schema definition of source_version
- **WHEN** the `DataSourceTencentCloudTeoConfigGroupVersionDetail()` schema is defined
- **THEN** the `config_group_version_info` block SHALL include a `source_version` field of type `schema.TypeString` with `Computed: true` and `Optional: true`
