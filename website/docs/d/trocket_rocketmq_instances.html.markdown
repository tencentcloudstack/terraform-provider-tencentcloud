---
subcategory: "TDMQ for RocketMQ(trocket)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_trocket_rocketmq_instances"
sidebar_current: "docs-tencentcloud-datasource-trocket_rocketmq_instances"
description: |-
  Use this data source to query detailed information of TROCKET rocketmq instances
---

# tencentcloud_trocket_rocketmq_instances

Use this data source to query detailed information of TROCKET rocketmq instances

## Example Usage

### Query all instances

```hcl
data "tencentcloud_trocket_rocketmq_instances" "example" {}
```

### Query instances by filters

```hcl
data "tencentcloud_trocket_rocketmq_instances" "example" {
  filters {
    name   = "InstanceId"
    values = ["rmq-1n58qbwg3"]
  }

  filters {
    name   = "InstanceName"
    values = ["tf-example"]
  }

  tag_filters {
    tag_key    = "createBy"
    tag_values = ["Terraform"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter query criteria list.
* `result_output_file` - (Optional, String) Used to save results.
* `tag_filters` - (Optional, List) Tag filters.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter values.

The `tag_filters` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_values` - (Required, Set) Tag values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Instance list.


