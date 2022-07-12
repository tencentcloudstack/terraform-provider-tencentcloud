---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_security_groups"
sidebar_current: "docs-tencentcloud-datasource-security_groups"
description: |-
  Use this data source to query detailed information of security groups.
---

# tencentcloud_security_groups

Use this data source to query detailed information of security groups.

## Example Usage

```hcl
data "tencentcloud_security_groups" "sglab" {
  security_group_id = tencentcloud_security_group.sglab.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Name of the security group to be queried. Conflict with `security_group_id`.
* `project_id` - (Optional, Int) Project ID of the security group to be queried. Conflict with `security_group_id`.
* `result_output_file` - (Optional, String) Used to save results.
* `security_group_id` - (Optional, String) ID of the security group to be queried. Conflict with `name` and `project_id`.
* `tags` - (Optional, Map) Tags of the security group to be queried. Conflict with `security_group_id`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `security_groups` - Information list of security group.
  * `be_associate_count` - Number of security group binding resources.
  * `create_time` - Creation time of security group.
  * `description` - Description of the security group.
  * `egress` - Egress rules set. For items like `[action]#[cidr_ip]#[port]#[protocol]`, it means a regular rule; for items like `sg-XXXX`, it means a nested security group.
  * `ingress` - Ingress rules set. For items like `[action]#[cidr_ip]#[port]#[protocol]`, it means a regular rule; for items like `sg-XXXX`, it means a nested security group.
  * `name` - Name of the security group.
  * `project_id` - Project ID of the security group.
  * `security_group_id` - ID of the security group.
  * `tags` - Tags of the security group.


