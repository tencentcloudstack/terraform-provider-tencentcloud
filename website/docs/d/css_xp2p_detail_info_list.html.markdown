---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_xp2p_detail_info_list"
sidebar_current: "docs-tencentcloud-datasource-css_xp2p_detail_info_list"
description: |-
  Use this data source to query detailed information of css xp2p_detail_info_list
---

# tencentcloud_css_xp2p_detail_info_list

Use this data source to query detailed information of css xp2p_detail_info_list

## Example Usage

```hcl
data "tencentcloud_css_xp2p_detail_info_list" "xp2p_detail_info_list" {
  query_time = "2023-11-01T14:55:01+08:00"
  type       = ["live"]
}
```

## Argument Reference

The following arguments are supported:

* `dimension` - (Optional, Set: [`String`]) The dimension parameter can be used to specify the dimension for the query. If this parameter is not passed, the query will default to stream-level data. If you pass this parameter, it will only retrieve data for the specified dimension. The available dimension currently supported is AppId dimension, which allows you to query data based on the application ID. Please note that the returned fields will be related to the specified dimension.
* `query_time` - (Optional, String) The UTC minute granularity query time for querying usage data for a specific minute is in the format: yyyy-mm-ddTHH:MM:00Z. Please refer to the link https://cloud.tencent.com/document/product/266/11732#I.For example, if the local time is 2019-01-08 10:00:00 in Beijing, the corresponding UTC time would be 2019-01-08T10:00:00+08:00.This query supports data from the past six months.
* `result_output_file` - (Optional, String) Used to save results.
* `stream_names` - (Optional, Set: [`String`]) The stream array can be used to specify the streams to be queried. If no stream is specified, the query will include all streams by default.
* `type` - (Optional, Set: [`String`]) The type array can be used to specify the type of media content to be queried. The two available options are live for live streaming and vod for video on demand. If no type is specified, the query will include both live and VOD content by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data_info_list` - P2P streaming statistical information.
  * `app_id` - AppId. Note: This field may return null, indicating that no valid value is available.
  * `cdn_bytes` - CDN traffic.
  * `online_people` - Online numbers.
  * `p2p_bytes` - P2P traffic.
  * `request_success` - Success numbers.
  * `request` - Request numbers.
  * `stream_name` - Stream ID.Note: This field may return null, indicating that no valid value is available.
  * `stuck_people` - People count.
  * `stuck_times` - Count.
  * `time` - The requested format for time in UTC with one-minute granularity is yyyy-mm-ddTHH:MM:SSZ. This format follows the ISO 8601 standard and is commonly used for representing timestamps in UTC. For more information and examples, you can refer to the link provided: https://cloud.tencent.com/document/product/266/11732#I.
  * `type` - Type, divided into two categories: live and vod.Note: This field may return null, indicating that no valid value is available.


