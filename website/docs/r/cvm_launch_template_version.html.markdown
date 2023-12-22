---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_launch_template_version"
sidebar_current: "docs-tencentcloud-resource-cvm_launch_template_version"
description: |-
  Provides a resource to create a cvm launch_template_version
---

# tencentcloud_cvm_launch_template_version

Provides a resource to create a cvm launch_template_version

## Example Usage

```hcl
resource "tencentcloud_cvm_launch_template_version" "foo" {
  placement {
    zone       = "ap-guangzhou-6"
    project_id = 0

  }
  launch_template_id                  = "lt-r9ajalbi"
  launch_template_version_description = "version description"
  disable_api_termination             = false
  instance_type                       = "S5.MEDIUM4"
  image_id                            = "img-9qrfy1xt"
}
```

## Argument Reference

The following arguments are supported:

* `launch_template_id` - (Required, String, ForceNew) Instance launch template ID. This parameter is used as a basis for creating new template versions.
* `placement` - (Required, List, ForceNew) Location of the instance. You can use this parameter to specify the attributes of the instance, such as its availability zone, project, and CDH (for dedicated CVMs).
* `action_timer` - (Optional, List, ForceNew) Scheduled tasks.
* `cam_role_name` - (Optional, String, ForceNew) The role name of CAM.
* `client_token` - (Optional, String, ForceNew) A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idem-potency of the request cannot be guaranteed.
* `data_disks` - (Optional, List, ForceNew) The configuration information of instance data disks. If this parameter is not specified, no data disk will be purchased by default.
* `disable_api_termination` - (Optional, Bool, ForceNew) Whether the termination protection is enabled. `TRUE`: Enable instance protection, which means that this instance can not be deleted by an API action.`FALSE`: Do not enable the instance protection. Default value: `FALSE`.
* `disaster_recover_group_ids` - (Optional, Set: [`String`], ForceNew) Placement group ID. You can only specify one.
* `dry_run` - (Optional, Bool, ForceNew) Whether the request is a dry run only.
* `enhanced_service` - (Optional, List, ForceNew) Enhanced service. You can use this parameter to specify whether to enable services such as Anti-DDoS and Cloud Monitor. If this parameter is not specified, Cloud Monitor and Anti-DDoS are enabled for public images by default.
* `host_name` - (Optional, String, ForceNew) Hostname of a CVM.
* `hpc_cluster_id` - (Optional, String, ForceNew) HPC cluster ID. The HPC cluster must and can only be specified for a high-performance computing instance.
* `image_id` - (Optional, String, ForceNew) Image ID.
* `instance_charge_prepaid` - (Optional, List, ForceNew) Describes the billing method of an instance.
* `instance_charge_type` - (Optional, String, ForceNew) The charge type of instance.
* `instance_count` - (Optional, Int, ForceNew) The number of instances to be purchased.
* `instance_market_options` - (Optional, List, ForceNew) Options related to bidding requests.
* `instance_name` - (Optional, String, ForceNew) Instance name to be displayed.
* `instance_type` - (Optional, String, ForceNew) The type of the instance. If this parameter is not specified, the system will dynamically specify the default model according to the resource sales in the current region.
* `internet_accessible` - (Optional, List, ForceNew) Describes the accessibility of an instance in the public network, including its network billing method, maximum bandwidth, etc.
* `launch_template_version_description` - (Optional, String, ForceNew) Description of instance launch template versions. This parameter can contain 2-256 characters.
* `launch_template_version` - (Optional, Int, ForceNew) This parameter, when specified, is used to create instance launch templates. If this parameter is not specified, the default version will be used.
* `login_settings` - (Optional, List, ForceNew) Describes login settings of an instance.
* `security_group_ids` - (Optional, Set: [`String`], ForceNew) Security groups to which the instance belongs. If this parameter is not specified, the instance will be associated with default security groups.
* `system_disk` - (Optional, List, ForceNew) System disk configuration information of the instance. If this parameter is not specified, it is assigned according to the system default.
* `tag_specification` - (Optional, List, ForceNew) Description of tags associated with resource instances during instance creation.
* `user_data` - (Optional, String, ForceNew) User data provided to the instance. This parameter needs to be encoded in base64 format with the maximum size of 16 KB.
* `virtual_private_cloud` - (Optional, List, ForceNew) Describes information on VPC, including subnets, IP addresses, etc.

