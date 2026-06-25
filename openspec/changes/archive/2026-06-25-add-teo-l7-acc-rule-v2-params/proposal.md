## Why

The `tencentcloud_teo_l7_acc_rule_v2` resource needs to be fully implemented as a RESOURCE_KIND_GENERAL resource for the TEO (TencentCloud EdgeOne) product. This resource manages L7 acceleration rules with full CRUD lifecycle operations (CreateL7AccRules, DescribeL7AccRules, ModifyL7AccRule, DeleteL7AccRules), enabling users to create, read, update, and delete L7 acceleration rules within a TEO zone.

## What Changes

- Add the `tencentcloud_teo_l7_acc_rule_v2` resource with full CRUD support:
  - `zone_id` (Required, ForceNew): Zone ID for the TEO site
  - `rules` (Required): Rule content list containing rule engine items with status, rule_name, description, and branches
  - `rule_id` (Computed): Rule ID returned after creation, used for update/delete/read operations
  - `status` (Optional): Rule status (enable/disable)
  - `rule_name` (Optional): Rule name with 255 character limit
  - `description` (Optional): Rule annotations, supports multiple entries
  - `branches` (Optional): Sub-rule branches for rule conditions and actions
  - `rule_priority` (Computed): Rule priority, output only
- Register the resource in `provider.go` and `provider.md`
- Add service layer function `DescribeTeoL7AccRuleById` for reading rule details
- Add documentation and unit tests

## Capabilities

### New Capabilities
- `teo-l7-acc-rule-v2-resource`: Full CRUD lifecycle management for TEO L7 acceleration rules via CreateL7AccRules, ModifyL7AccRule, DeleteL7AccRules, and DescribeL7AccRules APIs

### Modified Capabilities

## Impact

- **Code**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.go`, `tencentcloud/services/teo/service_tencentcloud_teo.go`, `tencentcloud/provider.go`, `tencentcloud/provider.md`
- **APIs**: Uses `teo/v20220901` package APIs: `CreateL7AccRules`, `ModifyL7AccRule`, `DeleteL7AccRules`, `DescribeL7AccRules`
- **Dependencies**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901` (already in vendor)
- **Documentation**: `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_v2.md`, `website/docs/r/teo_l7_acc_rule_v2.html.markdown`
