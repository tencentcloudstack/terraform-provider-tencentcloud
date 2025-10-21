---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_gateway_nodes"
sidebar_current: "docs-tencentcloud-datasource-tse_gateway_nodes"
description: |-
  Use this data source to query detailed information of tse gateway_nodes
---

# tencentcloud_tse_gateway_nodes

Use this data source to query detailed information of tse gateway_nodes

## Example Usage

```hcl
data "tencentcloud_tse_gateway_nodes" "gateway_nodes" {
  gateway_id = "gateway-ddbb709b"
  group_id   = "group-013c0d8e"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) gateway ID.
* `group_id` - (Optional, String) gateway group ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `node_list` - nodes information.
  * `group_id` - Group IDNote: This field may return null, indicating that a valid value is not available.
  * `group_name` - Group nameNote: This field may return null, indicating that a valid value is not available.
  * `node_id` - gateway node id.
  * `node_ip` - gateway node ip.
  * `status` - statusNote: This field may return null, indicating that a valid value is not available.
  * `zone_id` - Zone idNote: This field may return null, indicating that a valid value is not available.
  * `zone` - ZoneNote: This field may return null, indicating that a valid value is not available.


