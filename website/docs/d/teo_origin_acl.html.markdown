---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_origin_acl"
sidebar_current: "docs-tencentcloud-datasource-teo_origin_acl"
description: |-
  Use this data source to query detailed information of TEO origin acl
---

# tencentcloud_teo_origin_acl

Use this data source to query detailed information of TEO origin acl

## Example Usage

### Query origin acl by zone Id

```hcl
data "tencentcloud_teo_origin_acl" "example" {
  zone_id = "zone-3fkff38fyw8s"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Specifies the site ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `origin_acl_info` - Describes the binding relationship between the l7 acceleration domain/l4 proxy instance and the origin server IP range.


