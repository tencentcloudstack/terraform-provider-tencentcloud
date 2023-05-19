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
resource "tencentcloud_sqlserver_business_intelligence_instance" "business_intelligence_instance" {
  zone                = "ap-guangzhou-6"
  memory              = 4
  storage             = 20
  cpu                 = 2
  machine_type        = "CLOUD_PREMIUM"
  project_id          = 0
  subnet_id           = "subnet-dwj7ipnc"
  vpc_id              = "vpc-4owdpnwr"
  db_version          = "201603"
  security_group_list = []
  weekly              = [1, 2, 3, 4, 5, 6, 7]
  start_time          = "00:00"
  span                = 6
  instance_name       = "create_db_name"
}

resource "tencentcloud_sqlserver_business_intelligence_file" "business_intelligence_file" {
  instance_id = tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance.id
  file_url    = "https://keep-sqlserver-1308919341.cos.ap-guangzhou.myqcloud.com/test.xlsx"
  file_type   = "FLAT"
  remark      = "test case."
}
```

## Argument Reference

The following arguments are supported:

* `file_type` - (Required, String, ForceNew) File Type FLAT - Flat File as Data Source, SSIS - ssis project package.
* `file_url` - (Required, String, ForceNew) Cos Url.
* `instance_id` - (Required, String, ForceNew) instance id.
* `remark` - (Optional, String, ForceNew) remark.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver business_intelligence_file can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_business_intelligence_file.business_intelligence_file business_intelligence_file_id
```

