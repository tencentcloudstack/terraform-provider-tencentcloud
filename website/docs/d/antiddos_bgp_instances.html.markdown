---
subcategory: "Anti-DDoS(antiddos)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_bgp_instances"
sidebar_current: "docs-tencentcloud-datasource-antiddos_bgp_instances"
description: |-
  Use this data source to query detailed information of AntiDDoS bgp instances
---

# tencentcloud_antiddos_bgp_instances

Use this data source to query detailed information of AntiDDoS bgp instances

## Example Usage

```hcl
data "tencentcloud_antiddos_bgp_instances" "example" {
  filter_region = "ap-guangzhou"
  filter_instance_id_list = [
    "bgp-00000fv1",
    "bgp-00000fwx",
    "bgp-00000fwy",
  ]

  filter_tag {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter_region` - (Required, String) Region.
* `filter_instance_id_list` - (Optional, Set: [`String`]) Instance ID list.
* `filter_tag` - (Optional, List) Filter by tag key and value.
* `result_output_file` - (Optional, String) Used to save results.

The `filter_tag` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bgp_instance_list` - Returns purchased Anti-DDoS package information.
  * `enterprise_package_config` - Enterprise edition Anti-DDoS package configuration.
    * `bandwidth` - Business bandwidth scale.
    * `basic_protect_bandwidth` - Basic protection bandwidth.
    * `elastic_bandwidth_flag` - Whether to enable elastic business bandwidth.
Default is false.
    * `elastic_protect_bandwidth` - Elastic bandwidth in Gbps, selectable elastic bandwidth [0,400,500,600,800,1000].
Default is 0.
    * `protect_ip_count` - Number of protected IPs.
    * `region` - Region where the Anti-DDoS package is purchased.
  * `instance_charge_prepaid` - Renewal period related.
    * `period` - Purchase duration: unit in months.
    * `renew_flag` - NOTIFY_AND_MANUAL_RENEW: Notify expiration without automatic renewal.
NOTIFY_AND_AUTO_RENEW: Notify expiration and automatically renew.
DISABLE_NOTIFY_AND_MANUAL_RENEW: No notification and no automatic renewal.
Default: Notify expiration without automatic renewal.
  * `instance_charge_type` - Payment method.
  * `instance_id` - Instance ID.
  * `package_type` - Anti-DDoS package type.
  * `standard_package_config` - Standard edition Anti-DDoS package configuration.
    * `bandwidth` - Protection business bandwidth 50Mbps.
    * `elastic_bandwidth_flag` - Whether to enable elastic protection bandwidth. true: enable 
Default is false: disable.
    * `protect_ip_count` - Number of protected IPs.
    * `region` - Region where the Anti-DDoS package is purchased.
  * `standard_plus_package_config` - Standard edition 2.0 Anti-DDoS package configuration.
    * `bandwidth` - Protection bandwidth 50Mbps.
    * `elastic_bandwidth_flag` - Whether to enable elastic business bandwidth.
true: enable
false: disable 
Default is disable.
    * `protect_count` - Protection count: TWO_TIMES: two full protections, UNLIMITED: unlimited protections.
    * `protect_ip_count` - Number of protected IPs.
    * `region` - Region where the Anti-DDoS package is purchased.
  * `tag_info_list` - Tag information.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.


