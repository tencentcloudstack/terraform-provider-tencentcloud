## Why

The `tencentcloud_teo_customize_error_page` resource currently does not expose the `References` field from the `DescribeCustomErrorPages` API response. This field indicates which business rules reference the error page, which is useful for users to understand error page dependencies before modifying or deleting them.

## What Changes

- Add a new computed attribute `references` (TypeList of TypeString) at the top level of the `tencentcloud_teo_customize_error_page` resource schema
- The `references` attribute contains a list of `BusinessId` strings from the `ErrorPageReference` objects in the API response
- Update the `resourceTencentCloudTeoCustomizeErrorPageRead` function to populate the `references` attribute from the `DescribeCustomErrorPages` API response

## Capabilities

### New Capabilities
- `teo-customize-error-page-error-pages`: Expose the `references` computed attribute on the TEO customize error page resource, containing the list of business IDs that reference this error page

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_customize_error_page.go` - Add `references` schema field and populate it in the read function
- `tencentcloud/services/teo/resource_tc_teo_customize_error_page_test.go` - Add unit tests for the new attribute
