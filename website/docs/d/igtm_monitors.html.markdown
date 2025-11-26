---
subcategory: "Intelligent Global Traffic Manager(IGTM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_igtm_monitors"
sidebar_current: "docs-tencentcloud-datasource-igtm_monitors"
description: |-
  Use this data source to query detailed information of IGTM monitors
---

# tencentcloud_igtm_monitors

Use this data source to query detailed information of IGTM monitors

## Example Usage

### Query all igtm monitors

```hcl
data "tencentcloud_igtm_monitors" "example" {}
```

### Query igtm monitors by filter

```hcl
data "tencentcloud_igtm_monitors" "example" {
  filters {
    name  = "MonitorId"
    value = ["12383"]
    fuzzy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Query filter conditions.
* `is_detect_num` - (Optional, Int) Whether to query detection count, 0 for no, 1 for yes.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Filter field name, supported MonitorName: monitor name; MonitorId: monitor ID.
* `value` - (Required, Set) Filter field values.
* `fuzzy` - (Optional, Bool) Whether to enable fuzzy query, only supports filter field name as domain.
When fuzzy query is enabled, Value maximum length is 1, otherwise Value maximum length is 5. (Reserved field, currently unused).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `monitor_data_set` - Monitor list.
  * `check_interval` - Detection period.
  * `check_protocol` - Detection protocol PING TCP HTTP HTTPS.
  * `continue_period` - Continuous period count.
  * `created_on` - Creation time.
  * `detect_num` - Detection count.
  * `detector_group_ids` - Monitoring node ID group.
  * `detector_style` - Monitoring node type.
AUTO INTERNAL OVERSEAS IPV6 ALL.
  * `enable_redirect` - Whether to enable 3xx redirect following ENABLED DISABLED.
  * `enable_sni` - Whether to enable SNI.
ENABLED DISABLED.
  * `fail_rate` - Failure rate upper limit 100.
  * `fail_times` - Failure count.
  * `host` - Detection host.
  * `monitor_id` - Detection rule ID.
  * `monitor_name` - Monitor name.
  * `packet_loss_rate` - Packet loss rate upper limit.
  * `path` - Detection path.
  * `ping_num` - Packet count.
  * `return_code_threshold` - Return value threshold.
  * `tcp_port` - TCP port.
  * `timeout` - Detection timeout.
  * `uin` - Owner user.
  * `updated_on` - Update time.


