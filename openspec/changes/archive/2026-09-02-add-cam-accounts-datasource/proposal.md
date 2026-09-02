## Why

Terraform users managing TencentCloud CAM (Cloud Access Management) need a way to query the full list of CAM accounts (including sub-users, collaborators, message receivers, etc.) so that downstream resources can reference account attributes such as Uin, Uid, Name, and UserType. Currently there is no datasource that exposes the `ListAccounts` API, forcing users to hard-code account identifiers.

## What Changes

- Add a new Terraform data source `tencentcloud_cam_accounts` (RESOURCE_KIND_DATASOURCE) that wraps the CAM `ListAccounts` API.
- The data source exposes the optional query parameter `user_type` to filter accounts by type. Pagination is handled automatically in the service layer: `ListAccounts` returns at most 100 records per call, so the service loops using the response's `IsTruncated` field (and the returned `Marker`) until all pages are fetched. No paging-related parameters (`max_items`, `marker`, `is_truncated`) are exposed in the data source schema.
- The data source flattens the collected `Users` array into a top-level `users` list with per-account fields (`uin`, `name`, `uid`, `remark`, `console_login`, `phone_num`, `country_code`, `email`, `create_time`, `user_type`).

## Capabilities

### New Capabilities
- `cam-accounts-datasource`: Query the list of CAM accounts via the `ListAccounts` API and expose the results as a Terraform data source.

### Modified Capabilities
<!-- None -->

## Impact

- **New files:**
  - `tencentcloud/services/cam/data_source_tc_cam_accounts.go`
  - `tencentcloud/services/cam/data_source_tc_cam_accounts.md`
  - `tencentcloud/services/cam/data_source_tc_cam_accounts_test.go`
- **Modified files:**
  - `tencentcloud/services/cam/service_tencentcloud_cam.go` — append a `DescribeCamAccountsByFilter` method wrapping `ListAccounts` with automatic pagination (up to 100 records per request, looping on `IsTruncated`/`Marker` until all pages are fetched).
  - `tencentcloud/provider.go` — register the new `tencentcloud_cam_accounts` data source.
- **API dependency:** `ListAccounts` from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116` (already vendored).
- No breaking changes; this is purely additive.
