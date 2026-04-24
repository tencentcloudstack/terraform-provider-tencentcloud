## Why

The `tencentcloud_teo_customize_error_page` resource currently does not expose the `error_pages` computed attribute from the `DescribeCustomErrorPages` API response. This attribute contains the list of custom error pages including the `References` field (which indicates which business rules reference the error page), which is not available through any existing schema field. Users need visibility into which business rules reference each error page.

## What Changes

- Add a new computed attribute `error_pages` to the `tencentcloud_teo_customize_error_page` resource schema
- The `error_pages` attribute will be a TypeList of objects, each containing: `page_id`, `name`, `content_type`, `description`, `content`, and `references` (list of `business_id` strings)
- Update the `resourceTencentCloudTeoCustomizeErrorPageRead` function to populate the `error_pages` attribute from the `DescribeCustomErrorPages` API response

## Capabilities

### New Capabilities
- `teo-customize-error-page-error-pages`: Expose the `error_pages` computed attribute on the TEO customize error page resource, including the `References` sub-field from the API response

### Modified Capabilities

## Impact

- `tencentcloud/services/teo/resource_tc_teo_customize_error_page.go` - Add `error_pages` schema field and populate it in the read function
- `tencentcloud/services/teo/resource_tc_teo_customize_error_page_test.go` - Add unit tests for the new attribute
- `tencentcloud/services/teo/resource_tc_teo_customize_error_page.md` - Update documentation with the new attribute
