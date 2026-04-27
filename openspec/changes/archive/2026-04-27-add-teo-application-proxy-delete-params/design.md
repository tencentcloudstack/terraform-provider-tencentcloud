## Context

The `tencentcloud_teo_application_proxy` resource manages TEO (TencentCloud EdgeOne) application proxies. It uses a composite ID (`zone_id#proxy_id`) to uniquely identify each proxy instance. The current delete function `resourceTencentCloudTeoApplicationProxyDelete` parses `zone_id` and `proxy_id` from `d.Id()` by splitting on the `FILED_SP` separator, then passes them to `service.DeleteTeoApplicationProxyById()`.

Per the project coding guidelines, when a resource uses a composite ID with `FILED_SP` as separator, the read, update, and delete methods MUST obtain the individual ID component fields from `d.Get()` rather than parsing `d.Id()`. The `DeleteApplicationProxy` cloud API already supports `ZoneId` and `ProxyId` as request parameters in the SDK (`DeleteApplicationProxyRequest`), so this change only requires updating the Go resource code.

Current delete function flow:
1. Parse `d.Id()` → split by `FILED_SP` → get `zoneId`, `proxyId`
2. Call `service.DeleteTeoApplicationProxyById(ctx, zoneId, proxyId)`

Target delete function flow:
1. Get `zone_id` from `d.Get("zone_id")` and `proxy_id` from `d.Get("proxy_id")`
2. Construct `DeleteApplicationProxyRequest` directly with `ZoneId` and `ProxyId`
3. Call the API with retry logic using `tccommon.WriteRetryTimeout`

## Goals / Non-Goals

**Goals:**
- Modify the delete function to read `zone_id` and `proxy_id` from `d.Get()` instead of parsing from `d.Id()`
- Construct and call the `DeleteApplicationProxy` API request directly in the delete function (instead of delegating to the service helper method), following the standard CRUD pattern
- Add retry logic (`tccommon.WriteRetryTimeout`) for the delete API call
- Update unit tests to verify the new delete function behavior

**Non-Goals:**
- Do not modify the resource schema (both `zone_id` and `proxy_id` already exist)
- Do not modify the read, create, or update functions
- Do not modify the `DeleteTeoApplicationProxyById` service method (it can remain for backward compatibility)

## Decisions

1. **Use `d.Get()` over `d.Id()` parsing**: Per coding guidelines, obtain `zone_id` and `proxy_id` from `d.Get()` rather than splitting `d.Id()`. This is the standard pattern for composite ID resources.

2. **Call API directly in delete function**: Instead of delegating to `service.DeleteTeoApplicationProxyById()`, construct the `DeleteApplicationProxyRequest` and call the API directly within the delete function with proper retry logic. This aligns with the pattern where each CRUD function directly manages its API interaction.

3. **Keep the offline-then-delete flow**: The current delete function first sets the proxy status to `offline` via `ModifyApplicationProxyStatus`, then calls `DeleteApplicationProxy`. This two-step flow must be preserved.

## Risks / Trade-offs

- [Risk] `d.Get("proxy_id")` returns empty string if the state is corrupted → Mitigation: The `proxy_id` field is always set during the read operation and is part of the composite ID, so it should always be available in state.
- [Risk] Changing the delete function could break existing behavior → Mitigation: The delete logic remains the same; only the source of `zone_id` and `proxy_id` values changes (from `d.Id()` parsing to `d.Get()`).
