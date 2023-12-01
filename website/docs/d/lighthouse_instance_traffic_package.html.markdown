---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_instance_traffic_package"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_instance_traffic_package"
description: |-
  Use this data source to query detailed information of lighthouse instance_traffic_package
---

# tencentcloud_lighthouse_instance_traffic_package

Use this data source to query detailed information of lighthouse instance_traffic_package

## Example Usage

```hcl
data "tencentcloud_lighthouse_instance_traffic_package" "instance_traffic_package" {
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Optional, Set: [`String`]) Instance ID list.
* `limit` - (Optional, Int) Number of returned results. Default value is 20. Maximum value is 100.
* `offset` - (Optional, Int) Offset. Default value is 0.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_traffic_package_set` - List of details of instance traffic packages.
  * `instance_id` - Instance ID.
  * `traffic_package_set` - List of traffic package details.
    * `deadline` - The expiration time of the traffic package. Expressed according to the ISO8601 standard, and using UTC time. The format is YYYY-MM-DDThh:mm:ssZ..
    * `end_time` - The end time of the effective period of the traffic packet. Expressed according to the ISO8601 standard, and using UTC time. The format is YYYY-MM-DDThh:mm:ssZ.
    * `start_time` - The start time of the effective cycle of the traffic packet. Expressed according to the ISO8601 standard, and using UTC time. The format is YYYY-MM-DDThh:mm:ssZ.
    * `status` - Traffic packet status:- `NETWORK_NORMAL`: normal.- `OVERDUE_NETWORK_DISABLED`: network disconnection due to arrears.
    * `traffic_overflow` - The amount of traffic that exceeds the quota of the traffic packet during the effective period of the traffic packet, in bytes.
    * `traffic_package_id` - Traffic packet ID.
    * `traffic_package_remaining` - The remaining traffic during the effective period of the traffic packet, in bytes.
    * `traffic_package_total` - The total traffic in bytes during the effective period of the traffic packet.
    * `traffic_used` - Traffic has been used during the effective period of the traffic packet, in bytes.


