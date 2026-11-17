## Why

The `tencentcloud_teo_deploy_config_group_version` resource exposes the `config_group_version_infos` block to carry config-group version metadata when releasing versions via the EdgeOne `DeployConfigGroupVersion` API, and refreshes these fields from the `DescribeDeployHistory` (`response.Records.ConfigGroupVersionInfos`) response. The cloud API's `ConfigGroupVersionInfo` struct contains a `SourceVersion` field (the version ID from which a version was derived), but the Terraform schema currently omits it. As a result, users cannot read the source version of a deployed version through Terraform, losing information that is present in the cloud API response.

## What Changes

- Add a new `source_version` field (Computed) to the `config_group_version_infos` nested block of the `tencentcloud_teo_deploy_config_group_version` resource.
- In the resource Read function, populate `source_version` from `DescribeDeployHistory` response field `response.Records[].ConfigGroupVersionInfos[].SourceVersion`.

## Capabilities

### New Capabilities
- `teo-deploy-config-group-version-source-version`: Add the `source_version` Computed field to the `config_group_version_infos` nested block of the `tencentcloud_teo_deploy_config_group_version` resource, exposing the EdgeOne cloud API `ConfigGroupVersionInfo.SourceVersion` field.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version.go` — add the `source_version` field to the `config_group_version_infos` nested schema and map it in the Read function.
  - `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version.md` — update documentation.
  - `tencentcloud/services/teo/resource_tc_teo_deploy_config_group_version_test.go` — add unit test coverage for the new field.
- **Cloud APIs:**
  - Read populates the field from `DescribeDeployHistory` (`response.Records[].ConfigGroupVersionInfos[].SourceVersion`, struct `ConfigGroupVersionInfo.SourceVersion`).
- **SDK dependency:** No SDK update required — `ConfigGroupVersionInfo` already exposes `SourceVersion *string` in the vendored `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` package.
- **Backward compatibility:** Fully backward compatible. The new field is `Computed` only and does not change Create behavior; existing Terraform configurations and state continue to work unchanged.
