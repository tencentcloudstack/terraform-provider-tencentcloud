---
subcategory: "DNSPOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnspod_record_analytics"
sidebar_current: "docs-tencentcloud-datasource-dnspod_record_analytics"
description: |-
  Use this data source to query detailed information of dnspod record_analytics
---

# tencentcloud_dnspod_record_analytics

Use this data source to query detailed information of dnspod record_analytics

## Example Usage

```hcl
data "tencentcloud_dnspod_record_analytics" "record_analytics" {
  domain     = "iac-tf.cloud"
  start_date = "2023-09-07"
  end_date   = "2023-11-07"
  subdomain  = "www"
  dns_format = "HOUR"
  # domain_id = 123
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) The domain to query for resolution volume.
* `end_date` - (Required, String) The end date of the query, format: YYYY-MM-DD.
* `start_date` - (Required, String) The start date of the query, format: YYYY-MM-DD.
* `subdomain` - (Required, String) The subdomain to query for resolution volume.
* `dns_format` - (Optional, String) DATE: Statistics by day dimension, HOUR: Statistics by hour dimension.
* `domain_id` - (Optional, Int) Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `alias_data` - Subdomain alias resolution statistics information.
  * `data` - The subtotal of the resolution volume for the current statistical dimension.
    * `date_key` - For daily statistics, it is the statistical date.
    * `hour_key` - For hourly statistics, it is the hour of the current time for statistics (0-23), e.g., when HourKey is 23, the statistical period is the resolution volume from 22:00 to 23:00. Note: This field may return null, indicating that no valid value can be obtained.
    * `num` - The subtotal of the resolution volume for the current statistical dimension.
  * `info` - Subdomain resolution statistics query information.
    * `dns_format` - DATE: Daily statistics, HOUR: Hourly statistics.
    * `dns_total` - Total resolution count for the current statistical period.
    * `domain` - The domain currently being queried.
    * `end_date` - End date of the current statistical period.
    * `start_date` - Start date of the current statistical period.
    * `subdomain` - The subdomain currently being analyzed.
* `data` - The subtotal of the resolution volume for the current statistical dimension.
  * `date_key` - For daily statistics, it is the statistical date.
  * `hour_key` - For hourly statistics, it is the hour of the current time for statistics (0-23), e.g., when HourKey is 23, the statistical period is the resolution volume from 22:00 to 23:00. Note: This field may return null, indicating that no valid value can be obtained.
  * `num` - The subtotal of the resolution volume for the current statistical dimension.
* `info` - Subdomain resolution statistics query information.
  * `dns_format` - DATE: Daily statistics, HOUR: Hourly statistics.
  * `dns_total` - Total resolution count for the current statistical period.
  * `domain` - The domain currently being queried.
  * `end_date` - End date of the current statistical period.
  * `start_date` - Start date of the current statistical period.
  * `subdomain` - The subdomain currently being analyzed.


