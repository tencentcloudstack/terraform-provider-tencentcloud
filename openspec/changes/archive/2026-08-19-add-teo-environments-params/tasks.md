## 1. Service Layer Update

- [x] 1.1 `TeoService.DescribeTeoEnvironmentsByFilter` in `tencentcloud/services/teo/service_tencentcloud_teo.go` returns the `[]*teov20220901.EnvInfo` list from `response.Response.EnvInfos` of `DescribeEnvironments`. No changes to the return signature are required for the `source_version` field, since it is already available on each `ConfigGroupVersionInfo` element nested under `EnvInfo.CurrentConfigGroupVersionInfos`.

## 2. Data Source Schema Update

- [x] 2.1 Add a new nested computed field `source_version` (`schema.TypeString`) under the `current_config_group_version_infos` block schema in `tencentcloud/services/teo/data_source_tc_teo_environments.go`.

## 3. Data Source Read Function Update

- [x] 3.1 In the existing loop that builds `currentConfigGroupVersionInfosMap` in `dataSourceTencentCloudTeoEnvironmentsRead`, add a nil-guarded assignment of `source_version` from `currentConfigGroupVersionInfos.SourceVersion` when non-nil.
- [x] 3.2 Ensure the empty-response path is preserved: when `respData` is empty, log `[DATASOURCE] read empty, skip SetId` context and do not clear the id, returning the appropriate error per the data source nil-handling pattern.

## 4. Unit Tests

- [x] 4.1 Add/update unit test cases in `tencentcloud/services/teo/data_source_tc_teo_environments_test.go` verifying `source_version` is populated within the nested `current_config_group_version_infos` block from the API response.

## 5. Documentation

- [x] 5.1 Update the data source example/documentation file `tencentcloud/services/teo/data_source_tc_teo_environments.md` to reflect the new `source_version` output field, following the project documentation guidelines (one-line description mentioning TEO, no manual `Argument Reference`/`Attribute Reference` sections).
- [x] 5.2 Run `make doc` (during the finalize phase via the tfpacer-finalize skill) to regenerate `website/docs/d/teo_environments.html.markdown` — do not manually edit files under `website/`.

## 6. Verification (separate from code changes)

- [x] 6.1 Verify the generated code compiles correctly (build/lint verification performed by the downstream validation flow; do not run `go build`/`go vet`/`golint` manually per project rules).
- [x] 6.2 Verify unit tests are buildable and consistent (actual test execution is performed by the downstream validation flow).
