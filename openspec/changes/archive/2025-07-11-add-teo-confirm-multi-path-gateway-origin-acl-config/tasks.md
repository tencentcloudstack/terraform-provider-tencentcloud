## 1. Resource Implementation

- [x] 1.1 Create `tencentcloud/services/teo/resource_tc_teo_confirm_multi_path_gateway_origin_acl_config.go` with schema definition including zone_id (Required, ForceNew), gateway_id (Required, ForceNew), origin_acl_version (Optional), and multi_path_gateway_origin_acl_info (Computed) with nested structures for current and next ACL info
- [x] 1.2 Implement Create function: set composite ID (zone_id + gateway_id joined by FILED_SP), delegate to Update
- [x] 1.3 Implement Read function: parse composite ID, call DescribeMultiPathGatewayOriginACL via helper.Retry with ReadRetryTimeout, flatten MultiPathGatewayOriginACLInfo response into computed output fields, handle resource-not-found by clearing ID
- [x] 1.4 Implement Update function: when origin_acl_version is specified/changed, call ConfirmMultiPathGatewayOriginACL with ZoneId, GatewayId, OriginACLVersion, then call Read to refresh state
- [x] 1.5 Implement Delete function: no-op (CONFIG resource)
- [x] 1.6 Implement helper functions for flattening API response structures (MultiPathGatewayOriginACLInfo, MultiPathGatewayCurrentOriginACL, MultiPathGatewayNextOriginACL, Addresses) into Terraform state

## 2. Service Layer

- [x] 2.1 Add `DescribeTeoConfirmMultiPathGatewayOriginAclById` method to TeoService in `tencentcloud/services/teo/service_tencentcloud_teo.go` to call DescribeMultiPathGatewayOriginACL with ZoneId and GatewayId
- [x] 2.2 Add retry logic wrapping the DescribeMultiPathGatewayOriginACL API call with tccommon.ReadRetryTimeout, using tccommon.RetryError for error wrapping

## 3. Provider Registration

- [x] 3.1 Register `tencentcloud_teo_confirm_multi_path_gateway_origin_acl` resource in `tencentcloud/provider.go`
- [x] 3.2 Add resource entry in `tencentcloud/provider.md`

## 4. Unit Tests

- [x] 4.1 Create `tencentcloud/services/teo/resource_tc_teo_confirm_multi_path_gateway_origin_acl_config_test.go` with gomonkey-based unit tests for Read and Update operations
- [x] 4.2 Run unit tests with `go test -gcflags=all=-l` to verify they pass

## 5. Documentation

- [x] 5.1 Create `tencentcloud/services/teo/resource_tc_teo_confirm_multi_path_gateway_origin_acl_config.md` with one-line description, Example Usage section (with zone_id and gateway_id), and Import section (since this is RESOURCE_KIND_CONFIG)
