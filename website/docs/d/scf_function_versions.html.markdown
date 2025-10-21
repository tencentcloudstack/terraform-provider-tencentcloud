---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_function_versions"
sidebar_current: "docs-tencentcloud-datasource-scf_function_versions"
description: |-
  Use this data source to query detailed information of scf function_versions
---

# tencentcloud_scf_function_versions

Use this data source to query detailed information of scf function_versions

## Example Usage

```hcl
data "tencentcloud_scf_function_versions" "function_versions" {
  function_name = "keep-1676351130"
  namespace     = "default"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String) Function Name.
* `namespace` - (Optional, String) The namespace where the function locates.
* `order_by` - (Optional, String) It specifies the sorting order of the results according to a specified field, such as `AddTime`, `ModTime`.
* `order` - (Optional, String) It specifies whether to return the results in ascending or descending order. The value is `ASC` or `DESC`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `versions` - Function version listNote: This field may return null, indicating that no valid values is found.
  * `add_time` - The creation timeNote: This field may return null, indicating that no valid value was found.
  * `description` - Version descriptionNote: This field may return null, indicating that no valid values is found.
  * `mod_time` - Update timeNote: This field may return null, indicating that no valid value was found.
  * `status` - Version statusNote: this field may return `null`, indicating that no valid values can be obtained.
  * `version` - Function version name.


