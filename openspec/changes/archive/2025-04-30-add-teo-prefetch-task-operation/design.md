## Context

Terraform Provider for TencentCloud needs to support TEO (EdgeOne) prefetch cache tasks. Currently, there is no Terraform resource to manage TEO prefetch operations. The TEO SDK provides two relevant APIs:

- **CreatePrefetchTask**: Submits URL prefetch tasks to pre-warm cache on EdgeOne nodes. This is an asynchronous API that returns a `JobId`.
- **DescribePrefetchTasks**: Queries the status of prefetch tasks, returning task details including status (processing/success/failed/timeout/canceled).

The resource type is `RESOURCE_KIND_OPERATION` — a one-time operation that does not persist state. The Create method submits the task and polls for completion, while Read/Update/Delete methods are empty.

Reference implementation: `resource_tc_teo_import_zone_config_operation.go` in the same service package follows the same pattern.

## Goals / Non-Goals

**Goals:**
- Add `tencentcloud_teo_prefetch_task` resource supporting prefetch cache operations on TEO
- Implement Create method that calls `CreatePrefetchTask` API and polls `DescribePrefetchTasks` until task completes
- Support all CreatePrefetchTask input parameters: ZoneId, Targets, Mode, Headers, PrefetchMediaSegments (skip deprecated EncodeUrl)
- Support querying task results via DescribePrefetchTasks with filters (job-id filter for polling)
- Follow RESOURCE_KIND_OPERATION pattern: empty Read/Update/Delete methods
- Register resource in provider.go and provider.md
- Add unit tests with gomonkey mock
- Add resource documentation .md file

**Non-Goals:**
- This resource does not manage persistent state — it represents a one-time operation
- No support for updating or deleting prefetch tasks (API does not support these operations)
- No data source for listing prefetch tasks (out of scope for this change)

## Decisions

### 1. Resource ID format
**Decision**: Use `zoneId + FILED_SP + jobId` as composite resource ID.
**Rationale**: Both zoneId and jobId are needed to query task status. Using the standard FILED_SP separator follows existing patterns (e.g., `resource_tc_teo_import_zone_config_operation.go`).

### 2. Skip deprecated EncodeUrl parameter
**Decision**: Do not include `EncodeUrl` in the Terraform schema since it's marked as deprecated in the SDK.
**Rationale**: Deprecated API parameters should not be exposed in new Terraform resources.

### 3. Polling strategy for async CreatePrefetchTask
**Decision**: After calling CreatePrefetchTask, poll DescribePrefetchTasks using the job-id filter with `6 * tccommon.ReadRetryTimeout` total timeout.
**Rationale**: Follows the same pattern as `resource_tc_teo_import_zone_config_operation.go`. The polling checks task status and waits until status is not "processing". If status is "failed" or "timeout", return a NonRetryableError.

### 4. Schema design for Create input parameters
**Decision**:
- `zone_id` (Required, ForceNew): Zone ID
- `targets` (Required, ForceNew): List of URLs to prefetch
- `mode` (Optional, ForceNew): Prefetch mode (default/edge)
- `headers` (Optional, ForceNew): List of HTTP headers (name/value pairs)
- `prefetch_media_segments` (Optional, ForceNew): Media segment prefetch control (on/off)

**Rationale**: ZoneId and Targets are required by the API. Mode, Headers, and PrefetchMediaSegments are optional. All are ForceNew since this is an operation resource.

### 5. Task result computed fields
**Decision**: Add computed attributes from DescribePrefetchTasks response for task results:
- `job_id` (Computed): Task job ID
- `tasks` (Computed): List of task results with fields (job_id, target, type, method, status, create_time, update_time, fail_type, fail_message)

**Rationale**: Users need to see the results of their prefetch operation. These are set during the Create polling phase.

## Risks / Trade-offs

- **[Async task timeout]** → Mitigation: Use generous retry timeout (6 * ReadRetryTimeout) for polling. If the task takes too long, the user will see a timeout error and can check the task status manually.
- **[API rate limiting]** → Mitigation: Polling uses standard ReadRetryTimeout intervals, not aggressive polling.
- **[White-list features]** → Mode=edge and PrefetchMediaSegments are white-list features per API docs. No special handling needed in Terraform — the API will return errors if the user doesn't have access.
