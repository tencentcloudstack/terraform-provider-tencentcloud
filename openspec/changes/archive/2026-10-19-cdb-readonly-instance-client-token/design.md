## Context

The `tencentcloud_mysql_readonly_instance` resource creates a CDB read-only instance by calling `CreateDBInstance` (包年包月, prepaid) or `CreateDBInstanceHour` (按量计费, postpaid) inside two helper functions `mysqlCreateReadonlyInstancePayByMonth` and `mysqlCreateReadonlyInstancePayByUse`. Currently these helpers call the cloud API directly without a `ClientToken` and without a retry wrapper. If the cloud API responds with an empty `InstanceIds` list (transient backend behavior), the create fails immediately.

The master instance resource `tencentcloud_mysql_instance` already implements a more robust pattern in `mysqlCreateInstancePayByMonth` / `mysqlCreateInstancePayByUse`:
1. Generate a client token via `helper.BuildToken()`.
2. Set `request.ClientToken`.
3. Wrap the API call in `resource.Retry`.
4. Inside the retry, if `r.Response.InstanceIds` is empty AND client token was provided, return `resource.RetryableError(...)` so the call is retried. The idempotent client token guarantees the backend does not create a duplicate instance on retry.

Both `CreateDBInstanceRequest` and `CreateDBInstanceHourRequest` in the vendored SDK (`cdb/v20170320/models.go`) already expose the optional `ClientToken *string` field, so no SDK/vendor change is required.

## Goals / Non-Goals

**Goals:**
- Align the readonly instance create path with the master instance client-token polling mechanism.
- On a transient empty-`InstanceIds` response, retry idempotently using `ClientToken` instead of failing outright.
- Keep the change backward compatible: no schema changes, no state format changes, no behavior change on the success path.

**Non-Goals:**
- Do not modify the Read, Update, or Delete flows.
- Do not modify the master instance resource (`tencentcloud_mysql_instance`).
- Do not add `ClientToken` as a user-facing schema field; it is an internal implementation detail.
- Do not change the existing post-create status-polling retry (waiting for `MYSQL_STATUS_RUNNING`); only the API call that performs the create is wrapped.

## Decisions

### Decision 1: Reuse `helper.BuildToken()` to generate the client token

**Rationale:** This is exactly what the master instance uses. Reusing the same helper keeps consistency and avoids introducing a new token-generation utility.

**Alternative considered:** Hardcode a UUID via `helper.BuildUUID()`. Rejected because the master instance uses `BuildToken()` and the goal is to align with master logic.

### Decision 2: Wrap the create API call in `resource.Retry(2*tccommon.WriteRetryTimeout, ...)`

**Rationale:** Matches the master instance retry duration (`2*tccommon.WriteRetryTimeout`) and structure. The retry closure shadows the response variable and assigns the successful response to the outer variable, mirroring `mysqlCreateInstancePayByMonth`.

**Alternative considered:** Use the existing non-retried direct call. Rejected — it does not handle the empty-`InstanceIds` case.

### Decision 3: Retry condition — "empty InstanceIds AND client token provided"

**Rationale:** Directly mirrors master logic: `if len(r.Response.InstanceIds) < 1 && clientToken != "" { return resource.RetryableError(...) }`. The client token makes the retry idempotent, so retrying is safe. On a genuine API error (`inErr != nil`), the error is returned via `tccommon.RetryError(inErr)` (PayByMonth) / `tccommon.RetryError(inErr)` (PayByUse) to allow retry on retryable SDK errors, consistent with master.

**Alternative considered:** Always `NonRetryableError` on API error. Rejected — master uses `tccommon.RetryError` which respects retryable error codes.

### Decision 4: Apply the same pattern to both PayByMonth and PayByUse paths

**Rationale:** Both `CreateDBInstance` and `CreateDBInstanceHour` support `ClientToken`, and the master instance applies the mechanism to both. The readonly instance also has both paths, so both get aligned.

### Decision 5: Unit tests use gomonkey to mock the cloud API

**Rationale:** Per the project's code generation requirements, new terraform resources' tests should use mock (gomonkey) instead of the terraform test suite. Since this change modifies an existing resource's create path, the new client-token behavior will be covered by mock-based tests that exercise the retry-then-success and retry-on-empty scenarios without real cloud API calls.

## Risks / Trade-offs

- [Risk: retry could mask a persistent backend failure and delay the create error] → Mitigation: the retry is bounded by `2*tccommon.WriteRetryTimeout`, matching master; after exhaustion the underlying error surfaces.
- [Risk: client token collisions] → Mitigation: `helper.BuildToken()` generates a 16-byte random base64 token; collisions are negligible, and this is the same mechanism already used by master.
- [Risk: behavior divergence between prepaid/postpaid] → Mitigation: apply identical logic to both paths so they stay consistent.
- [Trade-off: slight increase in code complexity] → Acceptable, as it directly mirrors the proven master instance pattern.