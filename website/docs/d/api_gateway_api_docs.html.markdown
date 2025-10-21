---
subcategory: "API GateWay(apigateway)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_api_gateway_api_docs"
sidebar_current: "docs-tencentcloud-datasource-api_gateway_api_docs"
description: |-
  Use this data source to query list information of api_gateway api_doc
---

# tencentcloud_api_gateway_api_docs

Use this data source to query list information of api_gateway api_doc

## Example Usage

```hcl
data "tencentcloud_api_gateway_api_docs" "my_api_doc" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `api_doc_list` - List of ApiDocs.
  * `api_doc_id` - Api Doc ID.
  * `api_doc_name` - Api Doc Name.
  * `api_doc_status` - Api Doc Status.


