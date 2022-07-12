---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_enis"
sidebar_current: "docs-tencentcloud-datasource-enis"
description: |-
  Use this data source to query query ENIs.
---

# tencentcloud_enis

Use this data source to query query ENIs.

## Example Usage

```hcl
data "tencentcloud_enis" "name" {
  name = "test eni"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, String) Description of the ENI. Conflict with `ids`.
* `ids` - (Optional, Set: [`String`]) ID of the ENIs to be queried. Conflict with `vpc_id`,`subnet_id`,`instance_id`,`security_group`,`name`,`ipv4` and `tags`.
* `instance_id` - (Optional, String) ID of the instance which bind the ENI. Conflict with `ids`.
* `ipv4` - (Optional, String) Intranet IP of the ENI. Conflict with `ids`.
* `name` - (Optional, String) Name of the ENI to be queried. Conflict with `ids`.
* `result_output_file` - (Optional, String) Used to save results.
* `security_group` - (Optional, String) A set of security group IDs which bind the ENI. Conflict with `ids`.
* `subnet_id` - (Optional, String) ID of the subnet within this vpc to be queried. Conflict with `ids`.
* `tags` - (Optional, Map) Tags of the ENI. Conflict with `ids`.
* `vpc_id` - (Optional, String) ID of the vpc to be queried. Conflict with `ids`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `enis` - An information list of ENIs. Each element contains the following attributes:
  * `create_time` - Creation time of the ENI.
  * `description` - Description of the ENI.
  * `id` - ID of the ENI.
  * `instance_id` - ID of the instance which bind the ENI.
  * `ipv4s` - A set of intranet IPv4s.
    * `description` - Description of the IP.
    * `ip` - Intranet IP.
    * `primary` - Indicates whether the IP is primary.
  * `mac` - MAC address.
  * `name` - Name of the ENI.
  * `primary` - Indicates whether the IP is primary.
  * `security_groups` - A set of security group IDs which bind the ENI.
  * `state` - States of the ENI.
  * `subnet_id` - ID of the subnet within this vpc.
  * `tags` - Tags of the ENI.
  * `vpc_id` - ID of the vpc.


