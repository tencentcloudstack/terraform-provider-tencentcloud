## Context

The TencentCloud Terraform Provider currently supports various EdgeOne (TEO) resources and operations, but lacks a resource to trigger cache purge operations. The TEO SDK provides `CreatePurgeTask` API to trigger cache purges and `DescribePurgeTasks` API to query purge task status and results.

This is a RESOURCE_KIND_OPERATION type resource, meaning it performs a one-time action (cache purge) and does not need to track persistent state. The standard pattern for operation resources in this codebase includes:
- Create: Calls the action API (with optional polling for async tasks)
- Read: Empty (returns nil)
- Delete: Empty (returns nil)
- No Update method

Existing TEO operation resources like `tencentcloud_teo_identify_zone_operation` and `tencentcloud_teo_l7_acc_rule_priority_operation` follow this pattern.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_teo_purge_task` that triggers EdgeOne cache purge operations
- Support all purge types: `purge_url`, `purge_prefix`, `purge_host`, `purge_all`, `purge_cache_tag`
- After triggering the purge, poll `DescribePurgeTasks` to verify task completion and return task results as computed attributes
- Follow the established OPERATION resource pattern with Create/Read/Delete (RUD empty)
- All schema fields set as ForceNew since this is an operation resource

**Non-Goals:**
- This resource will NOT support updating an existing purge task (operation resources are immutable)
- This resource will NOT persist state beyond the purge operation lifecycle
- This resource will NOT support importing (one-time operation has no persistent state to import)

## Decisions

### 1. Create API: Use `CreatePurgeTask`
The `CreatePurgeTask` API triggers the actual cache purge. It accepts `ZoneId`, `Type`, `Method`, `Targets`, and `CacheTag` parameters. This is the primary action API for this operation resource.

**Rationale**: `CreatePurgeTask` is the designated API for triggering purge operations in the TEO SDK. No alternative exists.

### 2. Polling: Use `DescribePurgeTasks` to verify completion
After calling `CreatePurgeTask`, the returned `JobId` is used to poll `DescribePurgeTasks` until the task reaches a terminal state (`success`, `failed`, `timeout`, `canceled`).

**Rationale**: The purge operation is asynchronous. Polling ensures Terraform waits until the purge completes before considering the resource creation done. The `DescribePurgeTasks` API supports filtering by `job-id`, making it suitable for targeted status checks.

**Alternative considered**: Not polling and returning immediately. Rejected because users need to know if the purge succeeded or failed.

### 3. Resource ID: Use `helper.BuildToken()`
Since this is a one-time operation with no persistent entity to identify, use `helper.BuildToken()` to generate a unique ID, consistent with other operation resources like `tencentcloud_teo_identify_zone_operation`.

**Rationale**: The purge task is ephemeral; there is no persistent resource ID to track. A generated token provides a unique Terraform state identifier.

### 4. Schema fields: All ForceNew
All input parameters (`zone_id`, `type`, `method`, `targets`, `cache_tag`) have `ForceNew: true`. Computed output fields (`job_id`, `tasks`) are not ForceNew.

**Rationale**: Operation resources are immutable by design. Any parameter change triggers a new purge operation (destroy + create).

### 5. CacheTag as a list block with domains sub-field
The `CacheTag` SDK struct contains a `Domains` field. In Terraform schema, this maps to a `TypeList` block with a nested `domains` field.

**Rationale**: Matches the SDK structure and provides clear type safety for the cache tag configuration.

### 6. Tasks output as computed TypeList
The `DescribePurgeTasks` response returns a list of `Task` structs. These are exposed as a computed `TypeList` attribute so users can inspect purge task results.

**Rationale**: Users need to see purge task status, targets, and failure information. Computed output avoids requiring user input while providing full visibility.

## Risks / Trade-offs

- **[Purge quota limits]** → The TEO API has daily and batch quotas for purge operations. Terraform cannot enforce these limits; API errors will surface if quotas are exceeded. Mitigation: Document the quota limitation and surface API error messages clearly.
- **[Async polling timeout]** → If purge tasks take longer than the Terraform timeout, the resource creation will fail. Mitigation: Use `tccommon.ReadRetryTimeout` for polling, which provides adequate timeout for typical purge operations.
- **[No state persistence]** → Since Read is empty, `terraform refresh` will not update the `tasks` output. Mitigation: This is by design for operation resources; the `tasks` output is set during Create only.
