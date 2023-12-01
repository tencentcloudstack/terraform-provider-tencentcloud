---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_backup_stream"
sidebar_current: "docs-tencentcloud-datasource-css_backup_stream"
description: |-
  Use this data source to query detailed information of css backup_stream
---

# tencentcloud_css_backup_stream

Use this data source to query detailed information of css backup_stream

## Example Usage

```hcl
data "tencentcloud_css_backup_stream" "backup_stream" {
  stream_name = "live"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.
* `stream_name` - (Optional, String) Stream id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `stream_info_list` - Backup stream group info.
  * `backup_list` - Backup stream info.
    * `app_name` - Push path.
    * `domain_name` - Push domain.
    * `master_flag` - Master stream flag.
    * `publish_time` - UTC time, eg, 2018-06-29T19:00:00Z.
    * `source_from` - Source from.
    * `upstream_sequence` - Push stream sequence.
  * `host_group_name` - Group name.
  * `optimal_enable` - Optimal switch, 1-enable, 0-disable.
  * `stream_name` - Stream name.


