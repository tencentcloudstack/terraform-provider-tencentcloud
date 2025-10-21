---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instance_by_cert_id"
sidebar_current: "docs-tencentcloud-datasource-clb_instance_by_cert_id"
description: |-
  Use this data source to query detailed information of clb instance_by_cert_id
---

# tencentcloud_clb_instance_by_cert_id

Use this data source to query detailed information of clb instance_by_cert_id

## Example Usage

```hcl
data "tencentcloud_clb_instance_by_cert_id" "instance_by_cert_id" {
  cert_ids = ["3a6B5y8v"]
}
```

## Argument Reference

The following arguments are supported:

* `cert_ids` - (Required, Set: [`String`]) Server or client certificate ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cert_set` - Certificate ID and list of CLB instances associated with it.
  * `cert_id` - Certificate ID.
  * `load_balancers` - List of CLB instances associated with certificate. Note: this field may return null, indicating that no valid values can be obtained.
    * `address_i_p_version` - IP version. Valid values: ipv4, ipv6. Note: this field may return null, indicating that no valid values can be obtained.
    * `address_i_pv6` - IPv6 address of a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.
    * `anycast_zone` - Anycast CLB publishing region. For non-anycast CLB, this field returns an empty string. Note: This field may return null, indicating that no valid values can be obtained.
    * `attribute_flags` - Cluster ID.Note: This field may return null, indicating that no valid values can be obtained.
    * `backup_zone_set` - backup zone.
      * `edge_zone` - Whether the AZ is an edge zone. Values: true, false. Note: This field may return null, indicating that no valid values can be obtained.
      * `local_zone` - Whether the AZ is the LocalZone, e.g., false. Note: This field may return null, indicating that no valid values can be obtained.
      * `zone_id` - Unique AZ ID in a numeric form, such as 100001. Note: This field may return null, indicating that no valid values can be obtained.
      * `zone_name` - AZ name, such as Guangzhou Zone 1. Note: This field may return null, indicating that no valid values can be obtained.
      * `zone_region` - AZ region, e.g., ap-guangzhou. Note: This field may return null, indicating that no valid values can be obtained.
      * `zone` - Unique AZ ID in a numeric form, such as 100001. Note: This field may return null, indicating that no valid values can be obtained.
    * `charge_type` - Billing mode of CLB instance. Valid values: PREPAID (monthly subscription), POSTPAID_BY_HOUR (pay as you go). Note: this field may return null, indicating that no valid values can be obtained.
    * `cluster_ids` - Cluster ID. Note: This field may return null, indicating that no valid values can be obtained.
    * `cluster_tag` - Dedicated layer-7 tag. Note: this field may return null, indicating that no valid values can be obtained.
    * `config_id` - Custom configuration ID at the CLB instance level. Note: This field may return null, indicating that no valid values can be obtained.
    * `create_time` - CLB instance creation time. Note: This field may return null, indicating that no valid values can be obtained.
    * `domain` - Domain name of the CLB instance. It is only available for public classic CLBs. This parameter will be discontinued soon. Please use LoadBalancerDomain instead. Note: This field may return null, indicating that no valid values can be obtained.
    * `exclusive_cluster` - Private network dedicated cluster. Note: this field may return null, indicating that no valid values can be obtained.
      * `classical_cluster` - vpcgw cluster. Note: this field may return null, indicating that no valid values can be obtained.
        * `cluster_id` - Unique cluster ID.
        * `cluster_name` - Cluster name. Note: this field may return null, indicating that no valid values can be obtained.
        * `zone` - Cluster AZ, such as ap-guangzhou-1. Note: this field may return null, indicating that no valid values can be obtained.
      * `l4_clusters` - Layer-4 dedicated cluster list. Note: this field may return null, indicating that no valid values can be obtained.
        * `cluster_id` - Unique cluster ID.
        * `cluster_name` - Cluster name. Note: this field may return null, indicating that no valid values can be obtained.
        * `zone` - Cluster AZ, such as ap-guangzhou-1. Note: this field may return null, indicating that no valid values can be obtained.
      * `l7_clusters` - Layer-7 dedicated cluster list. Note: this field may return null, indicating that no valid values can be obtained.
        * `cluster_id` - Unique cluster ID.
        * `cluster_name` - Cluster name. Note: this field may return null, indicating that no valid values can be obtained.
        * `zone` - Cluster AZ, such as ap-guangzhou-1. Note: this field may return null, indicating that no valid values can be obtained.
    * `expire_time` - CLB instance expiration time, which takes effect only for prepaid instances. Note: This field may return null, indicating that no valid values can be obtained.
    * `extra_info` - Reserved field which can be ignored generally.Note: This field may return null, indicating that no valid values can be obtained.
      * `tgw_group_name` - TgwGroup name. Note: This field may return null, indicating that no valid values can be obtained.
      * `zhi_tong` - Whether to enable VIP direct connection. Note: This field may return null, indicating that no valid values can be obtained.
    * `forward` - CLB type identifier. Value range: 1 (CLB); 0 (classic CLB).
    * `health_log_set_id` - Health check logset ID of CLB CLS. Note: this field may return null, indicating that no valid values can be obtained.
    * `health_log_topic_id` - Health check log topic ID of CLB CLS. Note: this field may return null, indicating that no valid values can be obtained.
    * `ipv6_mode` - This field is meaningful only when the IP address version is ipv6. Valid values: IPv6Nat64, IPv6FullChain. Note: this field may return null, indicating that no valid values can be obtained.
    * `is_block_time` - Time blocked or unblocked. Note: this field may return null, indicating that no valid values can be obtained.
    * `is_block` - Whether VIP is blocked. Note: this field may return null, indicating that no valid values can be obtained.
    * `is_ddos` - Whether an Anti-DDoS Pro instance can be bound. Note: This field may return null, indicating that no valid values can be obtained.
    * `isolated_time` - CLB instance isolation time. Note: This field may return null, indicating that no valid values can be obtained.
    * `isolation` - 0: not isolated; 1: isolated. Note: This field may return null, indicating that no valid values can be obtained.
    * `load_balancer_domain` - Domain name of the CLB instance. Note: This field may return null, indicating that no valid values can be obtained.
    * `load_balancer_id` - CLB instance ID.
    * `load_balancer_name` - CLB instance name.
    * `load_balancer_pass_to_target` - Whether a real server opens the traffic from a CLB instance to the internet. Note: this field may return null, indicating that no valid values can be obtained.
    * `load_balancer_type` - CLB instance network type:OPEN: public network; INTERNAL: private network.
    * `load_balancer_vips` - List of VIPs of a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.
    * `local_bgp` - Whether the IP type is the local BGP.
    * `log_set_id` - Logset ID of CLB Log Service (CLS). Note: This field may return null, indicating that no valid values can be obtained.
    * `log_topic_id` - Log topic ID of CLB Log Service (CLS). Note: This field may return null, indicating that no valid values can be obtained.
    * `log` - Log information. Only the public network CLB that have HTTP or HTTPS listeners can generate logs. Note: This field may return null, indicating that no valid values can be obtained.
    * `master_zone` - Primary AZ. Note: This field may return null, indicating that no valid values can be obtained.
      * `edge_zone` - Whether the AZ is an edge zone. Values: true, false. Note: This field may return null, indicating that no valid values can be obtained.
      * `local_zone` - Whether the AZ is the LocalZone, e.g., false. Note: This field may return null, indicating that no valid values can be obtained.
      * `zone_id` - .
      * `zone_name` - AZ name, such as Guangzhou Zone 1. Note: This field may return null, indicating that no valid values can be obtained.
      * `zone_region` - AZ region, e.g., ap-guangzhou. Note: This field may return null, indicating that no valid values can be obtained.
      * `zone` - Unique AZ ID in a numeric form, such as 100001. Note: This field may return null, indicating that no valid values can be obtained.
    * `mix_ip_target` - If the layer-7 listener of an IPv6FullChain CLB instance is enabled, the CLB instance can be bound with an IPv4 and an IPv6 CVM instance simultaneously. Note: this field may return null, indicating that no valid values can be obtained.
    * `network_attributes` - CLB instance network attributes. Note: This field may return null, indicating that no valid values can be obtained.
      * `bandwidthpkg_sub_type` - Bandwidth package type, such as SINGLEISP. Note: This field may return null, indicating that no valid values can be obtained.
      * `internet_charge_type` - TRAFFIC_POSTPAID_BY_HOUR: hourly pay-as-you-go by traffic; BANDWIDTH_POSTPAID_BY_HOUR: hourly pay-as-you-go by bandwidth; BANDWIDTH_PACKAGE: billed by bandwidth package (currently, this method is supported only if the ISP is specified).
      * `internet_max_bandwidth_out` - Maximum outbound bandwidth in Mbps, which applies only to public network CLB. Value range: 0-65,535. Default value: 10.
    * `nfv_info` - Whether it is an NFV CLB instance. No returned information: no; l7nfv: yes. Note: this field may return null, indicating that no valid values can be obtained.
    * `numerical_vpc_id` - VPC ID in a numeric form. Note: This field may return null, indicating that no valid values can be obtained.
    * `open_bgp` - Protective CLB identifier. Value range: 1 (protective), 0 (non-protective). Note: This field may return null, indicating that no valid values can be obtained.
    * `prepaid_attributes` - Prepaid billing attributes of a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.
      * `period` - Cycle, indicating the number of months (reserved field). Note: This field may return null, indicating that no valid values can be obtained.
      * `renew_flag` - Renewal type. AUTO_RENEW: automatic renewal; MANUAL_RENEW: manual renewal. Note: This field may return null, indicating that no valid values can be obtained.
    * `project_id` - ID of the project to which a CLB instance belongs. 0: default project.
    * `secure_groups` - Security group of a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.
    * `sla_type` - Specification of the LCU-supported instance. Note: This field may return null, indicating that no valid values can be obtained.
    * `snat_ips` - SnatIp list after SnatPro load balancing is enabled. Note: this field may return null, indicating that no valid values can be obtained.
      * `ip` - IP address, such as 192.168.0.1.
      * `subnet_id` - Unique VPC subnet ID, such as subnet-12345678.
    * `snat_pro` - Whether to enable SnatPro. Note: this field may return null, indicating that no valid values can be obtained.
    * `snat` - SNAT is enabled for all private network classic CLB created before December 2016. Note: This field may return null, indicating that no valid values can be obtained.
    * `status_time` - Last status change time of a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.
    * `status` - CLB instance status, including:0: creating; 1: running. Note: This field may return null, indicating that no valid values can be obtained.
    * `subnet_id` - Subnet where a CLB instance resides (meaningful only for private network VPC CLB). Note: This field may return null, indicating that no valid values can be obtained.
    * `tags` - CLB instance tag information. Note: This field may return null, indicating that no valid values can be obtained.
      * `tag_key` - Tag key.
      * `tag_value` - Tag value.
    * `target_region_info` - Basic information of a backend server bound to a CLB instance. Note: This field may return null, indicating that no valid values can be obtained.
      * `region` - Region of the target, such as ap-guangzhou.
      * `vpc_id` - Network of the target, which is in the format of vpc-abcd1234 for VPC or 0 for basic network.
    * `vip_isp` - ISP to which a CLB IP address belongs. Note: This field may return null, indicating that no valid values can be obtained.
    * `vpc_id` - VPC ID Note: This field may return null, indicating that no valid values can be obtained.
    * `zones` - Availability zone of a VPC-based private network CLB instance. Note: this field may return null, indicating that no valid values can be obtained.


