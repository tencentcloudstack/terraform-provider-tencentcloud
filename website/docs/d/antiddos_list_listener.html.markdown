---
subcategory: "Anti-DDoS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_list_listener"
sidebar_current: "docs-tencentcloud-datasource-antiddos_list_listener"
description: |-
  Use this data source to query detailed information of antiddos list_listener
---

# tencentcloud_antiddos_list_listener

Use this data source to query detailed information of antiddos list_listener

## Example Usage

```hcl
data "tencentcloud_antiddos_list_listener" "list_listener" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `layer4_listeners` - L4 listener list.
  * `backend_port` - Origin port, value 1~65535.
  * `frontend_port` - Forwarding port, value 1~65535.
  * `instance_detail_rule` - Resource instance to which the rule belongs.
    * `cname` - Instance cname.
    * `eip_list` - Resource instance ip.
    * `instance_id` - InstanceId.
  * `instance_details` - Resource instance.
    * `eip_list` - Instance ip.
    * `instance_id` - InstanceId.
  * `protocol` - TCP or UDP.
  * `real_servers` - Source server list.
    * `port` - 0~65535.
    * `real_server` - Source server addr, ip or domain.
    * `rs_type` - 1: domain, 2: ip.
    * `weight` - The return weight of the source station, ranging from 1 to 100.
* `layer7_listeners` - Layer 7 forwarding listener list.
  * `domain` - Domain.
  * `instance_detail_rule` - Resource instance to which the rule belongs.
    * `cname` - Cname.
    * `eip_list` - Instance ip list.
    * `instance_id` - Instance id.
  * `instance_details` - InstanceDetails.
    * `eip_list` - Instance ip list.
    * `instance_id` - Instance id.
  * `protocol` - Protocol.
  * `proxy_type_list` - List of forwarding types.
    * `proxy_ports` - Forwarding listening port list, port value is 1~65535.
    * `proxy_type` - Http, https.
  * `real_servers` - Source server list.
    * `port` - 0-65535.
    * `real_server` - Source server list.
    * `rs_type` - 1: domain, 2: ip.
    * `weight` - Weight: 1-100.
  * `vport` - Port.


