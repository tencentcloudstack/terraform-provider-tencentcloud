## Why

Terraform Provider for TencentCloud currently has no datasource for DB Custom regions (dbdc product). Users need to query the list of supported DB Custom regions in their Terraform configurations to reference region availability and sale status when planning DB Custom clusters and nodes. Adding `tencentcloud_dbdc_db_custom_regions` datasource enables users to read and query DB Custom region attributes within Terraform workflows.

## What Changes

- Add new datasource `tencentcloud_dbdc_db_custom_regions` (RESOURCE_KIND_DATASOURCE) that calls `DescribeDBCustomRegions` API to query the DB Custom supported region list
- Add datasource schema with computed output field: `region_set` (list of region details, flattened fields: `region`, `region_state`)
- Add service layer method `DescribeDBCustomRegions` in the existing `service_tencentcloud_dbdc.go` file
- Register the datasource in `provider.go` and `provider.md`
- Add documentation file `data_source_tc_dbdc_db_custom_regions.md`
- Add unit test file `data_source_tc_dbdc_db_custom_regions_test.go`

## Capabilities

### New Capabilities
- `dbdc-db-custom-regions-datasource`: Datasource to query DB Custom supported region list from dbdc product. The `DescribeDBCustomRegions` API has no input parameters and returns a `RegionSet` list where each `RegionInfo` element contains `Region` (region name) and `RegionState` (sale status: SELL / SOLD_OUT). Returns region details flattened into the `region_set` list.

### Modified Capabilities
<!-- No existing capabilities are being modified -->

## Impact

- New files: `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_regions.go`, `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_regions.md`, `tencentcloud/services/dbdc/data_source_tc_dbdc_db_custom_regions_test.go`
- Modified files: `tencentcloud/services/dbdc/service_tencentcloud_dbdc.go` (add `DescribeDBCustomRegions` service method), `tencentcloud/provider.go` (register datasource), `tencentcloud/provider.md` (add datasource entry)
- API dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029.DescribeDBCustomRegions`
