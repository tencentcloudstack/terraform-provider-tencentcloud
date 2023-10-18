---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_domain_list"
sidebar_current: "docs-tencentcloud-datasource-dnspod_domain_list"
description: |-
  Use this data source to query detailed information of dnspod domain_list
---

# tencentcloud_dnspod_domain_list

Use this data source to query detailed information of dnspod domain_list

## Example Usage

```hcl
data "tencentcloud_dnspod_domain_list" "domain_list" {
  type               = "ALL"
  group_id           = [1]
  keyword            = ""
  sort_field         = "UPDATED_ON"
  sort_type          = "DESC"
  status             = ["PAUSE"]
  package            = [""]
  remark             = ""
  updated_at_begin   = "2021-05-01 03:00:00"
  updated_at_end     = "2024-05-10 20:00:00"
  record_count_begin = 0
  record_count_end   = 100
  project_id         = -1
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, String) Get domain names based on domain group type. Available values are ALL, MINE, SHARE, RECENT. ALL: All MINE: My domain names SHARE: Domain names shared with me RECENT: Recently operated domain names.
* `group_id` - (Optional, Set: [`Int`]) Get domain names based on domain group id, which can be obtained through the GroupId field in DescribeDomain or DescribeDomainList interface.
* `keyword` - (Optional, String) Get domain names based on keywords.
* `package` - (Optional, Set: [`String`]) Get domain names based on the package, which can be obtained through the Grade field in DescribeDomain or DescribeDomainList interface.
* `project_id` - (Optional, Int) Project ID.
* `record_count_begin` - (Optional, Int) The start point of the domain name&amp;#39;s record count query range.
* `record_count_end` - (Optional, Int) The end point of the domain name&amp;#39;s record count query range.
* `remark` - (Optional, String) Get domain names based on remark information.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_field` - (Optional, String) Sorting field. Available values are NAME, STATUS, RECORDS, GRADE, UPDATED_ON. NAME: Domain name STATUS: Domain status RECORDS: Number of records GRADE: Package level UPDATED_ON: Update time.
* `sort_type` - (Optional, String) Sorting type, ascending: ASC, descending: DESC.
* `status` - (Optional, Set: [`String`]) Get domain names based on domain status. Available values are ENABLE, LOCK, PAUSE, SPAM. ENABLE: Normal LOCK: Locked PAUSE: Paused SPAM: Banned.
* `updated_at_begin` - (Optional, String) The start time of the domain name&amp;#39;s update time to be obtained, such as &amp;#39;2021-05-01 03:00:00&amp;#39;.
* `updated_at_end` - (Optional, String) The end time of the domain name&amp;#39;s update time to be obtained, such as &amp;#39;2021-05-10 20:00:00&amp;#39;.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domain_list` - Domain list.
  * `cname_speedup` - Whether to enable CNAME acceleration, enabled: ENABLE, disabled: DISABLE.
  * `created_on` - Domain addition time.
  * `dns_status` - DNS settings status, error: DNSERROR, normal: empty string.
  * `domain_id` - Unique identifier assigned to the domain by the system.
  * `effective_dns` - Valid DNS assigned to the domain by the system.
  * `grade_level` - Sequence number corresponding to the domain package level.
  * `grade_title` - Package name.
  * `grade` - Domain package level code.
  * `group_id` - Group Id the domain belongs to.
  * `is_vip` - Whether it is a paid package.
  * `name` - Original format of the domain.
  * `owner` - Domain owner account.
  * `punycode` - Punycode encoded domain format.
  * `record_count` - Number of records under the domain.
  * `remark` - Domain remark description.
  * `search_engine_push` - Whether to enable search engine push optimization, YES: YES, NO: NO.
  * `status` - Domain status, normal: ENABLE, paused: PAUSE, banned: SPAM.
  * `tag_list` - Domain-related tag list Note: This field may return null, indicating that no valid value can be obtained.
    * `tag_key` - Tag key.
    * `tag_value` - Tag Value. Note: This field may return null, indicating that no valid value can be obtained.
  * `ttl` - Default TTL value for domain resolution records.
  * `updated_on` - Domain update time.
  * `vip_auto_renew` - Whether the domain has VIP auto-renewal enabled, YES: YES, NO: NO, DEFAULT: DEFAULT.
  * `vip_end_at` - Paid package expiration time.
  * `vip_start_at` - Paid package activation time.


