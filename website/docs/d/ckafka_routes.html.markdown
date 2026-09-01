---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_routes"
sidebar_current: "docs-tencentcloud-datasource-ckafka_routes"
description: |-
  Use this data source to query detailed information of CKafka routes
---

# tencentcloud_ckafka_routes

Use this data source to query detailed information of CKafka routes

## Example Usage

```hcl
data "tencentcloud_ckafka_routes" "example" {
  instance_id     = "ckafka-exampleabc"
  main_route_flag = true
  route_id        = 123
}

output "routes" {
  value = data.tencentcloud_ckafka_routes.example.routers
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Id of the ckafka instance.
* `main_route_flag` - (Optional, Bool) Whether to display the main route. When set to true, the main route created at instance creation will be additionally returned.
* `result_output_file` - (Optional, String) Used to save results.
* `route_id` - (Optional, Int) Route id, used to query a specific route.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `routers` - A list of ckafka routes. Each element contains the following attributes:
  * `access_type` - Instance access type. 0: PLAINTEXT, 1: SASL_PLAINTEXT, 2: SSL, 3: SASL_SSL.
  * `broker_vip_list` - Virtual IP list (1 to 1 broker nodes).
    * `vip` - Virtual IP.
    * `vport` - Virtual port.
  * `delete_timestamp` - Delete timestamp.
  * `domain_port` - Domain port.
  * `domain` - Domain.
  * `note` - Remark.
  * `route_id` - Route id.
  * `status` - Route status. 1: creating, 2: created, 3: create failed, 4: deleting, 6: delete failed.
  * `subnet` - Subnet id.
  * `vip_list` - Virtual IP list.
    * `vip` - Virtual IP.
    * `vport` - Virtual port.
  * `vip_type` - Routing network type. 3: vpc routing, 7: internal support routing, 1: public network routing.
  * `vpc_id` - Vpc id.


