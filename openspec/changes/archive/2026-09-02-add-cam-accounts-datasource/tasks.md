## 1. Service Layer

- [x] 1.1 Append `DescribeCamAccountsByFilter(ctx, paramMap)` to `tencentcloud/services/cam/service_tencentcloud_cam.go` — build `cam.NewListAccountsRequest()` from `paramMap["MaxItems"]`, `paramMap["Marker"]`, `paramMap["UserType"]`, call `me.client.UseCamClient().ListAccounts(request)`, return `([]*cam.ListAllUser, marker string, isTruncated bool, error)`

## 2. Data Source Implementation

- [x] 2.1 Create `tencentcloud/services/cam/data_source_tc_cam_accounts.go` with `DataSourceTencentCloudCamAccounts()` schema: optional inputs `max_items` (TypeInt), `marker` (TypeString, Optional+Computed), `user_type` (TypeString, no ValidateFunc); flattened `users` TypeList whose Elem schema.Resource contains `uin`, `name`, `uid`, `remark`, `console_login`, `phone_num`, `country_code`, `email`, `create_time`, `user_type`; computed outputs `marker` and `is_truncated`; optional `result_output_file`
- [x] 2.2 Implement `dataSourceTencentCloudCamAccountsRead` — build paramMap from optional inputs, wrap `DescribeCamAccountsByFilter` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`; inside retry return `NonRetryableError` on empty response (do NOT `d.SetId("")`); flatten `Users` into `users` list setting each field only when non-nil; set `marker` and `is_truncated`; on retry-exhausted path log `[DATASOURCE] read empty, skip SetId`; SetId via `helper.DataResourceIdsHash(ids)`
- [x] 2.3 Verify all returned errors are checked; functions that cannot fail use `_ = func()` for the err

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_cam_accounts` → `cam.DataSourceTencentCloudCamAccounts()` in the data source map of `tencentcloud/provider.go`

## 4. Documentation & Tests

- [x] 4.1 Create `tencentcloud/services/cam/data_source_tc_cam_accounts.md` with one-line description ("Use this data source to query ..."), Example Usage (including paging example using `max_items`/`marker`/`user_type`), no Argument/Attribute Reference sections (auto-generated), no Import section
- [x] 4.2 Create `tencentcloud/services/cam/data_source_tc_cam_accounts_test.go` using gomonkey mocks (no Terraform test suite) — mock `DescribeCamAccountsByFilter` to return sample `ListAllUser` data and verify the read handler flattens fields correctly
