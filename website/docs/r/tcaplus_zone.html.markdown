---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_zone"
sidebar_current: "docs-tencentcloud-resource-tcaplus_zone"
description: |-
  Use this resource to create tcaplus zone
---

# tencentcloud_tcaplus_zone

Use this resource to create tcaplus zone

## Example Usage

```hcl
resource "tencentcloud_tcaplus_application" "test" {
  idl_type                 = "PROTO"
  app_name                 = "tf_tcaplus_app_test"
  vpc_id                   = "vpc-7k6gzox6"
  subnet_id                = "subnet-akwgvfa3"
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_zone" "zone" {
  app_id    = tencentcloud_tcaplus_application.test.id
  zone_name = "tf_test_zone_name"
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) Application of the tcapplus zone belongs.
* `zone_name` - (Required) Name of the tcapplus zone. length should between 1 and 30.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the tcapplus zone.
* `table_count` - Number of tables.
* `total_size` - The total storage(MB).


