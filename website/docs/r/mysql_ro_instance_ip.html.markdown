---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_ro_instance_ip"
sidebar_current: "docs-tencentcloud-resource-mysql_ro_instance_ip"
description: |-
  Provides a resource to create a mysql ro_instance_ip
---

# tencentcloud_mysql_ro_instance_ip

Provides a resource to create a mysql ro_instance_ip

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_mysql_ro_instance_ip" "example" {
  instance_id    = "cdbro-bdlvcfpj"
  uniq_subnet_id = tencentcloud_subnet.subnet.id
  uniq_vpc_id    = tencentcloud_vpc.vpc.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Read-only instance ID, in the format: cdbro-3i70uj0k, which is the same as the read-only instance ID displayed on the cloud database console page.
* `uniq_subnet_id` - (Optional, String, ForceNew) Subnet descriptor, for example: subnet-1typ0s7d.
* `uniq_vpc_id` - (Optional, String, ForceNew) vpc descriptor, for example: vpc-a23yt67j, if this field is passed, UniqSubnetId must be passed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `ro_vip` - Intranet IP address of the read-only instance.
* `ro_vport` - Intranet port number of the read-only instance.


