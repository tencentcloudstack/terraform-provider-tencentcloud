## 1. Service Layer

- [x] 1.1 Add `DescribeTeoMultiPathGatewayOriginAclByFilter` method to `tencentcloud/services/teo/service_tencentcloud_teo.go` that calls `DescribeMultiPathGatewayOriginACL` API with paramMap containing ZoneId and GatewayId, returns `*teov20220901.DescribeMultiPathGatewayOriginACLResponseParams`

## 2. Data Source Schema and Read Function

- [x] 2.1 Create `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway_origin_acl.go` with `DataSourceTencentCloudTeoMultiPathGatewayOriginAcl()` function defining the schema (zone_id Required, gateway_id Required, result_output_file Optional, multi_path_gateway_origin_acl_info Computed with nested blocks)
- [x] 2.2 Implement `dataSourceTencentCloudTeoMultiPathGatewayOriginAclRead()` function with nil-safe nested response mapping, retry logic using `tccommon.ReadRetryTimeout`, composite ID using `tccommon.FILED_SP`, and result output file support

## 3. Provider Registration

- [x] 3.1 Add data source entry `"tencentcloud_teo_multi_path_gateway_origin_acl": teo.DataSourceTencentCloudTeoMultiPathGatewayOriginAcl()` to `tencentcloud/provider.go` dataSources map
- [x] 3.2 Add data source entry to `tencentcloud/provider.md` data source list

## 4. Documentation

- [x] 4.1 Create `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway_origin_acl.md` with description, example usage, and import section (following TEO doc pattern, no Argument/Attribute Reference sections)

## 5. Unit Tests

- [x] 5.1 Create `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway_origin_acl_test.go` with unit tests using gomonkey mock approach to test the Read function logic
- [x] 5.2 Run unit tests with `go test -gcflags=all=-l` to verify all test cases pass

## 6. Verification

- [x] 6.1 Verify all new and modified files compile correctly and follow the project patterns
