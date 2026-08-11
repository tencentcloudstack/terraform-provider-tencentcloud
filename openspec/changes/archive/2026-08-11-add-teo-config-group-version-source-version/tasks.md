## 1. Schema Changes

- [x] 1.1 Add `source_version` field (TypeString, Computed: true, Optional: true, with description) to the `config_group_version_info` block schema in `DataSourceTencentCloudTeoConfigGroupVersionDetail()` in `tencentcloud/services/teo/data_source_tc_teo_config_group_version_detail.go`

## 2. Read Function Changes

- [x] 2.1 In `dataSourceTencentCloudTeoConfigGroupVersionDetailRead`, add a nil-check for `respData.ConfigGroupVersionInfo.SourceVersion` and map it to `configGroupVersionInfoMap["source_version"]` within the existing `if respData.ConfigGroupVersionInfo != nil` block, placed alongside the sibling field mappings

## 3. Documentation

- [x] 3.1 Update `tencentcloud/services/teo/data_source_tc_teo_config_group_version_detail.md` to describe the new `source_version` field within the `config_group_version_info` block

## 4. Validation

- [x] 4.1 Verify the code compiles successfully (gofmt in finalize phase)
- [x] 4.2 Verify no lint errors
