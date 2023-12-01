---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_general_cloud_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_general_cloud_instance"
description: |-
  Provides a resource to create a sqlserver general_cloud_instance
---

# tencentcloud_sqlserver_general_cloud_instance

Provides a resource to create a sqlserver general_cloud_instance

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

resource "tencentcloud_sqlserver_general_cloud_instance" "example" {
  name                 = "tf_example"
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_HSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  db_version           = "2008R2"
  security_group_list  = [tencentcloud_security_group.security_group.id]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
}
```

## Argument Reference

The following arguments are supported:

* `cpu` - (Required, Int) Cpu, unit: CORE.
* `machine_type` - (Required, String) The host disk type of the purchased instance, CLOUD_HSSD-enhanced SSD cloud disk for virtual machines, CLOUD_TSSD-extremely fast SSD cloud disk for virtual machines, CLOUD_BSSD-universal SSD cloud disk for virtual machines.
* `memory` - (Required, Int) Memory, unit: GB.
* `name` - (Required, String) Name of the SQL Server instance.
* `storage` - (Required, Int) instance disk storage, unit: GB.
* `zone` - (Required, String) Instance AZ, such as ap-guangzhou-1 (Guangzhou Zone 1). Purchasable AZs for an instance can be obtained through the DescribeZones API.
* `auto_renew_flag` - (Optional, Int) Automatic renewal flag: 0-normal renewal 1-automatic renewal, the default is 1 automatic renewal. Valid only when purchasing a prepaid instance. Valid only when the 'instance_charge_type' parameter value is 'PREPAID'.
* `collation` - (Optional, String) System character set collation, default: Chinese_PRC_CI_AS.
* `db_version` - (Optional, String) sqlserver version, currently all supported versions are: 2008R2 (SQL Server 2008 R2 Enterprise), 2012SP3 (SQL Server 2012 Enterprise), 201202 (SQL Server 2012 Standard), 2014SP2 (SQL Server 2014 Enterprise), 201402 (SQL Server 2014 Standard), 2016SP1 (SQL Server 2016 Enterprise), 201602 (SQL Server 2016 Standard), 2017 (SQL Server 2017 Enterprise), 201702 (SQL Server 2017 Standard), 2019 (SQL Server 2019 Enterprise), 201902 (SQL Server 2019 Standard). Each region supports different versions for sale, and the version information that can be sold in each region can be pulled through the DescribeProductConfig interface. If left blank, the default version is 2008R2.
* `ha_type` - (Optional, String, **Deprecated**) It has been deprecated from version 1.81.2. Upgrade the high-availability architecture of sqlserver, upgrade from mirror disaster recovery to always on cluster disaster recovery, only support 2017 and above and support always on high-availability instances, do not support downgrading to mirror disaster recovery, CLUSTER-upgrade to always on capacity Disaster, if not filled, the high-availability architecture will not be modified.
* `instance_charge_type` - (Optional, String) Payment mode, the value supports PREPAID (prepaid), POSTPAID (postpaid).
* `period` - (Optional, Int) Purchase instance period, the default value is 1, which means one month. The value cannot exceed 48. Valid only when the 'instance_charge_type' parameter value is 'PREPAID'.
* `project_id` - (Optional, Int) project ID.
* `resource_tags` - (Optional, List) A collection of tags bound to the new instance.
* `security_group_list` - (Optional, Set: [`String`]) Security group list, fill in the security group ID in the form of sg-xxx.
* `span` - (Optional, Int) Maintainable time window configuration, duration, unit: hour.
* `start_time` - (Optional, String) Maintainable time window configuration, daily maintainable start time.
* `subnet_id` - (Optional, String) VPC subnet ID, in the form of subnet-bdoe83fa; SubnetId and VpcId need to be set at the same time or not set at the same time.
* `time_zone` - (Optional, String) System time zone, default: China Standard Time.
* `vpc_id` - (Optional, String) VPC network ID, in the form of vpc-dsp338hz; SubnetId and VpcId need to be set at the same time or not set at the same time.
* `weekly` - (Optional, Set: [`Int`]) Maintainable time window configuration, in weeks, indicates the days of the week that allow maintenance, 1-7 represent Monday to weekend respectively.

The `resource_tags` object supports the following:

* `tag_key` - (Optional, String) tag key.
* `tag_value` - (Optional, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver general_cloud_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_general_cloud_instance.example mssql-i9ma6oy7
```

