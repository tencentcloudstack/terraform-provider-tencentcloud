---
subcategory: "TencentDB for MongoDB(mongodb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_instance_urls"
sidebar_current: "docs-tencentcloud-datasource-mongodb_instance_urls"
description: |-
  Use this data source to query detailed information of mongodb instance urls
---

# tencentcloud_mongodb_instance_urls

Use this data source to query detailed information of mongodb instance urls

## Example Usage

```hcl
data "tencentcloud_mongodb_instance_urls" "mongodb_instance_urls" {
  instance_id = "cmgo-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `urls` - Example connection string access address in the form of an instance URI. Contains: URI type and connection string address.


