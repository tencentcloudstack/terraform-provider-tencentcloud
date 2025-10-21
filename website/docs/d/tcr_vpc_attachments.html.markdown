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

* `instance_id` - (Required, String) ID of the instance to query.
* `result_output_file` - (Optional, String) Used to save results.
* `subnet_id` - (Optional, String) ID of subnet to query.
* `vpc_id` - (Optional, String) ID of VPC to query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vpc_attachment_list` - Information list of the dedicated TCR namespaces.
  * `access_ip` - IP address of this VPC access.
  * `enable_public_domain_dns` - Whether to enable public domain dns.
  * `enable_vpc_domain_dns` - Whether to enable vpc domain dns.
  * `status` - Status of this VPC access.
  * `subnet_id` - ID of subnet.
  * `vpc_id` - ID of VPC.


