## Why

The `tencentcloud_teo_config_group_version_detail` data source currently reads the `ConfigGroupVersionInfo` from the `DescribeConfigGroupVersionDetail` API but does not expose the `SourceVersion` field. The SDK struct `ConfigGroupVersionInfo` already includes `SourceVersion` (the source version ID that the current version was derived from), but the Terraform schema omits it, so users cannot retrieve this information through the data source.

## What Changes

- Add a new computed field `source_version` to the `config_group_version_info` block of the `tencentcloud_teo_config_group_version_detail` data source schema, mapping to the cloud API response field `response.ConfigGroupVersionInfo.SourceVersion`.
- Populate `source_version` in the Read function from the API response when the field is non-nil.
- Update the data source documentation to describe the new field.

## Capabilities

### New Capabilities
- `teo-config-group-version-source-version`: Expose the `source_version` field (the source version ID that the config group version was derived from) in the `tencentcloud_teo_config_group_version_detail` data source.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/teo/data_source_tc_teo_config_group_version_detail.go` — add `source_version` schema field to the `config_group_version_info` block and map it from the API response in the Read function
  - `tencentcloud/services/teo/data_source_tc_teo_config_group_version_detail.md` — update documentation with the new field
  - `website/docs/d/teo_config_group_version_detail.html.markdown` — generated via `make doc` (not manually edited)
- **SDK dependency:** No SDK update needed — `ConfigGroupVersionInfo.SourceVersion` already exists in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901/models.go`
- **Backward compatibility:** fully backward compatible — the new field is Computed only and additive
- **API constraints:** `SourceVersion` is a read-only output field of `DescribeConfigGroupVersionDetail`; no API request changes required
