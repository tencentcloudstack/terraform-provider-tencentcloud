---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ha_vip_eip_attachments"
sidebar_current: "docs-tencentcloud-datasource-ha_vip_eip_attachments"
description: |-
  Use this data source to query detailed information of HA VIP EIP attachments
---

# tencentcloud_ha_vip_eip_attachments

Use this data source to query detailed information of HA VIP EIP attachments

## Example Usage

```hcl
data "tencentcloud_ha_vip_eip_attachments" "foo" {
  havip_id   = "havip-kjqwe4ba"
  address_ip = "1.1.1.1"
}
```

## Argument Reference

The following arguments are supported:

* `havip_id` - (Required, String) ID of the attached HA VIP to be queried.
* `address_ip` - (Optional, String) Public IP address of EIP to be queried.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ha_vip_eip_attachment_list` - A list of HA VIP EIP attachments. Each element contains the following attributes:
  * `address_ip` - Public IP address of EIP.
  * `havip_id` - ID of the attached HA VIP.


