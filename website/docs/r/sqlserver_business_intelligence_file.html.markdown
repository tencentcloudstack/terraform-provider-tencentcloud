---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_business_intelligence_file"
sidebar_current: "docs-tencentcloud-resource-sqlserver_business_intelligence_file"
description: |-
  Provides a resource to create a sqlserver business_intelligence_file
---

# tencentcloud_sqlserver_business_intelligence_file

Provides a resource to create a sqlserver business_intelligence_file

## Example Usage

```hcl
resource "tencentcloud_sqlserver_business_intelligence_file" "business_intelligence_file" {
  instance_id = "mssql-wu8wka8a"
  file_u_r_l  = ""
  file_type   = "SSIS"
  remark      = ""
}
```

## Argument Reference

The following arguments are supported:

* `file_type` - (Required, String) File Type FLAT - Flat File as Data Source, SSIS - ssis project package.
* `file_url` - (Required, String) Cos Url.
* `instance_id` - (Required, String) instance id.
* `remark` - (Optional, String) remark.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver business_intelligence_file can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_business_intelligence_file.business_intelligence_file business_intelligence_file_id
```