The `action_timer` object supports the following:

* `action_time` - (Optional, String, ForceNew) Execution time, displayed according to ISO8601 standard, and UTC time is used. The format is YYYY-MM-DDThh:mm:ssZ. For example, 2018-05-29T11:26:40Z, the execution must be at least 5 minutes later than the current time.
* `externals` - (Optional, List, ForceNew) Additional data.
* `timer_action` - (Optional, String, ForceNew) Timer name. Currently TerminateInstances is the only supported value.

The `automation_service` object of `enhanced_service` supports the following:

* `enabled` - (Optional, Bool, ForceNew) Whether to enable the TAT service.

The `data_disks` object supports the following:

* `disk_size` - (Required, Int, ForceNew) Data disk size (in GB). The minimum adjustment increment is 10 GB. The value range varies by data disk type.
* `cdc_id` - (Optional, String, ForceNew) ID of the dedicated cluster to which the instance belongs.
* `delete_with_instance` - (Optional, Bool, ForceNew) Whether to terminate the data disk when its CVM is terminated. Default value: `true`.
* `disk_id` - (Optional, String, ForceNew) System disk ID. System disks whose type is LOCAL_BASIC or LOCAL_SSD do not have an ID and do not support this parameter. It is only used as a response parameter for APIs such as DescribeInstances, and cannot be used as a request parameter for APIs such as RunInstances.
* `disk_type` - (Optional, String, ForceNew) The type of data disk.
* `encrypt` - (Optional, Bool, ForceNew) Specifies whether the data disk is encrypted.
* `kms_key_id` - (Optional, String, ForceNew) ID of the custom CMK in the format of UUID or `kms-abcd1234`.
* `snapshot_id` - (Optional, String, ForceNew) Data disk snapshot ID. The size of the selected data disk snapshot must be smaller than that of the data disk. Note: This field may return null, indicating that no valid value is found.
* `throughput_performance` - (Optional, Int, ForceNew) Cloud disk performance in MB/s.

The `enhanced_service` object supports the following:

* `automation_service` - (Optional, List, ForceNew) Whether to enable the TAT service. If this parameter is not specified, the TAT service is enabled for public images and disabled for other images by default.
* `monitor_service` - (Optional, List, ForceNew) Enables cloud monitor service. If this parameter is not specified, the cloud monitor service will be enabled by default.
* `security_service` - (Optional, List, ForceNew) Enables cloud security service. If this parameter is not specified, the cloud security service will be enabled by default.

The `externals` object of `action_timer` supports the following:

* `release_address` - (Optional, Bool, ForceNew) Release address.
* `storage_block_attr` - (Optional, List, ForceNew) Information on local HDD storage.
* `unsupport_networks` - (Optional, Set, ForceNew) Not supported network.

The `instance_charge_prepaid` object supports the following:

* `period` - (Required, Int, ForceNew) Subscription period; unit: month; valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60.
* `renew_flag` - (Optional, String, ForceNew) Auto renewal flag. Valid values: NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify upon expiration nor renew automatically &lt;br&gt;&lt;br&gt;Default value: NOTIFY_AND_MANUAL_RENEW. If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis if the account balance is sufficient.

The `instance_market_options` object supports the following:

* `spot_options` - (Required, List, ForceNew) Options related to bidding.
* `market_type` - (Optional, String, ForceNew) Market option type. Currently spot is the only supported value.

The `internet_accessible` object supports the following:

