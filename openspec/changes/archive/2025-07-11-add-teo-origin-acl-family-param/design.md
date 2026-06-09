## Context

The `tencentcloud_teo_origin_acl` resource manages TEO (TencentCloud EdgeOne) origin ACL protection for a zone. Currently, the resource only supports configuring L7 hosts and L4 proxy IDs but does not expose the `OriginACLFamily` parameter, which controls the geographic domain for origin ACL protection.

The cloud API already supports `OriginACLFamily` in:
- **EnableOriginACL** (Create): Sets the control domain when enabling origin ACL. Defaults to standard global (`gaz`) if not specified.
- **ModifyOriginACL** (Update): Changes the control domain. If not specified, the domain remains unchanged.
- **DescribeOriginACL** (Read): Returns the current `OriginACLFamily` value in the `OriginACLInfo` response struct.
- **DisableOriginACL** (Delete): Does not include `OriginACLFamily` — deletion is zone-level and does not need this parameter.

The parameter is of type `*string` with valid values: `gaz`, `mlc`, `emc`, `plat-gaz`, `plat-mlc`, `plat-emc`.

## Goals / Non-Goals

**Goals:**
- Add `origin_acl_family` as an Optional + Computed string parameter to the `tencentcloud_teo_origin_acl` resource schema
- Pass `OriginACLFamily` in Create (EnableOriginACL) and Update (ModifyOriginACL) API calls when the user specifies a value
- Read `OriginACLFamily` from the DescribeOriginACL response in the Read handler
- Add `origin_acl_family` to the data source `tencentcloud_teo_origin_acl` under the `origin_acl_info` block
- Maintain full backward compatibility — existing configurations without `origin_acl_family` continue to work

**Non-Goals:**
- Adding validation for the `origin_acl_family` valid values (the API will reject invalid values)
- Adding a separate data source for `DescribeAvailableOriginACLFamily`
- Changing the resource ID format or any existing schema fields

## Decisions

### Decision 1: Schema field type — Optional + Computed
**Choice**: `Optional: true, Computed: true` with `TypeString`
**Rationale**: The API defaults to `gaz` when not specified on Enable, and preserves the current value when not specified on Modify. Using `Computed: true` ensures the Read handler populates the actual value from the API response, so the state reflects the real configuration even when the user didn't explicitly set it.

### Decision 2: Update handler — separate ModifyOriginACL call for origin_acl_family
**Choice**: When `origin_acl_family` changes, make a separate ModifyOriginACL call with only `ZoneId` and `OriginACLFamily` set (without `OriginACLEntities`).
**Rationale**: The ModifyOriginACL API accepts both `OriginACLEntities` and `OriginACLFamily` in the same request. Since `origin_acl_family` is independent of the entity-level changes (l7_hosts/l4_proxy_ids), it can be sent in the same call or a separate one. Using a separate call keeps the logic clean and avoids mixing entity operations with family changes. However, the API allows both in one call, so we could also combine them — but separation is cleaner for error handling.

Actually, on further review, it's simpler and more consistent to set `OriginACLFamily` in ALL ModifyOriginACL calls when `d.HasChange("origin_acl_family")` is true. This way, the family change is propagated with the first ModifyOriginACL batch call, reducing the number of API calls.

**Final choice**: When `origin_acl_family` has changed, set the `OriginACLFamily` field on the first ModifyOriginACL request in the update handler. If there are no entity changes but only origin_acl_family change, make a standalone ModifyOriginACL call with just ZoneId and OriginACLFamily.

### Decision 3: Data source field placement
**Choice**: Add `origin_acl_family` as a Computed string field inside the `origin_acl_info` block of the data source.
**Rationale**: The `OriginACLFamily` field is part of the `OriginACLInfo` struct in the API response, so it naturally belongs in the `origin_acl_info` block of the data source.

## Risks / Trade-offs

- **[Risk] API rejects invalid OriginACLFamily values** → The API returns an error for unsupported values. Users will see the API error message. This is acceptable and consistent with how other string parameters work in the provider.
- **[Risk] Changing origin_acl_family might affect existing ACL rules** → The API documentation states this is a control domain setting. Changing it could alter which IP ranges are used. This is by design — users explicitly opt into this change.
- **[Trade-off] Not adding ValidateFunc for valid values** → We could add validation, but the API already validates and the valid values might change over time. Relying on API validation is more maintainable.
