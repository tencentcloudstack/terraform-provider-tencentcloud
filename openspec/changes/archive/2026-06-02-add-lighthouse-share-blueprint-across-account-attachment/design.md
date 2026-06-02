## Context

The TencentCloud Lighthouse service provides APIs for sharing blueprints (images) across different TencentCloud accounts. The SDK (already vendored) provides:

- `ShareBlueprintAcrossAccounts`: Share a blueprint to target accounts (bind)
- `DescribeBlueprintsShareAcrossAccountInfos`: Query sharing information (read)
- `CancelShareBlueprintAcrossAccounts`: Cancel sharing (unbind)

Currently there is no Terraform resource to manage these sharing relationships. This attachment resource will fill that gap, following the existing patterns in `tencentcloud_lighthouse_key_pair_attachment` and `tencentcloud_lighthouse_disk_attachment`.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource (`tencentcloud_lighthouse_share_blueprint_across_account_attachment`) to manage blueprint cross-account sharing
- Support Create (share), Read (query share status), and Delete (cancel share) operations
- Support Terraform import for existing sharing relationships
- Use composite ID format (`blueprint_id#account_id`) consistent with existing attachment resources

**Non-Goals:**
- No Update operation (sharing relationships are immutable - must delete and recreate to change)
- No support for listing all shared blueprints as a datasource (separate concern)
- No modification of existing Lighthouse resources or datasources

## Decisions

### Decision 1: Schema Design - All Fields ForceNew

**Choice**: All schema fields (`blueprint_id`, `account_ids`) are `Required` + `ForceNew: true`.

**Rationale**: This is an attachment resource managing a binding relationship. The Create API (`ShareBlueprintAcrossAccounts`) and Delete API (`CancelShareBlueprintAcrossAccounts`) only accept `BlueprintId` and `AccountIds`. There is no update/modify API for changing a sharing relationship. Users must destroy (`terraform destroy`) and recreate to change the relationship.

**Alternatives considered**:
- Making `account_ids` optional with `ForceNew: false` would imply updatability, which the API doesn't support.
- Adding an Update function that calls Delete+Create would be error-prone and inconsistent with the attachment resource pattern.

### Decision 2: Composite Resource ID

**Choice**: Use `blueprint_id` + `tccommon.FILED_SP` + `account_ids` as the resource ID (e.g., `lhbp-xxx#100012345678#100087654321`).

**Rationale**: The sharing relationship is defined by the combination of blueprint and target accounts. Following the existing pattern in `resource_tc_lighthouse_key_pair_attachment.go` which uses `keyId#instanceId` composite ID. Since `account_ids` is a list, we join the sorted account IDs using `tccommon.FILED_SP`.

**Alternatives considered**:
- Using only `blueprint_id` as ID: Would not support multiple attachment resources for the same blueprint with different account sets.
- Hashing the combined values: Less readable and harder to debug.

### Decision 3: Read Operation - Query by BlueprintId

**Choice**: In the Read function, call `DescribeBlueprintsShareAcrossAccountInfos` with `BlueprintIds` containing the `blueprint_id`, then filter results to find matching accounts.

**Rationale**: The describe API supports filtering by `BlueprintIds`. Since the resource tracks a specific blueprint sharing a specific set of accounts, querying by blueprint ID is the most direct approach. If the sharing info is not found (either blueprint not shared or accounts don't match), set `d.SetId("")` to indicate the resource no longer exists.

**Alternatives considered**:
- Using the `Filters` parameter with `account-id` filter: More complex, less direct. The `BlueprintIds` parameter is simpler and more reliable.

### Decision 4: account_ids as List with Sort-consistency

**Choice**: `account_ids` is defined as `TypeList` (ordered), and accounts are sorted in the Read function for consistent comparison.

**Rationale**: The API may return account IDs in any order. To ensure Terraform state consistency, the read function should sort the accounts before comparing with state. The `account_ids` type uses `TypeList` to maintain order.

## Risks / Trade-offs

- **[Low Risk] API Latency**: The `ShareBlueprintAcrossAccounts` and `CancelShareBlueprintAcrossAccounts` are synchronous operations. No async polling is needed based on the API documentation.
- **[Medium Risk] Partial Sharing**: If `account_ids` contains multiple accounts and the API partially succeeds (shares to some but not all), Terraform may be out of sync. Mitigation: The read function will detect the actual state and surface discrepancies.
- **[Low Risk] Idempotency**: Sharing an already-shared blueprint or canceling sharing that doesn't exist: These scenarios are handled by the API itself (it returns success or error), and Terraform's state management ensures consistency.