---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_instances_set"
sidebar_current: "docs-tencentcloud-datasource-instances_set"
description: |-
  Use this data source to query cvm instances in parallel.
---

# tencentcloud_instances_set

Use this data source to query cvm instances in parallel.

## Example Usage

```hcl
data "tencentcloud_instances_set" "foo" {
  vpc_id = "vpc-4owdpnwr"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, String) The available zone that the CVM instance locates at.
* `instance_id` - (Optional, String) ID of the instances to be queried.
* `instance_name` - (Optional, String) Name of the instances to be queried.
* `project_id` - (Optional, Int) The project CVM belongs to.
* `result_output_file` - (Optional, String) Used to save results.
* `subnet_id` - (Optional, String) ID of a vpc subnetwork.
* `tags` - (Optional, Map) Tags of the instance.
* `vpc_id` - (Optional, String) ID of the vpc to be queried.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - An information list of cvm instance. Each element contains the following attributes:
  * `allocate_public_ip` - Indicates whether public ip is assigned.
  * `availability_zone` - The available zone that the CVM instance locates at.
  * `cam_role_name` - CAM role name authorized to access.
  * `cpu` - The number of CPU cores of the instance.
  * `create_time` - Creation time of the instance.
  * `data_disks` - An information list of data disk. Each element contains the following attributes:
    * `data_disk_id` - Image ID of the data disk.
    * `data_disk_size` - Size of the data disk.
    * `data_disk_type` - Type of the data disk.
    * `delete_with_instance` - Indicates whether the data disk is destroyed with the instance.
  * `expired_time` - Expired time of the instance.
  * `image_id` - ID of the image.
  * `instance_charge_type_prepaid_renew_flag` - The way that CVM instance will be renew automatically or not when it reach the end of the prepaid tenancy.
  * `instance_charge_type` - The charge type of the instance.
  * `instance_id` - ID of the instances.
  * `instance_name` - Name of the instances.
  * `instance_type` - Type of the instance.
  * `internet_charge_type` - The charge type of the instance.
  * `internet_max_bandwidth_out` - Public network maximum output bandwidth of the instance.
  * `memory` - Instance memory capacity, unit in GB.
  * `private_ip` - Private IP of the instance.
  * `project_id` - The project CVM belongs to.
  * `public_ip` - Public IP of the instance.
  * `security_groups` - Security groups of the instance.
  * `status` - Status of the instance.
  * `subnet_id` - ID of a vpc subnetwork.
  * `system_disk_id` - Image ID of the system disk.
  * `system_disk_size` - Size of the system disk.
  * `system_disk_type` - Type of the system disk.
  * `tags` - Tags of the instance.
  * `vpc_id` - ID of the vpc.


