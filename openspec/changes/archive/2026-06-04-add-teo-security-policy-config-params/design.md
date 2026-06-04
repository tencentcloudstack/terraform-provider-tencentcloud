## Context

The `tencentcloud_teo_security_policy_config` resource manages TEO (TencentCloud EdgeOne) security policy configurations. It uses:
- `ModifySecurityPolicy` API for create/update operations
- `DescribeSecurityPolicy` API for read operations

The resource already exists at `tencentcloud/services/teo/resource_tc_teo_security_policy_config.go` with the following top-level schema fields:
- `zone_id` (Required, ForceNew) - Zone ID
- `entity` (Optional, ForceNew) - Security policy type (ZoneDefaultPolicy/Template/Host)
- `host` (Optional, ForceNew) - Domain name for Host entity type
- `template_id` (Optional, ForceNew) - Template ID for Template entity type
- `security_policy` (Optional, TypeList, MaxItems: 1) - Security policy configuration using expression grammar
- `security_config` (Optional, Computed, TypeList, MaxItems: 1) - Classic web protection settings

The `SecurityPolicy` parameter in the `ModifySecurityPolicy` API maps to the `security_policy` schema field and contains sub-structures for:
- `CustomRules` - Custom rule configuration
- `ManagedRules` - Managed rules configuration
- `HttpDDoSProtection` - HTTP DDoS protection configuration
- `RateLimitingRules` - Rate limiting rules configuration
- `ExceptionRules` - Exception rules configuration
- `BotManagement` - Bot management configuration
- `BotManagementLite` - Basic bot management configuration
- `DefaultDenySecurityActionParameters` - Default deny action parameters

The `DescribeSecurityPolicy` API returns `response.Response.SecurityPolicy` which is read back into the `security_policy` terraform attribute.

## Goals / Non-Goals

**Goals:**
- Ensure the `security_policy` parameter is properly supported in the `ModifySecurityPolicy` API call (create/update)
- Ensure the `DescribeSecurityPolicy` API properly reads back the `SecurityPolicy` response into the `security_policy` attribute
- Ensure proper input parameters (`ZoneId`, `Entity`, `Host`, `TemplateId`) are passed to `DescribeSecurityPolicy`
- Maintain backward compatibility with existing terraform configurations

**Non-Goals:**
- Modifying the `security_config` (classic) parameter handling
- Adding new sub-fields within the `SecurityPolicy` structure beyond what the vendor SDK supports
- Changing the resource ID format or import behavior

## Decisions

1. **Parameter mapping**: The `request.SecurityPolicy` field in `ModifySecurityPolicy` maps to the terraform schema field `security_policy`. This is a complex nested structure (TypeList with MaxItems: 1) containing all security policy sub-configurations.

2. **Read operation**: The `DescribeSecurityPolicy` API requires `ZoneId`, `Entity`, `Host`, and `TemplateId` as input parameters (extracted from the resource ID) and returns `SecurityPolicy` in the response which is flattened into the `security_policy` attribute.

3. **Resource lifecycle**: Since there is no `CreateSecurityPolicy` or `DeleteSecurityPolicy` API, the resource uses `ModifySecurityPolicy` for both create and update operations. Delete is a no-op (the security policy always exists for a zone).

4. **Retry handling**: API calls use `tccommon.ReadRetryTimeout` with `resource.RetryContext` for retry logic, wrapping errors with `tccommon.RetryError()`.

## Risks / Trade-offs

- [Risk] The `SecurityPolicy` structure is complex with deeply nested fields → Mitigation: Follow existing patterns in the resource code for flattening/expanding nested structures.
- [Risk] Backward compatibility with existing state files → Mitigation: Only add Optional fields, never modify existing Required fields or change field types.
