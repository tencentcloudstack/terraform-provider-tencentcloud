---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_security_ip_group_content"
sidebar_current: "docs-tencentcloud-datasource-teo_security_ip_group_content"
description: |-
  Use this data source to query the IP list within a specified TEO security IP group.
---

# tencentcloud_teo_security_ip_group_content

Use this data source to query the IP list within a specified TEO security IP group.

## Example Usage

```hcl
data "tencentcloud_teo_security_ip_group_content" "example" {
  zone_id  = "zone-3fkff38fyw8s"
  group_id = 33711
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, Int) IP group ID.
* `zone_id` - (Required, String) Site ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ip_list` - List of IPs or CIDR blocks in the IP group.
* `ip_total_count` - Total count of IPs or CIDR blocks in the IP group.


