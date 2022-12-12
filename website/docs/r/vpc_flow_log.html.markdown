---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_flow_log"
sidebar_current: "docs-tencentcloud-resource-vpc_flow_log"
description: |-
  Provides a resource to create a vpc flow_log
---

# tencentcloud_vpc_flow_log

Provides a resource to create a vpc flow_log

## Example Usage

```hcl
resource "tencentcloud_vpc_flow_log" "flow_log" {
  flow_log_name        = "foo"
  resource_type        = "NETWORKINTERFACE"
  resource_id          = "eni-xxxxxxxx"
  traffic_type         = "ALL"
  vpc_id               = "vpc-xxxxxxxx"
  flow_log_description = "My testing log"
  cloud_log_id         = "a1b2c3d4-e5f6a7b8-c9d0e1f2-a3b4c5d6"
  storage_type         = "cls"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `flow_log_name` - (Required, String) Specify flow log name.
* `resource_id` - (Required, String) Specify resource unique Id of `resource_type` configured.
* `resource_type` - (Required, String) Specify resource type. NOTE: Only support `NETWORKINTERFACE` for now. Values: `VPC`, `SUBNET`, `NETWORKINTERFACE`, `CCN`, `NAT`, `DCG`.
* `traffic_type` - (Required, String) Specify log traffic type, values: `ACCEPT`, `REJECT`, `ALL`.
* `cloud_log_id` - (Optional, String) Specify flow log storage id.
* `cloud_log_region` - (Optional, String) Specify flow log storage region, default using current.
* `flow_log_description` - (Optional, String) Specify flow Log description.
* `flow_log_storage` - (Optional, List) Specify consumer detail, required while `storage_type` is `ckafka`.
* `storage_type` - (Optional, String) Specify consumer type, values: `cls`, `ckafka`.
* `tags` - (Optional, Map) Tag description list.
* `vpc_id` - (Optional, String) Specify vpc Id, ignore while `resource_type` is `CCN` (unsupported) but required while other types.

The `flow_log_storage` object supports the following:

* `storage_id` - (Optional, String) Specify storage instance id, required while `storage_type` is `ckafka`.
* `storage_topic` - (Optional, String) Specify storage topic id, required while `storage_type` is `ckafka`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc flow_log can be imported using the flow log Id combine vpc Id, e.g.

```
$ terraform import tencentcloud_vpc_flow_log.flow_log flow_log_id fl-xxxx1234#vpc-yyyy5678
```

