---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_vpc_attachments"
sidebar_current: "docs-tencentcloud-datasource-tcr_vpc_attachments"
description: |-
  Use this data source to query detailed information of TCR VPC attachment.
---

# tencentcloud_tcr_vpc_attachments

Use this data source to query detailed information of TCR VPC attachment.

## Example Usage

```hcl
data "tencentcloud_tcr_vpc_attachments" "id" {
  instance_id = "cls-satg5125"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the instance to query.
* `result_output_file` - (Optional) Used to save results.
* `subnet_id` - (Optional) ID of subnet to query.
* `vpc_id` - (Optional) ID of VPC to query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vpc_attachment_list` - Information list of the dedicated TCR namespaces.
  * `access_ip` - IP address of this VPC access.
  * `status` - Status of this VPC access.
  * `subnet_id` - ID of subnet.
  * `vpc_id` - ID of VPC.


