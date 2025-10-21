---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_function"
sidebar_current: "docs-tencentcloud-resource-wedata_function"
description: |-
  Provides a resource to create a wedata function
---

# tencentcloud_wedata_function

Provides a resource to create a wedata function

## Example Usage

```hcl
resource "tencentcloud_wedata_function" "example" {
  type               = "HIVE"
  kind               = "ANALYSIS"
  name               = "tf_example"
  cluster_identifier = "emr-m6u3qgk0"
  db_name            = "tf_db_example"
  project_id         = "1612982498218618880"
  class_name         = "tf_class_example"
  resource_list {
    path = "/wedata-demo-1314991481/untitled3-1.0-SNAPSHOT.jar"
    name = "untitled3-1.0-SNAPSHOT.jar"
    id   = "5b28bcdf-a0e6-4022-927d-927d399c4593"
    type = "cos"
  }
  description = "description."
  usage       = "usage info."
  param_desc  = "param info."
  return_desc = "return value info."
  example     = "example info."
  comment     = "V1"
}
```

## Argument Reference

The following arguments are supported:

* `class_name` - (Required, String) Class name of function entry.
* `cluster_identifier` - (Required, String) Cluster ID.
* `comment` - (Required, String) Comment.
* `db_name` - (Required, String) Database name.
* `description` - (Required, String) Description of the function.
* `example` - (Required, String) Example of the function.
* `kind` - (Required, String) Function Kind, Enum: ANALYSIS, ENCRYPTION, AGGREGATE, LOGIC, DATE_AND_TIME, MATH, CONVERSION, STRING, IP_AND_DOMAIN, WINDOW, OTHER.
* `name` - (Required, String) Function Name.
* `param_desc` - (Required, String) Description of the Parameter.
* `project_id` - (Required, String) Project ID.
* `resource_list` - (Required, List) Resource of the function, stored in WeData COS(.jar,...).
* `return_desc` - (Required, String) Description of the Return value.
* `type` - (Required, String) Function Type, Enum: HIVE, SPARK, DLC.
* `usage` - (Required, String) Usage of the function.

The `resource_list` object supports the following:

* `name` - (Required, String) Resource Name.
* `path` - (Required, String) Resource Path.
* `id` - (Optional, String) Resource ID.
* `md5` - (Optional, String) Resource MD5 Value.
* `type` - (Optional, String) Resource Type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `function_id` - Function ID.


