---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lite_hbase_instance"
sidebar_current: "docs-tencentcloud-resource-lite_hbase_instance"
description: |-
  Provides a resource to create a emr lite_hbase_instance
---

# tencentcloud_lite_hbase_instance

Provides a resource to create a emr lite_hbase_instance

## Example Usage

```hcl
resource "tencentcloud_lite_hbase_instance" "lite_hbase_instance" {
  instance_name = "tf-test"
  pay_mode      = 0
  disk_type     = "CLOUD_HSSD"
  disk_size     = 100
  node_type     = "8C32G"
  zone_settings {
    zone = "ap-shanghai-2"
    vpc_settings {
      vpc_id    = "vpc-xxxxxx"
      subnet_id = "subnet-xxxxxx"
    }
    node_num = 3
  }
  tags {
    tag_key   = "test"
    tag_value = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `disk_size` - (Required, Int) Instance single-node disk capacity, in GB. The single-node disk capacity must be greater than or equal to 100 and less than or equal to 10000, with an adjustment step size of 20.
* `disk_type` - (Required, String) Instance disk type, fill in CLOUD_HSSD to indicate performance cloud storage.
* `instance_name` - (Required, String) Instance name. Length limit is 6-36 characters. Only Chinese characters, letters, numbers, -, and _ are allowed.
* `node_type` - (Required, String) Instance node type, can be filled in as 4C16G, 8C32G, 16C64G, 32C128G, case insensitive.
* `pay_mode` - (Required, Int) Instance pay mode. Value range: 0: indicates post pay mode, that is, pay-as-you-go.
* `zone_settings` - (Required, List) Detailed configuration of the instance availability zone, currently supports multiple availability zones, the number of availability zones can only be 1 or 3, including zone name, VPC information, and number of nodes. The total number of nodes across all zones must be greater than or equal to 3 and less than or equal to 50.
* `tags` - (Optional, List) List of tags to bind to the instance.

The `tags` object supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

The `vpc_settings` object of `zone_settings` supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) VPC ID.

The `zone_settings` object supports the following:

* `node_num` - (Required, Int) Number of nodes.
* `vpc_settings` - (Required, List) Private network related information configuration. This parameter can be used to specify the ID of the private network, subnet ID, and other information.
* `zone` - (Required, String) The availability zone to which the instance belongs, such as ap-guangzhou-1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

emr lite_hbase_instance can be imported using the id, e.g.

```
terraform import tencentcloud_lite_hbase_instance.lite_hbase_instance lite_hbase_instance_id
```

