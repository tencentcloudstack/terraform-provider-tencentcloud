---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_gateway_instances"
sidebar_current: "docs-tencentcloud-datasource-dc_gateway_instances"
description: |-
  Use this data source to query detailed information of direct connect gateway instances.
---

# tencentcloud_dc_gateway_instances

Use this data source to query detailed information of direct connect gateway instances.

## Example Usage

```hcl
resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_dc_gateway" "ccn_main" {
  name                = "ci-cdg-ccn-test"
  network_instance_id = "${tencentcloud_ccn.main.id}"
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}

#You need to sleep for a few seconds because there is a cache on the server
data "tencentcloud_dc_gateway_instances" "name_select"{
  name = "${tencentcloud_dc_gateway.ccn_main.name}"
}

data "tencentcloud_dc_gateway_instances"  "id_select" {
  dcg_id = "${tencentcloud_dc_gateway.ccn_main.id}"
}
```

## Argument Reference

The following arguments are supported:

* `dcg_id` - (Optional, ForceNew) ID of the DCG to be queried.
* `name` - (Optional, ForceNew) Name of the DCG to be queried.
* `result_output_file` - (Optional, ForceNew) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - Information list of the DCG.
  * `cnn_route_type` - Type of CCN route, the available value include 'BGP' and 'STATIC'.
  * `create_time` - Creation time of resource.
  * `dcg_id` - ID of the DCG
  * `dcg_ip` - IP of the DCG
  * `enable_bgp` - Indicates whether the BGP is enabled.
  * `gateway_type` - Type of the gateway, the available value include 'NORMAL' and 'NAT'. Default is 'NORMAL'.
  * `name` - Name of the DCG
  * `network_instance_id` - Type of associated network, the available value include 'VPC' and 'CCN'.
  * `network_type` - IP of the DCG


