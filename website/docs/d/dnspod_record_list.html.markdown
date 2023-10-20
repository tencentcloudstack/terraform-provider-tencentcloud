---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_record_list"
sidebar_current: "docs-tencentcloud-datasource-dnspod_record_list"
description: |-
  Use this data source to query detailed information of dnspod record_list
---

# tencentcloud_dnspod_record_list

Use this data source to query detailed information of dnspod record_list

## Example Usage

```hcl
data "tencentcloud_dnspod_record_list" "record_list" {
  domain = "iac-tf.cloud"
  # domain_id = 123
  # sub_domain = "www"
  record_type = ["A", "NS", "CNAME", "NS", "AAAA"]
  # record_line = [""]
  group_id            = []
  keyword             = ""
  sort_field          = "UPDATED_ON"
  sort_type           = "DESC"
  record_value        = "bicycle.dnspod.net"
  record_status       = ["ENABLE"]
  weight_begin        = 0
  weight_end          = 100
  mx_begin            = 0
  mx_end              = 10
  ttl_begin           = 1
  ttl_end             = 864000
  updated_at_begin    = "2021-09-07"
  updated_at_end      = "2023-12-07"
  remark              = ""
  is_exact_sub_domain = true
  # project_id = -1
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) The domain to which the resolution record belongs.
* `domain_id` - (Optional, Int) The domain ID to which the resolution record belongs. If DomainId is provided, the system will ignore the Domain parameter. You can find all Domain and DomainId through the DescribeDomainList interface.
* `group_id` - (Optional, Set: [`Int`]) When retrieving resolution records under certain groups, pass this group ID. You can obtain the GroupId field through the DescribeRecordGroupList interface.
* `is_exact_sub_domain` - (Optional, Bool) Whether to perform an exact search based on the SubDomain parameter.
* `keyword` - (Optional, String) Search for resolution records by keyword, currently supporting searching host headers and record values.
* `mx_begin` - (Optional, Int) The starting point of the resolution record MX priority query interval.
* `mx_end` - (Optional, Int) The endpoint of the resolution record MX priority query interval.
* `project_id` - (Optional, Int) Project ID.
* `record_line` - (Optional, Set: [`String`]) Retrieve resolution records for certain line IDs. You can view the allowed line information for the current domain through the DescribeRecordLineList interface.
* `record_status` - (Optional, Set: [`String`]) Get the resolution record based on the resolution record status. The possible values are ENABLE and DISABLE. ENABLE: Normal DISABLE: Paused.
* `record_type` - (Optional, Set: [`String`]) Retrieve certain types of resolution records, such as A, CNAME, NS, AAAA, explicit URL, implicit URL, CAA, SPF, etc.
* `record_value` - (Optional, String) Get the resolution record based on the resolution record value.
* `remark` - (Optional, String) Get the resolution record based on the resolution record remark.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_field` - (Optional, String) Sorting field, supporting NAME, LINE, TYPE, VALUE, WEIGHT, MX, TTL, UPDATED_ON fields. NAME: The host header of the resolution record LINE: The resolution record line TYPE: The resolution record type VALUE: The resolution record value WEIGHT: The weight MX: MX priority TTL: The resolution record cache time UPDATED_ON: The resolution record update time.
* `sort_type` - (Optional, String) Sorting method, ascending: ASC, descending: DESC. The default value is ASC.
* `sub_domain` - (Optional, String) Retrieve resolution records based on the host header of the resolution record. Fuzzy matching is used by default. You can set the IsExactSubdomain parameter to true for precise searching.
* `ttl_begin` - (Optional, Int) The starting point of the resolution record TTL query interval.
* `ttl_end` - (Optional, Int) The endpoint of the resolution record TTL query interval.
* `updated_at_begin` - (Optional, String) The starting point of the resolution record update time query interval.
* `updated_at_end` - (Optional, String) The endpoint of the resolution record update time query interval.
* `weight_begin` - (Optional, Int) The starting point of the resolution record weight query interval.
* `weight_end` - (Optional, Int) The endpoint of the resolution record weight query interval.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `record_count_info` - Statistics of the number of records.
  * `list_count` - Number of records returned in the list.
  * `subdomain_count` - Number of subdomains.
  * `total_count` - Total number of records.
* `record_list` - List of records.
  * `default_ns` - Whether it is the default NS record.
  * `line_id` - Line ID.
  * `line` - Record line.
  * `monitor_status` - Record monitoring status, normal: OK, alarm: WARN, downtime: DOWN, empty if monitoring is not set or paused.
  * `mx` - MX value, only available for MX records Note: This field may return null, indicating that no valid value can be obtained.
  * `name` - Host header.
  * `record_id` - Record ID.
  * `remark` - Record remark description.
  * `status` - Record status, enabled: ENABLE, paused: DISABLE.
  * `ttl` - Record cache time.
  * `type` - Record type.
  * `updated_on` - Update time.
  * `value` - Record value.
  * `weight` - Record weight, used for load balancing records. Note: This field may return null, indicating that no valid value can be obtained.


