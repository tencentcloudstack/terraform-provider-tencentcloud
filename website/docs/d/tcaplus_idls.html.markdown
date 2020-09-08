---
subcategory: "TcaplusDB"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcaplus_idls"
sidebar_current: "docs-tencentcloud-datasource-tcaplus_idls"
description: |-
  Use this data source to query  IDL information of the TcaplusDB table.
---

# tencentcloud_tcaplus_idls

Use this data source to query  IDL information of the TcaplusDB table.

## Example Usage

```hcl
data "tencentcloud_tcaplus_idls" "id_test" {
  cluster_id = "19162256624"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) Id of the TcaplusDB cluster to be query.
* `result_output_file` - (Optional) File for saving results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of TcaplusDB table IDL. Each element contains the following attributes.
  * `idl_id` - Id of the IDL.


