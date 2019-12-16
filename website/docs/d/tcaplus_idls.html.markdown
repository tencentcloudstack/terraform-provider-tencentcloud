---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_idls"
sidebar_current: "docs-tencentcloud-datasource-tcaplus_idls"
description: |-
  Use this data source to query tcaplus idl files
---

# tencentcloud_tcaplus_idls

Use this data source to query tcaplus idl files

## Example Usage

```hcl
data "tencentcloud_tcaplus_idls" "id_test" {
  app_id = "19162256624"
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) Id of the tcapplus application to be query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of tcaplus idls. Each element contains the following attributes.
  * `idl_id` - Id of this idl.


