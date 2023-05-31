---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_ckafka_zone"
sidebar_current: "docs-tencentcloud-datasource-ckafka_ckafka_zone"
description: |-
  Use this data source to query detailed information of ckafka ckafka_zone
---

# tencentcloud_ckafka_ckafka_zone

Use this data source to query detailed information of ckafka ckafka_zone

## Example Usage

```hcl
data "tencentcloud_ckafka_ckafka_zone" "ckafka_zone" {
}
```

## Argument Reference

The following arguments are supported:

* `cdc_id` - (Optional, String) cdc professional cluster business parameters.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - query result complex object entity.
  * `cluster_info` - User exclusive cluster information.
    * `available_band_width` - The current available bandwidth of the cluster in MBs.
    * `available_disk_size` - The current available disk of the cluster, in GB.
    * `cluster_id` - ClusterId.
    * `cluster_name` - ClusterName.
    * `max_band_width` - Maximum cluster bandwidth in MBs.
    * `max_disk_size` - The largest disk in the cluster, in GB.
    * `zone_id` - Availability zone to which the cluster belongs, indicating the availability zone to which the cluster belongs.
    * `zone_ids` - The availability zone where the cluster node is located. If the cluster is a cross-availability zone cluster, it includes multiple availability zones where the cluster node is located.
  * `max_bandwidth` - Maximum purchased bandwidth in Mbs.
  * `max_buy_instance_num` - The maximum number of purchased instances.
  * `message_price` - Postpaid message unit price.
    * `real_total_cost` - discount price.
    * `total_cost` - original price.
  * `physical` - Physical Exclusive Edition Configuration.
  * `profession` - Professional Edition configuration.
  * `public_network_limit` - Public network bandwidth configuration.
  * `public_network` - Public network bandwidth.
  * `standard_s2` - Standard Edition S2 configuration.
  * `standard` - Purchase Standard Edition Configuration.
  * `unit_price` - Postpaid unit price.
    * `real_total_cost` - discount price.
    * `total_cost` - original price.
  * `zone_list` - zone list.
    * `app_id` - app id.
    * `exflag` - extra flag.
    * `flag` - flag.
    * `is_internal_app` - internal APP.
    * `sales_info` - Standard Edition Sold Out Information.
      * `flag` - Manually set flags.
      * `platform` - Professional Edition, Standard Edition flag.
      * `sold_out` - sold out flag: true sold out.
      * `version` - ckakfa version(1.1.1/2.4.2/0.10.2).
    * `sold_out` - json object, key is model, value true is sold out, false is not sold out.
    * `zone_id` - zone id.
    * `zone_name` - zone name.
    * `zone_status` - zone status.


