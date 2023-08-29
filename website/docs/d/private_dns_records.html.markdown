---
subcategory: "PrivateDNS"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_private_dns_records"
sidebar_current: "docs-tencentcloud-datasource-private_dns_records"
description: |-
  Use this data source to query detailed information of private dns records
---

# tencentcloud_private_dns_records

Use this data source to query detailed information of private dns records

## Example Usage

```hcl
data "tencentcloud_private_dns_records" "private_dns_record" {
  zone_id = "zone-xxxxxx"
  filters {
    name   = "Value"
    values = ["8.8.8.8"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Private zone id: zone-xxxxxx.
* `filters` - (Optional, List) Filter parameters (Value and RecordType filtering are supported).
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) Parameter name.
* `values` - (Required, Set) Parameter values.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `record_set` - Parse record list.
  * `created_on` - Record creation time.
  * `enabled` - Enabled. 0 meaning paused, 1 meaning senabled.
  * `extra` - Additional information.
  * `mx` - MX priority: required if the record type is MX. Value range: 5,10,15,20,30,40,50.
  * `record_id` - Record sid.
  * `record_type` - Record type, optional record type are: A, AAAA, CNAME, MX, TXT, PTR.
  * `record_value` - Record value.
  * `status` - Record status.
  * `sub_domain` - Subdomain name.
  * `ttl` - Record cache time, the smaller the value, the faster it takes effect. The value is 1-86400s. The default is 600.
  * `updated_on` - Record update time.
  * `weight` - Record weight, value is 1-100.
  * `zone_id` - Private zone id: zone-xxxxxx.


