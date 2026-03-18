# Change: Add Tags support to tencentcloud_cam_policy resource

## Why
The `tencentcloud_cam_policy` resource currently does not support tags management, while the TencentCloud CAM CreatePolicy and GetPolicy APIs both support the Tags parameter. Users need the ability to add, update, and manage tags for CAM policies through Terraform to enable better resource organization, cost allocation, and access control.

## What Changes
- Add `tags` field to the `tencentcloud_cam_policy` resource schema
- Implement tags creation during policy creation using CreatePolicy API
- Implement tags reading during policy read using GetPolicy API  
- Implement tags update during policy update using the unified ModifyTags function
- Add tags handling in resource deletion to follow Terraform best practices

## Impact
- Affected specs: `cam-policy-resource`
- Affected code:
  - `tencentcloud/services/cam/resource_tc_cam_policy.go` - Add schema field and CRUD operations
  - `tencentcloud/services/cam/resource_tc_cam_policy_test.go` - Add test cases for tags
- No breaking changes - the `tags` field is optional
- Full backward compatibility maintained
