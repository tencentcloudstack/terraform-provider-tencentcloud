---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_firewall_rules_template"
sidebar_current: "docs-tencentcloud-datasource-lighthouse_firewall_rules_template"
description: |-
  Use this data source to query detailed information of lighthouse firewall_rules_template
---

# tencentcloud_lighthouse_firewall_rules_template

Use this data source to query detailed information of lighthouse firewall_rules_template

## Example Usage

```hcl
data "tencentcloud_lighthouse_firewall_rules_template" "firewall_rules_template" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `firewall_rule_set` - Firewall rule details list.
  * `action` - Valid values are (ACCEPT, DROP). Default value is ACCEPT.
  * `app_type` - Application type. Valid values are custom, HTTP (80), HTTPS (443), Linux login (22), Windows login (3389), MySQL (3306), SQL Server (1433), all TCP ports, all UDP ports, Ping-ICMP, ALL.
  * `cidr_block` - IP range or IP (mutually exclusive). Default value is 0.0.0.0/0, which indicates all sources.
  * `firewall_rule_description` - Firewall rule description.
  * `port` - Port. Valid values are ALL, one single port, multiple ports separated by commas, or port range indicated by a minus sign.
  * `protocol` - Protocol. Valid values are TCP, UDP, ICMP, ALL.


