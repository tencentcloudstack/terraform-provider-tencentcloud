## Context

Tencent Cloud EdgeOne (TEO) provides Multi-Path Gateway functionality that allows users to configure secure acceleration for origin servers. When the origin IP segments (回源 IP 网段) are updated by the cloud service, users need to confirm the update via the `ConfirmMultiPathGatewayOriginACL` API to acknowledge that they have updated their origin firewall configurations. This is a CONFIG-type resource where the resource (zone + gateway) exists independently and the config is always present — the Terraform resource manages reading the current ACL state and confirming version updates.

The resource follows the same pattern as other TEO config resources (e.g., `tencentcloud_teo_certificate_config`, `tencentcloud_teo_security_policy_config`):
- Create: Sets the composite ID and delegates to Update
- Read: Calls `DescribeMultiPathGatewayOriginACL` to fetch current ACL state
- Update: Calls `ConfirmMultiPathGatewayOriginACL` to confirm a version, then reads back
- Delete: No-op (config resource, cannot be truly deleted)

## Goals / Non-Goals

**Goals:**
- Implement a Terraform RESOURCE_KIND_CONFIG resource `tencentcloud_teo_confirm_multi_path_gateway_origin_acl`
- Support Read operation via `DescribeMultiPathGatewayOriginACL` to show current and pending ACL info
- Support Update operation via `ConfirmMultiPathGatewayOriginACL` to confirm origin ACL version updates
- Use `zone_id` and `gateway_id` as composite resource ID (joined by `tccommon.FILED_SP`)
- Expose the full ACL info structure (current + next) as computed output fields
- Support resource import
- Add unit tests using gomonkey mock approach
- Add documentation (.md file)

**Non-Goals:**
- This resource does NOT manage the creation or deletion of multi-path gateways themselves
- This resource does NOT manage the actual IP segments — those are controlled by the cloud service; this resource only confirms version updates
- No support for creating or deleting the underlying gateway

## Decisions

### 1. Resource Type: RESOURCE_KIND_CONFIG
**Decision**: Implement as a CONFIG resource with RU (Read-Update) pattern.
**Rationale**: The origin ACL config always exists as long as the zone and gateway exist. The `ConfirmMultiPathGatewayOriginACL` API confirms/updates the version, and `DescribeMultiPathGatewayOriginACL` reads the state. There is no Create or Delete API — this matches the CONFIG resource pattern perfectly.

### 2. Composite ID: zone_id + gateway_id
**Decision**: Use `zone_id` and `gateway_id` joined by `tccommon.FILED_SP` as the resource ID.
**Rationale**: Both APIs require these two parameters to identify the gateway. This follows the established pattern in the codebase (e.g., `tencentcloud_teo_certificate_config` uses `zone_id + host`).

### 3. Schema Design
**Decision**: Input parameters (`zone_id`, `gateway_id`, `origin_acl_version`) are set as user-configurable fields. The full `multi_path_gateway_origin_acl_info` is a computed output block.
**Rationale**:
- `zone_id` and `gateway_id` are required identifiers (ForceNew: true)
- `origin_acl_version` is an optional parameter for the Update/Confirm operation — when set, it confirms the specified version
- The response structure `MultiPathGatewayOriginACLInfo` (containing `MultiPathGatewayCurrentOriginACL` and `MultiPathGatewayNextOriginACL`) is a computed read-only output since it reflects the server-side state

### 4. Create and Delete Operations
**Decision**: Create sets the ID and delegates to Update (which may be a no-op if no `origin_acl_version` is specified). Delete is a no-op.
**Rationale**: CONFIG resources don't have real Create/Delete operations. The config exists as long as the underlying resource exists. The Delete being a no-op follows the pattern of other config resources.

### 5. Read Operation with Retry
**Decision**: Use `helper.Retry()` with `tccommon.ReadRetryTimeout` for the Read operation.
**Rationale**: Following the established pattern for cloud API calls with eventual consistency.

### 6. No Async Polling Needed
**Decision**: The `ConfirmMultiPathGatewayOriginACL` API does not appear to be asynchronous based on the SDK — it returns synchronously. No polling loop is needed.
**Rationale**: The API response only contains `RequestId`, indicating a synchronous operation. After calling Confirm, we immediately call Read to refresh the state.

## Risks / Trade-offs

- **[Risk] Origin ACL version mismatch**: If the user specifies an `origin_acl_version` that doesn't match the pending next version, the API may return an error. → **Mitigation**: Terraform will surface the API error, and the user can re-read to check the current version.
- **[Risk] External changes**: The ACL info can change externally (e.g., cloud service updates IP segments). → **Mitigation**: Terraform Read will detect the current state on every plan. The computed output fields will reflect the latest server state.
- **[Trade-off] No real delete**: As a CONFIG resource, Delete is a no-op. Removing the resource from Terraform state doesn't affect the actual cloud configuration. This is by design for CONFIG resources.
