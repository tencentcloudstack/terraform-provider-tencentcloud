---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_microservice_api_version"
sidebar_current: "docs-tencentcloud-datasource-tsf_microservice_api_version"
description: |-
  Use this data source to query detailed information of tsf microservice_api_version
---

# tencentcloud_tsf_microservice_api_version

Use this data source to query detailed information of tsf microservice_api_version

## Example Usage

```hcl
data "tencentcloud_tsf_microservice_api_version" "microservice_api_version" {
  microservice_id = "ms-yq3jo6jd"
  path            = ""
  method          = "get"
}
```

## Argument Reference

The following arguments are supported:

* `microservice_id` - (Required, String) Microservice ID.
* `method` - (Optional, String) request method.
* `path` - (Optional, String) api path.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - api version list.
  * `application_id` - Application ID.
  * `application_name` - Application Name.
  * `pkg_version` - application pkg version.


