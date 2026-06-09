---
subcategory: "KeeWiDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_keewidb_instances"
sidebar_current: "docs-tencentcloud-datasource-keewidb_instances"
description: |-
  Use this data source to query KeeWiDB instances.
---

# tencentcloud_keewidb_instances

Use this data source to query KeeWiDB instances.

## Example Usage

### Query all instances

```hcl
data "tencentcloud_keewidb_instances" "example" {}
```

### Query instances by filter

```hcl
data "tencentcloud_keewidb_instances" "example" {
  instance_id     = "kee-4nmzc0ul"
  instance_name   = "tf-example"
  uniq_vpc_ids    = ["vpc-mjwornzj"]
  uniq_subnet_ids = ["subnet-1ed4w7to"]
  billing_mode    = "postpaid"
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew` - (Optional, List: [`Int`]) Filter by renewal mode. 0: manual renewal; 1: auto-renewal; 2: no renewal on expiry.
* `billing_mode` - (Optional, String) Filter by billing mode. postpaid: pay-as-you-go; prepaid: prepaid.
* `instance_id` - (Optional, String) Filter by instance ID, e.g. `kee-6ubhg****`.
* `instance_name` - (Optional, String) Filter by instance name.
* `order_by` - (Optional, String) Sort field. Valid values: projectId, createtime, instancename, type, curDeadline.
* `order_type` - (Optional, Int) Sort direction. 1: descending (default); 0: ascending.
* `project_ids` - (Optional, List: [`Int`]) Filter by project IDs.
* `result_output_file` - (Optional, String) Used to save results.
* `search_key` - (Optional, String) Fuzzy search keyword. Supports instance ID or instance name.
* `search_keys` - (Optional, List: [`String`]) Search keywords. Supports instance ID, instance name, and private network IP.
* `status` - (Optional, List: [`Int`]) Filter by instance status. 0: pending init; 1: in process; 2: running; -2: isolated; -3: to be deleted.
* `subnet_ids` - (Optional, List: [`String`]) Filter by subnet ID (numeric format).
* `tag_keys` - (Optional, List: [`String`]) Filter by tag keys.
* `tag_list` - (Optional, List) Filter by tag key and value.
* `type` - (Optional, Int) Filter by instance type. 13: standard; 14: cluster.
* `uniq_subnet_ids` - (Optional, List: [`String`]) Filter by subnet ID (string format, e.g. subnet-xxx).
* `uniq_vpc_ids` - (Optional, List: [`String`]) Filter by VPC ID (string format, e.g. vpc-xxx).
* `vpc_ids` - (Optional, List: [`String`]) Filter by VPC ID (numeric format).

The `tag_list` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - List of KeeWiDB instances.
  * `auto_renew_flag` - Auto-renewal flag. 1: enabled; 0: disabled.
  * `billing_mode` - Billing mode. 0: pay-as-you-go; 1: prepaid.
  * `compression` - Data compression switch. ON or OFF.
  * `createtime` - Instance creation time.
  * `deadline_time` - Instance expiry time.
  * `disk_size` - Disk capacity (GB).
  * `engine` - Storage engine.
  * `instance_id` - Instance ID.
  * `instance_name` - Instance name.
  * `machine_memory` - Instance memory capacity (GB).
  * `no_auth` - Whether the instance is password-free.
  * `port` - Instance port.
  * `product_type` - Product type. standalone or cluster.
  * `project_id` - Project ID.
  * `project_name` - Project name.
  * `region_id` - Region ID.
  * `region` - Region, e.g. ap-guangzhou.
  * `size` - Total persistent memory capacity (MB).
  * `status` - Instance status.
  * `type` - Instance type. 13: standard; 14: cluster.
  * `uniq_subnet_id` - Subnet ID (string format).
  * `uniq_vpc_id` - VPC ID (string format).
  * `wan_ip` - Instance VIP.
  * `zone_id` - Availability zone ID.


