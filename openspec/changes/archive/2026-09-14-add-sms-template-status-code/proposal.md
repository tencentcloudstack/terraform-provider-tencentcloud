# Change: Add status_code field to tencentcloud_sms_template resource

## Why

The `tencentcloud_sms_template` resource currently supports creating, updating, reading, and deleting SMS templates, but does not expose the template's audit/review status as a Terraform attribute. The underlying `DescribeSmsTemplateList` API already returns a `StatusCode` field (via `DescribeTemplateListStatus.StatusCode`) indicating whether a template has been approved, is pending review, or was rejected. Without this field, users cannot determine the template's lifecycle state through Terraform and must check the TencentCloud console manually.

Adding this computed attribute enables users to:
- View the template's current audit status as part of their Terraform state
- Use Terraform outputs or `depends_on` to react to status changes
- Have complete visibility into the resource's cloud-side state

## What Changes

- **tencentcloud_sms_template** resource: Add a new `status_code` computed field to the resource schema
  - **Type**: `TypeInt`
  - **Attributes**: `Computed: true`, `Optional: false`, `Required: false`, `ForceNew: false`
  - **Description**: "Template status. 0: approved and effective, 1: pending review, 2: approved pending activation, -1: review failed or rejected"
- Update the `Read` function (`resourceTencentCloudSmsTemplateRead`) to read `StatusCode` from the `DescribeTemplateListStatus` struct and set it in the Terraform state
- This is a purely additive change - no modifications to existing fields, no new API calls, no changes to Create/Update/Delete operations

## Capabilities

### New Capabilities
- `sms-template-status-code`: Expose the `StatusCode` field from `DescribeSmsTemplateList` API response as a computed Terraform attribute on the `tencentcloud_sms_template` resource

### Modified Capabilities
None. No existing spec requirements are being changed.

## Impact

### Affected Code
- `tencentcloud/services/sms/resource_tc_sms_template.go` - Schema definition (add `status_code` field) and Read function (read and set field value)
- `tencentcloud/services/sms/resource_tc_sms_template_test.go` - Add test coverage for the new field (optional, follow existing test patterns)

### Backward Compatibility
- ✅ **Fully backward compatible** - New field is `Computed: true` (read-only output only), existing configurations are unaffected
- ✅ **No schema breaking changes** - All existing fields remain unchanged
- ✅ **No API changes** - The `StatusCode` field is already returned by the existing `DescribeSmsTemplateList` call in the Read path
- ✅ **No dependency changes** - Uses the existing SDK version, no new dependencies required

### API Compatibility
- ✅ `DescribeSmsTemplateList` API response already includes `StatusCode` in the `DescribeTemplateListStatus` struct
- ✅ The service layer function `DescribeSmsTemplate` in `service_tencentcloud_sms.go` already returns a `*sms.DescribeTemplateListStatus` which contains `StatusCode`
- ✅ Existing `DescribeTemplateListStatus` struct has `StatusCode *int64` field available