## 1. Schema Definition

- [x] 1.1 Add `references` computed attribute (TypeList of TypeString) at the top level of the `tencentcloud_teo_customize_error_page` resource schema in `tencentcloud/services/teo/resource_tc_teo_customize_error_page.go`

## 2. CRUD Function Update

- [x] 2.1 Update `resourceTencentCloudTeoCustomizeErrorPageRead` to flatten the `References` field (list of `ErrorPageReference` with `BusinessId`) from the API response into the top-level `references` attribute

## 3. Unit Tests

- [x] 3.1 Add unit tests in `tencentcloud/services/teo/resource_tc_teo_customize_error_page_test.go` to verify the `references` computed attribute is correctly populated, including empty references scenario
