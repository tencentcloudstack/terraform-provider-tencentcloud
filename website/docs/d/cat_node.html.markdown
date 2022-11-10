---
subcategory: "Cloud Automated Testing(CAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cat_node"
sidebar_current: "docs-tencentcloud-datasource-cat_node"
description: |-
  Use this data source to query detailed information of cat node
---

# tencentcloud_cat_node

Use this data source to query detailed information of cat node

## Example Usage

```hcl
data "tencentcloud_cat_node" "node" {
  node_type = 1
  location  = 2
  is_ipv6   = false
}
```

## Argument Reference

The following arguments are supported:

* `is_ipv6` - (Optional, Bool) is IPv6.
* `location` - (Optional, Int) Node area:1=Chinese Mainland,2=Hong Kong, Macao and Taiwan,3=Overseas.
* `node_name` - (Optional, String) Node name.
* `node_type` - (Optional, Int) Node type 1:IDC,2:LastMile,3:Mobile.
* `pay_mode` - (Optional, Int) Payment mode:1=Trial version,2=Paid version.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `node_define` - Probe node list.
  * `city` - City.
  * `code_type` - If the node type is base, it is an availability dial test point; if it is blank, it is an advanced dial test point.
  * `code` - Node ID.
  * `district` - District.
  * `ip_type` - IP type:1 = IPv4,2 = IPv6.
  * `location` - Node area:1=Chinese Mainland,2=Hong Kong, Macao and Taiwan,3=Overseas.
  * `name` - Node name.
  * `net_service` - Network service provider.
  * `node_define_status` - Node status: 1=running, 2=offline.
  * `type` - Node Type;1 = IDC,2 = LastMile,3 = Mobile.


