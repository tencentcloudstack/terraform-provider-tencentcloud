## 1. Service Layer

- [x] 1.1 Add `DescribeTeoMultiPathGatewayRegionByFilter` method to `tencentcloud/services/teo/service_tencentcloud_teo.go`, wrapping `DescribeMultiPathGatewayRegions` API call with retry handling using `tccommon.ReadRetryTimeout`

## 2. Data Source Implementation

- [x] 2.1 Create `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway_region.go` with schema definition (zone_id Required, gateway_regions Computed with nested region_id/cn_name/en_name, result_output_file Optional) and Read function
- [x] 2.2 Create `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway_region.md` documentation file with description, example usage

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_multi_path_gateway_region` data source in `tencentcloud/provider.go`
- [x] 3.2 Register `tencentcloud_teo_multi_path_gateway_region` data source in `tencentcloud/provider.md`

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway_region_test.go` with gomonkey mock tests for successful read and API error scenarios
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` to verify correctness
