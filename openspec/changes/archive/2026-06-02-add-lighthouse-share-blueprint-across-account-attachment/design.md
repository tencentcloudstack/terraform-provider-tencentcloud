## Context

The TencentCloud Lighthouse service provides APIs for sharing blueprints (images) across different TencentCloud accounts. The SDK (already vendored) provides:

- `ShareBlueprintAcrossAccounts`: Share a blueprint to target accounts (bind)
- `DescribeBlueprintsShareAcrossAccountInfos`: Query sharing information (read)
- `CancelShareBlueprintAcrossAccounts`: Cancel sharing (unbind)

Currently there is no Terraform resource to manage these sharing relationships. This attachment resource will fill that gap, following the existing patterns in `tencentcloud_lighthouse_key_pair_attachment` and `tencentcloud_lighthouse_disk_attachment`.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource (`tencentcloud_lighthouse_share_blueprint_across_account_attachment`) to manage blueprint cross-account sharing
- Support Create (share), Read (query share status), Update (incrementally modify shared accounts), and Delete (cancel share) operations
- Support Terraform import for existing sharing relationships using only `blueprint_id` as ID
- Use simple ID format (`blueprint_id`) to enable easy import without needing to specify account IDs

**Non-Goals:**
- No support for listing all shared blueprints as a datasource (separate concern)
- No modification of existing Lighthouse resources or datasources

## Decisions

### Decision 1: Schema Design - BlueprintId ForceNew, account_ids Updatable

**Choice**: 
- `blueprint_id` is `Required` + `ForceNew: true` (immutable - changing blueprint requires recreation)
- `account_ids` is `Required` + `ForceNew: false` (updatable - supports incremental add/remove)

**Rationale**: This design allows users to incrementally add or remove accounts from an existing sharing relationship without destroying and recreating the entire resource. The Update function implements incremental updates by calculating the diff between old and new account lists, then calling the appropriate APIs to remove deleted accounts and add new ones.

**Alternatives considered**:
- Making all fields `ForceNew: true` with no update API: Would require destroy+create for any account change, causing temporary removal of all sharing relationships.
- Implementing Update as Delete+Create: Would have the same issues as ForceNew approach.

### Decision 2: Simple Resource ID (BlueprintId Only)

**Choice**: Use `blueprint_id` as the resource ID (e.g., `lhbp-xxx`).

**Rationale**: Using only `blueprint_id` as the resource ID simplifies import operations significantly. Users can import using just the blueprint ID, and the Read function will query the actual account IDs from the API. This eliminates the need for users to know or specify the account list during import.

**Alternatives considered**:
- Composite ID (`blueprint_id#account_ids`): Would require users to know and specify all account IDs during import, which is inconvenient.
- Auto-generated UUID: Less meaningful and harder to debug.

### Decision 3: Incremental Update Strategy

**Choice**: The Update function implements incremental updates by:
1. Calculating the diff between old and new `account_ids`
2. Calling `CancelShareBlueprintAcrossAccounts` for removed accounts
3. Calling `ShareBlueprintAcrossAccounts` for added accounts
4. Calling Read to refresh state

**Rationale**: This approach provides several benefits:
- No intermediate state where all sharing is temporarily cancelled
- More efficient API usage (only operate on changed accounts)
- Consistent with the console behavior (incremental add/remove)
- Better user experience (no disruption to existing shares)

**Alternatives considered**:
- Destroy + Recreate: Simpler but causes temporary complete removal of sharing.

### Decision 4: Read Operation - Query by BlueprintId

**Choice**: In the Read function, call `DescribeBlueprintsShareAcrossAccountInfos` with `BlueprintIds` containing the `blueprint_id`, then return all currently shared account IDs.

**Rationale**: Since the resource ID is now just `blueprint_id`, the Read function should return the actual state of all accounts this blueprint is shared with. This ensures Terraform state stays in sync with reality.

**Alternatives considered**:
- Filtering by specific accounts: Not needed since we're tracking the full share state for the blueprint.

### Decision 5: Account IDs Sorting for State Consistency

**Choice**: Both Create and Read functions sort `account_ids` in ascending order before storing to state.

**Rationale**: The API may return account IDs in any order, and users may specify them in any order in their configuration. To prevent Terraform from showing spurious diffs during `plan`, we ensure consistent ordering:
1. **Create function**: Sorts `account_ids` after successful API call and before storing to state
2. **Read function**: Sorts `account_ids` returned from API before updating state

This guarantees that regardless of input order or API return order, the state always contains sorted account IDs. When a user runs `terraform plan`, if the set of accounts is the same (even if specified in different order), there will be no changes shown.

**Alternatives considered**:
- Using TypeSet instead of TypeList: Would handle ordering automatically but loses list semantics and position information.
- No sorting: Would cause unnecessary plan diffs when only order differs, confusing users.

## Risks / Trade-offs

- **[Low Risk] API Latency**: The `ShareBlueprintAcrossAccounts` and `CancelShareBlueprintAcrossAccounts` are synchronous operations. No async polling is needed based on the API documentation.
- **[Medium Risk] Partial Sharing**: If `account_ids` contains multiple accounts and one of the update APIs partially fails, Terraform may be out of sync. Mitigation: The read function will detect the actual state and surface discrepancies.
- **[Low Risk] Idempotency**: Sharing an already-shared blueprint or canceling sharing that doesn't exist: These scenarios are handled by the API itself (it returns success or error), and Terraform's state management ensures consistency.
