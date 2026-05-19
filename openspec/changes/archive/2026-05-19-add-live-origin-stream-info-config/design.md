## Context

The `live` service in the provider already has `service_tencentcloud_live.go` with existing query helpers. The SDK `live/v20180801` includes both `ModifyOriginStreamInfoRequest` and `DescribeOriginStreamInfoRequest`/`ResponseParams`. The `OriginStreamCustomizationRule` nested struct is also present. No SDK changes are needed.

The resource is a **config-type**: there is no dedicated create/delete API — the first `ModifyOriginStreamInfo` call configures the domain, and "delete" is a no-op (configuration cannot be truly deleted, only changed). ID is the `DomainName`.

## Goals / Non-Goals

**Goals:**
- Implement `tencentcloud_live_origin_stream_info_config` with Create/Read/Update/Delete lifecycle following the waf_owasp_rule_status_config style.
- Schema fields match `ModifyOriginStreamInfo` request parameters exactly.
- Async polling after Modify: retry until `DescribeOriginStreamInfo.Status` is `1` or `3`.
- Service layer method `DescribeLiveOriginStreamInfo` for Read.
- Acceptance test and .md doc.

**Non-Goals:**
- Supporting `StreamPackageRegion` (read-only context field not in ModifyOriginStreamInfo request).
- Pagination — `DescribeOriginStreamInfo` is a single-record query by domain.

## Decisions

### D1: Config resource pattern (no real delete)
`Create` sets ID and delegates to `Update`. `Delete` is a no-op returning nil. This matches other config-type resources in the provider.

### D2: Nested `customization_rules` as TypeList of Resource
`CustomizationRules` is `[]*OriginStreamCustomizationRule` in SDK. Expose as `TypeList` with sub-schema matching the struct fields. Each element has its own `origin_address` as `TypeList of TypeString`.

### D3: Async wait in WriteRetryTimeout loop
After `ModifyOriginStreamInfo` succeeds, poll `DescribeOriginStreamInfo` with `resource.Retry(tccommon.WriteRetryTimeout)`. Terminal states: `Status == 1` (success) or `Status == 3` (closed). Non-terminal: `Status == 0` or `2`. Any other value is a non-retryable error.

### D4: `status` as Computed field
`Status` from `DescribeOriginStreamInfo` is exposed as a Computed `TypeInt` field so users can observe the current configuration state.

## Risks / Trade-offs

- [Risk] `DescribeOriginStreamInfo` may return `ResourceNotFound` if the domain has never been configured → Mitigation: treat as `d.SetId("")` in Read.
- [Risk] Async timeout — WriteRetryTimeout (3 min) may not be enough for slow configurations → Mitigation: document in resource description; users can re-apply.
