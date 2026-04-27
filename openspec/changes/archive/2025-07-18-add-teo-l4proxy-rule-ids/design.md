## Context

The `tencentcloud_teo_l4_proxy_rule` resource manages Layer 4 proxy rules for TencentCloud EdgeOne (TEO). When creating L4 proxy rules via the `CreateL4ProxyRules` API, the response includes `L4ProxyRuleIds` — a list of IDs assigned to the newly created rules. Currently, the resource extracts the first rule ID to compose the composite resource ID (`zoneId#proxyId#ruleId`), but the full list of rule IDs is not persisted in the Terraform state.

The resource uses a composite ID format: `zoneId#proxyId#ruleId`, where `ruleId` is the first element of `L4ProxyRuleIds` from the create response.

## Goals / Non-Goals

**Goals:**
- Add `l4proxy_rule_ids` as a computed attribute (TypeList of TypeString) to the resource schema
- Persist the `L4ProxyRuleIds` values from the `CreateL4ProxyRules` API response into the Terraform state
- Ensure the `l4proxy_rule_ids` is populated in the read function for state consistency
- Maintain full backward compatibility with existing configurations

**Non-Goals:**
- Changing the existing composite ID format or the `l4_proxy_rules` nested block structure
- Adding any other new parameters beyond `l4proxy_rule_ids`
- Modifying the update or delete logic

## Decisions

1. **Schema type: TypeList of TypeString (Computed)**
   - The `L4ProxyRuleIds` from the API is `[]*string`, so TypeList of TypeString maps directly
   - Computed only — users cannot set this value; it is populated from the API response
   - Alternative considered: TypeSet — rejected because rule IDs are ordered and duplicates are not expected

2. **Populate in create and read functions**
   - In create: set `l4proxy_rule_ids` from `response.Response.L4ProxyRuleIds` after the API call succeeds
   - In read: reconstruct `l4proxy_rule_ids` from the rule ID in the composite ID (since the DescribeL4ProxyRules API returns individual rule details, not the full list of IDs from batch creation)
   - Actually, since the resource creates one rule at a time (takes `l4_proxy_rules` with MaxItems:1), the list will contain exactly one rule ID, which is the same as the third component of the composite ID

3. **No changes to update/delete**
   - `l4proxy_rule_ids` is computed and read-only, so it does not participate in update or delete operations

## Risks / Trade-offs

- **State consistency**: Since the resource only supports one rule (MaxItems:1 for `l4_proxy_rules`), `l4proxy_rule_ids` will always be a single-element list. This is consistent with the current behavior where only the first rule ID from the response is used.
- **Read reconstruction**: On read, the rule ID is reconstructed from the composite ID rather than from a separate API call, which is reliable since the ID was originally set from the create response.
