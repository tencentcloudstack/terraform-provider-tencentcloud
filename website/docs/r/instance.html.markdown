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

~> **NOTE:** At present, `PREPAID` instance cannot be deleted directly and must wait it to be outdated and released automatically.

~> **NOTE:** Currently, the `placement_group_id` field only supports setting and modification, but not deletion.

~> **NOTE:** When creating a CVM instance using a `launch_template_id`, if you set other parameter values ​​at the same time, the template definition values ​​will be overwritten.

~> **NOTE:** It is recommended to use resource `tencentcloud_eip` to create a AntiDDos Eip, and then call resource `tencentcloud_eip_association` to bind it to resource `tencentcloud_instance`.

~> **NOTE:** When creating a prepaid CVM instance and binding a data disk, you need to explicitly set `delete_with_instance` to `false`.

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

### Create CVM instance with placement_group_id

```hcl
resource "tencentcloud_instance" "example" {
  instance_name                    = "tf-example"
  availability_zone                = "ap-guangzhou-6"
  image_id                         = "img-eb30mz89"
  instance_type                    = "S5.MEDIUM4"
  system_disk_size                 = 50
  system_disk_name                 = "sys_disk_1"
  hostname                         = "user"
  project_id                       = 0
  vpc_id                           = "vpc-i5yyodl9"
  subnet_id                        = "subnet-hhi88a58"
  placement_group_id               = "ps-ejt4brtz"
  force_replace_placement_group_id = false

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 100
    encrypt        = false
    data_disk_name = "data_disk_1"
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

### Create CVM instance with template

```hcl
resource "tencentcloud_instance" "example" {
  launch_template_id      = "lt-b20scl2a"
  launch_template_version = 1
}
```

### Create CVM instance with AntiDDos Eip

```hcl
resource "tencentcloud_instance" "example" {
  instance_name              = "tf-example"
  availability_zone          = "ap-guangzhou-6"
  image_id                   = "img-eb30mz89"
  instance_type              = "S5.MEDIUM4"
  system_disk_type           = "CLOUD_HSSD"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = "vpc-i5yyodl9"
  subnet_id                  = "subnet-hhi88a58"
  orderly_security_groups    = ["sg-l222vn6w"]
  internet_charge_type       = "BANDWIDTH_PACKAGE"
  bandwidth_package_id       = "bwp-rp2nx3ab"
  ipv4_address_type          = "AntiDDoSEIP"
  anti_ddos_package_id       = "bgp-31400fvq"
  allocate_public_ip         = true
  internet_max_bandwidth_out = 100
  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 100
    encrypt        = false
  }
  tags = {
    tagKey = "tagValue"
  }
}
```

### Create CVM instance with setting running flag

```hcl
resource "tencentcloud_instance" "example" {
  instance_name           = "tf-example"
  availability_zone       = "ap-guangzhou-6"
  image_id                = "img-eb30mz89"
  instance_type           = "S5.MEDIUM4"
  system_disk_type        = "CLOUD_HSSD"
  system_disk_size        = 50
  hostname                = "user"
  project_id              = 0
  vpc_id                  = "vpc-i5yyodl9"
  subnet_id               = "subnet-hhi88a58"
  orderly_security_groups = ["sg-ma82yjwp"]
  running_flag            = false
  stop_type               = "SOFT_FIRST"
  stopped_mode            = "KEEP_CHARGING"
  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 100
    encrypt        = false
  }
  tags = {
    tagKey = "tagValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `allocate_public_ip` - (Optional, Bool, ForceNew) Associate a public IP address with an instance in a VPC or Classic. Boolean value, Default is false.
* `anti_ddos_package_id` - (Optional, String, ForceNew) Anti-DDoS service package ID. This is required when you want to request an AntiDDoS IP.
* `availability_zone` - (Optional, String, ForceNew) The available zone for the CVM instance.
* `bandwidth_package_id` - (Optional, String) bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.
* `cam_role_name` - (Optional, String) CAM role name authorized to access.
* `cdh_host_id` - (Optional, String, ForceNew) Id of cdh instance. Note: it only works when instance_charge_type is set to `CDHPAID`.
* `cdh_instance_type` - (Optional, String) Type of instance created on cdh, the value of this parameter is in the format of CDH_XCXG based on the number of CPU cores and memory capacity. Note: it only works when instance_charge_type is set to `CDHPAID`.
* `data_disks` - (Optional, List, ForceNew) Settings for data disks.
* `dedicated_cluster_id` - (Optional, String, ForceNew) Exclusive cluster id.
* `disable_api_termination` - (Optional, Bool) Whether the termination protection is enabled. Default is `false`. If set true, which means that this instance can not be deleted by an API action.
* `disable_automation_service` - (Optional, Bool) Disable enhance service for automation, it is enabled by default. When this options is set, monitor agent won't be installed. Modifications may lead to the reinstallation of the instance's operating system.
* `disable_monitor_service` - (Optional, Bool) Disable enhance service for monitor, it is enabled by default. When this options is set, monitor agent won't be installed. Modifications may lead to the reinstallation of the instance's operating system.
* `disable_security_service` - (Optional, Bool) Disable enhance service for security, it is enabled by default. When this options is set, security agent won't be installed. Modifications may lead to the reinstallation of the instance's operating system.
* `force_delete` - (Optional, Bool) Indicate whether to force delete the instance. Default is `false`. If set true, the instance will be permanently deleted instead of being moved into the recycle bin. Note: only works for `PREPAID` instance.
* `force_replace_placement_group_id` - (Optional, Bool) Whether to force the instance host to be replaced. Value range: true: Allows the instance to change the host and restart the instance. Local disk machines do not support specifying this parameter; false: Does not allow the instance to change the host and only join the placement group on the current host. This may cause the placement group to fail to change. Only useful for change `placement_group_id`, Default is false.
* `hostname` - (Optional, String) The hostname of the instance. Windows instance: The name should be a combination of 2 to 15 characters comprised of letters (case insensitive), numbers, and hyphens (-). Period (.) is not supported, and the name cannot be a string of pure numbers. Other types (such as Linux) of instances: The name should be a combination of 2 to 60 characters, supporting multiple periods (.). The piece between two periods is composed of letters (case insensitive), numbers, and hyphens (-). Changing the `hostname` will cause the instance system to restart.
* `hpc_cluster_id` - (Optional, String, ForceNew) High-performance computing cluster ID. If the instance created is a high-performance computing instance, you need to specify the cluster in which the instance is placed, otherwise it cannot be specified.
* `image_id` - (Optional, String) The image to use for the instance. Modifications may lead to the reinstallation of the instance's operating system.
* `instance_charge_type_prepaid_period` - (Optional, Int) The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`, `48`, `60`.
* `instance_charge_type_prepaid_renew_flag` - (Optional, String) Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.
* `instance_charge_type` - (Optional, String) The charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, `SPOTPAID`, `CDHPAID` and `CDCPAID`. The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR` and `CDHPAID`. `PREPAID` instance may not allow to delete before expired. `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time. `CDHPAID` instance must set `cdh_instance_type` and `cdh_host_id`.
* `instance_name` - (Optional, String) The name of the instance. The max length of instance_name is 128, and default value is `Terraform-CVM-Instance`.
* `instance_type` - (Optional, String) The type of the instance.
* `internet_charge_type` - (Optional, String) Internet charge type of the instance, Valid values are `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`. If not set, internet charge type are consistent with the cvm charge type by default. This value takes NO Effect when changing and does not need to be set when `allocate_public_ip` is false.
* `internet_max_bandwidth_out` - (Optional, Int) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bits per second). This value does not need to be set when `allocate_public_ip` is false.
* `ipv4_address_type` - (Optional, String, ForceNew) AddressType. Default value: WanIP. For beta users of dedicated IP. the value can be: HighQualityEIP: Dedicated IP. Note that dedicated IPs are only available in partial regions. For beta users of Anti-DDoS IP, the value can be: AntiDDoSEIP: Anti-DDoS EIP. Note that Anti-DDoS IPs are only available in partial regions.
* `ipv6_address_count` - (Optional, Int, ForceNew) Specify the number of randomly generated IPv6 addresses for the Elastic Network Interface.
* `ipv6_address_type` - (Optional, String, ForceNew) IPv6 AddressType. Default value: WanIP. EIPv6: Elastic IPv6; HighQualityEIPv6: Premium IPv6, only China Hong Kong supports premium IPv6. To allocate IPv6 addresses to resources, please specify the Elastic IPv6 type.
* `keep_image_login` - (Optional, Bool) Whether to keep image login or not, default is `false`. When the image type is private or shared or imported, this parameter can be set `true`. Modifications may lead to the reinstallation of the instance's operating system..
* `key_ids` - (Optional, Set: [`String`]) The key pair to use for the instance, it looks like `skey-16jig7tx`. Modifications may lead to the reinstallation of the instance's operating system.
* `key_name` - (Optional, String, **Deprecated**) Please use `key_ids` instead. The key pair to use for the instance, it looks like `skey-16jig7tx`. Modifications may lead to the reinstallation of the instance's operating system.
* `launch_template_id` - (Optional, String, ForceNew) Instance launch template ID. This parameter allows you to create an instance using the preset parameters in the instance template.
* `launch_template_version` - (Optional, Int, ForceNew) The instance launch template version number. If given, a new instance launch template will be created based on the given version number.
* `orderly_security_groups` - (Optional, List: [`String`]) A list of orderly security group IDs to associate with.
* `password` - (Optional, String) Password for the instance. In order for the new password to take effect, the instance will be restarted after the password change. Modifications may lead to the reinstallation of the instance's operating system.
* `placement_group_id` - (Optional, String) The ID of a placement group.
* `private_ip` - (Optional, String) The private IP to be assigned to this instance, must be in the provided subnet and available.
* `project_id` - (Optional, Int) The project the instance belongs to, default to 0.
* `release_address` - (Optional, Bool) Release elastic IP. Under EIP 2.0, only the first EIP under the primary network card is provided, and the EIP types are limited to HighQualityEIP, AntiDDoSEIP, EIPv6, and HighQualityEIPv6. Default behavior is not released.
* `running_flag` - (Optional, Bool) Set instance to running or stop. Default value is true, the instance will shutdown when this flag is false.
* `security_groups` - (Optional, Set: [`String`], **Deprecated**) It will be deprecated. Use `orderly_security_groups` instead. A list of security group IDs to associate with.
* `spot_instance_type` - (Optional, String) Type of spot instance, only support `ONE-TIME` now. Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `spot_max_price` - (Optional, String, ForceNew) Max price of a spot instance, is the format of decimal string, for example "0.50". Note: it only works when instance_charge_type is set to `SPOTPAID`.
* `stop_type` - (Optional, String) Instance shutdown mode. Valid values: SOFT_FIRST: perform a soft shutdown first, and force shut down the instance if the soft shutdown fails; HARD: force shut down the instance directly; SOFT: soft shutdown only. Default value: SOFT.
* `stopped_mode` - (Optional, String) Billing method of a pay-as-you-go instance after shutdown. Available values: `KEEP_CHARGING`,`STOP_CHARGING`. Default `KEEP_CHARGING`.
* `subnet_id` - (Optional, String) The ID of a VPC subnet. If you want to create instances in a VPC network, this parameter must be set.
* `system_disk_id` - (Optional, String) System disk snapshot ID used to initialize the system disk. When system disk type is `LOCAL_BASIC` and `LOCAL_SSD`, disk id is not supported.
* `system_disk_name` - (Optional, String) Name of the system disk.
* `system_disk_resize_online` - (Optional, Bool) Resize online.
* `system_disk_size` - (Optional, Int) Size of the system disk. unit is GB, Default is 50GB. If modified, the instance may force stop.
* `system_disk_type` - (Optional, String) System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_BASIC`: cloud disk, `CLOUD_SSD`: cloud SSD disk, `CLOUD_PREMIUM`: Premium Cloud Storage, `CLOUD_BSSD`: Basic SSD, `CLOUD_HSSD`: Enhanced SSD, `CLOUD_TSSD`: Tremendous SSD. NOTE: If modified, the instance may force stop.
* `tags` - (Optional, Map) A mapping of tags to assign to the resource. For tag limits, please refer to [Use Limits](https://intl.cloud.tencent.com/document/product/651/13354).
* `user_data_raw` - (Optional, String) The user data to be injected into this instance, in plain text. Conflicts with `user_data`. Up to 16 KB after base64 encoded. If `user_data_replace_on_change` is set to `true`, updates to this field will trigger the destruction and recreation of the CVM instance.
* `user_data_replace_on_change` - (Optional, Bool) When used in combination with `user_data` or `user_data_raw` will trigger a destroy and recreate of the CVM instance when set to `true`. Default is `false`.
* `user_data` - (Optional, String) The user data to be injected into this instance. Must be base64 encoded and up to 16 KB. If `user_data_replace_on_change` is set to `true`, updates to this field will trigger the destruction and recreation of the CVM instance.
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
* `kms_key_id` - (Optional, String, ForceNew) Optional parameters. When purchasing an encryption disk, customize the key. When this parameter is passed in, the `encrypt` parameter need be set.
* `throughput_performance` - (Optional, Int, ForceNew) Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cpu` - The number of CPU cores of the instance.
* `create_time` - Create time of the instance.
* `expired_time` - Expired time of the instance.
* `instance_status` - Current status of the instance.
* `ipv6_addresses` - IPv6 address of the instance.
* `memory` - Instance memory capacity, unit in GB.
* `os_name` - Instance os name.
* `public_ip` - Public IP of the instance.
* `public_ipv6_addresses` - The public IPv6 address to which the instance is bound.
* `uuid` - Globally unique ID of the instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `15m`) Used when creating the resource.

## Import

CVM instance can be imported using the id, e.g.

```
terraform import tencentcloud_instance.example ins-2qol3a80
```

