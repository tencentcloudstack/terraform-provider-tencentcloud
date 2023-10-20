---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_record_line_list"
sidebar_current: "docs-tencentcloud-datasource-dnspod_record_line_list"
description: |-
  Use this data source to query detailed information of dnspod record_line_list
---

# tencentcloud_dnspod_record_line_list

Use this data source to query detailed information of dnspod record_line_list

## Example Usage

```hcl
data "tencentcloud_dnspod_record_line_list" "record_line_list" {
  domain       = "iac-tf.cloud"
  domain_grade = "DP_FREE"
  domain_id    = 123
}
```

## Argument Reference

The following arguments are supported:

* `domain_grade` - (Required, String) Domain level. + Old packages: D_FREE, D_PLUS, D_EXTRA, D_EXPERT, D_ULTRA correspond to free package, personal luxury, enterprise 1, enterprise 2, enterprise 3. + New packages: DP_FREE, DP_PLUS, DP_EXTRA, DP_EXPERT, DP_ULTRA correspond to new free, personal professional, enterprise basic, enterprise standard, enterprise flagship.
* `domain` - (Required, String) Domain.
* `domain_id` - (Optional, Int) Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `line_group_list` - Line group list.
  * `line_id` - Line group ID.
  * `line_list` - List of lines included in the line group.
  * `name` - Line group name.
  * `type` - Group type.
* `line_list` - Line list.
  * `line_id` - Line ID.
  * `name` - Line name.


