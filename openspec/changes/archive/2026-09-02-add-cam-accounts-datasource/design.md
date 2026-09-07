## Context

The Terraform TencentCloud provider already exposes several CAM data sources (`tencentcloud_cam_users`, `tencentcloud_cam_sub_accounts`, `tencentcloud_cam_roles`, etc.), but none of them wraps the `ListAccounts` API. `ListAccounts` differs from `ListUsers` in that it returns **all** account types (Owner, SubUser, CICUser, WechatCorpUser, AgentIdentity, Collaborator, MessageReceiver) in a single call and supports server-side paging via `MaxItems` / `Marker` plus filtering by `UserType`. This makes it the appropriate API for a general-purpose "list all CAM accounts" data source.

The existing `tencentcloud_cam_users` data source uses the older `ListUsers` API and performs client-side filtering, whereas the new `tencentcloud_cam_accounts` uses `ListAccounts` server-side filtering (`UserType`) and relies on automatic pagination implemented in the service layer to collect every page of results.

## Goals / Non-Goals

**Goals:**
- Provide a read-only data source `tencentcloud_cam_accounts` that calls the `ListAccounts` API.
- Expose the optional `user_type` request parameter so users can filter by account type.
- Implement automatic full pagination in the service layer: `ListAccounts` returns at most 100 records per call (MaxItems range [1, 100]), so the service requests pages of 100 and follows the response `Marker` while `IsTruncated` is true until every page is fetched.
- Flatten the collected `Users` response array into a top-level `users` list, with each element's fields spread directly (no extra nesting layer), per provider conventions for RESOURCE_KIND_DATASOURCE.
- Register the data source in `provider.go` and produce docs/tests.

**Non-Goals:**
- Not exposing any paging-related parameter in the data source schema (`max_items`, `marker`, `is_truncated` are all omitted) — pagination is an internal service-layer concern and the data source always returns the complete account list.
- Not replacing the existing `tencentcloud_cam_users` or `tencentcloud_cam_sub_accounts` data sources.

## Decisions

### Decision 1: Use `ListAccounts` (not `ListUsers`)
`ListAccounts` returns all account types and supports `MaxItems`, `Marker`, and `UserType` natively, matching the required parameter mapping. `ListUsers` only returns sub-users and lacks server-side paging. Therefore `ListAccounts` is the correct API.

### Decision 2: Auto-paginate in the service layer instead of exposing paging parameters
`ListAccounts` returns at most 100 records per request (`MaxItems` valid range [1, 100]) and reports truncation with `IsTruncated` plus a `Marker` to continue from. Instead of exposing `max_items` / `marker` inputs and `marker` / `is_truncated` outputs in the schema, the service method `DescribeCamAccountsByFilter` loops internally: it sets `MaxItems = 100` per request, appends `response.Response.Users` to the result, and keeps calling the API with `request.Marker = response.Response.Marker` while `IsTruncated` is true (and `Marker` is non-nil), so the data source always receives the complete account list.

### Decision 3: Flatten `users` list — no extra nesting layer
Per provider rule: "资源参数 schema 中禁止创建'该资源列表型数据'这一层 schema". The `users` field is a `TypeList` whose `Elem` is a `schema.Resource` containing the per-account fields directly. Each field (`uin`, `name`, `uid`, `remark`, `console_login`, `phone_num`, `country_code`, `email`, `create_time`, `user_type`) sits inside that resource elem — this is the standard flattened list pattern used by `tencentcloud_igtm_instance_list` and `tencentcloud_cam_sub_accounts`.

### Decision 4: Field types match the SDK
The `ListAccounts` response (`ListAllUser` struct) types:
- `Uin` *int64 → schema `uin` TypeInt
- `Name` *string → schema `name` TypeString
- `Uid` *int64 → schema `uid` TypeInt
- `Remark` *string → schema `remark` TypeString
- `ConsoleLogin` *int64 → schema `console_login` TypeInt (keep as int, matching the raw API; do not convert to bool to avoid ambiguity)
- `PhoneNum` *string → schema `phone_num` TypeString
- `CountryCode` *string → schema `country_code` TypeString
- `Email` *string → schema `email` TypeString
- `CreateTime` *string → schema `create_time` TypeString
- `UserType` *string → schema `user_type` TypeString (the output `UserType` is a string enum in `ListAllUser`, unlike `SubAccountUser` where it is int)

Request parameter types:
- `user_type` TypeString (maps to `UserType *string`; no validate func — per user requirement, the enum range is not enforced in schema)
- No `max_items` / `marker` / `is_truncated` schema fields — pagination is handled entirely inside the service layer.

### Decision 5: Service-layer method
Add `DescribeCamAccountsByFilter(ctx, paramMap)` to `service_tencentcloud_cam.go`. It builds a `ListAccountsRequest` from `paramMap["UserType"]`, sets `MaxItems = 100` per page, and loops: each iteration calls `me.client.UseCamClient().ListAccounts(request)`, appends `response.Response.Users`, and continues with `request.Marker = response.Response.Marker` while `response.Response.IsTruncated` is true and `Marker` is non-nil. It returns `([]*cam.ListAllUser, error)`. The data source read handler wraps the call in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.

### Decision 6: DATASOURCE empty-response handling
Consistent with the existing CAM data sources (`tencentcloud_cam_users`, `tencentcloud_cam_sub_accounts`), the read handler does not special-case an empty `Users` list: the retry block assigns the result and the handler proceeds, producing an empty `users` list and a `SetId` computed from the (empty) ids. On the retry-exhausted path the handler logs `[DATASOURCE] read empty, skip SetId` and returns the error.

### Decision 7: SetId
Use `helper.DataResourceIdsHash(ids)` where `ids` is the list of `uin` strings, consistent with `tencentcloud_cam_users`.

## Risks / Trade-offs

- **[Risk] Auto-pagination issues multiple API calls** → Accounts with more than 100 entries trigger several sequential `ListAccounts` requests (100 per page) on every refresh. Accepted because the data source must return the complete account list; `ratelimit.Check` throttles each call.
- **[Risk] No paging control for users** → Users cannot limit the result size or resume from a marker via the data source. Accepted per requirement: the schema must not expose paging parameters; the service fetches everything using `IsTruncated`.
- **[Risk] `ConsoleLogin` int vs bool** → Kept as int to match the raw `ListAllUser` struct (`*int64`) and avoid a lossy conversion. Documented in the field description.
