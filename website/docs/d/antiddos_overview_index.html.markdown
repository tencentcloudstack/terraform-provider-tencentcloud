---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_overview_index"
sidebar_current: "docs-tencentcloud-datasource-antiddos_overview_index"
description: |-
  Use this data source to query detailed information of antiddos overview index
---

# tencentcloud_antiddos_overview_index

Use this data source to query detailed information of antiddos overview index

## Example Usage

```hcl
data "tencentcloud_antiddos_overview_index" "overview_index" {
  start_time = "2023-11-20 12:32:12"
  end_time   = "2023-11-21 12:32:12"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) EndTime.
* `start_time` - (Required, String) StartTime.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `all_ip_count` - ip count.
* `antiddos_domain_count` - AntiddosDomainCount.
* `antiddos_ip_count` - Total number of advanced defense IPs (including advanced defense packets and advanced defense IPs).
* `attack_domain_count` - AttackDomainCount.
* `attack_ip_count` - AttackIpCount.
* `block_ip_count` - BlockIpCount.
* `max_attack_flow` - MaxAttackFlow.
* `new_attack_ip` - The IP address in the most recent attack.
* `new_attack_time` - The time in the most recent attack.
* `new_attack_type` - The type in the most recent attack.


