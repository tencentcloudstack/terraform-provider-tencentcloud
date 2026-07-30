## Context

TencentCloud VOD provides AIGC quota management through four cloud APIs:
- `CreateAigcQuota`: Create a quota for a sub-application (no resource ID returned)
- `DescribeAigcQuotas`: List quotas for a sub-application with filtering by quota type and API token
- `ModifyAigcQuota`: Update the quota limit value
- `DeleteAigcQuota`: Remove a quota

The existing `tencentcloud_vod_aigc_api_token` resource serves as a reference for code patterns (composite ID, service layer retry, async consistency polling). The new resource will be created in the `tencentcloud/services/vod/` directory alongside existing VOD resources.

The cloud SDK (`tencentcloud-sdk-go/tencentcloud/vod/v20180717`) is already vendored and includes all required types.

## Goals / Non-Goals

**Goals:**
- Provide a Terraform resource `tencentcloud_vod_aigc_quota` with full CRUD lifecycle (Create, Read, Update, Delete)
- Support import via composite ID `{sub_app_id}#{quota_type}#{api_token}`
- Handle async consistency (Create may have a brief sync delay before Describe reflects changes)
- Follow existing VOD service patterns (e.g., `tencentcloud_vod_aigc_api_token`)

**Non-Goals:**
- No datasource for listing quotas (this is a resource-only change)
- No batch operations or complex quota management beyond single-resource CRUD

## Decisions

### Decision 1: Composite ID Format

**Choice**: `{sub_app_id}#{quota_type}#{api_token}` using `tccommon.FILED_SP` separator.

**Rationale**: CreateAigcQuota does not return a unique resource ID. The tuple (SubAppId, QuotaType, ApiToken) uniquely identifies a quota. When ApiToken is empty (for Image/Video quotas), the ID will be `{sub_app_id}#{quota_type}#`. This follows the existing pattern used by `tencentcloud_vod_aigc_api_token` (which uses `{sub_app_id}#{api_token}`).

**Alternative considered**: Using a generated UUID — rejected because there's no server-side ID to anchor to, and the tuple is the natural key.

### Decision 2: Schema Field Design (Flat, Not Nested)

**Choice**: Flatten AigcQuotaItem fields directly into the resource schema top level.

**Rationale**: Following rule #13 — "禁止创建该资源列表型数据这一层 schema". The Describe response returns `QuotaSet []*AigcQuotaItem`, but for a single-resource CRUD, we always get 0 or 1 item. We expand the item's fields (QuotaLimit, Usage) into top-level schema fields rather than creating a `quota_set` nested block.

Fields:
| Schema Name | Type | Mode | SDK Source |
|---|---|---|---|
| `sub_app_id` | TypeInt | Required, ForceNew | Request input |
| `quota_type` | TypeString | Required, ForceNew | Request input |
| `quota_limit` | TypeInt | Required | Request input / AigcQuotaItem.QuotaLimit |
| `api_token` | TypeString | Optional, ForceNew | Request input |
| `usage` | TypeInt | Computed | AigcQuotaItem.Usage |

### Decision 3: Read Flow

**Choice**: Call DescribeAigcQuotas with (SubAppId, QuotaType, ApiToken), iterate QuotaSet for the matching item.

**Rationale**: The Describe endpoint returns a paginated list. With the specific filters applied, we expect 0 or 1 results. If 0 results, set `d.SetId("")` to signal resource not found. We also need to pass Limit=100 (the SDK max) to ensure we catch the quota if it's not on the first page.

### Decision 4: Create/Update Consistency Polling

**Choice**: After Create, poll DescribeAigcQuotas until the quota appears (similar to `tencentcloud_vod_aigc_api_token` pattern). No explicit polling needed for Update — the Modify API is synchronous.

**Rationale**: The VOD AIGC API token resource documentation mentions ~30s sync delay after Create. The same pattern applies to quotas. The existing `resource.Retry(tccommon.ReadRetryTimeout, ...)` pattern is proven.

### Decision 5: Service Layer Methods

**Choice**: Add four new methods to `VodService`:
- `CreateVodAigcQuota(ctx, subAppId, quotaType, quotaLimit, apiToken) error`
- `DescribeVodAigcQuotaById(ctx, subAppId, quotaType, apiToken) (*vod.AigcQuotaItem, error)` — returns nil if not found
- `ModifyVodAigcQuota(ctx, subAppId, quotaType, quotaLimit, apiToken) error`
- `DeleteVodAigcQuota(ctx, subAppId, quotaType, apiToken) error`

**Rationale**: Follows existing VOD service conventions. Each method wraps the cloud API call with retry logic and proper error handling.

## Risks / Trade-offs

- **[Risk] Create returns no ID** → Mitigation: Use composite ID tuple. Document the import format clearly in the `.md` file.
- **[Risk] ApiToken is only meaningful for Text quotas but optional in schema** → Mitigation: Accept the field for all quota types; the cloud API handles validation. For Image/Video quotas, the token part of the composite ID will be empty.
- **[Risk] Describe returns a list, not a single item** → Mitigation: Filter by the exact (SubAppId, QuotaType, ApiToken) tuple; if multiple items match (should not happen), treat as not found and let retry handle it.