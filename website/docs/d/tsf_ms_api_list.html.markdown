---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_ms_api_list"
sidebar_current: "docs-tencentcloud-datasource-tsf_ms_api_list"
description: |-
  Use this data source to query detailed information of tsf ms_api_list
---

# tencentcloud_tsf_ms_api_list

Use this data source to query detailed information of tsf ms_api_list

## Example Usage

```hcl
data "tencentcloud_tsf_ms_api_list" "ms_api_list" {
  microservice_id = "ms-yq3jo6jd"
  search_word     = "echo"
}
```

## Argument Reference

The following arguments are supported:

* `microservice_id` - (Required, String) Microservice Id.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) search word, support  service name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - result list.
  * `content` - api list.
    * `description` - Method description. Note: This field may return null, indicating that no valid value was found.
    * `method` - api method.
    * `path` - api path.
    * `status` - API status. 0: offline, 1: online.Note: This field may return null, indicating that no valid value was found.
  * `total_count` - Quantity.


