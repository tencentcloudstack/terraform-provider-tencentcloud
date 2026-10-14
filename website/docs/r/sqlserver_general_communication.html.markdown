---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_general_communication"
sidebar_current: "docs-tencentcloud-resource-sqlserver_general_communication"
description: |-
  Provides a resource to create a sqlserver general_communication
---

# tencentcloud_sqlserver_general_communication

Provides a resource to create a sqlserver general_communication

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

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_basic_instance" "example" {
  name                   = "tf-example"
  availability_zone      = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = tencentcloud_vpc.vpc.id
  subnet_id              = tencentcloud_subnet.subnet.id
  project_id             = 0
  memory                 = 4
  storage                = 100
  cpu                    = 2
  machine_type           = "CLOUD_PREMIUM"
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  security_groups        = [tencentcloud_security_group.security_group.id]

  tags = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_general_communication" "example" {
  instance_id = tencentcloud_sqlserver_basic_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) ID of instances.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeFlowStatus](https://www.tencentcloud.com/document/product/238/32113) - This API is used to query flow status.
* Create: [OpenInterCommunication](https://www.tencentcloud.com/document/product/238/50605) - This API is used to enable instance interconnection, which can interconnect business intelligence services.
* Read: [DescribeDBInstanceInter](https://www.tencentcloud.com/document/product/238/50606) - This API is used to query the information of the interconnected instances.
* Delete: [CloseInterCommunication](https://www.tencentcloud.com/document/product/238/50611) - This API is used to disable instance interconnection.
* Delete: [DescribeFlowStatus](https://www.tencentcloud.com/document/product/238/32113) - This API is used to query flow status.

## Import

sqlserver general_communication can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_communication.example mssql-hlh6yka1
```

