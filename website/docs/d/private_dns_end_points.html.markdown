---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_end_points"
sidebar_current: "docs-tencentcloud-datasource-private_dns_end_points"
description: |-
  Use this data source to query detailed information of Private Dns end points
---

# tencentcloud_private_dns_end_points

Use this data source to query detailed information of Private Dns end points

## Example Usage

### Query all private dns end points

```hcl
data "tencentcloud_private_dns_end_points" "example" {}
```

### Query all private dns end points by filters

```hcl
data "tencentcloud_private_dns_end_points" "example" {
  filters {
    name   = "EndPointName"
    values = ["tf-example"]
  }

  filters {
    name   = "EndPointId"
    values = ["eid-72dc11b8f3"]
  }

  filters {
    name   = "EndPointServiceId"
    values = ["vpcsvc-61wcwmar"]
  }

  filters {
    name = "EndPointVip"
    values = [
      "172.10.10.1"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter parameters. Valid values: EndPointName, EndPointId, EndPointServiceId, and EndPointVip.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Parameter name.
* `values` - (Required, Set) Array of parameter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `end_point_set` - Endpoint list.
Note: This field may return null, indicating that no valid values can be obtained.


