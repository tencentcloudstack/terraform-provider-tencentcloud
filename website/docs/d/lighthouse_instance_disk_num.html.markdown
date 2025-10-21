---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_instance_disk_num"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_instance_disk_num"
description: |-
  Use this data source to query detailed information of lighthouse instance_disk_num
---

# tencentcloud_lighthouse_instance_disk_num

Use this data source to query detailed information of lighthouse instance_disk_num

## Example Usage

```hcl
data "tencentcloud_lighthouse_instance_disk_num" "instance_disk_num" {
  instance_ids = ["lhins-xxxxxx"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Required, Set: [`String`]) List of instance IDs.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `attach_detail_set` - Mount information list.
  * `attached_disk_count` - Number of elastic cloud disks mounted to the instance.
  * `instance_id` - Instance Id.
  * `max_attach_count` - Number of elastic cloud disks that can be mounted.


