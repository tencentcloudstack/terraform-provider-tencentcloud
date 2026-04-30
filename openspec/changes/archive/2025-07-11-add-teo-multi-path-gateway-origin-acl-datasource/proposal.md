## Why

TEO (EdgeOne) multi-path gateway supports source origin ACL protection, but there is currently no Terraform data source to query the origin ACL details of a multi-path gateway. Users need to read the multi-path gateway origin ACL information (including current and next version ACL details, IPv4/IPv6 CIDR lists, version info, and update confirmation status) via Terraform to integrate with their infrastructure-as-code workflows.

## What Changes

- Add a new data source `tencentcloud_teo_multi_path_gateway_origin_acl` to query multi-path gateway origin ACL details using the `DescribeMultiPathGatewayOriginACL` API
- The data source accepts `zone_id` and `gateway_id` as required input parameters
- The data source outputs `multi_path_gateway_origin_acl_info` containing:
  - `multi_path_gateway_current_origin_acl`: Current effective origin ACL (with `entire_addresses`, `version`, `is_planed`)
  - `multi_path_gateway_next_origin_acl`: Next version origin ACL (with `entire_addresses`, `added_addresses`, `removed_addresses`, `no_change_addresses`, `version`)
  - Each address block contains `ipv4` and `ipv6` CIDR lists
- Register the new data source in `provider.go` and `provider.md`
- Add corresponding documentation `.md` file

## Capabilities

### New Capabilities
- `teo-multi-path-gateway-origin-acl-datasource`: Data source for querying TEO multi-path gateway origin ACL details via DescribeMultiPathGatewayOriginACL API

### Modified Capabilities

## Impact

- New files: `tencentcloud/services/teo/data_source_tc_teo_multi_path_gateway_origin_acl.go`, test file, and `.md` documentation
- Modified files: `tencentcloud/services/teo/service_tencentcloud_teo.go` (add service method), `tencentcloud/provider.go` (register data source), `tencentcloud/provider.md` (add data source entry)
- API dependency: `DescribeMultiPathGatewayOriginACL` from `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
