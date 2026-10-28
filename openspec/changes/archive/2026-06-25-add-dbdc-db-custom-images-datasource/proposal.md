## Why

Terraform users managing DBDC (Database Dedicated Cluster) resources need to discover available custom OS images before creating DB Custom clusters or adding nodes to them. The cloud API `DescribeDBCustomImages` provides this information, but there is currently no Terraform data source to query it. Adding `tencentcloud_dbdc_db_custom_images` enables users to reference valid image IDs in their Terraform configurations without hardcoding them.

## What Changes

- Add new data source `tencentcloud_dbdc_db_custom_images` (RESOURCE_KIND_DATASOURCE) to query available DB Custom system images via the `DescribeDBCustomImages` API
- The data source will expose the full `ImageSet` response as `image_set`, flattening each `DBCustomImage` element's fields (`image_id`, `os_name`, `image_type`, `architecture`) to the schema top level
- Register the new data source in `provider.go` and `provider.md`
- Add unit tests using gomonkey mock pattern (not Terraform acceptance test suite)
- Add `.md` documentation file for the data source

## Capabilities

### New Capabilities
- `dbdc-db-custom-images-datasource`: Data source to query DB Custom available OS images via DescribeDBCustomImages API, exposing image_id, os_name, image_type, and architecture fields

### Modified Capabilities
- (none)

## Impact

- New files in `tencentcloud/services/dbdc/`: data_source_tc_dbdc_db_custom_images.go, data_source_tc_dbdc_db_custom_images_test.go, data_source_tc_dbdc_db_custom_images.md
- Modified files: `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` (add DescribeDBCustomImagesByFilter method), `tencentcloud/provider.go` (add data source registration), `tencentcloud/provider.md` (add documentation entry)
- Cloud API dependency: `DescribeDBCustomImages` from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029`
