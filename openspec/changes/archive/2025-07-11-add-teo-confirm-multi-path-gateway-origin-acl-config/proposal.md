## Why

Users of Tencent Cloud EdgeOne (TEO) need to manage the confirmation of Multi-Path Gateway origin ACL updates as infrastructure code. When the origin IP segments (回源 IP 网段) for a multi-path gateway are updated by the cloud service, users must confirm the update to acknowledge they've updated their origin firewall configurations. Currently, there is no Terraform resource to manage this confirmation process, requiring manual intervention through the console or API.

## What Changes

- Add a new Terraform resource `tencentcloud_teo_confirm_multi_path_gateway_origin_acl` of type RESOURCE_KIND_CONFIG
- The resource supports Read and Update operations:
  - **Read**: Query the current multi-path gateway origin ACL details via `DescribeMultiPathGatewayOriginACL`, including current active ACL info and pending next-version ACL info
  - **Update**: Confirm the origin ACL version update via `ConfirmMultiPathGatewayOriginACL`, which acknowledges that the user has updated their origin firewall to the latest IP segments
- The resource uses `zone_id` and `gateway_id` as composite identifiers (separated by `tccommon.FILED_SP`)
- The resource file name follows the pattern: `resource_tc_teo_confirm_multi_path_gateway_origin_acl_config.go`

## Capabilities

### New Capabilities
- `teo-confirm-multi-path-gateway-origin-acl`: Manages the confirmation of multi-path gateway origin ACL updates in TEO, including reading current/pending ACL info and confirming version updates

### Modified Capabilities
<!-- No existing capabilities are modified -->

## Impact

- New files in `tencentcloud/services/teo/`:
  - `resource_tc_teo_confirm_multi_path_gateway_origin_acl_config.go` (resource implementation)
  - `resource_tc_teo_confirm_multi_path_gateway_origin_acl_config_test.go` (unit tests)
  - `resource_tc_teo_confirm_multi_path_gateway_origin_acl_config.md` (documentation)
- Modified files:
  - `tencentcloud/provider.go` (register new resource)
  - `tencentcloud/provider.md` (add resource documentation entry)
- Cloud APIs used:
  - `DescribeMultiPathGatewayOriginACL` (Read)
  - `ConfirmMultiPathGatewayOriginACL` (Update)
- SDK package: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901`
