---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_instance"
sidebar_current: "docs-tencentcloud-resource-instance"
description: |-
  Provides a CVM instance resource.
---

# tencentcloud_instance

Provides a CVM instance resource.

~> **NOTE:** You can launch an CVM instance for a VPC network via specifying parameter `vpc_id`. One instance can only belong to one VPC.

~> **NOTE:** At present, 'PREPAID' instance cannot be deleted and must wait it to be outdated and released automatically.

## Example Usage

```hcl
data "tencentcloud_images" "my_favorate_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "my_favorate_zones" {
}

// Create VPC resource
resource "tencentcloud_vpc" "app" {
  cidr_block = "10.0.0.0/16"
  name       = "awesome_app_vpc"
}

resource "tencentcloud_subnet" "app" {
  vpc_id            = "${tencentcloud_vpc.app.id}"
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  name              = "awesome_app_subnet"
  cidr_block        = "10.0.1.0/24"
}

// Create 2 CVM instances to host awesome_app
resource "tencentcloud_instance" "my_awesome_app" {
  instance_name              = "awesome_app"
  availability_zone          = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id                   = "${data.tencentcloud_images.my_favorate_image.images.0.image_id}"
  instance_type              = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = "${tencentcloud_vpc.app.id}"
  subnet_id                  = "${tencentcloud_subnet.app.id}"
  internet_max_bandwidth_out = 20
  count                      = 2

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The available zone that the CVM instance locates at.
* `image_id` - (Required, ForceNew) The Image to use for the instance. Change 'image_id' will case instance destroy and re-created.
* `allocate_public_ip` - (Optional, ForceNew) Associate a public ip address with an instance in a VPC or Classic. Boolean value, Default is false.
* `data_disks` - (Optional) Settings for data disk.
* `disable_monitor_service` - (Optional) Disable enhance service for monitor, it is enabled by default. When this options is set, monitor agent won't be installed.
* `disable_security_service` - (Optional) Disable enhance service for security, it is enabled by default. When this options is set, security agent won't be installed.
* `hostname` - (Optional, ForceNew) The hostname of CVM. Windows instance: The name should be a combination of 2 to 15 characters comprised of letters (case insensitive), numbers, and hyphens (-). Period (.) is not supported, and the name cannot be a string of pure numbers. Other types (such as Linux) of instances: The name should be a combination of 2 to 60 characters, supporting multiple periods (.). The piece between two periods is composed of letters (case insensitive), numbers, and hyphens (-).
* `instance_charge_type_prepaid_period` - (Optional) The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36.
* `instance_charge_type_prepaid_renew_flag` - (Optional) When enabled, the CVM instance will be renew automatically when it reach the end of the prepaid tenancy. Valid values are `NOTIFY_AND_AUTO_RENEW`, `NOTIFY_AND_MANUAL_RENEW` and `DISABLE_NOTIFY_AND_MANUAL_RENEW`. NOTE: it only works when instance_charge_type is set to `PREPAID`.
* `instance_charge_type` - (Optional, ForceNew) The charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR` and `SPOTPAID`, The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`. `PREPAID` instance may not allow to delete before expired. `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time.
* `instance_name` - (Optional) The name of the CVM. The max length of instance_name is 60, and default value is `Terrafrom-CVM-Instance`.
* `instance_type` - (Optional, ForceNew) The type of instance to start.
* `internet_charge_type` - (Optional, ForceNew) Internet charge type of the instance, Valid values are `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`. The default is `TRAFFIC_POSTPAID_BY_HOUR`.
* `internet_max_bandwidth_out` - (Optional, ForceNew) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bit per second). If this value is not specified, then automatically sets it to 0 Mbps.
* `key_name` - (Optional) The key pair to use for the instance, it looks like skey-16jig7tx.
* `password` - (Optional) Password to an instance. In order to take effect new password, the instance will be restarted after modifying the password.
* `placement_group_id` - (Optional, ForceNew) The id of a placement group.
* `private_ip` - (Optional) The private ip to be assigned to this instance, must be in the provided subnet and available.
* `project_id` - (Optional) The project CVM belongs to, default to 0.
* `running_flag` - (Optional) Set instance to running or stop. Default value is true, the instance will shutdown when flag is false.
* `security_groups` - (Optional) A list of security group ids to associate with.
* `spot_instance_type` - (Optional) Type of spot instance, only support `ONE-TIME` now. Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `spot_max_price` - (Optional, ForceNew) Max price of spot instance, is the format of decimal string, for example "0.50". Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `subnet_id` - (Optional) The id of a VPC subnetwork. If you want to create instances in VPC network, this parameter must be set.
* `system_disk_id` - (Optional) System disk snapshot ID used to initialize the system disk. When system disk type is `LOCAL_BASIC` and `LOCAL_SSD`, disk id is not supported.
* `system_disk_size` - (Optional) Size of the system disk. Value range: [50, 1000], and unit is GB. Default is 50GB.
* `system_disk_type` - (Optional, ForceNew) Type of the system disk. Valid values are `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_SSD` and `CLOUD_PREMIUM`, default value is `CLOUD_BASIC`. NOTE: `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.
* `tags` - (Optional) A mapping of tags to assign to the resource. For tag limits, please refer to [Use Limits](https://intl.cloud.tencent.com/document/product/651/13354).
* `user_data_raw` - (Optional, ForceNew) The user data to be specified into this instance, plain text. Conflicts with `user_data`. Limited in 16 KB after encrypted in base64 format.
* `user_data` - (Optional, ForceNew) The user data to be specified into this instance. Must be encrypted in base64 format and limited in 16 KB.
* `vpc_id` - (Optional) The id of a VPC network. If you want to create instances in VPC network, this parameter must be set.

The `data_disks` object supports the following:

* `data_disk_size` - (Required) Size of the data disk, and unit is GB. If disk type is `CLOUD_SSD`, the size range is [100, 16000], and the others are [10-16000].
* `data_disk_type` - (Required) Type of the data disk. Valid values are `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_SSD` and `CLOUD_PREMIUM`. NOTE: `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.
* `data_disk_id` - (Optional) Data disk snapshot ID used to initialize the data disk. When data disk type is `LOCAL_BASIC` and `LOCAL_SSD`, disk id is not supported.
* `delete_with_instance` - (Optional) Decides whether the disk is deleted with instance(only applied to cloud disk), default to true.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Create time of the instance.
* `expired_time` - Expired time of the instance.
* `instance_status` - Current status of the instance.
* `public_ip` - Public ip of the instance.


## Import

CVM instance can be imported using the id, e.g.

```
terraform import tencentcloud_instance.foo ins-2qol3a80
```

