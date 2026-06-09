## Context

The `tencentcloud_teo_l4_proxy` resource manages TEO (TencentCloud EdgeOne) Layer 4 proxy instances. Currently, the resource captures `ProxyId` from the CreateL4Proxy API response and uses it internally as part of the composite resource ID (`zoneId#proxyId`), but does not expose it as a user-accessible schema attribute.

The `ProxyId` is also returned by the DescribeL4Proxy API in the `L4Proxy` struct, making it available during the read operation.

## Goals / Non-Goals

**Goals:**
- Expose `proxy_id` as a computed attribute on the `tencentcloud_teo_l4_proxy` resource so users can reference the L4 proxy instance ID in their Terraform configurations
- Mark `ddos_protection_config` as deprecated
- Maintain full backward compatibility with existing configurations and state

**Non-Goals:**
- Changing the resource ID format (it remains `zoneId#proxyId`)
- Adding any other new parameters beyond `proxy_id`
- Removing `ddos_protection_config` entirely (it remains for backward compatibility)

## Decisions

1. **Add `proxy_id` as a Computed attribute (not Optional+Computed)**
   - Rationale: `proxy_id` is assigned by the cloud API upon creation and cannot be user-specified. Using `Computed: true` alone clearly communicates this is a read-only server-assigned value.

2. **Set `proxy_id` in the read function from DescribeL4Proxy response**
   - Rationale: The `L4Proxy` struct returned by `DescribeTeoL4ProxyById` already contains `ProxyId`. Setting it in the read function follows the standard pattern and ensures consistency after import operations.

3. **Set `proxy_id` explicitly in the create function before calling Read**
   - Rationale: After the `CreateL4Proxy` API call, `proxy_id` is set from `response.Response.ProxyId` via `d.Set("proxy_id", proxyId)` in the create function. The subsequent `resourceTencentCloudTeoL4ProxyRead` call will also set it from the DescribeL4Proxy response, providing redundancy.

4. **`proxy_id` is set from `respData.ProxyId` in the read function**
   - Rationale: When importing or refreshing, the read function retrieves the L4Proxy from the DescribeL4Proxy API. The `ProxyId` field is set from `respData.ProxyId` when not nil, which is the canonical source.

5. **Mark `ddos_protection_config` as deprecated with `Deprecated` and `Computed: true`**
   - Rationale: The cloud API has deprecated `DDosProtectionConfig`. Adding `Deprecated: "It has been deprecated from version 1.82.90."` follows the provider's standard deprecation pattern. Adding `Computed: true` to the field and all sub-fields (`level_mainland`, `max_bandwidth_mainland`, `level_overseas`) ensures existing state is preserved and API-returned values can be stored without requiring user input.

## Risks / Trade-offs

- **[Risk] Adding a computed attribute changes the state shape** → Mitigation: Computed attributes are added to state automatically on the next read/refresh. No state migration is needed. This is backward compatible.
- **[Risk] Deprecating `ddos_protection_config` may affect existing users** → Mitigation: The field remains functional; only a deprecation warning is shown. Users can continue using it but are advised to stop. The `Computed: true` addition ensures the field can hold API-returned values without user input.
