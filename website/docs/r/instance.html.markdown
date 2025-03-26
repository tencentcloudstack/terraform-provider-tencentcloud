---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_instance"
sidebar_current: "docs-tencentcloud-resource-instance"
description: |-
  Provides a CVM instance resource.
---

# tencentcloud_instance

Provides a CVM instance resource.

~> **NOTE:** You can launch an CVM instance for a VPC network via specifying parameter `vpc_id`. One instance can only belong to one VPC.

~> **NOTE:** At present, 'PREPAID' instance cannot be deleted directly and must wait it to be outdated and released automatically.

## Example Usage

### Create a general POSTPAID_BY_HOUR CVM instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

data "tencentcloud_images" "images" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "OpenCloudOS Server"
}

data "tencentcloud_instance_types" "types" {
  filter {
    name   = "instance-family"
    values = ["S1", "S2", "S3", "S4", "S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
}

// create CVM instance
resource "tencentcloud_instance" "example" {
  instance_name     = "tf-example"
  availability_zone = var.availability_zone
  image_id          = data.tencentcloud_images.images.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.types.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  hostname          = "user"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

### Create a general PREPAID CVM instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

data "tencentcloud_images" "images" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "OpenCloudOS Server"
}

data "tencentcloud_instance_types" "types" {
  filter {
    name   = "instance-family"
    values = ["S1", "S2", "S3", "S4", "S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
}

// create CVM instance
resource "tencentcloud_instance" "example" {
  instance_name                           = "tf-example"
  availability_zone                       = var.availability_zone
  image_id                                = data.tencentcloud_images.images.images.0.image_id
  instance_type                           = data.tencentcloud_instance_types.types.instance_types.0.instance_type
  system_disk_type                        = "CLOUD_PREMIUM"
  system_disk_size                        = 50
  hostname                                = "user"
  project_id                              = 0
  vpc_id                                  = tencentcloud_vpc.vpc.id
  subnet_id                               = tencentcloud_subnet.subnet.id
  instance_charge_type                    = "PREPAID"
  instance_charge_type_prepaid_period     = 1
  instance_charge_type_prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  force_delete                            = true
  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }

  timeouts {
    create = "30m"
  }
}
```

### Create a dedicated cluster CVM instance

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

data "tencentcloud_images" "images" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "OpenCloudOS Server"
}

data "tencentcloud_instance_types" "types" {
  filter {
    name   = "instance-family"
    values = ["S1", "S2", "S3", "S4", "S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
  cdc_id            = "cluster-262n63e8"
  is_multicast      = false
}

// create CVM instance
resource "tencentcloud_instance" "example" {
  instance_name        = "tf-example"
  availability_zone    = var.availability_zone
  image_id             = data.tencentcloud_images.images.images.0.image_id
  instance_type        = data.tencentcloud_instance_types.types.instance_types.0.instance_type
  dedicated_cluster_id = "cluster-262n63e8"
  instance_charge_type = "CDCPAID"
  system_disk_type     = "CLOUD_SSD"
  system_disk_size     = 50
  hostname             = "user"
  project_id           = 0
  vpc_id               = tencentcloud_vpc.vpc.id
  subnet_id            = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_SSD"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, String, ForceNew) The available zone for the CVM instance.
* `image_id` - (Required, String) The image to use for the instance. Changing `image_id` will cause the instance reset.
* `allocate_public_ip` - (Optional, Bool, ForceNew) Associate a public IP address with an instance in a VPC or Classic. Boolean value, Default is false.
* `bandwidth_package_id` - (Optional, String) bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.
* `cam_role_name` - (Optional, String) CAM role name authorized to access.
* `cdh_host_id` - (Optional, String, ForceNew) Id of cdh instance. Note: it only works when instance_charge_type is set to `CDHPAID`.
* `cdh_instance_type` - (Optional, String) Type of instance created on cdh, the value of this parameter is in the format of CDH_XCXG based on the number of CPU cores and memory capacity. Note: it only works when instance_charge_type is set to `CDHPAID`.
* `data_disks` - (Optional, List, ForceNew) Settings for data disks.
* `dedicated_cluster_id` - (Optional, String, ForceNew) Exclusive cluster id.
* `disable_api_termination` - (Optional, Bool) Whether the termination protection is enabled. Default is `false`. If set true, which means that this instance can not be deleted by an API action.
* `disable_automation_service` - (Optional, Bool) Disable enhance service for automation, it is enabled by default. When this options is set, monitor agent won't be installed. Modifying will cause the instance reset.
* `disable_monitor_service` - (Optional, Bool) Disable enhance service for monitor, it is enabled by default. When this options is set, monitor agent won't be installed. Modifying will cause the instance reset.
* `disable_security_service` - (Optional, Bool) Disable enhance service for security, it is enabled by default. When this options is set, security agent won't be installed. Modifying will cause the instance reset.
* `force_delete` - (Optional, Bool) Indicate whether to force delete the instance. Default is `false`. If set true, the instance will be permanently deleted instead of being moved into the recycle bin. Note: only works for `PREPAID` instance.
* `hostname` - (Optional, String) The hostname of the instance. Windows instance: The name should be a combination of 2 to 15 characters comprised of letters (case insensitive), numbers, and hyphens (-). Period (.) is not supported, and the name cannot be a string of pure numbers. Other types (such as Linux) of instances: The name should be a combination of 2 to 60 characters, supporting multiple periods (.). The piece between two periods is composed of letters (case insensitive), numbers, and hyphens (-). Modifying will cause the instance reset.
* `instance_charge_type_prepaid_period` - (Optional, Int) The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`, `48`, `60`.
* `instance_charge_type_prepaid_renew_flag` - (Optional, String) Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.
* `instance_charge_type` - (Optional, String) The charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, `SPOTPAID`, `CDHPAID` and `CDCPAID`. The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR` and `CDHPAID`. `PREPAID` instance may not allow to delete before expired. `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time. `CDHPAID` instance must set `cdh_instance_type` and `cdh_host_id`.
* `instance_count` - (Optional, Int, **Deprecated**) It has been deprecated from version 1.59.18. Use built-in `count` instead. The number of instances to be purchased. Value range:[1,100]; default value: 1.
* `instance_name` - (Optional, String) The name of the instance. The max length of instance_name is 128, and default value is `Terraform-CVM-Instance`.
* `instance_type` - (Optional, String) The type of the instance.
* `internet_charge_type` - (Optional, String) Internet charge type of the instance, Valid values are `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`. If not set, internet charge type are consistent with the cvm charge type by default. This value takes NO Effect when changing and does not need to be set when `allocate_public_ip` is false.
* `internet_max_bandwidth_out` - (Optional, Int) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bits per second). This value does not need to be set when `allocate_public_ip` is false.
* `keep_image_login` - (Optional, Bool) Whether to keep image login or not, default is `false`. When the image type is private or shared or imported, this parameter can be set `true`. Modifying will cause the instance reset.
* `key_ids` - (Optional, Set: [`String`]) The key pair to use for the instance, it looks like `skey-16jig7tx`. Modifying will cause the instance reset.
* `key_name` - (Optional, String, **Deprecated**) Please use `key_ids` instead. The key pair to use for the instance, it looks like `skey-16jig7tx`. Modifying will cause the instance reset.
* `orderly_security_groups` - (Optional, List: [`String`]) A list of orderly security group IDs to associate with.
* `password` - (Optional, String) Password for the instance. In order for the new password to take effect, the instance will be restarted after the password change. Modifying will cause the instance reset.
* `placement_group_id` - (Optional, String, ForceNew) The ID of a placement group.
* `private_ip` - (Optional, String) The private IP to be assigned to this instance, must be in the provided subnet and available.
* `project_id` - (Optional, Int) The project the instance belongs to, default to 0.
* `running_flag` - (Optional, Bool) Set instance to running or stop. Default value is true, the instance will shutdown when this flag is false.
* `security_groups` - (Optional, Set: [`String`], **Deprecated**) It will be deprecated. Use `orderly_security_groups` instead. A list of security group IDs to associate with.
* `spot_instance_type` - (Optional, String) Type of spot instance, only support `ONE-TIME` now. Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `spot_max_price` - (Optional, String, ForceNew) Max price of a spot instance, is the format of decimal string, for example "0.50". Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `stopped_mode` - (Optional, String) Billing method of a pay-as-you-go instance after shutdown. Available values: `KEEP_CHARGING`,`STOP_CHARGING`. Default `KEEP_CHARGING`.
* `subnet_id` - (Optional, String) The ID of a VPC subnet. If you want to create instances in a VPC network, this parameter must be set.
* `system_disk_id` - (Optional, String) System disk snapshot ID used to initialize the system disk. When system disk type is `LOCAL_BASIC` and `LOCAL_SSD`, disk id is not supported.
* `system_disk_name` - (Optional, String) Name of the system disk.
* `system_disk_resize_online` - (Optional, Bool) Resize online.
* `system_disk_size` - (Optional, Int) Size of the system disk. unit is GB, Default is 50GB. If modified, the instance may force stop.
* `system_disk_type` - (Optional, String) System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_BASIC`: cloud disk, `CLOUD_SSD`: cloud SSD disk, `CLOUD_PREMIUM`: Premium Cloud Storage, `CLOUD_BSSD`: Basic SSD, `CLOUD_HSSD`: Enhanced SSD, `CLOUD_TSSD`: Tremendous SSD. NOTE: If modified, the instance may force stop.
* `tags` - (Optional, Map) A mapping of tags to assign to the resource. For tag limits, please refer to [Use Limits](https://intl.cloud.tencent.com/document/product/651/13354).
* `user_data_raw` - (Optional, String) The user data to be injected into this instance, in plain text. Conflicts with `user_data`. Up to 16 KB after base64 encoded.
* `user_data` - (Optional, String) The user data to be injected into this instance. Must be base64 encoded and up to 16 KB.
* `vpc_id` - (Optional, String) The ID of a VPC network. If you want to create instances in a VPC network, this parameter must be set.

The `data_disks` object supports the following:

* `data_disk_size` - (Required, Int) Size of the data disk, and unit is GB.
* `data_disk_type` - (Required, String, ForceNew) Data disk type. For more information about limits on different data disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: LOCAL_BASIC: local disk, LOCAL_SSD: local SSD disk, LOCAL_NVME: local NVME disk, specified in the InstanceType, LOCAL_PRO: local HDD disk, specified in the InstanceType, CLOUD_BASIC: HDD cloud disk, CLOUD_PREMIUM: Premium Cloud Storage, CLOUD_SSD: SSD, CLOUD_HSSD: Enhanced SSD, CLOUD_TSSD: Tremendous SSD, CLOUD_BSSD: Balanced SSD.
* `data_disk_id` - (Optional, String) Data disk ID used to initialize the data disk. When data disk type is `LOCAL_BASIC` and `LOCAL_SSD`, disk id is not supported.
* `data_disk_name` - (Optional, String) Name of data disk.
* `data_disk_snapshot_id` - (Optional, String, ForceNew) Snapshot ID of the data disk. The selected data disk snapshot size must be smaller than the data disk size.
* `delete_with_instance_prepaid` - (Optional, Bool, ForceNew) Decides whether the disk is deleted with instance(only applied to `CLOUD_BASIC`, `CLOUD_SSD` and `CLOUD_PREMIUM` disk with `PREPAID` instance), default is false.
* `delete_with_instance` - (Optional, Bool, ForceNew) Decides whether the disk is deleted with instance(only applied to `CLOUD_BASIC`, `CLOUD_SSD` and `CLOUD_PREMIUM` disk with `POSTPAID_BY_HOUR` instance), default is true.
* `encrypt` - (Optional, Bool, ForceNew) Decides whether the disk is encrypted. Default is `false`.
* `throughput_performance` - (Optional, Int, ForceNew) Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cpu` - The number of CPU cores of the instance.
* `create_time` - Create time of the instance.
* `expired_time` - Expired time of the instance.
* `instance_status` - Current status of the instance.
* `memory` - Instance memory capacity, unit in GB.
* `os_name` - Instance os name.
* `public_ip` - Public IP of the instance.
* `uuid` - Globally unique ID of the instance.


## Import

CVM instance can be imported using the id, e.g.

```
terraform import tencentcloud_instance.example ins-2qol3a80
```

