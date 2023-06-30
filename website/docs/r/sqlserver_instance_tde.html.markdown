---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_instance_tde"
sidebar_current: "docs-tencentcloud-resource-sqlserver_instance_tde"
description: |-
  Provides a resource to create a sqlserver instance_tde
---

# tencentcloud_sqlserver_instance_tde

Provides a resource to create a sqlserver instance_tde

## Example Usage

```hcl
resource "tencentcloud_sqlserver_instance_tde" "instance_tde" {
  instance_id             = "mssql-qelbzgwf"
  certificate_attribution = "self"
}
```

## Argument Reference

The following arguments are supported:

* `certificate_attribution` - (Required, String) Certificate attribution. self- means to use the account's own certificate, others- means to refer to the certificate of other accounts, and the default is self.
* `instance_id` - (Required, String) Instance ID.
* `quote_uin` - (Optional, String) Other referenced main account IDs, required when CertificateAttribute is others.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver instance_tde can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_instance_tde.instance_tde instance_tde_id
```

