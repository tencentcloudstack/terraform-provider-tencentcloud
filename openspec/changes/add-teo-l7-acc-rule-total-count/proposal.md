## Why

The `DescribeL7AccRules` API from TencentCloud Edge Optimization (TEO) service returns a `TotalCount` field in the response, which indicates the total number of rules available. Currently, this field is not exposed in the Terraform provider's `tencentcloud_teo_l7_acc_rule` resource schema. Exposing this field allows users to access the total count of L7 access rules, which is useful for:
- Displaying rule statistics in dashboards
- Implementing pagination-aware data retrieval
- Providing better visibility into the total number of configured rules

This is a non-breaking enhancement that adds read-only access to existing API data.

## What Changes

- Add `total_count` field to the `tencentcloud_teo_l7_acc_rule` resource schema
  - Type: `Int`
  - Computed: `true` (read-only)
  - Description: "Total number of L7 access rules"
- Update the `resourceTencentCloudTeoL7AccRuleRead` function to populate this field from the API response
- The `DescribeL7AccRules` API already returns `TotalCount`, so no backend changes required
- This is a schema-only addition that does not affect resource CRUD operations

## Capabilities

### New Capabilities
- `teo-l7-rule-total-count`: Expose the TotalCount field from DescribeL7AccRules API response in the tencentcloud_teo_l7_acc_rule resource schema

### Modified Capabilities

## Impact

- **Affected Files:**
  - `/repo/tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go` - Add TotalCount field to schema and read logic
  - `/repo/website/docs/r/teo_l7_acc_rule.html.markdown` - Add documentation for the new field

- **API Impact:** None - the DescribeL7AccRules API already returns TotalCount

- **Backward Compatibility:** Fully backward compatible - this is a computed field only, no required or optional fields added

- **Dependencies:** No new dependencies required

- **Testing Impact:**
  - Update acceptance tests to verify TotalCount field is correctly populated
  - No changes to existing test logic required (just add assertions for new field)
