---
subcategory: "Cloud Automated Testing(CAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cat_node_groups"
sidebar_current: "docs-tencentcloud-datasource-cat_node_groups"
description: |-
  Use this data source to query detailed information of cat node groups
---

# tencentcloud_cat_node_groups

Use this data source to query detailed information of cat node groups

## Example Usage

```hcl
data "tencentcloud_cat_node_groups" "groups" {
  node_type       = [1]
  ip_type         = 1
  node_group_type = 1
}
```

## Argument Reference

The following arguments are supported:

* `district_id` - (Optional, Int) Province or country ID. 0 means all. Defaults to 0 if not set.
* `ip_type` - (Optional, Int) IP type. 0: all, 1: IPv4, 2: IPv6. Defaults to 0 if not set.
* `name` - (Optional, String) Probe node description keyword.
* `net_service_id` - (Optional, Int) ISP ID. 0: all, 1: China Telecom, 2: China Unicom, 3: China Mobile, 99: others. Defaults to 0 if not set.
* `node_group_type` - (Optional, Int) Node group type. 0: advanced probe group, 1: availability node, 2: my probe group. Defaults to 0 if not set.
* `node_type` - (Optional, List: [`Int`]) Node type. 0: all, 1: IDC, 2: LastMile, 3: Mobile. Defaults to 0 if not set.
* `probe_type` - (Optional, Int) Test type, including scheduled test and instant test. 0-scheduled probe, others mean instant probe.
* `region_id` - (Optional, Int) Region ID. 0: selected probe points, 1: Chinese Mainland, 2: Hong Kong, Macao and Taiwan, 3: Asia Pacific, 4: Europe and America, 5: Africa and Oceania. Defaults to 0 if not set.
* `result_output_file` - (Optional, String) Used to save results.
* `task_category` - (Optional, Int) Node category. 0: all, 1: PC, 2: Mobile. Defaults to 0 if not set.
* `task_type` - (Optional, Int) Task type, such as 1, 2, 3, 4, 5, 6, 7. 1-page performance, 2-file upload, 3-file download, 4-port performance, 5-network quality, 6-audio and video experience, 7-domain whois. Defaults to 0 if not set, no filtering by task type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `district_list` - Province or country list.
  * `id` - Province (international) or ISP ID.
  * `name` - Name.
* `net_service_list` - ISP list.
  * `id` - Province (international) or ISP ID.
  * `name` - Name.
* `node_list` - Tree node list, two levels in total.
  * `children` - Child nodes.
    * `children` - Node list.
      * `content` - Node name.
      * `id` - Node code.
    * `content` - Child node name.
    * `id` - Child node ID.
  * `content` - Node name.
  * `id` - Node ID.


