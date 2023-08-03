---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_instance_tde"
sidebar_current: "docs-tencentcloud-resource-sqlserver_instance_tde"
description: |-
  Provides a resource to create a sqlserver instance_tde
---

# tencentcloud_sqlserver_instance_tde

Provides a resource to create a sqlserver instance_tde

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_sqlserver_instance" "example" {
  name              = "tf-example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type       = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  project_id        = 0
  memory            = 16
  storage           = 40
}

resource "tencentcloud_sqlserver_instance_tde" "instance_tde" {
  instance_id             = tencentcloud_sqlserver_instance.example.id
  certificate_attribution = "self"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_attribution` - (Required, String) Certificate attribution. self- means to use the account's own certificate, others- means to refer to the certificate of other accounts, and the default is self.
* `instance_id` - (Required, String) Instance ID.
* `quote_uin` - (Optional, String) Other referenced main account IDs, required when CertificateAttribute is others.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver instance_tde can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_instance_tde.instance_tde instance_tde_id
```

