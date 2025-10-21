---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_gateway_services"
sidebar_current: "docs-tencentcloud-datasource-tse_gateway_services"
description: |-
  Use this data source to query detailed information of tse gateway_services
---

# tencentcloud_tse_gateway_services

Use this data source to query detailed information of tse gateway_services

## Example Usage

```hcl
data "tencentcloud_tse_gateway_services" "gateway_services" {
  gateway_id = "gateway-ddbb709b"
  filters {
    key   = "name"
    value = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String) gateway ID.
* `filters` - (Optional, List) filter conditions, valid value:name,upstreamType.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `key` - (Optional, String) filter name.
* `value` - (Optional, String) filter value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - result.
  * `service_list` - service list.
    * `created_time` - created time.
    * `editable` - editable status.
    * `id` - service ID.
    * `name` - service name.
    * `tags` - tag list.
    * `upstream_info` - upstream information.
      * `algorithm` - load balance algorithm,default:round-robin,least-connections and consisten_hashing also support.
      * `auto_scaling_cvm_port` - auto scaling group port of cvm.
      * `auto_scaling_group_id` - auto scaling group ID of cvm.
      * `auto_scaling_hook_status` - hook status in auto scaling group of cvm.
      * `auto_scaling_tat_cmd_status` - tat cmd status in auto scaling group of cvm.
      * `host` - an IP address or domain name.
      * `namespace` - namespace.
      * `port` - port.
      * `real_source_type` - exact source service type.
      * `scf_lambda_name` - scf lambda name.
      * `scf_lambda_qualifier` - scf lambda version.
      * `scf_namespace` - scf lambda namespace.
      * `scf_type` - scf lambda type.
      * `service_name` - the name of the service in registry or kubernetes.
      * `slow_start` - slow start time, unit:second,when it&#39;s enabled, weight of the node is increased from 1 to the target value gradually.
      * `source_id` - service source ID.
      * `source_name` - the name of source service.
      * `source_type` - source service type.
      * `targets` - provided when service type is IPList.
        * `created_time` - created time.
        * `health` - health.
        * `host` - Host.
        * `port` - port.
        * `source` - source of target.
        * `weight` - weight.
    * `upstream_type` - service type.
  * `total_count` - total count.


