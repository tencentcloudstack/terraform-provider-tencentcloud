---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_ssl"
sidebar_current: "docs-tencentcloud-resource-mysql_ssl"
description: |-
  Provides a resource to create a MySQL SSL
---

# tencentcloud_mysql_ssl

Provides a resource to create a MySQL SSL

## Example Usage

### For mysql instance SSL

```hcl
resource "tencentcloud_mysql_ssl" "example" {
  instance_id = "cdb-j5rprr8n"
  status      = "OFF"
}
```

### For mysql RO group SSL

```hcl
resource "tencentcloud_mysql_ssl" "example" {
  ro_group_id = "cdbrg-k9a6gup3"
  status      = "ON"
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Required, String) Whether to enable SSL. `ON` means enabled, `OFF` means not enabled.
* `instance_id` - (Optional, String, ForceNew) Instance ID. Example value: cdb-c1nl9rpv.
* `ro_group_id` - (Optional, String, ForceNew) RO group ID. Example value: cdbrg-k9a6gup3.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `url` - The certificate download link. Example value: http://testdownload.url.


## Import

MySQL SSL can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_ssl.example cdb-j5rprr8n
```

Or

```
terraform import tencentcloud_mysql_ssl.example cdbrg-k9a6gup3
```

