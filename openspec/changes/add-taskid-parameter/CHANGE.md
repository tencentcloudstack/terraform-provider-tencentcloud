# Change: add-taskid-parameter

## Summary
Added support for the `task_id` field in the `tencentcloud_teo_l7_acc_rule` resource to read the TaskId returned by the ImportZoneConfig API.

## Impact

### Code Changes
- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule.go`
  - Added `task_id` field to resource schema (Computed + Optional)
  - Updated Update function to save TaskId from ImportZoneConfig API response to state

- `tencentcloud/services/teo/resource_tc_teo_l7_acc_rule_test.go`
  - Added verification of `task_id` field in acceptance tests

### Behavior Changes
- Users can now read the `task_id` field from the `tencentcloud_teo_l7_acc_rule` resource
- The `task_id` field is populated automatically when creating or updating the resource
- No impact on existing configurations due to backward compatibility (Computed + Optional)

## Migration Notes
No migration required. The `task_id` field is optional and computed, so existing resources will continue to work without changes.

## Testing
- Added acceptance test verification for `task_id` field
- Field is set to Computed + Optional to ensure backward compatibility
