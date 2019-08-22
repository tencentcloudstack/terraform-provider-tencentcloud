---
layout: "tencentcloud"
page_title: "tencentcloud: tencentcloud_instance"
sidebar_current: "docs-tencentcloud-resource-cvm-instance"
description: |-
  Provides a CVM instance resource.
---

# tencentcloud_instance

Provides a CVM instance resource.

~> **NOTE:** You can launch an CVM instance for a VPC network via specifying parameter `vpc_id`. One instance can only belong to one VPC.

~> **NOTE:** At present, 'PREPAID' instance cannot be deleted and must wait it to be outdated and released automatically.

## Example Usage

```hcl
data "tencentcloud_image" "my_favorate_image" {
  os_name = "ubuntu"

  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S4"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "my_favorate_zones" {}

// Create Security Group with 2 rules
resource "tencentcloud_security_group" "app" {
  name        = "web accessibility"
  description = "make it accessible for both production and stage ports"
}

resource "tencentcloud_security_group_rule" "web" {
  security_group_id = "${tencentcloud_security_group.app.id}"
  type              = "ingress"
  cidr_ip           = "0.0.0.0/0"
  ip_protocol       = "tcp"
  port_range        = "80,3000,8080"
  policy            = "accept"
}

resource "tencentcloud_security_group_rule" "ssh" {
  security_group_id = "${tencentcloud_security_group.app.id}"
  type              = "ingress"
  cidr_ip           = "202.119.230.10/32"
  ip_protocol       = "tcp"
  port_range        = "22"
  policy            = "accept"
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

// Create 10 CVM instances to host awesome_app
resource "tencentcloud_instance" "my_awesome_app" {
  instance_name     = "awesome_app"
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type     = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  key_name          = "${tencentcloud_key_pair.random_key.id}"
  hostname          = "awesome_app"
  project_id        = 0

  tags = {
    tagKey = "tagValue"
  }

  security_groups = [
    "${tencentcloud_security_group.app.id}",
  ]

  vpc_id    = "${tencentcloud_vpc.app.id}"
  subnet_id = "${tencentcloud_subnet.app.id}"

  internet_max_bandwidth_out = 20
  count                      = 10
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required,ForceNew) The Image to use for the instance. Change 'image_id' will case instance destroy and re-created.

* `availability_zone` - (Required) The available zone that the CVM instance locates at.

* `instance_name` - (Optional) The name of the CVM. This instance_name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://. If not specified, Terraform will autogenerate a default name is `CVM-Instance`.

* `instance_type` - (Required) The type of instance to start.

* `hostname` - (Optional) The hostname of CVM.

* `project_id` - (Optional) The project CVM belongs to, default to 0.

* `instance_charge_type` - (Optional) Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, The default is `POSTPAID_BY_HOUR`.

* `instance_charge_type_prepaid_period` - (Optional) The tenancy (time unit is month) of the prepaid instance, **NOTE**: it only works when `instance_charge_type` is set to `PREPAID`.

* `instance_charge_type_prepaid_renew_flag` - (Optional) Auto renewal flag. Value range:
NOTIFY_AND_AUTO_RENEW: notify expiry and renew automatically
NOTIFY_AND_MANUAL_RENEW: notify expiry but not renew automatically
DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify expiry nor renew automatically

If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis when the account balance is sufficient. **NOTE**: it only works when `instance_charge_type` is set to `PREPAID`.

* `internet_charge_type` - (Optional) Internet charge type of the instance, Valid values are `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`. Default is `TRAFFIC_POSTPAID_BY_HOUR`.

* `internet_max_bandwidth_out` - (Optional) Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bit per second). Value range:  [0, 200], If this value is not specified, then automatically sets it to 0 Mbps.

* `allocate_public_ip` - (Optional) Associate a public ip address with an instance in a VPC or Classic. Boolean value, Default is false.

* `vpc_id` - (Optional) The id of a VPC network. If you want to create instances in VPC network, this parameter must be set.

* `subnet_id` - (Optional) The id of a VPC subnetwork. If you want to create instances in VPC network, this parameter must be set.

* `private_ip` - (Optional) The private ip to be assigned to this instance, must be in the provided subnet and available.

* `security_groups` - (Optional)  A list of security group ids to associate with.

* `system_disk_type` - (Optional) Valid values are `LOCAL_BASIC`, `LOCAL_SSD`,  `CLOUD_BASIC`, `CLOUD_SSD` and `CLOUD_PREMIUM`. **NOTE**: LOCAL_BASIC and LOCAL_SSD are deprecated.

* `system_disk_size` - (Optional) Size of the system disk, value range: 50GB ~ 1TB. Default is 50GB.

* `data_disks` - (Optional) Settings for data disk. In each disk, `data_disk_type` indicates the disk type, valid values are `LOCAL_BASIC`, `LOCAL_SSD`,  `CLOUD_BASIC`, `CLOUD_SSD` and `CLOUD_PREMIUM`. **NOTE**, it must follow the system_disk_type, and all disks must be the same type. `data_disk_size` is the size of the data disk, value range: 100GB~1.6TB.
`data_disk_size` is the size of the data disk, value range: 60GB~1.6TB. `delete_with_instance` decides whether the disk is deleted with instance(only applied to `POSTPAID_BY_HOUR` cloud disk), default to true.

* `disable_security_service` - (Optional) Disable enhance service for security, it is enabled by default. When this options is set, security agent won't be installed.

* `disable_monitor_service` - (Optional) Disable enhance service for monitor, it is enabled by default. When this options is set, monitor agent won't be installed.

* `key_name` - (Optional) The key pair to use for the instance, it looks like `skey-16jig7tx`.

* `password` - (Optional) Password to an instance. In order to take effect new password, the instance will be restarted after modifying the password.

* `user_data` - (Optional) The user data to be specified into this instance. Must be encrypted in base64 format and limited in 16 KB.

* `user_data_raw` - (Optional) The user data to be specified into this instance, plain text. Conflicts with `user_data`. Limited in 16 KB after encrypted in base64 format.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The instance ID, something looks like `ins-xxxxxx`.
* `instance_status` - The Status of the instance.
* `private_ip` - The Local IP Address of the instance.
* `public_ip` - The instance public ip.
* `vpc_id` - The VPC Id associated with the instance.
* `subnet_id` - The Subnet Id associated with the instance.
* `system_disk_type` - The system disk type on the instance.
* `system_disk_size` - The system disk type on the instance.
* `data_disks` - The data disks info. In each data disk, `data_disk_type` is the disk type. `data_disk_size` is the size of the disk.
* `key_name` - The key pair id of the instance.

## Import

CVM instance can be imported using the id, e.g.

```
terraform import tencentcloud_instance.foo ins-2qol3a80
```
