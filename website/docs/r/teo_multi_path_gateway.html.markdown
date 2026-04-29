---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_multi_path_gateway"
sidebar_current: "docs-tencentcloud-resource-teo_multi_path_gateway"
description: |-
  Provides a resource to create a TEO multi-path gateway.
---

# tencentcloud_teo_multi_path_gateway

Provides a resource to create a TEO multi-path gateway.

## Example Usage

### Cloud gateway type

```hcl
resource "tencentcloud_teo_multi_path_gateway" "cloud_example" {
  zone_id      = "zone-3fkff38fyw8s"
  gateway_type = "cloud"
  gateway_name = "tf-cloud-gateway"
  gateway_port = 8080
  region_id    = "ap-guangzhou"
}
```

### Private gateway type

```hcl
resource "tencentcloud_teo_multi_path_gateway" "private_example" {
  zone_id      = "zone-3fkff38fyw8s"
  gateway_type = "private"
  gateway_name = "tf-private-gateway"
  gateway_port = 9090
  gateway_ip   = "10.0.0.1"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_name` - (Required, String) Gateway name, up to 16 characters, available characters (a-z, A-Z, 0-9, -, _).
* `gateway_port` - (Required, Int) Gateway port, range 1~65535 (except 8888).
* `gateway_type` - (Required, String, ForceNew) Gateway type. Valid values: `cloud` (cloud gateway managed by Tencent Cloud), `private` (private gateway deployed by user).
* `zone_id` - (Required, String, ForceNew) Site ID.
* `gateway_ip` - (Optional, String) Gateway IP address, required when GatewayType is private. Please ensure the address has been registered in Tencent Cloud Multi-Path Gateway system.
* `region_id` - (Optional, String, ForceNew) Gateway region, required when GatewayType is cloud. You can get RegionId list from DescribeMultiPathGatewayRegions API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `gateway_id` - Gateway ID.
* `need_confirm` - Whether the gateway origin IP list needs to be confirmed. Valid values: `true` (origin IP list changed, need confirmation), `false` (no change, no need to confirm).
* `status` - Gateway status. Valid values: `creating`, `online`, `offline`, `disable`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `10m`) Used when creating the resource.
* `update` - (Defaults to `10m`) Used when updating the resource.

## Import

TEO multi-path gateway can be imported using the zoneId#gatewayId, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway.example zone-3fkff38fyw8s#gw-abc123
```

