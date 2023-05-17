---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_business_intelligence_instance"
sidebar_current: "docs-tencentcloud-resource-sqlserver_business_intelligence_instance"
description: |-
  Provides a resource to create a sqlserver business_intelligence_instance
---

# tencentcloud_sqlserver_business_intelligence_instance

Provides a resource to create a sqlserver business_intelligence_instance

## Example Usage

```hcl
resource "tencentcloud_sqlserver_business_intelligence_instance" "business_intelligence_instance" {
  zone                = "ap-guangzhou-1"
  memory              = 10
  storage             = 100
  cpu                 = 2
  machine_type        = "CLOUD_SSD"
  project_id          = 0
  goods_num           = 1
  subnet_id           = "subnet-bdoe83fa"
  vpc_id              = "vpc-dsp338hz"
  d_b_version         = ""
  security_group_list =
  weekly              =
  start_time          = ""
  span                =
  resource_tags {
    tag_key   = ""
    tag_value = ""

  }
}
```

## Argument Reference

The following arguments are supported:

* `cpu` - (Required, Int) The number of CPU cores of the instance you want to purchase.
* `instance_name` - (Required, String) Instance Name.
* `machine_type` - (Required, String) The host type of purchased instance. Valid values: CLOUD_PREMIUM (virtual machine with premium cloud disk), CLOUD_SSD (virtual machine with SSD).
* `memory` - (Required, Int) Instance memory size in GB.
* `storage` - (Required, Int) Instance disk size in GB.
* `zone` - (Required, String) Instance AZ, such as ap-guangzhou-1 (Guangzhou Zone 1). Purchasable AZs for an instance can be obtained through theDescribeZones API.
* `db_version` - (Optional, String) Supported versions of business intelligence server. Valid values: 201603 (SQL Server 2016 Integration Services), 201703 (SQL Server 2017 Integration Services), 201903 (SQL Server 2019 Integration Services). Default value: 201903. As the purchasable versions are region-specific, you can use the DescribeProductConfig API to query the information of purchasable versions in each region.
* `project_id` - (Optional, Int) Project ID.
* `resource_tags` - (Optional, List) Tags associated with the instances to be created.
* `security_group_list` - (Optional, List: [`String`]) Security group list, which contains security group IDs in the format of sg-xxx.
* `span` - (Optional, Int) Configuration of the maintenance window, which specifies the maintenance duration in hours.
* `start_time` - (Optional, String) Configuration of the maintenance window, which specifies the start time of daily maintenance.
* `subnet_id` - (Optional, String) VPC subnet ID in the format of subnet-bdoe83fa. Both SubnetId and VpcId need to be set or unset at the same time.
* `vpc_id` - (Optional, String) VPC ID in the format of vpc-dsp338hz. Both SubnetId and VpcId need to be set or unset at the same time.
* `weekly` - (Optional, List: [`Int`]) Configuration of the maintenance window, which specifies the day of the week when maintenance can be performed. Valid values: 1 (Monday), 2 (Tuesday), 3 (Wednesday), 4 (Thursday), 5 (Friday), 6 (Saturday), 7 (Sunday).

The `resource_tags` object supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver business_intelligence_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance business_intelligence_instance_id
```

