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
2. Calling `CancelShareBlueprintAcrossAccounts` for removed accounts (in batches)
3. Calling `ShareBlueprintAcrossAccounts` for added accounts (in batches)
4. Calling Read to refresh state

**Rationale**: This approach provides several benefits:
- No intermediate state where all sharing is temporarily cancelled
- More efficient API usage (only operate on changed accounts)
- Consistent with the console behavior (incremental add/remove)
- Better user experience (no disruption to existing shares)

**Alternatives considered**:
- Destroy + Recreate: Simpler but causes temporary complete removal of sharing.

### Decision 4: Read Operation - Query by BlueprintId with Pagination

**Choice**: In the Read function, call `DescribeBlueprintsShareAcrossAccountInfos` with `BlueprintIds` containing the `blueprint_id`, then paginate through results using `Offset`/`Limit` (page size 100) until all shared account IDs are collected.

**Rationale**: Since the resource ID is just `blueprint_id`, the Read function should return the actual state of all accounts this blueprint is shared with. The API supports pagination via `Limit`/`Offset` parameters and returns `TotalCount`. If there are many shared accounts, a single page may not be sufficient, so we must loop through all pages.

**Alternatives considered**:
- Single query without pagination: Would miss accounts if more than one page of data exists.
- Filtering by specific accounts: Not needed since we're tracking the full share state for the blueprint.

### Decision 5: Use TypeSet for account_ids

**Choice**: Define `account_ids` as `schema.TypeSet` (not `schema.TypeList`) to handle ordering automatically.

**Rationale**: Using TypeSet means Terraform treats `account_ids` as an unordered collection. Since the API returns account IDs in non-deterministic order and users may specify them in any order, TypeSet eliminates spurious plan diffs caused by ordering differences. No manual sorting logic is needed — the framework handles this natively.

**Alternatives considered**:
- Using TypeList with manual sorting: Adds unnecessary sorting complexity; TypeSet handles this automatically.
- No ordering control: Would cause unnecessary plan diffs when only order differs, confusing users.

### Decision 6: Batch Processing for Write APIs

**Choice**: Both `ShareBlueprintAcrossAccounts` and `CancelShareBlueprintAcrossAccounts` APIs have a maximum limit of 10 account IDs per request. All write operations (Create, Update, Delete) split the account ID list into batches of 10 and make sequential API calls.

**Rationale**: This ensures the resource works correctly regardless of how many account IDs the user provides. Without batching, providing more than 10 accounts would cause API errors. The batching is transparent to users — they simply provide any number of account IDs and the provider handles splitting internally.

**Implementation details**:
- Constant `shareBlueprintBatchSize = 10` defines the batch size
- Create loops through batches calling `ShareBlueprintAcrossAccounts`
- Update splits `toRemove` into batches for `CancelShareBlueprintAcrossAccounts` and `toAdd` into batches for `ShareBlueprintAcrossAccounts`
- Delete loops through batches calling `CancelShareBlueprintAcrossAccounts`

## Risks / Trade-offs

- **[Low Risk] API Latency**: The `ShareBlueprintAcrossAccounts` and `CancelShareBlueprintAcrossAccounts` are synchronous operations. No async polling is needed based on the API documentation.
- **[Medium Risk] Partial Sharing**: If `account_ids` contains multiple accounts and one of the update APIs partially fails, Terraform may be out of sync. Mitigation: The read function will detect the actual state and surface discrepancies.
- **[Low Risk] Idempotency**: Sharing an already-shared blueprint or canceling sharing that doesn't exist: These scenarios are handled by the API itself (it returns success or error), and Terraform's state management ensures consistency.
