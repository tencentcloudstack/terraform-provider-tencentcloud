---
subcategory: "Anti-DDoS(Dayu)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dayu_ddos_policy_attachments"
sidebar_current: "docs-tencentcloud-datasource-dayu_ddos_policy_attachments"
description: |-
  Use this data source to query detailed information of dayu DDoS policy attachments
---

# tencentcloud_dayu_ddos_policy_attachments

Use this data source to query detailed information of dayu DDoS policy attachments

## Example Usage

```hcl
data "tencentcloud_dayu_ddos_policy_attachments" "foo_type" {
  resource_type = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_type
}
data "tencentcloud_dayu_ddos_policy_attachments" "foo_resource" {
  resource_id   = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_id
  resource_type = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_type
}
data "tencentcloud_dayu_ddos_policy_attachments" "foo_policy" {
  resource_type = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.resource_type
  policy_id     = tencentcloud_dayu_ddos_policy_attachment.dayu_ddos_policy_attachment.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_type` - (Required) Type of the resource that the DDoS policy works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.
* `policy_id` - (Optional) Id of the policy to be queried.
* `resource_id` - (Optional) Id of the attached resource to be queried.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `dayu_ddos_policy_attachment_list` - A list of dayu DDoS policy attachments. Each element contains the following attributes:
  * `policy_id` - Id of the policy.
  * `resource_id` - Id of the attached resource.
  * `resource_type` - Type of the resource that the DDoS policy works for.


