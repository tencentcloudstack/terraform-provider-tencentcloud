## Context

The `tencentcloud_teo_security_policy_config` resource manages TEO (TencentCloud EdgeOne) security policy configurations. It needs to support the full parameter set of the `ModifySecurityPolicy` and `DescribeSecurityPolicy` cloud APIs from the `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` package.

The resource is a RESOURCE_KIND_GENERAL type that uses:
- `ModifySecurityPolicy` for Create and Update operations
- `DescribeSecurityPolicy` for Read operations
- No dedicated Delete API (delete is a no-op or resets to defaults)

The resource supports three entity types for targeting policies:
- `ZoneDefaultPolicy`: Site-level policy (ID format: `{zone_id}#ZoneDefaultPolicy`)
- `Host`: Domain-level policy (ID format: `{zone_id}#Host#{host}`)
- `Template`: Template-level policy (ID format: `{zone_id}#Template#{template_id}`)

## Goals / Non-Goals

**Goals:**
- Support all parameters of `ModifySecurityPolicy` API: `zone_id`, `entity`, `host`, `template_id`, `security_config`, `security_policy`
- Support reading back `security_policy` from `DescribeSecurityPolicy` response
- Use composite ID with `tccommon.FILED_SP` separator based on entity type
- Implement retry logic with `tccommon.ReadRetryTimeout` for API calls
- Maintain backward compatibility with existing configurations

**Non-Goals:**
- Adding new sub-fields within `SecurityPolicy` or `SecurityConfig` structs beyond what the SDK provides
- Implementing a dedicated Delete API (the resource uses a no-op delete)
- Modifying the existing resource ID format

## Decisions

1. **Resource ID format**: Use composite ID `{zone_id}#{entity}[#{host_or_template_id}]` with `tccommon.FILED_SP` as separator. This allows import support and uniquely identifies the policy target.

2. **ForceNew fields**: `zone_id`, `entity`, `host`, and `template_id` are ForceNew since changing them targets a different policy entirely.

3. **security_config vs security_policy**: Both are Optional. `security_config` is the legacy configuration format, while `security_policy` is the newer expression-based format. The API handles precedence between them.

4. **Read behavior**: `DescribeSecurityPolicy` returns `SecurityPolicy` in the response. The `security_config` field is write-only since the Describe API does not return it.

## Risks / Trade-offs

- [Risk] `security_config` is write-only (not returned by DescribeSecurityPolicy) → Mark as Computed to avoid perpetual diffs; document this behavior.
- [Risk] Complex nested structures in `SecurityPolicy` → Use careful nil checks in Read to avoid panics on optional fields.
