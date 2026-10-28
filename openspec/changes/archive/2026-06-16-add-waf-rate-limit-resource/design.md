## Context

TencentCloud WAF provides rate limiting APIs (`CreateRateLimitV2`, `DescribeRateLimitsV2`, `UpdateRateLimitV2`, `DeleteRateLimitsV2`) in the `waf/v20180125` SDK package. The existing WAF service directory (`tencentcloud/services/waf/`) already contains multiple resources and data sources. This change adds a new RESOURCE_KIND_GENERAL resource `tencentcloud_waf_rate_limit` following the established patterns.

The resource manages rate limit rules scoped to a WAF domain. Each rule has a unique `LimitRuleID` returned by the Create API. The composite resource ID will be `domain#limit_rule_id` using `tccommon.FILED_SP` as separator.

## Goals / Non-Goals

**Goals:**
- Provide full CRUD lifecycle management for WAF rate limit rules via Terraform
- Support all configurable parameters exposed by the WAF rate limit APIs (domain, name, priority, status, limit_window, limit_object, limit_strategy, limit_method, limit_paths, limit_headers, limit_header_name, get/post params, ip_location, redirect_info, block_page, object_src, quota_share, paths_option, order)
- Support resource import using composite ID `domain#limit_rule_id`
- Follow existing provider patterns for retry, error handling, and logging

**Non-Goals:**
- Batch management of multiple rate limit rules in a single resource
- Data source for querying rate limit rules (can be added separately)
- Managing WAF domain lifecycle (handled by existing resources)

## Decisions

### 1. Composite ID: `domain#limit_rule_id`
**Rationale**: The rate limit rule is scoped to a domain. Both `domain` and `limit_rule_id` are required for Read/Update/Delete operations. Using `tccommon.FILED_SP` (`#`) as separator follows the established pattern in this provider.

### 2. Read implementation using DescribeRateLimitsV2 with ID filter
**Rationale**: The `DescribeRateLimitsV2` API accepts an `Id` parameter to filter by rule ID. In the Read function, we pass `Domain` and `Id` to retrieve the specific rule. The response returns a `RateLimits` array; we match the first element. The API also has `Offset`/`Limit` for pagination, but since we filter by specific ID, pagination is not needed for the Read path.

### 3. Nested struct types as sub-blocks in Terraform schema
**Rationale**: The API uses several nested struct types:
- `LimitWindow` (second, minute, hour, quota_share) → `limit_window` block (List, MaxItems: 1)
- `LimitMethod` (method, type) → `limit_method` block (List, MaxItems: 1)
- `LimitPath` (path, type) → `limit_paths` block (List, MaxItems: 1)
- `LimitHeader` (key, value, type) → `limit_headers` block (List)
- `LimitHeaderName` (params_name, type) → `limit_header_name` block (List, MaxItems: 1)
- `MatchOption` (params, func, content) → used for `get_params_name`, `get_params_value`, `post_params_name`, `post_params_value`, `ip_location` blocks (List, MaxItems: 1)
- `RedirectInfo` (protocol, domain, url) → `redirect_info` block (List, MaxItems: 1)
- `PathItem` (path, method) → `paths_option` block (List)

### 4. Unit tests with gomonkey mocking
**Rationale**: Per project requirements, new resources use gomonkey-based mocking for unit tests rather than Terraform acceptance test suites. This avoids dependency on real cloud credentials during testing.

### 5. Delete uses array parameter
**Rationale**: `DeleteRateLimitsV2` accepts `LimitRuleIds` as an array of int64. For single resource deletion, we pass a single-element array containing the rule ID.

## Risks / Trade-offs

- [Risk] `DescribeRateLimitsV2` returns a list; if the specific rule ID is not found, the list may be empty → Mitigation: Log the ID before clearing state, return appropriate error for data consistency.
- [Risk] Complex nested schema with many optional fields may be hard to maintain → Mitigation: Follow existing patterns in the provider; each nested type maps directly to an API struct.
- [Risk] The `Order` field in Create/Update is `*int64` but in Describe request it is `*string` (sort order) → Mitigation: These are different parameters with the same name in different APIs. In the Terraform schema, `order` maps to the Create/Update `Order` (int64, execution order). The Describe `Order` (string, sort direction) is only used internally for API calls, not exposed to users.
