---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_ssl"
sidebar_current: "docs-tencentcloud-resource-mysql_ssl"
description: |-
  Provides a resource to create a mysql ssl
---

# tencentcloud_mysql_ssl

Provides a resource to create a mysql ssl

## Example Usage

```hcl
resource "tencentcloud_mysql_ssl" "ssl" {
  instance_id = "cdb-j5rprr8n"
  status      = "OFF"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID. Example value: cdb-c1nl9rpv.
* `status` - (Required, String) Whether to enable SSL. `ON` means enabled, `OFF` means not enabled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `url` - The certificate download link. Example value: http://testdownload.url.


## Import

mysql ssl can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_ssl.ssl instanceId
```

