---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_domain_analytics"
sidebar_current: "docs-tencentcloud-datasource-dnspod_domain_analytics"
description: |-
  Use this data source to query detailed information of dnspod domain_analytics
---

# tencentcloud_dnspod_domain_analytics

Use this data source to query detailed information of dnspod domain_analytics

## Example Usage

```hcl
data "tencentcloud_dnspod_domain_analytics" "domain_analytics" {
  domain     = "dnspod.cn"
  start_date = "2023-10-07"
  end_date   = "2023-10-12"
  dns_format = "HOUR"
  # domain_id = 123
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) The domain name to query for resolution volume.
* `end_date` - (Required, String) The end date of the query, format: YYYY-MM-DD.
* `start_date` - (Required, String) The start date of the query, format: YYYY-MM-DD.
* `dns_format` - (Optional, String) DATE: Statistics by day dimension HOUR: Statistics by hour dimension.
* `domain_id` - (Optional, Int) Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alias_data` - Domain alias resolution volume statistics information.
  * `data` - Subtotal of resolution volume for the current statistical dimension.
    * `date_key` - For daily statistics, it is the statistical date.
    * `hour_key` - For hourly statistics, it is the hour of the current time (0-23), for example, when HourKey is 23, the statistical period is the resolution volume from 22:00 to 23:00. Note: This field may return null, indicating that no valid value can be obtained.
    * `num` - Subtotal of resolution volume for the current statistical dimension.
  * `info` - Domain resolution volume statistics query information.
    * `dns_format` - DATE: Statistics by day dimension HOUR: Statistics by hour dimension.
    * `dns_total` - Total resolution volume for the current statistical period.
    * `domain` - The domain name currently being queried.
    * `end_date` - End time of the current statistical period.
    * `start_date` - Start time of the current statistical period.
* `data` - Subtotal of resolution volume for the current statistical dimension.
  * `date_key` - For daily statistics, it is the statistical date.
  * `hour_key` - For hourly statistics, it is the hour of the current time (0-23), for example, when HourKey is 23, the statistical period is the resolution volume from 22:00 to 23:00. Note: This field may return null, indicating that no valid value can be obtained.
  * `num` - Subtotal of resolution volume for the current statistical dimension.
* `info` - Domain resolution volume statistics query information.
  * `dns_format` - DATE: Statistics by day dimension HOUR: Statistics by hour dimension.
  * `dns_total` - Total resolution volume for the current statistical period.
  * `domain` - The domain name currently being queried.
  * `end_date` - End time of the current statistical period.
  * `start_date` - Start time of the current statistical period.