* `bandwidth_package_id` - (Optional, String, ForceNew) Bandwidth package ID.
* `internet_charge_type` - (Optional, String, ForceNew) Network connection billing plan.
* `internet_max_bandwidth_out` - (Optional, Int, ForceNew) The maximum outbound bandwidth of the public network, in Mbps. The default value is 0 Mbps.
* `public_ip_assigned` - (Optional, Bool, ForceNew) Whether to assign a public IP.

The `login_settings` object supports the following:

* `keep_image_login` - (Optional, String, ForceNew) Whether to keep the original settings of an image.
* `key_ids` - (Optional, Set, ForceNew) List of key IDs. After an instance is associated with a key, you can access the instance with the private key in the key pair.
* `password` - (Optional, String, ForceNew) Login password of the instance.

The `monitor_service` object of `enhanced_service` supports the following:

* `enabled` - (Optional, Bool, ForceNew) Whether to enable Cloud Monitor.

The `placement` object supports the following:

* `zone` - (Required, String, ForceNew) ID of the availability zone where the instance resides. You can call the DescribeZones API and obtain the ID in the returned Zone field.
* `host_ids` - (Optional, Set, ForceNew) ID list of CDHs from which the instance can be created. If you have purchased CDHs and specify this parameter, the instances you purchase will be randomly deployed on the CDHs.
* `host_ips` - (Optional, Set, ForceNew) IPs of the hosts to create CVMs.
* `project_id` - (Optional, Int, ForceNew) ID of the project to which the instance belongs. This parameter can be obtained from the projectId returned by DescribeProject. If this is left empty, the default project is used.

The `security_service` object of `enhanced_service` supports the following:

* `enabled` - (Optional, Bool) Whether to enable Cloud Security.

The `spot_options` object of `instance_market_options` supports the following:

* `max_price` - (Required, String, ForceNew) Bidding price.
* `spot_instance_type` - (Optional, String, ForceNew) Bidding request type. Currently only one-time is supported.

The `storage_block_attr` object of `externals` supports the following:

* `max_size` - (Required, Int, ForceNew) Maximum capacity of local HDD storage.
* `min_size` - (Required, Int, ForceNew) Minimum capacity of local HDD storage.
* `type` - (Required, String, ForceNew) Local HDD storage type. Value: LOCAL_PRO.

The `system_disk` object supports the following:

* `cdc_id` - (Optional, String, ForceNew) ID of the dedicated cluster to which the instance belongs.
* `disk_id` - (Optional, String, ForceNew) System disk ID. System disks whose type is LOCAL_BASIC or LOCAL_SSD do not have an ID and do not support this parameter. It is only used as a response parameter for APIs such as DescribeInstances, and cannot be used as a request parameter for APIs such as RunInstances.
* `disk_size` - (Optional, Int, ForceNew) System disk size; unit: GB; default value: 50 GB.
* `disk_type` - (Optional, String, ForceNew) The type of system disk. Default value: the type of hard disk currently in stock.

The `tag_specification` object supports the following:

* `resource_type` - (Required, String, ForceNew) The type of resource that the tag is bound to.
* `tags` - (Required, List) List of tags.

The `tags` object of `tag_specification` supports the following:

* `key` - (Required, String, ForceNew) Tag key.
* `value` - (Required, String, ForceNew) Tag value.

The `virtual_private_cloud` object supports the following:

* `subnet_id` - (Required, String, ForceNew) VPC subnet ID in the format subnet-xxx, if you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.
* `vpc_id` - (Required, String, ForceNew) VPC ID in the format of vpc-xxx, if you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.
* `as_vpc_gateway` - (Optional, Bool, ForceNew) Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC.
* `ipv6_address_count` - (Optional, Int, ForceNew) Number of IPv6 addresses randomly generated for the ENI.
* `private_ip_addresses` - (Optional, Set, ForceNew) Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cvm launch_template_version can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_launch_template_version.launch_template_version ${launch_template_id}#${launch_template_version}
```

