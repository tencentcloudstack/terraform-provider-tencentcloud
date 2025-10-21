---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_connections"
sidebar_current: "docs-tencentcloud-datasource-mongodb_instance_connections"
description: |-
  Use this data source to query detailed information of mongodb instance_connections
---

# tencentcloud_mongodb_instance_connections

Use this data source to query detailed information of mongodb instance_connections

## Example Usage

```hcl
data "tencentcloud_mongodb_instance_connections" "instance_connections" {
  instance_id = "cmgo-9d0p6umb"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `clients` - Client connection info.
  * `count` - client connection count.
  * `internal_service` - is internal.
  * `ip` - client connection ip.


