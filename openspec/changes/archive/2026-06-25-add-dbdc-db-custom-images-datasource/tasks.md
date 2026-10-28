## 1. Service Layer

- [x] 1.1 Add `DescribeDBCustomImagesByFilter` method to `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` with automatic pagination (Limit=100), retry logic, ratelimit check, and NonRetryableError for nil/empty response

## 2. Data Source Implementation

- [x] 2.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_images.go` with schema definition (result_output_file + image_set with flattened fields: image_id, os_name, image_type, architecture)
- [x] 2.2 Implement Read function with defer logging, paramMap construction, retry-wrapped service call, nil-check field flattening, d.SetId(helper.BuildToken()), and result_output_file handling

## 3. Provider Registration

- [x] 3.1 Add `"tencentcloud_dbdc_db_custom_images": dbdc.DataSourceTencentCloudDbdcDbCustomImages()` entry in `tencentcloud/provider.go` DataSourcesMap
- [x] 3.2 Add `tencentcloud_dbdc_db_custom_images` entry in `tencentcloud/provider.md`

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_images_test.go` with gomonkey mock tests: basic read, nil-field handling, empty response error, schema structure validation
- [x] 4.2 Run unit tests with `go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomImagesDS" -v -count=1 -gcflags="all=-l"` and verify all pass

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_images.md` following gendoc format: one-line description mentioning DBDC, Example Usage with HCL examples, no Argument/Attribute Reference sections
