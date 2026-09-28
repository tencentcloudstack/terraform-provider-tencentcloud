## Context

The `tencentcloud_clb_log_topic` resource (`tencentcloud/services/clb/resource_tc_clb_log_topic.go`) manages CLB log topics. It currently exposes `log_set_id`, `topic_name`, `status`, `tags`, and the computed `create_time`.

The resource is unusual because it spans two Tencent Cloud SDK packages:
- **CLB SDK** (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317`) is used for **Create** via `ClbService.CreateTopic`, which wraps `clb.CreateTopicRequest`.
- **CLS SDK** (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016`) is used for **Modify** (`cls.NewModifyTopicRequest()`) and **Read** (via `svccls.ClsService.DescribeClsTopicById`, which calls `DescribeTopics` and returns `*cls.TopicInfo`).

**Current state — API behavior analysis (verified against vendored SDK):**

| API | SDK package | `Period` in request | `Period` in response | Go type |
|-----|-------------|---------------------|----------------------|---------|
| `CreateTopic` | `clb/v20180317` | Yes (`CreateTopicRequest.Period`) | N/A | `*uint64` |
| `ModifyTopic` | `cls/v20201016` | Yes (`ModifyTopicRequest.Period`) | N/A | `*int64` |
| `DescribeTopics` | `cls/v20201016` | N/A | Yes (`TopicInfo.Period`) | `*int64` |

**SDK constraint — type mismatch:** `Period` is `*uint64` on the CLB `CreateTopicRequest` but `*int64` on the CLS `ModifyTopicRequest` and `TopicInfo`. The Terraform schema will use `TypeInt`; values must be cast to the correct SDK type at each call site (`uint64(period)` for Create, `int64(period)` for Modify).

**API semantics:** `Period` is the log storage lifecycle in days. Standard storage supports 1-3600 days; 3640 means permanent retention. The API applies a default of 30 days when `Period` is not supplied. `Period` is accepted by both Create and Modify, so it is **updatable in-place** (not `ForceNew`).

## Goals / Non-Goals

**Goals:**
- Add an optional `period` (TypeInt) parameter to the `tencentcloud_clb_log_topic` resource schema that is not `ForceNew` (updatable in-place).
- Pass `period` to the `CreateTopic` API (CLB SDK) as `Period *uint64`.
- Pass `period` to the `ModifyTopic` API (CLS SDK) as `Period *int64` when it changes during Update.
- Read `period` back from the `DescribeTopics` API response (`TopicInfo.Period *int64`) so state refresh and import populate the field.
- Update the resource `.md` example with `period` usage.
- Add mock-based unit tests (gomonkey) for Create, Read, and Update of the `period` parameter.
- Maintain full backward compatibility (optional field, defaults unset).

**Non-Goals:**
- Adding validation constraints (range 1-3600/3640) beyond relying on the API as the final authority — the API will reject invalid values and surface a clear error.
- Changing the existing `status`, `tags`, or `create_time` behavior.
- Adding `period` to any CLB datasource (out of scope).
- Splitting the resource into CLB-/CLS-specific resources.

## Decisions

### Decision 1: `period` is updatable (not `ForceNew`)
**Rationale:** The `ModifyTopic` API (CLS SDK) accepts `Period`, so the value can be changed in-place after creation. Marking it `ForceNew` would unnecessarily destroy and recreate log topics (and their data) on every retention change.

### Decision 2: `period` is `Optional` and `Computed`
**Rationale:** The user can explicitly configure retention; marking it `Computed` as well lets the provider surface the API-applied default (30 days) in state after Read when the user leaves it unset, following the same pattern as the existing `status` field on this resource. When unset, the provider does not send `Period` on Create and the API applies its own default.

### Decision 3: Service-layer `CreateTopic` receives `period` via the existing `params` map
**Rationale:** `ClbService.CreateTopic(ctx, params map[string]interface{})` already accepts a generic params map (used for `topic_name`, `partition_count`, `tags`). Adding `params["period"]` avoids changing the function signature and is consistent with the existing pattern. Inside the service function, `Period` is set on `clb.CreateTopicRequest` as `*uint64` only when present.

### Decision 4: Update flow builds a single `ModifyTopic` request when any updatable field changes
**Rationale:** `ModifyTopic` (CLS SDK) accepts `Status`, `Tags`, and `Period` in the same request, so the Update function builds one `cls.NewModifyTopicRequest()` whenever any of `status`, `tags`, or `period` changes, populates only the changed fields, and issues a single retried `ModifyTopic` call for all of them. This reuses one call path for all updatable fields (consistent with the previous `status`/`tags` handling) and avoids redundant API calls when multiple fields change at once.

### Decision 5: Type casting at call sites
**Rationale:** Because the CLB Create SDK uses `*uint64` and the CLS Modify/Describe SDK uses `*int64`, the resource code casts the schema `int` value to `uint64` before Create and to `int64` before Modify/Read-set. This is localized to each call site and avoids introducing a shared type abstraction for a single field.

## Risks / Trade-offs

- **[Risk] Type mismatch between CLB and CLS SDKs could cause silent truncation** → Mitigation: explicit `uint64(period)` for Create and `int64(period)` for Modify; the value range (1-3640) fits comfortably in both types, so no overflow. Read sets `*TopicInfo.Period` (`*int64`) into the `TypeInt` schema directly.
- **[Risk] Drift between user-configured `period` and API default (30) when unset** → Mitigation: `period` is `Optional` (not `Computed`), so Terraform does not attempt to manage the value when the user leaves it blank. Users who want to pin retention simply set `period` explicitly.
- **[Risk] Update path must still call `ModifyTopic` even if only `period` changes** → Mitigation: restructure the Update function so a single `ModifyTopic` request is built when any of `status`, `tags`, or `period` changes, preserving existing retry behavior and the `cls.NewModifyTopicRequest()` pattern.
