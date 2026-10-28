## Why

The `tencentcloud_teo_environments` data source currently exposes environment list information from the `DescribeEnvironments` API, but omits the source version of each config group version (`SourceVersion`) returned by the cloud API. Users need `SourceVersion` to understand which version a deployed config group version was derived from. Adding this field improves observability and aligns the Terraform data source with the full cloud API response.

## What Changes

- Add a new nested computed field `source_version` (string) under `env_infos.current_config_group_version_infos` of the `tencentcloud_teo_environments` data source schema, mapped from `response.EnvInfos[].CurrentConfigGroupVersionInfos[].SourceVersion` of `DescribeEnvironments`.
- Update the data source `Read` function to populate the new field from the API response when it is non-nil.
- Update the `data_source_tc_teo_environments.md` documentation file to reflect the new field.
- Add unit test cases in `data_source_tc_teo_environments_test.go` covering the new field.

## Capabilities

### New Capabilities
- `teo-environments-params`: Adds the `source_version` output field to the `tencentcloud_teo_environments` data source, surfacing the config group version source version from the `DescribeEnvironments` API.

### Modified Capabilities
<!-- None. This is a new capability for the teo environments data source parameters. -->

## Impact

- **Affected code**:
  - `tencentcloud/services/teo/data_source_tc_teo_environments.go` (schema + Read function)
  - `tencentcloud/services/teo/data_source_tc_teo_environments_test.go` (unit tests)
  - `tencentcloud/services/teo/data_source_tc_teo_environments.md` (documentation)
- **APIs**: `DescribeEnvironments` (package `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`) — `SourceVersion` (on `ConfigGroupVersionInfo`, nested under `EnvInfo.CurrentConfigGroupVersionInfos`) is already present in the vendored SDK, confirming the cloud API supports this field.
- **Compatibility**: Backward compatible. The new field is `Computed` (output-only) addition; existing Terraform configurations and state are unaffected.
- **Dependencies**: No new dependencies; uses the existing vendored `tencentcloud-sdk-go` teo module.
