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

### Cloud type gateway

```hcl
resource "tencentcloud_teo_multi_path_gateway" "example" {
  zone_id      = "zone-3fkff38fyw8s"
  gateway_type = "cloud"
  gateway_name = "tf-example-cloud"
  gateway_port = 8080
  region_id    = "ap-guangzhou"
}
```

### Private type gateway

```hcl
resource "tencentcloud_teo_multi_path_gateway" "example" {
  zone_id      = "zone-3fkff38fyw8s"
  gateway_type = "private"
  gateway_name = "tf-example-private"
  gateway_port = 8080
  gateway_ip   = "10.0.0.1"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_name` - (Required, String) Gateway name, up to 16 characters, available characters (a-z, A-Z, 0-9, -, _).
* `gateway_port` - (Required, Int) Gateway port, range 1-65535 (except 8888).
* `gateway_type` - (Required, String, ForceNew) Gateway type. Valid values: `cloud` (cloud gateway managed by Tencent Cloud), `private` (self-deployed private gateway).
* `zone_id` - (Required, String, ForceNew) Site ID.
* `gateway_ip` - (Optional, String) Gateway address, required when GatewayType is private.
* `region_id` - (Optional, String, ForceNew) Gateway region, required when GatewayType is cloud.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `gateway_id` - Gateway ID.
* `lines` - Line information list.
  * `line_address` - Line address, format is host:port.
  * `line_id` - Line ID.
  * `line_type` - Line type. Valid values: `direct`, `proxy`, `custom`.
  * `proxy_id` - L4 proxy instance ID, returned when LineType is proxy.
  * `rule_id` - Forwarding rule ID, returned when LineType is proxy.
* `need_confirm` - Whether the gateway origin IP list change needs confirmation.
* `status` - Gateway status. Valid values: `creating`, `online`, `offline`, `disable`.


## Import

TEO multi-path gateway can be imported using the zoneId#gatewayId, e.g.

```
terraform import tencentcloud_teo_multi_path_gateway.example zone-3fkff38fyw8s#gw-2qrk328yw8s
```

