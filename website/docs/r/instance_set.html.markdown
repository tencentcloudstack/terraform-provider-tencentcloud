---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_instance_set"
sidebar_current: "docs-tencentcloud-resource-instance_set"
description: |-
  Provides a CVM instance set resource.
---

# tencentcloud_instance_set

Provides a CVM instance set resource.

~> **NOTE:** You can launch an CVM instance for a VPC network via specifying parameter `vpc_id`. One instance can only belong to one VPC.

~> **NOTE:** This resource is designed to cater for the scenario of creating CVM in large batches.

~> **NOTE:** After run command `terraform apply`, must wait all cvms is ready, then run command `terraform plan`, either it will casue state change.

## Example Usage

```hcl
data "tencentcloud_images" "my_favorite_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "Tencent Linux release 3.2 (Final)"
}

data "tencentcloud_instance_types" "my_favorite_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "my_favorite_zones" {
}

// Create VPC resource
resource "tencentcloud_vpc" "app" {
  cidr_block = "10.0.0.0/16"
  name       = "awesome_app_vpc"
}

resource "tencentcloud_subnet" "app" {
  vpc_id            = tencentcloud_vpc.app.id
  availability_zone = data.tencentcloud_availability_zones.my_favorite_zones.zones.0.name
  name              = "awesome_app_subnet"
  cidr_block        = "10.0.1.0/24"
}

// Create 10 CVM instances to host awesome_app
resource "tencentcloud_instance_set" "my_awesome_app" {
  timeouts {
    create = "5m"
    read   = "20s"
    delete = "1h"
  }

  instance_count    = 10
  instance_name     = "awesome_app"
  availability_zone = data.tencentcloud_availability_zones.my_favorite_zones.zones.0.name
  image_id          = data.tencentcloud_images.my_favorite_image.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  hostname          = "user"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.app.id
  subnet_id         = tencentcloud_subnet.app.id
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The available zone for the CVM instance.
* `image_id` - (Required, ForceNew) The image to use for the instance. Changing `image_id` will cause the instance reset.
* `allocate_public_ip` - (Optional, ForceNew) Associate a public IP address with an instance in a VPC or Classic. Boolean value, Default is false.
* `bandwidth_package_id` - (Optional) bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.
* `cam_role_name` - (Optional, ForceNew) CAM role name authorized to access.
* `disable_monitor_service` - (Optional) Disable enhance service for monitor, it is enabled by default. When this options is set, monitor agent won't be installed. Modifying will cause the instance reset.
* `disable_security_service` - (Optional) Disable enhance service for security, it is enabled by default. When this options is set, security agent won't be installed. Modifying will cause the instance reset.
* `hostname` - (Optional) The hostname of the instance. Windows instance: The name should be a combination of 2 to 15 characters comprised of letters (case insensitive), numbers, and hyphens (-). Period (.) is not supported, and the name cannot be a string of pure numbers. Other types (such as Linux) of instances: The name should be a combination of 2 to 60 characters, supporting multiple periods (.). The piece between two periods is composed of letters (case insensitive), numbers, and hyphens (-). Modifying will cause the instance reset.
* `instance_charge_type` - (Optional) The charge type of instance. Only support `POSTPAID_BY_HOUR`.
* `instance_count` - (Optional, ForceNew) The number of instances to be purchased. Value range:[1,100]; default value: 1.
* `instance_name` - (Optional) The name of the instance. The max length of instance_name is 60, and default value is `Terraform-CVM-Instance`.
* `instance_type` - (Optional) The type of the instance.
* `internet_charge_type` - (Optional, ForceNew) Internet charge type of the instance, Valid values are `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`. This value does not need to be set when `allocate_public_ip` is false.
* `internet_max_bandwidth_out` - (Optional) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bits per second). This value does not need to be set when `allocate_public_ip` is false.
* `keep_image_login` - (Optional) Whether to keep image login or not, default is `false`. When the image type is private or shared or imported, this parameter can be set `true`. Modifying will cause the instance reset.
* `key_name` - (Optional) The key pair to use for the instance, it looks like `skey-16jig7tx`. Modifying will cause the instance reset.
* `password` - (Optional) Password for the instance. In order for the new password to take effect, the instance will be restarted after the password change. Modifying will cause the instance reset.
* `placement_group_id` - (Optional, ForceNew) The ID of a placement group.
* `private_ip` - (Optional) The private IP to be assigned to this instance, must be in the provided subnet and available.
* `project_id` - (Optional) The project the instance belongs to, default to 0.
* `security_groups` - (Optional) A list of security group IDs to associate with.
* `subnet_id` - (Optional) The ID of a VPC subnet. If you want to create instances in a VPC network, this parameter must be set.
* `system_disk_id` - (Optional) System disk snapshot ID used to initialize the system disk. When system disk type is `LOCAL_BASIC` and `LOCAL_SSD`, disk id is not supported.
* `system_disk_size` - (Optional) Size of the system disk. Valid value ranges: (50~1000). and unit is GB. Default is 50GB. If modified, the instance may force stop.
* `system_disk_type` - (Optional) System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_SSD`: SSD, `CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: 1. `CLOUD_BASIC`, `LOCAL_BASIC` and `LOCAL_SSD` are deprecated; 2. If modified, the instance may force stop.
* `user_data_raw` - (Optional, ForceNew) The user data to be injected into this instance, in plain text. Conflicts with `user_data`. Up to 16 KB after base64 encoded.
* `user_data` - (Optional, ForceNew) The user data to be injected into this instance. Must be base64 encoded and up to 16 KB.
* `vpc_id` - (Optional) The ID of a VPC network. If you want to create instances in a VPC network, this parameter must be set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the instance.
* `expired_time` - Expired time of the instance.
* `instance_ids` - instance id list.
* `instance_status` - Current status of the instance.
* `public_ip` - Public IP of the instance.


