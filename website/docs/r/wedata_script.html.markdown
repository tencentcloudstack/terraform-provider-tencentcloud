---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_script"
sidebar_current: "docs-tencentcloud-resource-wedata_script"
description: |-
  Provides a resource to create a wedata script
---

# tencentcloud_wedata_script

Provides a resource to create a wedata script

## Example Usage

```hcl
resource "tencentcloud_wedata_script" "example" {
  file_path           = "/datastudio/project/tf_example.sql"
  project_id          = "1470575647377821696"
  bucket_name         = "wedata-demo-1257305158"
  region              = "ap-guangzhou"
  file_extension_type = "sql"
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Optional, String) Cos bucket name.
* `file_extension_type` - (Optional, String) File Extension Type:jar, sql, zip, py, sh, txt, di, dg, pyspark, kjb, ktr, csv.
* `file_path` - (Optional, String) Cos file path:/datastudio/project/projectId/.
* `project_id` - (Optional, String) Project id.
* `region` - (Optional, String) Cos region.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `resource_id` - Resource ID.


## Import

wedata script can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_script.example 1470575647377821696#/datastudio/project/tf_example.sql#4147824b-7ba2-432b-8a8b-7e747594c926
```

