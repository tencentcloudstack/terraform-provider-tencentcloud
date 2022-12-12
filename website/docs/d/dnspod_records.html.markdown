---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_records"
sidebar_current: "docs-tencentcloud-datasource-dnspod_records"
description: |-
  Use this data source to query dnspod record list.
---

# tencentcloud_dnspod_records

Use this data source to query dnspod record list.

## Example Usage

```hcl
data "tencentcloud_dnspod_records" "record" {
  domain    = "example.com"
  subdomain = "www"
}

output "result" {
  value = data.tencentcloud_dnspod_records.record.result
}
```

Use verbose filter

```hcl
data "tencentcloud_dnspod_records" "record" {
  domain      = "example.com"
  subdomain   = "www"
  limit       = 100
  record_type = "TXT"
  sort_field  = "updated_on"
  sort_type   = "DESC"
}

output "result" {
  value = data.tencentcloud_dnspod_records.record.result
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Optional, String) The ID of the domain for which DNS records are to be obtained. If DomainId is passed in, the system will omit the parameter domain.
* `domain` - (Optional, String) The domain for which DNS records are to be obtained.
* `group_id` - (Optional, String) The group ID.
* `keyword` - (Optional, String) The keyword for searching for DNS records. Host headers and record values are supported.
* `limit` - (Optional, Int) The limit. It defaults to 100 and can be up to 3,000.
* `offset` - (Optional, Int) The offset. Default value: 0.
* `record_line_id` - (Optional, String) The split zone ID. If `record_line_id` is passed in, the system will omit the parameter `record_line`.
* `record_line` - (Optional, String) The split zone name.
* `record_type` - (Optional, String) The type of DNS record, such as A, CNAME, NS, AAAA, explicit URL, implicit URL, CAA, or SPF record.
* `result_output_file` - (Optional, String) Used for store query result as JSON.
* `sort_field` - (Optional, String) The sorting field. Available values: name, line, type, value, weight, mx, and ttl,updated_on.
* `sort_type` - (Optional, String) The sorting type. Valid values: ASC (ascending, default), DESC (descending).
* `subdomain` - (Optional, String) The host header of a DNS record. If this parameter is passed in, only the DNS record corresponding to this host header will be returned.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `record_count_info` - Count info of the queried record list.
  * `list_count` - The count of records returned in the list.
  * `subdomain_count` - The subdomain count.
  * `total_count` - The total record count.
* `result` - The record list result.
  * `line_id` - The split zone ID.
  * `line` - The record split zone.
  * `monitor_status` - The monitoring status of the record. Valid values: OK (normal), WARN (warning), and DOWN (downtime). It is empty if no monitoring is set or the monitoring is suspended.
  * `mx` - The MX value, applicable to the MX record only.
Note: This field may return null, indicating that no valid values can be obtained.
  * `name` - The host name.
  * `record_id` - Record ID.
  * `remark` - The record remarks.
  * `status` - The record status. Valid values: ENABLE (enabled), DISABLE (disabled).
  * `ttl` - The record cache time.
  * `type` - The record type.
  * `updated_on` - The update time.
  * `value` - The record value.
  * `weight` - The record weight, which is required for round-robin DNS records.


