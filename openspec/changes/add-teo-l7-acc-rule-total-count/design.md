## Context

The Terraform Provider for TencentCloud currently exposes the `tencentcloud_teo_l7_acc_rule` resource, which manages L7 (Layer 7) access rules for the Edge Optimization (TEO) service. The resource interacts with the TencentCloud API using the `DescribeL7AccRules` endpoint.

Currently, the API response includes a `TotalCount` field that indicates the total number of rules available, but this field is not exposed in the Terraform resource schema. This means users cannot access this potentially useful information through the Terraform provider.

The resource implementation follows standard Terraform provider patterns:
- Located at `/repo/tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go`
- Uses the Terraform Plugin SDK v2
- Calls TencentCloud SDK (`teov20220901` package) via the service layer
- Implements standard CRUD operations (Create, Read, Update, Delete)
- Uses `helper.Retry()` for consistency and `defer tccommon.LogElapsed()` for logging

## Goals / Non-Goals

**Goals:**
- Add the `total_count` field to the `tencentcloud_teo_l7_acc_rule` resource schema
- Ensure the field is correctly populated from the API response during read operations
- Maintain full backward compatibility with existing Terraform configurations
- Update resource documentation to reflect the new field

**Non-Goals:**
- Modify any existing resource behavior or schema (only add new field)
- Change any API interactions or calls
- Modify any resource lifecycle operations (Create, Update, Delete)
- Add pagination or filtering functionality to the resource

## Decisions

### Field Type and Schema

**Decision:** Add `total_count` as a `TypeInt` with `Computed: true` and no `Optional` flag.

**Rationale:**
- `TypeInt` matches the API response type (`*int64`)
- `Computed: true` indicates it's read-only and populated by the provider
- No `Optional` flag ensures it's not settable by users (comes from API only)
- This follows the standard pattern for API-derived fields in the provider

**Alternative considered:** Using `TypeString` - rejected because the API returns an integer.

### Implementation Location

**Decision:** Add the field to the resource's `Schema` map and populate it in `resourceTencentCloudTeoL7AccRuleRead`.

**Rationale:**
- The field is part of the resource's state and should be defined in the Schema
- The Read operation is where API data is retrieved and populated
- This follows the standard pattern for all computed fields in the provider
- The API call `DescribeTeoL7AccRuleById` already returns the response with `TotalCount`

**Alternative considered:** Adding to separate data source - rejected because the data is already available in the resource read operation.

### API Response Parsing

**Decision:** Use `d.Set("total_count", *respData.TotalCount)` after the existing schema population code.

**Rationale:**
- The `TotalCount` field is already available in the `respData` object
- This approach is consistent with how other fields are set
- Minimal code changes required (one line addition)
- No need to modify the service layer (`DescribeTeoL7AccRuleById` already returns the full response)

**Alternative considered:** Modifying service layer to extract TotalCount separately - rejected as it adds unnecessary complexity.

## Risks / Trade-offs

**Risk:** The API might return `nil` for `TotalCount` in certain scenarios (e.g., empty rule list).

→ **Mitigation:** The SDK returns `*int64` (pointer type), so we need to check for nil before dereferencing. Use helper function or explicit nil check:
```go
if respData.TotalCount != nil {
    _ = d.Set("total_count", *respData.TotalCount)
}
```

**Risk:** Existing acceptance tests might fail if they expect the resource state to not include the new field.

→ **Mitigation:** Since this is a `Computed` field, it will be automatically populated during state refresh and will not cause test failures. Existing tests will simply see an additional field in the state, which is acceptable.

**Trade-off:** Adding the field increases the resource state size, but the field value is a single integer, so the impact is negligible.

## Migration Plan

**Deployment Steps:**
1. Update the resource schema in `resource_tc_teo_l7_acc_rule.go`
2. Update the Read function to populate the new field
3. Update the documentation file `website/docs/r/teo_l7_acc_rule.html.markdown`
4. Update or add acceptance tests to verify the field is populated
5. Run existing acceptance tests to ensure backward compatibility

**Rollback Strategy:**
- Simply remove the new field from the schema and Read function
- Since this is a `Computed` field, removing it will not corrupt existing Terraform states (users just won't see the field anymore)
- No state migration is required because the field is read-only

**State Migration:** Not required. The field is `Computed` only, so it doesn't exist in user configurations and will be populated fresh on each `terraform apply` or `terraform refresh`.

## Open Questions

None. The implementation is straightforward with no architectural decisions or dependencies that require further investigation.

The only remaining question is whether to add assertions to existing acceptance tests or create new tests for the `total_count` field. This will be determined during the implementation phase.
