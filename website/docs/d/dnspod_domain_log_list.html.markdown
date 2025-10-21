---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_domain_log_list"
sidebar_current: "docs-tencentcloud-datasource-dnspod_domain_log_list"
description: |-
  Use this data source to query detailed information of dnspod domain_log_list
---

# tencentcloud_dnspod_domain_log_list

Use this data source to query detailed information of dnspod domain_log_list

## Example Usage

```hcl
data "tencentcloud_dnspod_domain_log_list" "domain_log_list" {
  domain    = "iac-tf.cloud"
  domain_id = 123
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain.
* `domain_id` - (Optional, Int) Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `log_list` - Domain Operation Log List. Note: This field may return null, indicating that no valid value can be obtained.


