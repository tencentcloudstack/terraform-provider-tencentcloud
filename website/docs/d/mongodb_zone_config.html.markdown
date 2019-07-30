---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mongodb_zone_config"
sidebar_current: "docs-tencentcloud-datasource-mongodb_zone_config"
description: |-
  Use this data source to query the available mongodb specifications for different zone.
---

# tencentcloud_mongodb_zone_config

Use this data source to query the available mongodb specifications for different zone.

## Example Usage

```hcl
data "tencentcloud_mongodb_zone_config" "mongodb" {
  available_zone = "ap-guangzhou-2"
}
```

## Argument Reference

The following arguments are supported:

* `available_zone` - (Optional) The available zone of the Mongodb.
* `result_output_file` - (Optional) Used to store results.


