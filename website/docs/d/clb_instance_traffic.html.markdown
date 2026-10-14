---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_instance_traffic"
sidebar_current: "docs-tencentcloud-datasource-clb_instance_traffic"
description: |-
  Use this data source to query detailed information of clb instance_traffic
---

# tencentcloud_clb_instance_traffic

Use this data source to query detailed information of clb instance_traffic

## Example Usage

```hcl
data "tencentcloud_clb_instance_traffic" "instance_traffic" {
  load_balancer_region = "ap-guangzhou"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_region` - (Optional, String) CLB instance region. If this parameter is not passed in, CLB instances in all regions will be returned.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `load_balancer_traffic` - Information of CLB instances sorted by outbound bandwidth from highest to lowest. Note: This field may return null, indicating that no valid values can be obtained.
  * `domain` - CLB domain name. Note: This field may return null, indicating that no valid values can be obtained.
  * `load_balancer_id` - CLB instance ID.
  * `load_balancer_name` - CLB instance name.
  * `out_bandwidth` - Maximum outbound bandwidth in Mbps.
  * `region` - CLB instance region.
  * `vip` - CLB instance VIP.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeLoadBalancerTraffic](https://cloud.tencent.com/document/api/214/50774) - This API is used to query CLB instances with high traffic under the current account, and return the top 10 results. For queries using a sub-account, only the CLB instances authorized to the sub-account will be returned.

