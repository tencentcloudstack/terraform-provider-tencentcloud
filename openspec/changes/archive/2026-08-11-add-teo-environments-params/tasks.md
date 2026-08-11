## 1. Service Layer Update

- [x] 1.1 Modify `TeoService.DescribeTeoEnvironmentsByFilter` in `tencentcloud/services/teo/service_tencentcloud_teo.go` to also return the `TotalCount` value (`*uint64`) from `response.Response.TotalCount` alongside the existing `[]*teov20220901.EnvInfo` return value. Update the return signature to `(ret []*teov20220901.EnvInfo, totalCount *uint64, errRet error)`.
- [x] 1.2 Guard the `TotalCount` extraction with a nil check on `response.Response` and `response.Response.TotalCount`; return `nil` for the count when absent.

## 2. Data Source Schema Update

- [x] 2.1 Add a new top-level computed field `total_count` (`schema.TypeInt`) to the `tencentcloud_teo_environments` data source schema in `tencentcloud/services/teo/data_source_tc_teo_environments.go`.
- [x] 2.2 Add a new nested computed field `source_version` (`schema.TypeString`) under the `current_config_group_version_infos` block schema in `data_source_tc_teo_environments.go`.

## 3. Data Source Read Function Update

- [x] 3.1 Update `dataSourceTencentCloudTeoEnvironmentsRead` in `data_source_tc_teo_environments.go` to consume the new `totalCount` return value from `DescribeTeoEnvironmentsByFilter` and set the `total_count` schema field (with a nil check) using `d.Set("total_count", ...)`.
- [x] 3.2 In the existing loop that builds `currentConfigGroupVersionInfosMap`, add a nil-guarded assignment of `source_version` from `currentConfigGroupVersionInfos.SourceVersion` when non-nil.
- [x] 3.3 Ensure the empty-response path is preserved: when `respData` is empty, log `[DATASOURCE] read empty, skip SetId` context and do not clear the id, returning the appropriate error per the data source nil-handling pattern.

## 4. Unit Tests

- [x] 4.1 Add/update unit test cases in `tencentcloud/services/teo/data_source_tc_teo_environments_test.go` using the gomonkey mock approach (mock `TeoService.DescribeTeoEnvironmentsByFilter`) to verify `total_count` is populated from a mocked response.
- [x] 4.2 Add unit test cases verifying `source_version` is populated within the nested `current_config_group_version_infos` block from a mocked response.
- [x] 4.3 Add unit test cases verifying nil `TotalCount` and nil `SourceVersion` are handled without error (fields left unset).

## 5. Documentation

- [x] 5.1 Update the data source example/documentation file `tencentcloud/services/teo/data_source_tc_teo_environments.md` to reflect the new `total_count` and `source_version` output fields, following the project documentation guidelines (one-line description mentioning TEO, no manual `Argument Reference`/`Attribute Reference` sections).
- [x] 5.2 Run `make doc` (during the finalize phase via the tfpacer-finalize skill) to regenerate `website/docs/d/teo_environments.html.markdown` — do not manually edit files under `website/`.

## 6. Verification (separate from code changes)

- [x] 6.1 Verify the generated code compiles correctly (build/lint verification performed by the downstream validation flow; do not run `go build`/`go vet`/`golint` manually per project rules).
- [x] 6.2 Verify unit tests are buildable and consistent with the gomonkey mock pattern (actual test execution is performed by the downstream validation flow).
