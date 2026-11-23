---
subcategory: "DNSPod"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_domain_instances"
sidebar_current: "docs-tencentcloud-datasource-dnspod_domain_instances"
description: |-
  Use this data source to query detailed information of dnspod domain instances
---

# tencentcloud_dnspod_domain_instances

Use this data source to query detailed information of dnspod domain instances

## Example Usage

```hcl
data "tencentcloud_dnspod_domain_instances" "example" {
  domain = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_list` - Domain list.
  * `create_time` - Create time of the domain.
  * `domain_id` - ID of the domain.
  * `domain` - The Domain.
  * `grade` - The DNS plan/package grade of the domain (e.g., DP_Free, DP_Plus).
  * `group_id` - The Group Id of Domain.
  * `id` - ID of the domain.
  * `is_mark` - Whether to Mark the Domain.
  * `record_count` - Number of DNS records under this domain.
  * `remark` - The remark of Domain.
  * `slave_dns` - Is secondary DNS enabled.
  * `status` - The status of Domain.
  * `updated_on` - Last modification time of the domain.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeDomain](https://www.tencentcloud.com/document/api/1157/49038)

