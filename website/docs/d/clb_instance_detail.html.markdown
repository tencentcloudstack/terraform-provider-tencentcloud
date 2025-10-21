---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instance_detail"
sidebar_current: "docs-tencentcloud-datasource-clb_instance_detail"
description: |-
  Use this data source to query detailed information of clb instance_detail
---

# tencentcloud_clb_instance_detail

Use this data source to query detailed information of clb instance_detail

## Example Usage

```hcl
data "tencentcloud_clb_instance_detail" "instance_detail" {
  target_type = "NODE"
}
```

## Argument Reference

The following arguments are supported:

* `fields` - (Optional, Set: [`String`]) List of fields. Only fields specified will be returned. If it's left blank, `null` is returned. The fields `LoadBalancerId` and `LoadBalancerName` are added by default. For details about fields.
* `filters` - (Optional, List) Filter condition of querying lists describing CLB instance details:loadbalancer-id - String - Required: no - (Filter condition) CLB instance ID, such as lb-12345678; project-id - String - Required: no - (Filter condition) Project ID, such as 0 and 123; network - String - Required: no - (Filter condition) Network type of the CLB instance, such as Public and Private.&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt; vip - String - Required: no - (Filter condition) CLB instance VIP, such as 1.1.1.1 and 2204::22:3; target-ip - String - Required: no - (Filter condition) Private IP of the target real servers, such as1.1.1.1 and 2203::214:4; vpcid - String - Required: no - (Filter condition) Identifier of the VPC instance to which the CLB instance belongs, such as vpc-12345678; zone - String - Required: no - (Filter condition) Availability zone where the CLB instance resides, such as ap-guangzhou-1; tag-key - String - Required: no - (Filter condition) Tag key of the CLB instance, such as name; tag:* - String - Required: no - (Filter condition) CLB instance tag, followed by tag key after the colon. For example, use {Name: tag:name,Values: [zhangsan, lisi]} to filter the tag key `name` with the tag value `zhangsan` and `lisi`; fuzzy-search - String - Required: no - (Filter condition) Fuzzy search for CLB instance VIP and CLB instance name, such as 1.
* `result_output_file` - (Optional, String) Used to save results.
* `target_type` - (Optional, String) Target type. Valid values: NODE and GROUP. If the list of fields contains `TargetId`, `TargetAddress`, `TargetPort`, `TargetWeight` and other fields, `Target` of the target group or non-target group must be exported.

The `filters` object supports the following:

