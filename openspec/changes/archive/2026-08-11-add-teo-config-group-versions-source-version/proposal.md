## Why

The `tencentcloud_teo_config_group_versions` data source (RESOURCE_KIND_DATASOURCE) exposes version information for an EdgeOne (TEO) configuration group via the `DescribeConfigGroupVersions` API. The API response's `ConfigGroupVersionInfo` structure includes a `SourceVersion` field that indicates the source version ID from which a version was derived, but this field is currently not mapped into the Terraform schema. Users querying config group versions cannot determine the lineage/source of each version, which limits traceability when managing version-controlled configurations.

## What Changes

- Add a new computed output field `source_version` to the `config_group_version_infos` nested schema of the `tencentcloud_teo_config_group_versions` data source.
- The `source_version` field maps from the cloud API response `ConfigGroupVersionInfos.SourceVersion` and represents the source version ID that a version was derived from (format: `ver-xxxxxxxx`, e.g., `ver-2kplomhisdcb`).
- Populate the `source_version` field in the data source Read function by reading `SourceVersion` from the `ConfigGroupVersionInfo` SDK structure when it is non-nil.
- Update the data source documentation to include the new `source_version` field.

## Capabilities

### New Capabilities
- `teo-config-group-versions-source-version`: Expose the `SourceVersion` field from the `DescribeConfigGroupVersions` API response as a computed `source_version` attribute within the `config_group_version_infos` block of the `tencentcloud_teo_config_group_versions` data source.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/teo/data_source_tc_teo_config_group_versions.go` — add `source_version` computed field to the `config_group_version_infos` nested schema and map `SourceVersion` from the API response in the Read function.
  - `tencentcloud/services/teo/data_source_tc_teo_config_group_versions_test.go` — add unit test coverage asserting `source_version` is populated from the mocked API response.
  - `tencentcloud/services/teo/data_source_tc_teo_config_group_versions.md` — document the new `source_version` output field (regenerated via `make doc` during finalization).
- **SDK dependency:** No SDK update needed. The `ConfigGroupVersionInfo` structure in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` already includes the `SourceVersion *string` field.
- **Backward compatibility:** Fully backward compatible. The change only adds a new computed output field; no existing schema fields are modified or removed, and no breaking changes are introduced to existing Terraform configurations or state.
