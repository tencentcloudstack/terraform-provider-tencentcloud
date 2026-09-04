## 1. Service Layer

- [x] 1.1 Append `DescribeCamAccountsByFilter(ctx, paramMap)` to `tencentcloud/services/cam/service_tencentcloud_cam.go` — build `cam.NewListAccountsRequest()` from `paramMap["UserType"]`, set `MaxItems = 100` per page (the API returns at most 100 records per call), and loop: call `me.client.UseCamClient().ListAccounts(request)`, append `response.Response.Users`, and continue with `request.Marker = response.Response.Marker` while the response's `IsTruncated` is true and `Marker` is non-nil; return `([]*cam.ListAllUser, error)`

## 2. Data Source Implementation

- [x] 2.1 Create `tencentcloud/services/cam/data_source_tc_cam_accounts.go` with `DataSourceTencentCloudCamAccounts()` schema: optional input `user_type` (TypeString, no ValidateFunc); flattened `users` TypeList whose Elem schema.Resource contains `uin`, `name`, `uid`, `remark`, `console_login`, `phone_num`, `country_code`, `email`, `create_time`, `user_type`; optional `result_output_file`. No paging-related parameters (`max_items`, `marker`, `is_truncated`) are exposed — pagination is automatic in the service layer
- [x] 2.2 Implement `dataSourceTencentCloudCamAccountsRead` — build paramMap from `user_type`, wrap `DescribeCamAccountsByFilter` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`; log `[DATASOURCE] read empty, skip SetId` on the retry-exhausted path; an empty `Users` list is not an error (consistent with `tencentcloud_cam_users`); flatten `Users` into `users` list setting each field only when non-nil; SetId via `helper.DataResourceIdsHash(ids)`
- [x] 2.3 Verify all returned errors are checked; functions that cannot fail use `_ = func()` for the err

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_cam_accounts` → `cam.DataSourceTencentCloudCamAccounts()` in the data source map of `tencentcloud/provider.go`

## 4. Documentation & Tests

- [x] 4.1 Create `tencentcloud/services/cam/data_source_tc_cam_accounts.md` with one-line description ("Use this data source to query ..."), Example Usage (querying all accounts and filtering by `user_type`), no Argument/Attribute Reference sections (auto-generated), no Import section
- [x] 4.2 Create `tencentcloud/services/cam/data_source_tc_cam_accounts_test.go` using gomonkey mocks (no Terraform test suite) — mock `DescribeCamAccountsByFilter` (returning `([]*cam.ListAllUser, error)`) with sample `ListAllUser` data and verify the read handler flattens fields correctly, handles nil fields and empty responses, and the schema exposes no paging parameters