* `name` - (Required, String) Filter name.
* `values` - (Required, Set) Filter value array.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `load_balancer_detail_set` - List of CLB instance details.Note: this field may return null, indicating that no valid values can be obtained.
  * `address_ip_version` - IP version of the CLB instance. Valid values: IPv4, IPv6.Note: this field may return null, indicating that no valid values can be obtained.
  * `address_ipv6` - IPv6 VIP address of the CLB instance.Note: this field may return null, indicating that no valid values can be obtained.
  * `address_isp` - ISP to which the CLB IP address belongs.Note: this field may return null, indicating that no valid values can be obtained.
  * `address` - CLB instance VIP.Note: this field may return null, indicating that no valid values can be obtained.
  * `charge_type` - CLB instance billing mode.Note: this field may return null, indicating that no valid values can be obtained.
  * `config_id` - Custom configuration IDs of CLB instances. Multiple IDs must be separated by commas (,).Note: This field may return null, indicating that no valid values can be obtained.
  * `create_time` - CLB instance creation time.Note: this field may return null, indicating that no valid values can be obtained.
  * `domain` - Domain name of the forwarding rule.Note: this field may return null, indicating that no valid values can be obtained.
  * `domains` - List o domain names associated with the forwarding ruleNote: This field may return `null`, indicating that no valid values can be obtained.
  * `extra_info` - Reserved field, which can be ignored generally.Note: this field may return null, indicating that no valid values can be obtained.
    * `tgw_group_name` - TgwGroup nameNote: This field may return null, indicating that no valid values can be obtained.
    * `zhi_tong` - Whether to enable VIP direct connectionNote: This field may return null, indicating that no valid values can be obtained.
  * `ipv6_mode` - IPv6 address type of the CLB instance. Valid values: IPv6Nat64, IPv6FullChain.Note: this field may return null, indicating that no valid values can be obtained.
  * `isolation` - 0: not isolated; 1: isolated.Note: this field may return null, indicating that no valid values can be obtained.
  * `listener_id` - CLB listener ID.Note: this field may return null, indicating that no valid values can be obtained.
  * `load_balancer_domain` - Domain name of the CLB instance.Note: This field may return null, indicating that no valid values can be obtained.
  * `load_balancer_id` - CLB instance ID.
  * `load_balancer_name` - CLB instance name.
  * `load_balancer_pass_to_target` - Whether the CLB instance is billed by IP.Note: this field may return `null`, indicating that no valid values can be obtained.
  * `load_balancer_type` - CLB instance network type:Public: public network; Private: private network.Note: this field may return null, indicating that no valid values can be obtained.
  * `location_id` - Forwarding rule ID.Note: this field may return null, indicating that no valid values can be obtained.
  * `network_attributes` - CLB instance network attribute.Note: this field may return null, indicating that no valid values can be obtained.
    * `bandwidth_pkg_sub_type` - Bandwidth package type, such as SINGLEISPNote: This field may return null, indicating that no valid values can be obtained.
    * `internet_charge_type` - TRAFFIC_POSTPAID_BY_HOUR: hourly pay-as-you-go by traffic; BANDWIDTH_POSTPAID_BY_HOUR: hourly pay-as-you-go by bandwidth;BANDWIDTH_PACKAGE: billed by bandwidth package (currently, this method is supported only if the ISP is specified).
    * `internet_max_bandwidth_out` - Maximum outbound bandwidth in Mbps, which applies only to public network CLB. Value range: 0-65,535. Default value: 10.
  * `port` - Listener port.Note: this field may return null, indicating that no valid values can be obtained.
  * `prepaid_attributes` - Pay-as-you-go attribute of the CLB instance.Note: this field may return null, indicating that no valid values can be obtained.
    * `period` - Cycle, indicating the number of months (reserved field)Note: This field may return null, indicating that no valid values can be obtained.
    * `renew_flag` - Renewal type. AUTO_RENEW: automatic renewal; MANUAL_RENEW: manual renewalNote: This field may return null, indicating that no valid values can be obtained.
  * `project_id` - ID of the project to which the CLB instance belongs. 0: default project.Note: this field may return null, indicating that no valid values can be obtained.
  * `protocol` - Listener protocol.Note: this field may return null, indicating that no valid values can be obtained.
  * `security_group` - List of the security groups bound to the CLB instance.Note: this field may return `null`, indicating that no valid values can be obtained.
  * `slave_zone` - The secondary zone of multi-AZ CLB instanceNote: This field may return `null`, indicating that no valid values can be obtained.
  * `sni_switch` - Whether SNI is enabled. This parameter is only meaningful for HTTPS listeners.Note: This field may return `null`, indicating that no valid values can be obtained.
  * `status` - CLB instance status, including:0: creating; 1: running.Note: this field may return null, indicating that no valid values can be obtained.
  * `tags` - CLB instance tag information.Note: this field may return null, indicating that no valid values can be obtained.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `target_address` - Address of target real servers.Note: this field may return null, indicating that no valid values can be obtained.
  * `target_health` - Health status of the target real server.Note: this field may return `null`, indicating that no valid values can be obtained.
  * `target_id` - ID of target real servers.Note: this field may return null, indicating that no valid values can be obtained.
  * `target_port` - Listening port of target real servers.Note: this field may return null, indicating that no valid values can be obtained.
  * `target_weight` - Forwarding weight of target real servers.Note: this field may return null, indicating that no valid values can be obtained.
  * `url` - Forwarding rule path.Note: this field may return null, indicating that no valid values can be obtained.
  * `vpc_id` - ID of the VPC instance to which the CLB instance belongs.Note: this field may return null, indicating that no valid values can be obtained.
  * `zone` - Availability zone where the CLB instance resides.Note: this field may return null, indicating that no valid values can be obtained.
  * `zones` - The AZ of private CLB instance. This is only available for beta users.Note: This field may return `null`, indicating that no valid values can be obtained.


