---
subcategory: "MapReduce(EMR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_serverless_hbase_instance"
sidebar_current: "docs-tencentcloud-resource-serverless_hbase_instance"
description: |-
  Provides a resource to create a emr serverless_hbase_instance
---

# tencentcloud_serverless_hbase_instance

Provides a resource to create a emr serverless_hbase_instance

## Example Usage

```hcl
resource "tencentcloud_serverless_hbase_instance" "serverless_hbase_instance" {
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

* `disk_size` - (Required, Int) Instance single-node disk capacity, in GB. The single-node disk capacity must be greater than or equal to 100 and less than or equal to 250 times the number of CPU cores. The capacity adjustment step is 100.
* `disk_type` - (Required, String) Instance disk type, Value range: CLOUD_HSSD: indicate performance cloud storage(ESSD). CLOUD_BSSD: indicate standard cloud storage(SSD).
* `instance_name` - (Required, String) Instance name. Length limit is 6-36 characters. Only Chinese characters, letters, numbers, -, and _ are allowed.
* `pay_mode` - (Required, Int) Instance pay mode. Value range: 0: indicates post-pay mode, that is, pay-as-you-go. 1: indicates pre-pay mode, that is, monthly subscription.
* `zone_settings` - (Required, List) Detailed configuration of the instance availability zone, currently supports multiple availability zones, the number of availability zones can only be 1 or 3, including zone name, VPC information, and number of nodes. The total number of nodes across all zones must be greater than or equal to 3 and less than or equal to 50.
* `auto_renew_flag` - (Optional, Int) AutoRenewFlag, Value range: 0: indicates NOTIFY_AND_MANUAL_RENEW; 1: indicates NOTIFY_AND_AUTO_RENEW; 2: indicates DISABLE_NOTIFY_AND_MANUAL_RENEW.
* `node_type` - (Optional, String) Instance node type, can be filled in as 4C16G, 8C32G, 16C64G, 32C128G, case insensitive.
* `tags` - (Optional, Set) List of tags to bind to the instance.
* `time_span` - (Optional, Int) Time span.
* `time_unit` - (Optional, String) Time unit, fill in m which means month.

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

emr serverless_hbase_instance can be imported using the id, e.g.

```
terraform import tencentcloud_serverless_hbase_instance.serverless_hbase_instance serverless_hbase_instance_id
```

