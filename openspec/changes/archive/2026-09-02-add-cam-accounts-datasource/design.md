## Context

The Terraform TencentCloud provider already exposes several CAM data sources (`tencentcloud_cam_users`, `tencentcloud_cam_sub_accounts`, `tencentcloud_cam_roles`, etc.), but none of them wraps the `ListAccounts` API. `ListAccounts` differs from `ListUsers` in that it returns **all** account types (Owner, SubUser, CICUser, WechatCorpUser, AgentIdentity, Collaborator, MessageReceiver) in a single call and supports server-side paging via `MaxItems` / `Marker` plus filtering by `UserType`. This makes it the appropriate API for a general-purpose "list all CAM accounts" data source.

The existing `tencentcloud_cam_users` data source uses the older `ListUsers` API and performs client-side filtering, whereas the new `tencentcloud_cam_accounts` will rely on server-side paging/filtering provided by `ListAccounts`.

## Goals / Non-Goals

**Goals:**
- Provide a read-only data source `tencentcloud_cam_accounts` that calls the `ListAccounts` API.
- Expose the optional request parameters `max_items`, `marker`, and `user_type` so users can control paging and filter by account type.
- Flatten the `Users` response array into a top-level `users` list, with each element's fields spread directly (no extra nesting layer), per provider conventions for RESOURCE_KIND_DATASOURCE.
- Surface the paging outputs `marker` and `is_truncated` so callers can detect truncation and chain subsequent calls.
- Register the data source in `provider.go` and produce docs/tests.

**Non-Goals:**
- Not implementing automatic full-pagination loop that fetches every page internally. The data source passes `MaxItems` / `Marker` through to the API and returns the single page plus the `marker` / `is_truncated` fields so the caller decides whether to continue. This mirrors the API's design.
- Not replacing the existing `tencentcloud_cam_users` or `tencentcloud_cam_sub_accounts` data sources.

## Decisions

### Decision 1: Use `ListAccounts` (not `ListUsers`)
`ListAccounts` returns all account types and supports `MaxItems`, `Marker`, and `UserType` natively, matching the required parameter mapping. `ListUsers` only returns sub-users and lacks server-side paging. Therefore `ListAccounts` is the correct API.

### Decision 2: Pass paging parameters through instead of auto-paging
The `ListAccounts` API is explicitly designed around `MaxItems` / `Marker` paging with an `IsTruncated` flag. The requirement maps `request.MaxItems → max_items` and `request.Marker → marker` as optional inputs, and `response.Response.Marker → marker` and `response.Response.IsTruncated → is_truncated` as outputs. This indicates the data source should expose paging control to the user rather than silently looping. We follow the requirement mapping directly.

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
- `max_items` TypeInt (maps to `MaxItems *int64`)
- `marker` TypeString (maps to `Marker *string`)
- `user_type` TypeString (maps to `UserType *string`; no validate func — per user requirement, the enum range is not enforced in schema)

### Decision 5: Service-layer method
Add `DescribeCamAccountsByFilter(ctx, paramMap)` to `service_tencentcloud_cam.go`. It builds a `ListAccountsRequest` from `paramMap`, calls `me.client.UseCamClient().ListAccounts(request)`, and returns `([]*cam.ListAllUser, marker string, isTruncated bool, error)`. The data source read handler wraps the call in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.

### Decision 6: DATASOURCE empty-response handling
Per provider rule for RESOURCE_KIND_DATASOURCE: inside the retry block, if the response is nil / `response.Response` is nil / `len(Users) == 0`, return `NonRetryableError` (do **not** `d.SetId("")`). On the outer retry-exhausted path, log `[DATASOURCE] read empty, skip SetId`.

### Decision 7: SetId
Use `helper.DataResourceIdsHash(ids)` where `ids` is the list of `uin` strings, consistent with `tencentcloud_cam_users`.

## Risks / Trade-offs

- **[Risk] `marker` is both an input and an output field** → The same schema key `marker` is Optional (input) and also set from the response. Terraform allows a field to be both Optional and Computed, which is the correct approach here so the user can supply an initial marker and read back the next marker.
- **[Risk] Paging is not automatic** → Users who have more accounts than `MaxItems` must issue multiple data source blocks with chained markers. This is a conscious trade-off to honor the API's native paging design and the required parameter mapping. Mitigated by exposing `is_truncated` so users know when to stop.
- **[Risk] `ConsoleLogin` int vs bool** → Kept as int to match the raw `ListAllUser` struct (`*int64`) and avoid a lossy conversion. Documented in the field description.
