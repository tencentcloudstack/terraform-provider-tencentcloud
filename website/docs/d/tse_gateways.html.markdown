---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_gateways"
sidebar_current: "docs-tencentcloud-datasource-tse_gateways"
description: |-
  Use this data source to query detailed information of tse gateways
---

# tencentcloud_tse_gateways

Use this data source to query detailed information of tse gateways

## Example Usage

```hcl
data "tencentcloud_tse_gateways" "gateways" {
  filters {
    name   = "GatewayId"
    values = ["gateway-ddbb709b"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) filter conditions, valid value:Type,Name,GatewayId,Tag,TradeType,InternetPaymode,Region.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) filter name.
* `values` - (Required, Set) filter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - gateways information.
  * `gateway_list` - gateway list.
    * `auto_renew_flag` - auto renew flag, `0`: default status, `1`: auto renew, `2`: auto not renew.
    * `create_time` - create time.
    * `cur_deadline` - expire date, for prepaid type.Note: This field may return null, indicating that a valid value is not available.
    * `description` - description of gateway.
    * `enable_cls` - whether to enable CLS log.
    * `enable_internet` - whether to open the public network of client.Note: This field may return null, indicating that a valid value is not available.
    * `engine_region` - engine region of gateway.
    * `feature_version` - product version. `TRIAL`, `STANDARD`(default value), `PROFESSIONAL`.
    * `gateway_id` - gateway ID.
    * `gateway_minor_version` - minor version of gateway.
    * `gateway_version` - gateway version. Reference value: `2.4.1`, `2.5.1`.
    * `ingress_class_name` - ingress class name.
    * `instance_port` - the port information that the instance monitors.
      * `http_port` - http port.
      * `https_port` - https port.
    * `internet_max_bandwidth_out` - public network outbound traffic bandwidth.
    * `internet_pay_mode` - trade type of internet. `BANDWIDTH`, `TRAFFIC`.
    * `isolate_time` - isolation time, used when the gateway is isolated.
    * `load_balancer_type` - load balance type of public internet.
    * `name` - gateway name.
    * `node_config` - original node config.
      * `number` - node number, 2-50.
      * `specification` - specification, 1c2g|2c4g|4c8g|8c16g.
    * `public_ip_addresses` - addresses of public internet.
    * `status` - status of gateway. May return values: `Creating`, `CreateFailed`, `Running`, `Modifying`, `UpdatingSpec`, `UpdateFailed`, `Deleting`, `DeleteFailed`, `Isolating`.
    * `tags` - tags infomation of gatewayNote: This field may return null, indicating that a valid value is not available.
      * `tag_key` - tag key.
      * `tag_value` - tag value.
    * `trade_type` - trade type. `0`: postpaid, `1`: Prepaid.
    * `type` - gateway type.
    * `vpc_config` - vpc infomation.
      * `subnet_id` - subnet ID. Assign an IP address to the engine in the VPC subnet.
      * `vpc_id` - subnet ID. Assign an IP address to the engine in the VPC subnet.
  * `total_count` - total count.


