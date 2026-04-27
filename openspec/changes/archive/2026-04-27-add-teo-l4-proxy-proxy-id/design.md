## Context

The `tencentcloud_teo_l4_proxy` resource manages TEO (TencentCloud EdgeOne) Layer 4 proxy instances. Currently, the resource captures `ProxyId` from the CreateL4Proxy API response and uses it internally as part of the composite resource ID (`zoneId#proxyId`), but does not expose it as a user-accessible schema attribute.

The `ProxyId` is also returned by the DescribeL4Proxy API in the `L4Proxy` struct, making it available during the read operation.

## Goals / Non-Goals

**Goals:**
- Expose `proxy_id` as a computed attribute on the `tencentcloud_teo_l4_proxy` resource so users can reference the L4 proxy instance ID in their Terraform configurations
- Maintain full backward compatibility with existing configurations and state

**Non-Goals:**
- Changing the resource ID format (it remains `zoneId#proxyId`)
- Modifying any existing schema attributes
- Adding any other new parameters beyond `proxy_id`

## Decisions

1. **Add `proxy_id` as a Computed attribute (not Optional+Computed)**
   - Rationale: `proxy_id` is assigned by the cloud API upon creation and cannot be user-specified. Using `Computed: true` alone clearly communicates this is a read-only server-assigned value.

2. **Set `proxy_id` in the read function from DescribeL4Proxy response**
   - Rationale: The `L4Proxy` struct returned by `DescribeTeoL4ProxyById` already contains `ProxyId`. Setting it in the read function follows the standard pattern and ensures consistency after import operations.

3. **No changes needed in create function for setting `proxy_id`**
   - Rationale: The create function already calls `resourceTencentCloudTeoL4ProxyRead` at the end, which will set `proxy_id`. No separate `d.Set("proxy_id", ...)` call is needed in the create function.

4. **`proxy_id` should also be set from the ID split in the read function**
   - Rationale: When importing, the read function parses the composite ID to get the proxyId. The `ProxyId` field should be set from `respData.ProxyId` when available, which is the canonical source.

## Risks / Trade-offs

- **[Risk] Adding a computed attribute changes the state shape** → Mitigation: Computed attributes are added to state automatically on the next read/refresh. No state migration is needed. This is backward compatible.
