## Why

The `tencentcloud_teo_environments` data source currently exposes environment list information from the `DescribeEnvironments` API, but omits two useful fields returned by the cloud API: the total environment count (`TotalCount`) and the source version of each config group version (`SourceVersion`). Users need `TotalCount` to know how many environments exist for a zone, and `SourceVersion` to understand which version a deployed config group version was derived from. Adding these fields improves observability and aligns the Terraform data source with the full cloud API response.

## What Changes

- Add a new top-level computed field `total_count` (int) to the `tencentcloud_teo_environments` data source schema, mapped from `response.TotalCount` of `DescribeEnvironments`.
- Add a new nested computed field `source_version` (string) under `env_infos.current_config_group_version_infos` of the `tencentcloud_teo_environments` data source schema, mapped from `response.EnvInfos[].CurrentConfigGroupVersionInfos[].SourceVersion` of `DescribeEnvironments`.
- Update the data source `Read` function to populate the new fields from the API response when they are non-nil.
- Update the `data_source_tc_teo_environments.md` documentation file to reflect the new fields.
- Add unit test cases in `data_source_tc_teo_environments_test.go` covering the new fields using the gomonkey mock approach.

## Capabilities

### New Capabilities
- `teo-environments-params`: Adds `total_count` and `source_version` output fields to the `tencentcloud_teo_environments` data source, surfacing the total environment count and config group version source version from the `DescribeEnvironments` API.

### Modified Capabilities
<!-- None. This is a new capability for the teo environments data source parameters. -->

## Impact

- **Affected code**:
  - `tencentcloud/services/teo/data_source_tc_teo_environments.go` (schema + Read function)
  - `tencentcloud/services/teo/data_source_tc_teo_environments_test.go` (unit tests)
  - `tencentcloud/services/teo/data_source_tc_teo_environments.md` (documentation)
- **APIs**: `DescribeEnvironments` (package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`) — both `TotalCount` (on `DescribeEnvironmentsResponseParams`) and `SourceVersion` (on `ConfigGroupVersionInfo`, nested under `EnvInfo.CurrentConfigGroupVersionInfos`) are already present in the vendored SDK, confirming the cloud API supports these fields.
- **Compatibility**: Backward compatible. Both new fields are `Computed` (output-only) additions; existing Terraform configurations and state are unaffected.
- **Dependencies**: No new dependencies; uses the existing vendored `tencentcloud-sdk-go` teo module.
