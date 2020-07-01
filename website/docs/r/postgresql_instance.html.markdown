---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_instance"
sidebar_current: "docs-tencentcloud-resource-postgresql_instance"
description: |-
  Use this resource to create postgresql instance.
---

# tencentcloud_postgresql_instance

Use this resource to create postgresql instance.

## Example Usage

```hcl
resource "tencentcloud_postgresql_instance" "foo" {
  name              = "example"
  availability_zone = var.availability_zone
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = "vpc-409mvdvv"
  subnet_id         = "subnet-nf9n81ps"
  engine_version    = "9.3.5"
  root_password     = "1qaA2k1wgvfa3ZZZ"
  charset           = "UTF8"
  project_id        = 0
  memory            = 2
  storage           = 10
}
```

## Argument Reference

The following arguments are supported:

* `memory` - (Required) Memory size(in GB). Allowed value must be larger than `memory` that data source `tencentcloud_postgresql_specinfos` provides.
* `name` - (Required) Name of the postgresql instance.
* `root_password` - (Required) Password of root account. This parameter can be specified when you purchase master instances, but it should be ignored when you purchase read-only instances or disaster recovery instances.
* `storage` - (Required) Volume size(in GB). Allowed value must be a multiple of 10. The storage must be set with the limit of `storage_min` and `storage_max` which data source `tencentcloud_postgresql_specinfos` provides.
* `availability_zone` - (Optional, ForceNew) Availability zone.
* `charge_type` - (Optional, ForceNew) Pay type of the postgresql instance. For now, only `POSTPAID_BY_HOUR` is valid.
* `charset` - (Optional, ForceNew) Charset of the root account. Valid values are `UTF8`,`LATIN1`.
* `engine_version` - (Optional, ForceNew) Version of the postgresql database engine. Allowed values are `9.3.5`, `9.5.4`, `10.4`.
* `project_id` - (Optional) Project id, default value is 0.
* `public_access_switch` - (Optional) Indicates whether to enable the access to an instance from public network or not.
* `subnet_id` - (Optional, ForceNew) ID of subnet.
* `vpc_id` - (Optional, ForceNew) ID of VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the postgresql instance.
* `private_access_ip` - IP for private access.
* `private_access_port` - Port for private access.
* `public_access_host` - Host for public access.
* `public_access_port` - Port for public access.


## Import

postgresql instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_postgresql_instance.foo postgres-cda1iex1
```

