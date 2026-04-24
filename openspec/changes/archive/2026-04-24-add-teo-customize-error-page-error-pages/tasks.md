## 1. Schema Definition

- [x] 1.1 Add `error_pages` computed attribute (TypeList of nested Resource) to the `tencentcloud_teo_customize_error_page` resource schema in `tencentcloud/services/teo/resource_tc_teo_customize_error_page.go`, with nested fields: `page_id` (string), `name` (string), `content_type` (string), `description` (string), `content` (string), `references` (TypeList of TypeString)

## 2. CRUD Function Update

- [x] 2.1 Update `resourceTencentCloudTeoCustomizeErrorPageRead` to flatten the `error_pages` data from the `DescribeCustomErrorPages` API response, including mapping the `References` field (list of `ErrorPageReference` with `BusinessId`) into the `references` sub-field

## 3. Unit Tests

- [x] 3.1 Add unit tests in `tencentcloud/services/teo/resource_tc_teo_customize_error_page_test.go` to verify the `error_pages` computed attribute is correctly populated, including the `references` sub-field

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/teo/resource_tc_teo_customize_error_page.md` to document the new `error_pages` computed attribute with its sub-fields
