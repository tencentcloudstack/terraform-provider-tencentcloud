---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instances"
sidebar_current: "docs-tencentcloud-datasource-clb_instances"
description: |-
  Use this data source to query detailed information of CLB
---

# tencentcloud_clb_instances

Use this data source to query detailed information of CLB

## Example Usage

```hcl
data "tencentcloud_clb_instances" "example" {
  clb_id             = "lb-k2zjp9lv"
  network_type       = "OPEN"
  clb_name           = "tf-example"
  project_id         = 0
  result_output_file = "myOutputPath"
}

# Parse JSON fields
output "exclusive_cluster_info" {
  value = jsondecode(data.tencentcloud_clb_instances.example.clb_list[0].exclusive_cluster)
}
```

## Argument Reference

The following arguments are supported:

* `clb_id` - (Optional, String) ID of the CLB to be queried.
* `clb_name` - (Optional, String) Name of the CLB to be queried.
* `master_zone` - (Optional, String) Master available zone id.
* `network_type` - (Optional, String) Type of CLB instance, and available values include `OPEN` and `INTERNAL`.
* `project_id` - (Optional, Int) Project ID of the CLB.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `clb_list` - A list of cloud load balancers. Each element contains the following attributes:
  * `address_ip_version` - IP version, only applicable to open CLB. Valid values are `IPV4`, `IPV6` and `IPv6FullChain`.
  * `address_ipv6` - IPv6 address of the CLB instance.
  * `anycast_zone` - Anycast CLB publishing region, returns empty string for non-anycast CLB.
  * `associate_endpoint` - Endpoint ID associated with the CLB instance.
  * `attribute_flags` - CLB attribute flags array.
  * `available_zone_affinity_info` - Available zone forwarding affinity information (JSON format).
  * `backup_zone_set` - Backup zone list, each element contains zone_id/zone/zone_name/zone_region/local_zone.
    * `local_zone` - Whether this available zone is local zone.
    * `zone_id` - Available zone unique id (numerical representation).
    * `zone_name` - Available zone name.
    * `zone_region` - Region that this available zone belongs to.
    * `zone` - Available zone unique id (string representation).
  * `charge_type` - Billing type, PREPAID: Prepaid, POSTPAID_BY_HOUR: Pay-as-you-go.
  * `clb_id` - ID of CLB.
  * `clb_name` - Name of CLB.
  * `clb_vips` - The virtual service address table of the CLB.
  * `cluster_id` - ID of the cluster.
  * `cluster_ids` - Cluster ID list.
  * `cluster_tag` - Layer-7 exclusive tag.
  * `config_id` - CLB dimension personalized configuration ID.
  * `create_time` - Create time of the CLB.
  * `domain` - CLB domain (only for public network Classic CLB), gradually deprecated.
  * `egress` - Network egress.
  * `exclusive_cluster` - Internal exclusive cluster information (JSON format).
  * `exclusive` - Whether the instance type is exclusive, 1: Exclusive, 0: Not exclusive.
  * `expire_time` - Expiration time of the CLB instance, only for prepaid CLB, format: YYYY-MM-DD HH:mm:ss.
  * `extra_info` - Reserved field, generally no need to pay attention (JSON format).
  * `forward` - CLB type identifier, 1: CLB, 0: Classic CLB.
  * `health_log_set_id` - Log service (CLS) health check log set ID.
  * `health_log_topic_id` - Log service (CLS) health check log topic ID.
  * `internet_bandwidth_max_out` - Max bandwidth out, only applicable to open CLB. Valid value ranges is [1, 2048]. Unit is MB.
  * `internet_charge_type` - Internet charge type, only applicable to open CLB. Valid values are `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.
  * `ipv6_mode` - IPv6 mode when IP version is ipv6, IPv6Nat64 or IPv6FullChain.
  * `is_block_time` - Time of blocking or unblocking, format: YYYY-MM-DD HH:mm:ss.
  * `is_block` - Whether the VIP is blocked.
  * `is_ddos` - Whether Anti-DDoS Pro can be bound.
  * `isolated_time` - Time when the CLB instance was isolated, format: YYYY-MM-DD HH:mm:ss.
  * `isolation` - Whether isolated, 0: Not isolated, 1: Isolated.
  * `load_balancer_domain` - Domain of the CLB instance.
  * `load_balancer_pass_to_target` - Whether backend services allow traffic from CLB.
  * `local_bgp` - Whether the IP type is local BGP.
  * `local_zone` - Whether this available zone is local zone, This field maybe null, means cannot get a valid value.
  * `log_set_id` - Log service (CLS) log set ID.
  * `log_topic_id` - Log service (CLS) log topic ID.
  * `mix_ip_target` - IPv6FullChain CLB layer-7 listener supports mixed binding of IPv4/IPv6 targets.
  * `network_type` - Types of CLB.
  * `nfv_info` - Whether CLB is NFV, empty: No, l7nfv: Layer-7 is NFV.
  * `numerical_vpc_id` - VPC ID in a numeric form. Note: This field may return null, indicating that no valid values can be obtained.
  * `open_bgp` - Anti-DDoS Pro LB identifier, 1: Anti-DDoS Pro, 0: Not Anti-DDoS Pro.
  * `prepaid_period` - Prepaid purchase period, unit: month.
  * `prepaid_renew_flag` - Prepaid renewal flag, NOTIFY_AND_AUTO_RENEW: Notify and auto-renew, NOTIFY_AND_MANUAL_RENEW: Notify but not auto-renew, DISABLE_NOTIFY_AND_MANUAL_RENEW: No notification and not auto-renew.
  * `project_id` - ID of the project.
  * `security_groups` - ID set of the security groups.
  * `sla_type` - Performance capacity type specification (clb.c1.small/clb.c2.medium/clb.c3.small/clb.c3.medium/clb.c4.small/clb.c4.medium/clb.c4.large/clb.c4.xlarge or empty string).
  * `snat_ips` - SnatIp list after enabling SnatPro (JSON format).
  * `snat_pro` - Whether SnatPro is enabled.
  * `snat` - Whether SNAT is enabled.
  * `status_time` - Latest state transition time of CLB.
  * `status` - The status of CLB.
  * `subnet_id` - ID of the subnet.
  * `tags` - The available tags within this CLB.
  * `target_count` - Number of bound backend services.
  * `target_region_info_region` - Region information of backend service are attached the CLB.
  * `target_region_info_vpc_id` - VpcId information of backend service are attached the CLB.
  * `vip_isp` - Network operator, only applicable to open CLB. Valid values are `CMCC`(China Mobile), `CTCC`(Telecom), `CUCC`(China Unicom) and `BGP`. If this ISP is specified, network billing method can only use the bandwidth package billing (BANDWIDTH_PACKAGE).
  * `vpc_id` - ID of the VPC.
  * `zone_id` - Available zone unique id(numerical representation), This field maybe null, means cannot get a valid value.
  * `zone_name` - Available zone name, This field maybe null, means cannot get a valid value.
  * `zone_region` - Region that this available zone belong to, This field maybe null, means cannot get a valid value.
  * `zone` - Available zone unique id(string representation), This field maybe null, means cannot get a valid value.
  * `zones` - Zones where rules are deployed for VPC internal load balancers with nearby access mode. Note: This field may return null, indicating no valid values can be obtained.


