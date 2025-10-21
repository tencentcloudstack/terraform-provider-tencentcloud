---
subcategory: "VPN Connections(VPN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_customer_gateways"
sidebar_current: "docs-tencentcloud-datasource-vpn_customer_gateways"
description: |-
  Use this data source to query detailed information of VPN customer gateways.
---

# tencentcloud_vpn_customer_gateways

Use this data source to query detailed information of VPN customer gateways.

## Example Usage

### Query all customer gateways

```hcl
data "tencentcloud_vpn_customer_gateways" "example" {}
```

### Query customer gateways by filters

```hcl
data "tencentcloud_vpn_customer_gateways" "example" {
  name              = "tf-example"
  id                = "cgw-r1g6c8fr"
  public_ip_address = "1.1.1.1"
  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) ID of the VPN customer gateway.
* `name` - (Optional, String) Name of the customer gateway. The length of character is limited to 1-60.
* `public_ip_address` - (Optional, String) Public ip address of the VPN customer gateway.
* `result_output_file` - (Optional, String) Used to save results.
* `tags` - (Optional, Map) Tags of the VPN customer gateway to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `gateway_list` - Information list of the dedicated gateways.
  * `create_time` - Create time of the VPN customer gateway.
  * `id` - ID of the VPN customer gateway.
  * `name` - Name of the VPN customer gateway.
  * `public_ip_address` - Public ip address of the VPN customer gateway.
  * `tags` - Tags of the VPN customer gateway.


