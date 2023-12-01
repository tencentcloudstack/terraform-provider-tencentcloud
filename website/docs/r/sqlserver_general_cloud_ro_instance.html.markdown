---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_general_cloud_ro_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_general_cloud_ro_instance"
description: |-
  Provides a resource to create a sqlserver general_cloud_ro_instance
---

# tencentcloud_sqlserver_general_cloud_ro_instance

Provides a resource to create a sqlserver general_cloud_ro_instance

## Example Usage

### If read_only_group_type value is 1 - Ship according to one instance and one read-only group:

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

resource "tencentcloud_sqlserver_general_cloud_ro_instance" "example" {
  instance_id          = tencentcloud_sqlserver_general_cloud_instance.example.id
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  read_only_group_type = 1
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_BSSD"
  instance_charge_type = "POSTPAID"
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  security_group_list  = [tencentcloud_security_group.security_group.id]
  collation            = "Chinese_PRC_CI_AS"
  time_zone            = "China Standard Time"
  resource_tags = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
}
```

### If read_only_group_type value is 2 - Ship after creating a read-only group, all instances are under this read-only group:

```hcl
resource "tencentcloud_sqlserver_general_cloud_ro_instance" "example" {
  instance_id                      = tencentcloud_sqlserver_general_cloud_instance.example.id
  zone                             = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  read_only_group_type             = 2
  read_only_group_name             = "test-ro-group"
  read_only_group_is_offline_delay = 1
  read_only_group_max_delay_time   = 10
  read_only_group_min_in_group     = 1
  memory                           = 4
  storage                          = 100
  cpu                              = 2
  machine_type                     = "CLOUD_BSSD"
  instance_charge_type             = "POSTPAID"
  subnet_id                        = tencentcloud_subnet.subnet.id
  vpc_id                           = tencentcloud_vpc.vpc.id
  security_group_list              = [tencentcloud_security_group.security_group.id]
  collation                        = "Chinese_PRC_CI_AS"
  time_zone                        = "China Standard Time"
  resource_tags = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
}
```

### If read_only_group_type value is 3 - All instances shipped are in the existing Some read-only groups below:

```hcl
resource "tencentcloud_sqlserver_general_cloud_ro_instance" "example" {
  instance_id          = tencentcloud_sqlserver_general_cloud_instance.example.id
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  read_only_group_type = 3
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_BSSD"
  read_only_group_id   = "mssqlrg-clboghrj"
  instance_charge_type = "POSTPAID"
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  security_group_list  = [tencentcloud_security_group.security_group.id]
  collation            = "Chinese_PRC_CI_AS"
  time_zone            = "China Standard Time"
  resource_tags = {
    test-key1 = "test-value1"
    test-key2 = "test-value2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cpu` - (Required, Int) Number of instance cores.
* `instance_id` - (Required, String) Primary instance ID, in the format: mssql-3l3fgqn7.
* `machine_type` - (Required, String) The host disk type of the purchased instance, CLOUD_HSSD-enhanced SSD cloud disk for virtual machines, CLOUD_TSSD-extremely fast SSD cloud disk for virtual machines, CLOUD_BSSD-universal SSD cloud disk for virtual machines.
* `memory` - (Required, Int) Instance memory size, in GB.
* `read_only_group_type` - (Required, Int) Read-only group type option, 1- Ship according to one instance and one read-only group, 2 - Ship after creating a read-only group, all instances are under this read-only group, 3 - All instances shipped are in the existing Some read-only groups below.
* `storage` - (Required, Int) Instance disk size, in GB.
* `zone` - (Required, String) Instance Availability Zone, similar to ap-guangzhou-1 (Guangzhou District 1); the instance sales area can be obtained through the interface DescribeZones.
* `collation` - (Optional, String) System character set collation, default: Chinese_PRC_CI_AS.
* `instance_charge_type` - (Optional, String) Payment mode, the value supports PREPAID (prepaid), POSTPAID (postpaid).
* `period` - (Optional, Int) Purchase instance period, the default value is 1, which means one month. The value cannot exceed 48.
* `read_only_group_id` - (Optional, String) Required when ReadOnlyGroupType=3, existing read-only group ID.
* `read_only_group_is_offline_delay` - (Optional, Int) Required when ReadOnlyGroupType=2, whether to enable the delayed elimination function for the newly created read-only group, 1-on, 0-off. When the delay between the read-only replica and the primary instance is greater than the threshold, it will be automatically removed.
* `read_only_group_max_delay_time` - (Optional, Int) Mandatory when ReadOnlyGroupType=2 and ReadOnlyGroupIsOfflineDelay=1, the threshold for delay culling of newly created read-only groups.
* `read_only_group_min_in_group` - (Optional, Int) Required when ReadOnlyGroupType=2 and ReadOnlyGroupIsOfflineDelay=1, the newly created read-only group retains at least the number of read-only replicas after delay elimination.
* `read_only_group_name` - (Optional, String) Required when ReadOnlyGroupType=2, the name of the newly created read-only group.
* `resource_tags` - (Optional, Map) Tag description list.
* `security_group_list` - (Optional, Set: [`String`]) Security group list, fill in the security group ID in the form of sg-xxx.
* `subnet_id` - (Optional, String) VPC subnet ID, in the form of subnet-bdoe83fa; SubnetId and VpcId need to be set at the same time or not set at the same time.
* `time_zone` - (Optional, String) System time zone, default: China Standard Time.
* `vpc_id` - (Optional, String) VPC network ID, in the form of vpc-dsp338hz; SubnetId and VpcId need to be set at the same time or not set at the same time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `ro_instance_id` - Primary read only instance ID, in the format: mssqlro-lbljc5qd.


