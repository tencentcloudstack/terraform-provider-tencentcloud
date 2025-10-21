---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_api_detail"
sidebar_current: "docs-tencentcloud-datasource-tsf_api_detail"
description: |-
  Use this data source to query detailed information of tsf api_detail
---

# tencentcloud_tsf_api_detail

Use this data source to query detailed information of tsf api_detail

## Example Usage

```hcl
data "tencentcloud_tsf_api_detail" "api_detail" {
  microservice_id = "ms-yq3jo6jd"
  path            = "/printRequest"
  method          = "GET"
  pkg_version     = "20210625192923"
  application_id  = "application-a24x29xv"
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required, String) application id.
* `method` - (Required, String) request method.
* `microservice_id` - (Required, String) microservice id.
* `path` - (Required, String) api path.
* `pkg_version` - (Required, String) pkg version.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - api detail.
  * `can_run` - can debug or not.
  * `definitions` - api data struct.
    * `name` - object name.
    * `properties` - object property list.
      * `description` - property description.
      * `name` - property name.
      * `type` - property type.
  * `description` - API description. Note: This field may return null, indicating that no valid value can be obtained.
  * `request_content_type` - api content type.
  * `request` - api request description.
    * `default_value` - default value.
    * `description` - param description.
    * `in` - param position.
    * `name` - param name.
    * `required` - require or not.
    * `type` - type.
  * `response` - api response.
    * `description` - param description.
    * `name` - param description.
    * `type` - param type.
  * `status` - API status 0: offline 1: online, default 0. Note: This section may return null, indicating that no valid value can be obtained.


