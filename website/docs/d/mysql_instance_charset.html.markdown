---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_instance_charset"
sidebar_current: "docs-tencentcloud-datasource-mysql_instance_charset"
description: |-
  Use this data source to query detailed information of mysql instance_charset
---

# tencentcloud_mysql_instance_charset

Use this data source to query detailed information of mysql instance_charset

## Example Usage

```hcl
data "tencentcloud_mysql_instance_charset" "instance_charset" {
  instance_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page, and you can use the [query instance list] (https://cloud.tencent.com/document/api/236/15872) interface Gets the value of the field InstanceId in the output parameter.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `charset` - The default character set of the instance, such as `latin1`, `utf8` etc.


