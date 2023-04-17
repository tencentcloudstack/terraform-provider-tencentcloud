---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_chc_hosts"
sidebar_current: "docs-tencentcloud-datasource-cvm_chc_hosts"
description: |-
  Use this data source to query detailed information of cvm chc_hosts
---

# tencentcloud_cvm_chc_hosts

Use this data source to query detailed information of cvm chc_hosts

## Example Usage

```hcl
data "tencentcloud_cvm_chc_hosts" "chc_hosts" {
  chc_ids = ["chc-xxxxxx"]
  filters {
    name   = "zone"
    values = ["ap-guangzhou-7"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `chc_ids` - (Optional, Set: [`String`]) CHC host ID. Up to 100 instances per request is allowed. ChcIds and Filters cannot be specified at the same time.
* `filters` - (Optional, List) - `zone` Filter by the availability zone, such as ap-guangzhou-1. Valid values: See [Regions and Availability Zones](https://www.tencentcloud.com/document/product/213/6091?from_cn_redirect=1).
- `instance-name` Filter by the instance name.
- `instance-state` Filter by the instance status. For status details, see [InstanceStatus](https://www.tencentcloud.com/document/api/213/15753?from_cn_redirect=1#InstanceStatus).
- `device-type` Filter by the device type.
- `vpc-id` Filter by the unique VPC ID.
- `subnet-id` Filter by the unique VPC subnet ID.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `chc_host_set` - List of returned instances.
  * `bmc_ip` - Out-of-band network IPNote: This field may return null, indicating that no valid values can be obtained.
  * `bmc_mac` - MAC address assigned under the out-of-band networkNote: This field may return null, indicating that no valid values can be obtained.
  * `bmc_security_group_ids` - Out-of-band network security group IDNote: This field may return null, indicating that no valid values can be obtained.
  * `bmc_virtual_private_cloud` - Out-of-band networkNote: This field may return null, indicating that no valid values can be obtained.
    * `as_vpc_gateway` - Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: no&lt;br&gt;&lt;br&gt;Default: FALSE.
    * `ipv6_address_count` - Number of IPv6 addresses randomly generated for the ENI.
    * `private_ip_addresses` - Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.
    * `subnet_id` - VPC subnet ID in the format subnet-xxx. To obtain valid subnet IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call DescribeSubnets and look for the unSubnetId fields in the response. If you specify DEFAULT for both SubnetId and VpcId when creating an instance, the default VPC will be used.
    * `vpc_id` - VPC ID in the format of vpc-xxx. To obtain valid VPC IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call the DescribeVpcEx API and look for the unVpcId fields in the response. If you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.
  * `chc_id` - CHC host ID.
  * `cpu` - CPU cores of the CHC hostNote: This field may return null, indicating that no valid values can be obtained.
  * `created_time` - Server creation time.
  * `cvm_instance_id` - ID of the associated CVMNote: This field may return null, indicating that no valid values can be obtained.
  * `deploy_ip` - Deployment network IPNote: This field may return null, indicating that no valid values can be obtained.
  * `deploy_mac` - MAC address assigned under the deployment networkNote: This field may return null, indicating that no valid values can be obtained.
  * `deploy_security_group_ids` - Deployment network security group IDNote: This field may return null, indicating that no valid values can be obtained.
  * `deploy_virtual_private_cloud` - Deployment networkNote: This field may return null, indicating that no valid values can be obtained.
    * `as_vpc_gateway` - Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC. Valid values:&lt;br&gt;&lt;li&gt;TRUE: yes;&lt;br&gt;&lt;li&gt;FALSE: no&lt;br&gt;&lt;br&gt;Default: FALSE.
    * `ipv6_address_count` - Number of IPv6 addresses randomly generated for the ENI.
    * `private_ip_addresses` - Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.
    * `subnet_id` - VPC subnet ID in the format subnet-xxx. To obtain valid subnet IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call DescribeSubnets and look for the unSubnetId fields in the response. If you specify DEFAULT for both SubnetId and VpcId when creating an instance, the default VPC will be used.
    * `vpc_id` - VPC ID in the format of vpc-xxx. To obtain valid VPC IDs, you can log in to the [console](https://console.tencentcloud.com/vpc/vpc?rid=1) or call the DescribeVpcEx API and look for the unVpcId fields in the response. If you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.
  * `device_type` - Device typeNote: This field may return null, indicating that no valid values can be obtained.
  * `disk` - Disk capacity of the CHC hostNote: This field may return null, indicating that no valid values can be obtained.
  * `hardware_description` - Instance hardware description, including CPU cores, memory capacity and disk capacity.Note: This field may return null, indicating that no valid values can be obtained.
  * `instance_name` - Instance name.
  * `instance_state` - CHC host status&lt;br/&gt;&lt;ul&gt;&lt;li&gt;REGISTERED: The CHC host is registered, but the out-of-band network and deployment network are not configured.&lt;/li&gt;&lt;li&gt;VPC_READY: The out-of-band network and deployment network are configured.&lt;/li&gt;&lt;li&gt;PREPARED: It&#39;s ready and can be associated with a CVM.&lt;/li&gt;&lt;li&gt;ONLINE: It&#39;s already associated with a CVM.&lt;/li&gt;&lt;/ul&gt;.
  * `memory` - Memory capacity of the CHC host (unit: GB)Note: This field may return null, indicating that no valid values can be obtained.
  * `placement` - Availability zone.
    * `host_id` - The ID of the CDH to which the instance belongs, only used as an output parameter.
    * `host_ids` - ID list of CDHs from which the instance can be created. If you have purchased CDHs and specify this parameter, the instances you purchase will be randomly deployed on the CDHs.
    * `host_ips` - IPs of the hosts to create CVMs.
    * `project_id` - ID of the project to which the instance belongs. This parameter can be obtained from the projectId returned by DescribeProject. If this is left empty, the default project is used.
    * `zone` - ID of the availability zone where the instance resides. You can call the [DescribeZones](https://www.tencentcloud.com/document/product/213/35071) API and obtain the ID in the returned Zone field.
  * `serial_number` - Server serial number.
  * `tenant_type` - Management typeHOSTING: HostingTENANT: LeasingNote: This field may return null, indicating that no valid values can be obtained.


