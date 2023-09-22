---
subcategory: "Waf"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_waf_domains"
sidebar_current: "docs-tencentcloud-datasource-waf_domains"
description: |-
  Use this data source to query detailed information of waf domains
---

# tencentcloud_waf_domains

Use this data source to query detailed information of waf domains

## Example Usage

### Find all domains

```hcl
data "tencentcloud_waf_domains" "example" {}
```

### Find domains by filter

```hcl
data "tencentcloud_waf_domains" "example" {
  instance_id = "waf_2kxtlbky01b3wceb"
  domain      = "tf.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Optional, String) Domain name.
* `instance_id` - (Optional, String) Unique ID of Instance.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domains` - Domain info list.
  * `alb_type` - Traffic Source: clb represents Tencent Cloud clb, apisix represents apisix gateway, tsegw represents Tencent Cloud API gateway, default clbNote: This field may return null, indicating that a valid value cannot be obtained.
  * `api_status` - API security switch status, 0 off, 1 onNote: This field may return null, indicating that a valid value cannot be obtained.
  * `app_id` - User appid.
  * `bot_status` - BOT switch status, 0 off, 1 on.
  * `cc_list` - Waf sandbox export addresses, should be added to the whitelist by the upstreams.
  * `cdc_clusters` - Cdc clustersNote: This field may return null, indicating that a valid value cannot be obtained.
  * `cls_status` - Whether to enable access logs, 1 enable, 0 disable.
  * `cname` - Cname address, used for dns access.
  * `create_time` - Create time.
  * `domain_id` - Domain unique ID.
  * `domain` - Domain name.
  * `edition` - Instance type, sparta-waf represents SAAS WAF, clb-waf represents CLB WAF.
  * `engine` - Rule and AI Defense Mode, 10 Rule Engine Observation&amp;amp;&amp;amp;AI Engine Shutdown Mode 11 Rule Engine Observation&amp;amp;&amp;amp;AI Engine Observation Mode 12 Rule Engine Observation&amp;amp;&amp;amp;AI Engine Interception Mode 20 Rule Engine Interception&amp;amp;&amp;amp;AI Engine Shutdown Mode 21 Rule Engine Interception&amp;amp;&amp;amp;AI Engine Observation Mode 22 Rule Engine Interception&amp;amp;&amp;amp;AI Engine Interception Mode.
  * `flow_mode` - CLBWAF traffic mode, 1 cleaning mode, 0 mirroring mode.
  * `instance_id` - Instance unique ID.
  * `instance_name` - Instance name.
  * `ipv6_status` - Ipv6 switch status, 0 off, 1 on.
  * `level` - Instance level.
  * `load_balancer_set` - List of bound LB.
    * `listener_id` - Listener unique IDNote: This field may return null, indicating that a valid value cannot be obtained.
    * `listener_name` - Listener nameNote: This field may return null, indicating that a valid value cannot be obtained.
    * `load_balancer_id` - LoadBalancer IDNote: This field may return null, indicating that a valid value cannot be obtained.
    * `load_balancer_name` - LoadBalancer nameNote: This field may return null, indicating that a valid value cannot be obtained.
    * `load_balancer_type` - Loadbalancer typeNote: This field may return null, indicating that a valid value cannot be obtained.
    * `numerical_vpc_id` - VPCID for load balancer, public network is -1, and internal network is filled in according to actual conditionsNote: This field may return null, indicating that a valid value cannot be obtained.
    * `protocol` - Listener protocolNote: This field may return null, indicating that a valid value cannot be obtained.
    * `region` - RegionNote: This field may return null, indicating that a valid value cannot be obtained.
    * `vip` - LoadBalancer ipNote: This field may return null, indicating that a valid value cannot be obtained.
    * `vport` - Listener portNote: This field may return null, indicating that a valid value cannot be obtained.
    * `zone` - Loadbalancer zoneNote: This field may return null, indicating that a valid value cannot be obtained.
  * `mode` - Rule defense mode, 0 observation mode, 1 interception mode.
  * `ports` - Listening ports.
    * `nginx_server_id` - Nginx server ID.
    * `port` - Listening port.
    * `protocol` - The listening protocol of listening port.
    * `upstream_port` - The upstream port for listening port.
    * `upstream_protocol` - The upstream protocol for listening port.
  * `post_ckafka_status` - Whether to enable the delivery of CKafka function, 0 off, 1 on.
  * `post_cls_status` - Whether to enable the delivery CLS function, 0 off, 1 on.
  * `region` - Region.
  * `rs_list` - Waf engine export addresses, should be added to the whitelist by the upstreams.
  * `sg_detail` - Detailed explanation of security group statusNote: This field may return null, indicating that a valid value cannot be obtained.
  * `sg_state` - Security group status, 0 does not display, 1 non Tencent cloud source site, 2 security group binding failed, 3 security group changedNote: This field may return null, indicating that a valid value cannot be obtained.
  * `state` - Clbwaf domain name listener status, 0 operation successful, 4 binding LB, 6 unbinding LB, 7 unbinding LB failed, 8 binding LB failed, 10 internal error.
  * `status` - Waf switch,0 off 1 on.


