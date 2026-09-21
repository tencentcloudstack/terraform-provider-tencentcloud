---
subcategory: "Backup and Disaster Recovery Center(BDRC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bdrc_instance_copy_pair"
sidebar_current: "docs-tencentcloud-resource-bdrc_instance_copy_pair"
description: |-
  Provides a resource to create a BDRC instance copy pair
---

# tencentcloud_bdrc_instance_copy_pair

Provides a resource to create a BDRC instance copy pair

## Example Usage

```hcl
resource "tencentcloud_bdrc_disaster_recovery_site_pair" "example" {
  disaster_recovery_type = "CROSS_ZONE"
  source_region          = "ap-shanghai"
  source_zone            = "ap-shanghai-3"
  target_region          = "ap-shanghai"
  target_zone            = "ap-shanghai-4"
  source_vpc             = "vpc-lx6q09ji"
  target_vpc             = "vpc-jktad5e6"
  site_pair_product_type = "INSTANCE"
  site_pair_name         = "tf-example"
  copy_type              = "ASY"
}

resource "tencentcloud_bdrc_disaster_recovery_protect_group" "example" {
  site_pair_id             = tencentcloud_bdrc_disaster_recovery_site_pair.example.site_pair_id
  protect_group_name       = "tf-example"
  data_direction           = "POSITIVE"
  protect_group_type       = "INSTANCE"
  recovery_point_objective = 15
}

resource "tencentcloud_bdrc_disaster_recovery_vpc_mapping" "example" {
  site_pair_id     = tencentcloud_bdrc_disaster_recovery_site_pair.example.site_pair_id
  source_vpc_id    = "vpc-lx6q09ji"
  source_subnet_id = "subnet-nyrg9pkl"
  target_vpc_id    = "vpc-jktad5e6"
  target_subnet_id = "subnet-jdzuvvvb"
}

resource "tencentcloud_bdrc_security_group_mapping" "example" {
  site_pair_id             = tencentcloud_bdrc_disaster_recovery_site_pair.example.site_pair_id
  src_security_group_id    = "sg-ool7tmf8"
  target_security_group_id = "sg-jfy3gi92"
}

resource "tencentcloud_bdrc_instance_copy_pair" "example" {
  protect_group_id         = tencentcloud_bdrc_disaster_recovery_protect_group.example.protect_group_id
  instance_copy_pair_name  = "tf-example"
  recovery_point_objective = 15

  create_target_instance_parameters {
    source_instance_id   = "ins-0ryh0vvn"
    instance_charge_type = "POSTPAID_BY_HOUR"
    instance_type        = "S5.MEDIUM2"
    image_id             = "img-9qrfy1xt"
    instance_name        = "tf-example-instance"

    placement {
      zone = "ap-shanghai-4"
    }

    system_disk {
      disk_type = "CLOUD_BSSD"
      disk_size = 50
    }

    virtual_private_cloud {
      vpc_id    = "vpc-jktad5e6"
      subnet_id = "subnet-jdzuvvvb"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `create_target_instance_parameters` - (Required, List, ForceNew) Target CVM creation parameters list.
* `protect_group_id` - (Required, String, ForceNew) Protect group ID.
* `delete_target_resource` - (Optional, Bool) Whether to delete the disaster-recovery site disk on destroy.
* `instance_copy_pair_name` - (Optional, String) Copy pair name.
* `recovery_point_objective` - (Optional, Int, ForceNew) User-desired RPO in minutes.

The `automation_service` object of `enhanced_service` supports the following:

* `enabled` - (Optional, Bool) Whether to enable.

The `basic_service` object of `enhanced_service` supports the following:

* `enabled` - (Optional, Bool) Whether to enable.

The `create_target_instance_parameters` object supports the following:

* `image_id` - (Required, String) Image ID.
* `instance_charge_type` - (Required, String) Instance billing mode.
* `placement` - (Required, List) Instance placement.
* `source_instance_id` - (Required, String) Source CVM ID.
* `system_disk` - (Required, List) System disk.
* `copy_pair_id` - (Optional, String) Copy pair ID for drill scenarios.
* `data_disks` - (Optional, List) Data disk list.
* `disaster_recover_group_ids` - (Optional, List) Placement group ID list.
* `enhanced_service` - (Optional, List) Enhanced service config.
* `host_name` - (Optional, String) Instance hostname.
* `instance_charge_prepaid` - (Optional, List) Prepaid billing settings.
* `instance_name` - (Optional, String) Instance display name.
* `instance_type` - (Optional, String) Instance type.
* `internet_accessible` - (Optional, List) Public bandwidth config.
* `login_settings` - (Optional, List) Login settings.
* `recovery_time` - (Optional, String) Recovery time point for drill scenarios.
* `spot_price` - (Optional, String) Spot instance max bid.
* `stopped_mode` - (Optional, String) Shutdown billing mode.
* `user_data` - (Optional, String) User data for the instance.
* `virtual_private_cloud` - (Optional, List) VPC config.

The `data_disks` object of `create_target_instance_parameters` supports the following:

* `delete_with_instance` - (Optional, Bool) Delete with instance.
* `disk_size` - (Optional, Int) Cloud disk size in GB.
* `disk_type` - (Optional, String) Cloud disk type.

The `enhanced_service` object of `create_target_instance_parameters` supports the following:

* `automation_service` - (Optional, List) Automation service.
* `basic_service` - (Optional, List) Basic service.
* `monitor_service` - (Optional, List) Monitor service.
* `security_service` - (Optional, List) Security service.

The `instance_charge_prepaid` object of `create_target_instance_parameters` supports the following:

* `period` - (Required, Int) Purchase duration in months.
* `renew_flag` - (Optional, String) Auto renewal flag.

The `internet_accessible` object of `create_target_instance_parameters` supports the following:

* `internet_charge_type` - (Optional, String) Network billing type.
* `internet_max_bandwidth_out` - (Optional, Int) Public network outbound bandwidth cap in Mbps.
* `internet_service_provider` - (Optional, String) Network service provider.
* `public_ip_assigned` - (Optional, Bool) Whether to assign a public IP.

The `login_settings` object of `create_target_instance_parameters` supports the following:

* `keep_image_login` - (Optional, String) Keep image login settings.
* `key_ids` - (Optional, List) Key ID list.
* `password` - (Optional, String) Instance login password.

The `monitor_service` object of `enhanced_service` supports the following:

* `enabled` - (Optional, Bool) Whether to enable.

The `placement` object of `create_target_instance_parameters` supports the following:

* `zone` - (Required, String) Zone ID.
* `host_id` - (Optional, String) Dedicated host ID.
* `host_ids` - (Optional, List) Dedicated host ID list.
* `project_id` - (Optional, Int) Project ID.
* `project_name` - (Optional, String) Project name.

The `security_service` object of `enhanced_service` supports the following:

* `enabled` - (Optional, Bool) Whether to enable.

The `system_disk` object of `create_target_instance_parameters` supports the following:

* `disk_size` - (Optional, Int) Cloud disk size in GB.
* `disk_type` - (Optional, String) Cloud disk type.

The `virtual_private_cloud` object of `create_target_instance_parameters` supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) VPC ID.
* `as_vpc_gateway` - (Optional, Bool) Used as public network gateway.
* `ipv6_address_count` - (Optional, Int) Number of IPv6 addresses.
* `private_ip_addresses` - (Optional, List) Private IP address list.
* `subnet_name` - (Optional, String) Subnet name.
* `vpc_name` - (Optional, String) VPC name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `account_uin` - Account Uin.
* `copy_pair_id` - Copy pair ID.
* `copy_pair_name` - Copy pair name from the describe response.
* `copy_pair_state` - Copy pair state.
* `copy_pair_type` - Copy pair type.
* `create_from` - Creation source.
* `create_time` - Creation time.
* `cvm_create_params` - CVM creation params JSON string.
* `data_direction` - Data direction.
* `deferred_create` - Whether deferred creation mode.
* `disaster_recovery_type` - Disaster recovery type.
* `disk_copy_pair_set` - Disk copy pair list for CVM.
  * `copy_pair_id` - Disk copy pair ID.
  * `copy_pair_name` - Disk copy pair name.
  * `create_time` - Creation time.
  * `source_resource_id` - Source resource ID.
  * `target_resource_id` - Target resource ID.
* `drill_group_id` - Drill group ID.
* `instance_copy_pair_id` - CVM copy pair ID.
* `instance_id` - Instance ID.
* `latest_protection_time` - Latest protection time.
* `peer_cloud_name` - Peer cloud name.
* `percent` - Replication progress percent.
* `protect_group_name` - Protect group name.
* `protection_time_set` - Protection time points.
* `rollback_percent` - Rollback progress.
* `rollbacking` - Whether rollback is in progress.
* `site_pair_id` - Site pair ID.
* `site_pair_name` - Site pair name.
* `source_region` - Source region.
* `source_resource_id` - Source resource ID.
* `source_vpc` - Source VPC.
* `source_zone` - Source zone.
* `sub_account_uin` - Sub-account Uin.
* `target_cvm_created` - Whether target CVM is actually created.
* `target_region` - Target region.
* `target_resource_id` - Target resource ID.
* `target_vpc` - Target VPC.
* `target_zone` - Target zone.


