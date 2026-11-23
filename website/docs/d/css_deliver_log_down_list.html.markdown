---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_deliver_log_down_list"
sidebar_current: "docs-tencentcloud-datasource-css_deliver_log_down_list"
description: |-
  Use this data source to query detailed information of css deliver_log_down_list
---

# tencentcloud_css_deliver_log_down_list

Use this data source to query detailed information of css deliver_log_down_list

## Example Usage



## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `log_info_list` - List of log information.
  * `file_size` - File size, in bytes.
  * `log_name` - Log name.
  * `log_time` - Log time. UTC format, for example: 2018-11-29T19:00:00Z.Note:Beijing time is UTC time + 8 hours, formatted according to the ISO 8601 standard, see ISO date format description for details.
  * `log_url` - Log download address.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeDeliverLogDownList](https://www.tencentcloud.com/search/?wd=DescribeDeliverLogDownList&type=document&lang=en)

